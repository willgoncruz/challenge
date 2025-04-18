package model

type MinHeap []Order

func (h MinHeap) Len() int {
	return len(h)
}

func (h MinHeap) Less(i, j int) bool {
	return h[i].TTL < h[j].TTL
}

func (h MinHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *MinHeap) Push(val interface{}) {
	*h = append(*h, val.(Order))
}

func (h *MinHeap) Pop() interface{} {
	heapDerefrenced := *h

	size := len(heapDerefrenced)
	val := heapDerefrenced[size-1]
	*h = heapDerefrenced[:size-1]

	return val
}

func (h MinHeap) Find(order Order) int {
	for i, o := range h {
		if o.ID == order.ID {
			return i
		}
	}

	return -1
}
