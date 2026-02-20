package inverted_index

type InvertedIndex struct {
	postings map[string]map[string]float64
}

func New() *InvertedIndex {
	return &InvertedIndex{
		postings: make(map[string]map[string]float64),
	}
}

func (i *InvertedIndex) AddDocument(docID string, freq map[string]int) {

	// compute total terms in doc
	total := 0
	for _, c := range freq {
		total += c
	}

	if total == 0 {
		return
	}

	// convert counts → TF
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

func (i *InvertedIndex) All() map[string]map[string]float64 {
	return i.postings
}