package model

import (
	"container/heap"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_MinHeap_Pop(t *testing.T) {
	minHeap := &MinHeap{}
	heap.Init(minHeap)

	for i := range 10 {
		heap.Push(minHeap, CommonOrder(strconv.Itoa(i), Cold))
	}

	item := heap.Pop(minHeap).(Order)
	assert.Equal(t, "0", item.ID)

	item = heap.Pop(minHeap).(Order)
	assert.Equal(t, "1", item.ID)

	item = heap.Pop(minHeap).(Order)
	assert.Equal(t, "2", item.ID)
}

func CommonOrder(id string, temp Temperature) Order {
	return Order{
		ID:        id,
		Name:      "test order",
		Temp:      temp,
		Freshness: 10,
		TTL:       time.Now().UnixMicro() + 10*time.Second.Abs().Microseconds(),
	}
}
