package inverted_index

import (
	"math"
	"sync"
)

type InvertedIndex struct {
	Postings   map[string]map[string]*Posting
	DocCount   int
	DocLengths map[string]int
	TotalTerms int
	Embeddings map[string][]float64
	mu         sync.RWMutex // Protects all fields
}

func New() *InvertedIndex {
	return &InvertedIndex{
		Postings: make(map[string]map[string]*Posting),
		DocLengths: make(map[string]int),
	}
}
func (i *InvertedIndex) AddEmbedding(docId string, embedding []float64) {
	i.mu.Lock()
	defer i.mu.Unlock()
	
	if i.Embeddings == nil {
		i.Embeddings = make(map[string][]float64)
	}
	i.Embeddings[docId] = embedding
}
func (i *InvertedIndex) GetEmbedding(docId string) []float64 {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.Embeddings[docId]
}
func (i *InvertedIndex) AddDocument(docID string, tokens []string) {
	i.mu.Lock()
	defer i.mu.Unlock()
	
	i.DocLengths[docID] = len(tokens)
	i.TotalTerms += len(tokens)
	for pos, token := range tokens {
		if _, ok := i.Postings[token]; !ok {
			i.Postings[token] = make(map[string]*Posting)
		}
		if _, ok := i.Postings[token][docID]; !ok {
			i.Postings[token][docID] = &Posting{}
		}
		p := i.Postings[token][docID]
		p.TF++
		p.Positions = append(p.Positions, pos)
	}
	i.DocCount++
}
func (i *InvertedIndex) DocLength(docId string) int {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.DocLengths[docId]
}
func (i *InvertedIndex) Get(term string) map[string]*Posting {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.Postings[term]
}
func (i *InvertedIndex) AvgDocLength() float64 {
	i.mu.RLock()
	defer i.mu.RUnlock()
	if i.DocCount == 0 {
		return 0
	}
	return float64(i.TotalTerms) / float64(i.DocCount)
}
func (i *InvertedIndex) GetIdf(term string) float64 {
	i.mu.RLock()
	defer i.mu.RUnlock()
	
	df := float64(len(i.Postings[term]))
	N := float64(i.DocCount)

	if df == 0 {
		return 0
	}

	return math.Log((N - df + 0.5) / (df + 0.5) + 1)
}

func (i *InvertedIndex) All() map[string]map[string]*Posting {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.Postings
}
func (i *InvertedIndex) Vocabulary() []string {
	i.mu.RLock()
	defer i.mu.RUnlock()
	
	terms := make([]string, 0, len(i.Postings))
	for t := range i.Postings {
		terms = append(terms, t)
	}
	return terms
}