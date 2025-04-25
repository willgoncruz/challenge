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
var auditActionChannel chan (model.Action)

func init() {
	reset()
}

func reset() {
	auditActionChannel = make(chan model.Action)
	book = &ledger{
		actions: []model.Action{},
	}
	writeLoop()
}

func writeLoop() {
	go func() { // Async write to ledger, in the order of channel writing
		for newAction := range auditActionChannel {
			book.actions = append(book.actions, newAction)
		}
	}()
}

func Clear() {
	book.actions = book.actions[:0]
}

func Audit(order model.Order, action model.ActionType) {
	log.Printf("New Audit Action: %+v %s", order, action)
	newAction := model.Action{
		ID:        order.ID,
		Action:    action,
		Timestamp: time.Now().UnixMicro(),
	}

	auditActionChannel <- newAction
}

// Return all the saved actions
func Retrieve() []model.Action {
	return book.actions
}
