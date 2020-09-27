package search

import "github.com/datainq/teryt"

type LocationWrapper struct {
	*teryt.Location
	Score      int
	SearchText []rune
}

type Heap []*LocationWrapper

func (h Heap) Len() int {
	return len(h)
}

func (h Heap) Less(i, j int) bool {
	a, b := h[i], h[j]
	return a.Score > b.Score || a.Score == b.Score && (len(a.Parts) > len(b.Parts))
}

func (h Heap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *Heap) Push(x interface{}) {
	*h = append(*h, x.(*LocationWrapper))
}

func (h *Heap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // avoid memory leak
	*h = old[0 : n-1]
	return item
}

type ReverseHeap []*LocationWrapper

func (h ReverseHeap) Len() int {
	return len(h)
}

func (h ReverseHeap) Less(i, j int) bool {
	a, b := h[i], h[j]
	return a.Score < b.Score || a.Score == b.Score && (len(a.Parts) < len(b.Parts))
}

func (h ReverseHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *ReverseHeap) Push(x interface{}) {
	*h = append(*h, x.(*LocationWrapper))
}

func (h *ReverseHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // avoid memory leak
	*h = old[0 : n-1]
	return item
}

type SearchResult struct {
	Location *teryt.Location
	Score    int
}
