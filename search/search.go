package search

import (
	"github.com/Karsod58/search_engine/documents"
	inverted_index "github.com/Karsod58/search_engine/index"
	"github.com/Karsod58/search_engine/processor"
)

type Searcher struct {
	idx       *inverted_index.InvertedIndex
	processor *processor.Processor
	docs      []documents.Document
	auto      *AutoComplete
}

func New(
	idx *inverted_index.InvertedIndex,
	p *processor.Processor,
	docs []documents.Document,
) *Searcher {

	ac := NewAutoComplete(idx.Vocabulary())

	return &Searcher{
		idx:       idx,
		processor: p,
		docs:      docs,
		auto:      ac,
	}
}

func (s *Searcher) Suggest(prefix string) []string {
	return s.auto.Suggest(prefix, 5)
}

func (s *Searcher) Search(query string, k int) []Result {

	terms, _ := s.processor.Process(query)
    root:=parse(query)
	allowed:=s.evaluate(root)
	nearDocs:=s.evaluateProximity(query)
	scores := make(map[string]float64)

	for _, term := range terms {
		postings := s.idx.Get(term)
		idf := s.idx.GetIdf(term)

		for docID, posting := range postings {
			if allowed!=nil {
				if !allowed[docID] {
					continue
				}
			}
			k1 := 1.5
b := 0.75
avgDocLen := s.idx.AvgDocLength()

tf := posting.TF
docLen := float64(s.idx.DocLength(docID))

numerator := tf * (k1 + 1)
denominator := tf + k1*(1-b+b*(docLen/avgDocLen))

bm25 := idf * (numerator / denominator)

scores[docID] += bm25
		}
		if nearDocs!=nil{
			for docId:=range nearDocs {
				if !nearDocs[docId] {
					delete(scores,docId)
				}
			}
		}
	}

postingsList := []map[string]*inverted_index.Posting{}

for _, term := range terms {
	postingsList = append(postingsList, s.idx.Get(term))
}

phraseDocs := inverted_index.HasPhrase(postingsList)

for docID := range phraseDocs {
	scores[docID] += 2.0  // phrase boost
}
	top := TopK(scores, k)

	results := make([]Result, 0, len(top))

	for _, t := range top {
		doc := documents.GetByID(s.docs, t.DocID)

	snippet := s.bestSnippet(doc.Text, terms, 25)
snippet = Highlight(snippet, terms)	

		results = append(results, Result{
			DocID:   t.DocID,
			Score:   t.Score,
			Snippet: snippet,
		})
	}

	return results
}