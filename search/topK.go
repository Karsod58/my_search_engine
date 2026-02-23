package search

import "container/heap"

type Result struct {
	DocID string
	Score float64
}
type MinHeap []Result

func (h MinHeap) Len() int { return len(h) }

func (h MinHeap) Less(i, j int) bool {
	return h[i].Score < h[j].Score 
}

func (h MinHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *MinHeap) Push(x interface{}) {
	*h = append(*h, x.(Result))
}

func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[:n-1]
	return item
}
func TopK(scores map[string]float64, k int) []Result {
	h := &MinHeap{}
	heap.Init(h)

	for docID, score := range scores {
		if h.Len() < k {
			heap.Push(h, Result{docID, score})
			continue
		}

		if score > (*h)[0].Score {
			heap.Pop(h)
			heap.Push(h, Result{docID, score})
		}
	}

	results := make([]Result, 0, h.Len())
	for h.Len() > 0 {
		results = append(results, heap.Pop(h).(Result))
	}

	for i, j := 0, len(results)-1; i < j; i, j = i+1, j-1 {
		results[i], results[j] = results[j], results[i]
	}

	return results
}