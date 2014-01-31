package matching

import (
	"database/sql"
	"fmt"
	"github.com/bitnel/bitnel-api/model"
	"github.com/bitnel/bitnel-api/money"
	"log"
	"time"
)

// the matching order thing!!!!!!!!!!!

type matchingError struct {
	order *model.Order
}

func (me *matchingError) Error() string {
	return fmt.Sprintf("matching: unable to match order uuid %s", me.order.Uuid)
}

const (
	feeFraction = 0.02
)

type Engine struct {
	orderNotifier chan *model.Order
	database      *sql.DB
}

func NewEngine(db *sql.DB, bufferSize int) *Engine {
	return &Engine{
		orderNotifier: make(chan *model.Order, bufferSize),
		database:      db,
	}
}

// starts the matching engine
func (m *Engine) Start() {
	go m.listen()
}

// adds an order to be matched
// returns no error on succesfull add
// otherwise fail miserably
func (m *Engine) Add(o *model.Order) *matchingError {
	select {
	case m.orderNotifier <- o:
		return nil
	default:
		return &matchingError{o}
	}
}

func (m *Engine) match(o *model.Order) *matchingError {
	tx, err := m.database.Begin()
	if err != nil {
		return &matchingError{o}
	}

	var stmt *sql.Stmt

	// maybe have more orders later on
	if *o.Side == model.BidSide {
		stmt, err = tx.Prepare(`SELECT uuid, price, size, initial_size, status, side FROM orders
			WHERE (status = $1 OR status = $2) AND side = $3 AND price <= $4
			ORDER BY price ASC, created_at ASC`)
	} else if *o.Side == model.AskSide {
		stmt, err = tx.Prepare(`SELECT uuid, price, size,initial_size, status, side FROM orders
			WHERE (status = $1 OR status = $2) AND side = $3 AND price >= $4
			ORDER BY price DESC, created_at ASC`)
	} else {
		return &matchingError{o}
	}
	if err != nil {
		return &matchingError{o}
	}

	rows, err := stmt.Query(model.OpenStatus, model.PartiallyFilledStatus, (*o.Side).CounterSide(), o.Price)
	if err != nil {
		return &matchingError{o}
	}

	var counterOrders []model.Order
	var trades map[string]model.Trade = make(map[string]model.Trade) // The counter order uuid is the key
	for rows.Next() && *o.Size > money.Unit(0) {
		var counterOrder model.Order
		err = rows.Scan(&counterOrder.Uuid, &counterOrder.Price, &counterOrder.Size, &counterOrder.InitialSize, &counterOrder.Status, &counterOrder.Side)
		if err != nil {
			return &matchingError{o}
		}
		var price money.Unit
		if *o.Side == model.AskSide {
			if *counterOrder.Price > *o.Price {
				price = *counterOrder.Price
			} else {
				price = *o.Price
			}
		} else {
			if *counterOrder.Price < *o.Price {
				price = *counterOrder.Price
			} else {
				price = *o.Price
			}
		}
		var trade model.Trade
		var amount money.Unit
		if *o.Size > *counterOrder.Size {
			counterOrder.Status = model.CompletedStatus
			o.Status = model.PartiallyFilledStatus
			amount = *counterOrder.Size
			*o.Size = *o.Size - *counterOrder.Size
			*counterOrder.Size = money.Unit(0)
		} else { // matched order gets totally filled
			o.Status = model.CompletedStatus
			counterOrder.Status = model.PartiallyFilledStatus
			// however
			if *counterOrder.Size = *counterOrder.Size - *o.Size; *counterOrder.Size == money.Unit(0) {
				counterOrder.Status = model.CompletedStatus
			}
			amount = *o.Size
			*o.Size = money.Unit(0)
		}
		trade.Amount = amount
		trade.Price = price
		trades[counterOrder.Uuid] = trade
		counterOrders = append(counterOrders, counterOrder)
	}
	rows.Close()

	orderStmt, err := tx.Prepare(`UPDATE orders
		SET size = $1, status = $2
		WHERE uuid = $3`)
	if err != nil {
		return &matchingError{o}
	}
	tradeStmt, err := tx.Prepare(`INSERT INTO trades (amount,price) RETURNING uuid`)
	if err != nil {
		return &matchingError{o}
	}
	transactionStmt, err := tx.Prepare(`INSERT INTO transactions (balance_uuid,type,amount,fee_amount, trade)  VALUES($1,$2,$3,$4,$5)`)
	if err != nil {
		return &matchingError{o}
	}
	reservedBalanceStmt, err := tx.Prepare(`UPDATE balances SET reserved_balance = reserved_balance+$1 WHERE user_uuid = $2 AND currency = $3 RETURNING uuid`)
	if err != nil {
		return &matchingError{o}
	}
	balanceStmt, err := tx.Prepare(`UPDATE balances SET balance = balance+$1 WHERE user_uuid = $2 AND currency = $3 RETURNING uuid`)
	if err != nil {
		return &matchingError{o}
	}
	var market model.Market
	marketStmt, err := tx.Prepare(`SELECT base_currency,quote_currency,currency_pair FROM markets WHERE uuid = $1`)
	if err != nil {
		return &matchingError{o}
	}
	if err = marketStmt.QueryRow(o.MarketUuid).Scan(&market.BaseCurrency, &market.QuoteCurrency, &market.CurrencyPair); err != nil {
		return &matchingError{o}
	}
	for _, counterOrder := range counterOrders {
		_, err = orderStmt.Exec(*counterOrder.Size, counterOrder.Status, counterOrder.Uuid)
		if err != nil {
			return &matchingError{o}
		}
		trade := trades[counterOrder.Uuid]
		if err = tradeStmt.QueryRow(trade.Amount, trade.Price).Scan(&trade.Uuid); err != nil {
			return &matchingError{o}
		}
		var transaction model.Transaction
		if *counterOrder.Side == model.AskSide {
			transaction.Amount = trade.Amount * trade.Price
		} else {
			transaction.Amount = -1 * trade.Amount * trade.Price
		}
		transaction.Type = model.TradeTransaction
		transaction.Trade = &counterOrder.Uuid
		transaction.FeeAmount = (transaction.Amount / (1 / feeFraction))
		var buyerAccountUuid string
		var sellerAccountUuid string
		if *o.Side == model.AskSide {
			sellerAccountUuid = *o.AccountUuid
			buyerAccountUuid = *counterOrder.AccountUuid
		} else {
			sellerAccountUuid = *counterOrder.AccountUuid
			buyerAccountUuid = *o.AccountUuid
		}
		var sellerBalance model.Balance
		if err = balanceStmt.QueryRow(transaction.Amount, sellerAccountUuid, market.QuoteCurrency).Scan(&sellerBalance.Uuid); err != nil {
			return &matchingError{o}
		}
		_, err := transactionStmt.Exec(sellerBalance.Uuid, model.TradeTransaction, transaction.Amount, transaction.FeeAmount, trade.Uuid)
		if err != nil {
			return &matchingError{o}
		}
		var buyerBalance model.Balance
		if err = reservedBalanceStmt.QueryRow(transaction.GetTotalAmount(), buyerAccountUuid, market.BaseCurrency).Scan(&buyerBalance.Uuid); err != nil {
			return &matchingError{o}
		}
		_, err = transactionStmt.Exec(buyerBalance.Uuid, model.TradeTransaction, transaction.Amount, transaction.FeeAmount, trade.Uuid)
		if err != nil {
			return &matchingError{o}
		}
	}

	if *o.Size == o.InitialSize { // order did not get filled
		o.Status = model.OpenStatus
	}

	_, err = orderStmt.Exec(*o.Size, o.Status, o.Uuid)
	if err != nil {
		return &matchingError{o}
	}

	if err = tx.Commit(); err != nil {
		return &matchingError{o}
	}

	return nil
}

func (m *Engine) listen() {
	for {
		order := <-m.orderNotifier

		err := m.match(order)
		for err != nil {
			log.Println(err)

			// wait 100 miliseconds and try again
			time.Sleep(time.Millisecond * 100)
			err = m.match(order)
		}

		// yay succesfull match
	}
}
