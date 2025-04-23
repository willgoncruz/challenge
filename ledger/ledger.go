package ledger

import (
	"challenge/model"
	"log"
	"time"
)

type ledger struct {
	actions []model.Action
}

var book *ledger

func init() {
	reset()
}

func reset() {
	book = &ledger{
		actions: []model.Action{},
	}
}

func Audit(order model.Order, action model.ActionType) {
	log.Printf("New Audit Action: %+v %s", order, action)
	newAction := model.Action{
		ID:        order.ID,
		Action:    action,
		Timestamp: time.Now().UnixMicro(),
	}

	book.actions = append(book.actions, newAction)
}

// Return all the saved actions
func Retrieve() []model.Action {
	return book.actions
}
