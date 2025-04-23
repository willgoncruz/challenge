package ledger

import (
	"challenge/model"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_AsyncAuditLedgerWriting(t *testing.T) {
	reset()

	Audit(CommonOrder("new-id", model.Room), model.Move)
	time.Sleep(time.Millisecond)
	assert.Equal(t, 1, len(book.actions))
}

func CommonOrder(id string, temp model.Temperature) model.Order {
	return model.Order{
		ID:        id,
		Name:      "test order",
		Temp:      temp,
		Freshness: 120,
	}
}
