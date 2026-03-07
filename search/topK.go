package search

import "container/heap"

type SearchStats struct {
	QueryTime    string
	DocsSearched int
	TermsMatched int
	TotalResults int
}

type Result struct {
	DocID       string
	Score       float64
	Snippet     string
	Title       string
	URL         string
	Corrections map[string]string
	Stats       *SearchStats
    Summary string
    Insights []string
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
    if k <= 0 {
        return nil
    }

    h := &MinHeap{}
    heap.Init(h)

    for docID, score := range scores {
        r := Result{
            DocID: docID,
            Score: score,
        }

        if h.Len() < k {
            heap.Push(h, r)
            continue
        }

        if score > (*h)[0].Score {
            heap.Pop(h)
            heap.Push(h, r)
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