package search

import (
    inverted_index "github.com/Karsod58/search_engine/index"
    "github.com/Karsod58/search_engine/processor"
)

type Searcher struct {
    idx       *inverted_index.InvertedIndex
    processor *processor.Processor
}

func New(idx *inverted_index.InvertedIndex, p *processor.Processor) *Searcher {
    return &Searcher{
        idx:       idx,
        processor: p,
    }
}

func (s *Searcher) Search(query string, k int) []Result {
    docScores := make(map[string]float64)

    terms, _ := s.processor.Process(query)

    for _, term := range terms {
        postings := s.idx.Get(term)
        idf := s.idx.GetIdf(term)

        for docID, tf := range postings {
            docScores[docID] += tf * idf
        }
    }

    return TopK(docScores, k)
}