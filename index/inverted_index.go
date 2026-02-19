package index

type InvertedIndex struct {
	postings map[string]map[string]int
}

func New() *InvertedIndex {
	return &InvertedIndex{
		postings: make(map[string]map[string]int),
	}
}

func (i *InvertedIndex) AddDocument(docID string, freq map[string]int) {
	for term, count := range freq {
		if _, ok := i.postings[term]; !ok {
			i.postings[term] = make(map[string]int)
		}
		i.postings[term][docID] = count
	}
}

func (i *InvertedIndex) Get(term string) map[string]int {
	return i.postings[term]
}

func (i *InvertedIndex) All() map[string]map[string]int {
	return i.postings
}