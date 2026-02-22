package search

import inverted_index "github.com/Karsod58/search_engine/index"

type Searcher struct {
	docScore map[string]float64
	idx      *inverted_index.InvertedIndex
}

func New(idx *inverted_index.InvertedIndex) *Searcher {
	return &Searcher{docScore: make(map[string]float64),
	idx:idx}
}

func (s *Searcher) Search() {
  for term,_:=range s.idx.All() {
	docs:=s.idx.Get(term)
	idf:=s.idx.GetIdf(term)
	for doc:= range docs {
		tf:=docs[doc]
      s.docScore[doc]+=(idf*tf)
	} 
  }
}
func(s *Searcher) Get(docsId string) float64{
	return  s.docScore[docsId]
}