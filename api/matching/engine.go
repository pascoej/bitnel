package matching

import (
	"database/sql"
	"fmt"
	"github.com/bitnel/api/model"
	"github.com/bitnel/api/money"
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
	requestNotifier chan interface{}
	database        *sql.DB
}

func NewEngine(db *sql.DB, bufferSize int) *Engine {
	return &Engine{
		requestNotifier: make(chan interface{}, bufferSize),
		database:        db,
	}
}

// starts the matching engine
func (m *Engine) Start() {
	go m.listen()
}

type CancelRequest struct {
	order *model.Order
}

type AddRequest struct {
	order *model.Order
}

// adds an order to be matched
// returns no error on succesfull add
// otherwise fail miserably
func (m *Engine) Add(o *model.Order) *matchingError {
	select {
	case m.requestNotifier <- &AddRequest{o}:
		return nil
	default:
		return &matchingError{o}
	}
}

func (m *Engine) Cancel(o *model.Order) *matchingError {
	select {
	case m.requestNotifier <- &CancelRequest{o}:
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
		stmt, err = tx.Prepare(`SELECT uuid, price, size, initial_size, status, side, account_uuid FROM orders
			WHERE (status = $1 OR status = $2) AND side = $3 AND price <= $4
			ORDER BY price ASC, created_at ASC`)
	} else if *o.Side == model.AskSide {
		stmt, err = tx.Prepare(`SELECT uuid, price, size,initial_size, status, side,account_uuid FROM orders
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
		err = rows.Scan(&counterOrder.Uuid, &counterOrder.Price, &counterOrder.Size, &counterOrder.InitialSize, &counterOrder.Status, &counterOrder.Side, &counterOrder.AccountUuid)
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
	tradeStmt, err := tx.Prepare(`INSERT INTO trades (amount,price) VALUES($1,$2) RETURNING uuid`)
	if err != nil {
		return &matchingError{o}
	}
	transactionStmt, err := tx.Prepare(`INSERT INTO transactions (balance_uuid,type,amount,fee_amount, trade)  VALUES($1,$2,$3,$4,$5)`)
	if err != nil {
		return &matchingError{o}
	}
	reservedBalanceStmt, err := tx.Prepare(`UPDATE balances SET reserved_balance = reserved_balance+$1 WHERE account_uuid = $2 AND currency = $3 RETURNING uuid`)
	if err != nil {
		return &matchingError{o}
	}
	balanceStmt, err := tx.Prepare(`UPDATE balances SET available_balance = available_balance+$1 WHERE account_uuid = $2 AND currency = $3 RETURNING uuid`)
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
		//The seller GETTING the quote currency.
		var sellerQuoteTransaction model.Transaction
		sellerQuoteTransaction.Amount = trade.Amount * trade.Price
		sellerQuoteTransaction.Type = model.TradeTransaction
		sellerQuoteTransaction.Trade = &counterOrder.Uuid
		sellerQuoteTransaction.FeeAmount = (sellerQuoteTransaction.Amount / (1 / feeFraction))
		//The seller LOSING the base currency
		var sellerBaseTransaction model.Transaction
		sellerBaseTransaction.Amount = trade.Amount
		sellerBaseTransaction.Type = model.TradeTransaction
		sellerBaseTransaction.Trade = &counterOrder.Uuid
		sellerBaseTransaction.FeeAmount = 0
		//The buying getting the base currency
		var buyerBaseTransaction model.Transaction
		buyerBaseTransaction.Amount = trade.Amount
		buyerBaseTransaction.Type = model.TradeTransaction
		buyerBaseTransaction.Trade = &counterOrder.Uuid
		buyerBaseTransaction.FeeAmount = (buyerBaseTransaction.Amount / (1 / feeFraction))
		//The buyer losing the quote currency
		var buyerQuoteTransaction model.Transaction
		buyerQuoteTransaction.Amount = trade.Amount * trade.Price
		buyerQuoteTransaction.Type = model.TradeTransaction
		buyerQuoteTransaction.Trade = &counterOrder.Uuid
		buyerQuoteTransaction.FeeAmount = 0
		var buyerAccountUuid string
		var sellerAccountUuid string
		if *o.Side == model.AskSide {
			sellerAccountUuid = *o.AccountUuid
			buyerAccountUuid = *counterOrder.AccountUuid
		} else {
			sellerAccountUuid = *counterOrder.AccountUuid
			buyerAccountUuid = *o.AccountUuid
		}

		//Editng seller's quote balance (FEE HERE)
		var sellerQuoteBalance model.Balance
		if err = balanceStmt.QueryRow(sellerQuoteTransaction.GetAmountAfterFee(), sellerAccountUuid, market.QuoteCurrency).Scan(&sellerQuoteBalance.Uuid); err != nil {
			return &matchingError{o}
		}
		_, err := transactionStmt.Exec(sellerQuoteBalance.Uuid, model.TradeTransaction, sellerQuoteTransaction.Amount, sellerQuoteTransaction.FeeAmount, trade.Uuid)
		if err != nil {
			return &matchingError{o}
		}
		//Editing seller's base (reserved) balance (NO FEE)
		var sellerBaseBalance model.Balance
		if err = reservedBalanceStmt.QueryRow(sellerBaseTransaction.Amount*-1, sellerAccountUuid, market.BaseCurrency).Scan(&sellerBaseBalance.Uuid); err != nil {
			return &matchingError{o}
		}
		_, err = transactionStmt.Exec(sellerBaseBalance.Uuid, model.TradeTransaction, sellerBaseTransaction.Amount, 0, trade.Uuid)
		if err != nil {
			return &matchingError{o}
		}
		//Editng bufer's base balance (FEE HERE)
		var buyerBaseBalance model.Balance
		if err = balanceStmt.QueryRow(buyerBaseTransaction.GetAmountAfterFee(), buyerAccountUuid, market.BaseCurrency).Scan(&buyerBaseBalance.Uuid); err != nil {
			return &matchingError{o}
		}
		_, err = transactionStmt.Exec(buyerBaseBalance.Uuid, model.TradeTransaction, buyerBaseTransaction.Amount, buyerBaseTransaction.FeeAmount, trade.Uuid)
		if err != nil {
			return &matchingError{o}
		}
		//editing buyer's quote balance (NO FEE)
		var buyerQuoteBalance model.Balance
		if err = reservedBalanceStmt.QueryRow(buyerQuoteTransaction.Amount*-1, buyerAccountUuid, market.QuoteCurrency).Scan(&buyerQuoteBalance.Uuid); err != nil {
			return &matchingError{o}
		}
		_, err = transactionStmt.Exec(buyerQuoteBalance.Uuid, model.TradeTransaction, buyerQuoteTransaction.Amount, 0, trade.Uuid)
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

func (m *Engine) cancel(o *model.Order) *matchingError {
	tx, err := m.database.Begin()
	if err != nil {
		return &matchingError{o}
	}

	stmt, err := tx.Prepare(`UPDATE orders SET status = $1 WHERE uuid = $2`)
	if err != nil {
		return &matchingError{o}
	}

	if _, err = stmt.Exec(model.CanceledStatus, o.Uuid); err != nil {
		return &matchingError{o}
	}

	if err = tx.Commit(); err != nil {
		return &matchingError{o}
	}

	return nil
}

func (m *Engine) listen() {
	for {
		req := <-m.requestNotifier

		switch req := req.(type) {
		case *AddRequest:
			err := m.match(req.order)
			for err != nil {
				log.Println(err)

				// wait 100 miliseconds and try again
				time.Sleep(time.Millisecond * 100)
				err = m.match(req.order)
			}
		case *CancelRequest:
			err := m.cancel(req.order)
			for err != nil {
				log.Println(err)

				// wait 100 miliseconds and try again
				time.Sleep(time.Millisecond * 100)
				err = m.match(req.order)
			}
		}
	}
}
