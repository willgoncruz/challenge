package model

type MinHeap []Order // Order min heap for deciding discard candidates

func (h MinHeap) Len() int {
	return len(h)
}

func (h MinHeap) Less(i, j int) bool {
	return h[i].TTL.Milliseconds() < h[j].TTL.Milliseconds()
}

func (h MinHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *MinHeap) Push(x interface{}) {
	*h = append(*h, x.(Order))
}

func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func (h MinHeap) Find(order Order) int {
	bottom := 0
	top := h.Len() - 1

	for bottom <= top {
		mid := (bottom + top) / 2
		if h[mid].ID == order.ID {
			return mid
		} else if h.Less(mid, bottom) {
			top = mid
		} else {
			bottom = mid + 1
		}
	}

	return -1
}
