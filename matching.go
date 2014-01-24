package main

import (
	"github.com/bitnel/bitnel-api/model"
	"github.com/bitnel/bitnel-api/money"
	"fmt"
	"time"
)

// the matching order thing!!!!!!!!!!!

type matchingError struct {
	order *model.Order
}

func (me *matchingError) Error() string {
	return fmt.Sprintln
}

type matchingEngine struct {
	orderNotifer chan *model.Order
}

func newMatchingEngine() *matchingEngine {
	return &matchingEngine{
		orderNotifer: make(chan *model.Order)
	}
}

func (m *matchingEngine) add(o *model.Order) {
	m.orderNotifier <- o
}

func (m *matchingEngine) start() {
	go m.listen()
}

func (m *matchingEngine) match(o *model.Order) *matchingError {
		tx, err := db.Begin()
		if err != nil {
			return &matchingError{o}
		}

		stmt, err := tx.Prepare(``)
		if err != nil {
			return &matchingError{o}
		}
}

func (m *matchingEngine) listen() {
	for {
		order := <-m.order

		// rerun if error
		for err := m.match(order); err != nil {
			time.Sleep(time.Millisecond * 100)
		}
	}
}