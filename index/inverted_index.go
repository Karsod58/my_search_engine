package inverted_index
 import "math"
type InvertedIndex struct {
	postings map[string]map[string]float64
	idf      map[string]float64
	docCount int
}

func New() *InvertedIndex {
	return &InvertedIndex{
		postings: make(map[string]map[string]float64),
		idf: make(map[string]float64),
	}
}

func (i *InvertedIndex) AddDocument(docID string, freq map[string]int) {
    i.docCount++
	total := 0
	for _, c := range freq {
		total += c
	}

	if total == 0 {
		return
	}

	for term, count := range freq {
		tf := float64(count) / float64(total)

		if _, ok := i.postings[term]; !ok {
			i.postings[term] = make(map[string]float64)
		}

		i.postings[term][docID] = tf
	}
}

func (i *InvertedIndex) Get(term string) map[string]float64 {
	return i.postings[term]
}
func (i *InvertedIndex) Finalize() {
	N:=float64(i.docCount)
	for term,docs:= range i.postings {
		df:=float64(len(docs))
		idfval:=math.Log(N/df)
		i.idf[term]=idfval
	}
}
func (i *InvertedIndex)  GetIdf(term string) float64{
	return i.idf[term]
}
func (i *InvertedIndex) All() map[string]map[string]float64 {
	return i.postings
}