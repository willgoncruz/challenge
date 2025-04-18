package ledger

import (
	"challenge/model"
	"time"
)

type ledger struct {
	actions []model.Action
}

var book *ledger

func init() {
	book = &ledger{
		actions: []model.Action{},
	}
}

func Audit(order model.Order, action model.ActionType) {
	// Ledger audit can happen async
	go func() {
		newAction := model.Action{
			ID:        order.ID,
			Action:    action,
			Timestamp: time.Now().UnixMicro(),
		}

		book.actions = append(book.actions, newAction)
	}()
}

// Return all the saved actions
func Retrieve() []model.Action {
	return book.actions
}
