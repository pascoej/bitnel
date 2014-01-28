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
		stmt, err = tx.Prepare(`SELECT uuid, price, size, initial_size FROM orders
			WHERE status = $1 OR status = $2 AND side = $3 AND price <= $4
			ORDER BY price ASC, created_at ASC`)
	} else {
		stmt, err = tx.Prepare(`SELECT uuid, price, size,initial_size FROM orders
			WHERE status = $1 OR status = $2  AND side = $3 AND price >= $4
			ORDER BY price ASC, created_at ASC`)
	}

	if err != nil {
		return &matchingError{o}
	}

	rows, err := stmt.Query(model.OpenStatus, model.PartiallyFilledStatus, (*o.Side).CounterSide(), o.Price)
	if err != nil {
		return &matchingError{o}
	}
	var counterOrders []model.Order
	// Keep matching until no more orders
	for rows.Next() && *o.Size > money.Unit(0) {
		var counterOrder model.Order
		err = rows.Scan(&counterOrder.Uuid, &counterOrder.Price, &counterOrder.Size, &counterOrder.InitialSize)
		if err != nil {
			return &matchingError{o}
		}
		if *o.Size > *counterOrder.Size {
			counterOrder.Status = model.CompletedStatus

			o.Status = model.PartiallyFilledStatus
			*o.Size = *o.Size - *counterOrder.Size
			*counterOrder.Size = money.Unit(0)
		} else { // matched order gets totally filled
			o.Status = model.CompletedStatus
			counterOrder.Status = model.PartiallyFilledStatus
			// however
			if *counterOrder.Size = *counterOrder.Size - *o.Size; *counterOrder.Size == money.Unit(0) {
				counterOrder.Status = model.CompletedStatus
			}
			*o.Size = money.Unit(0)
		}
		counterOrders = append(counterOrders, counterOrder)
	}
	rows.Close()
	stmt, err = tx.Prepare(`UPDATE orders
		SET size = $1, status = $2
		WHERE uuid = $3`)
	if err != nil {
		return &matchingError{o}
	}
	for _, counterOrder := range counterOrders {
		_, err = stmt.Exec(*counterOrder.Size, counterOrder.Status, counterOrder.Uuid)
		if err != nil {
			return &matchingError{o}
		}
	}
	if *o.Size == o.InitialSize { // order did not get filled
		o.Status = model.OpenStatus
	}
	_, err = stmt.Exec(*o.Size, o.Status, o.Uuid)
	if err != nil {
		return &matchingError{o}
	}
	tx.Commit()
	return nil
}

/*

	stmt err := tx.Prepare(`UPDATE orders
		SET size = 0, status = $1
		WHERE uuid = $2`)
	if err != nil {
		return &matchingError{o}
	}

	// our original order eats the smaller order
	if o.Size > counterOrder.Size {
		counterOrder.Status = model.CompletedStatus

		_, err := stmt.Exec(model.CompletedStatus, order.Uuid)
		if err != nil {
			return &matchingError{o}
		}

		o.Size = o.Size - counterOrder.Size
	} else { // smaller order is too bigggg
		o.Status = model.CompletedStatus

		_, err := stmt.Exec(model.CompletedStatus, o.Uuid)
		if err != nil {
			return &matchingError{o}
		}

		if counterOrder.Size = counterOrder.Size - o.Size; counterOrder.Size == 0 {
			counterOrder.Status = model.CompletedStatus

			_, err := stmt.Exec(model.CompletedStatus, counterOrder.Uuid)
			if err != nil {
				return &matchingError{o}
			}
		}

	}*/

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
