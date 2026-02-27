package inverted_index
 import "math"
type InvertedIndex struct {
	postings map[string]map[string]*Posting
	idf      map[string]float64
	docCount int
}

func New() *InvertedIndex {
	return &InvertedIndex{
		postings: make(map[string]map[string]*Posting),
		idf: make(map[string]float64),
	}
}

func (i *InvertedIndex) AddDocument(docID string, tokens []string) {
   for pos,token:= range  tokens{
	 if _,ok:=i.postings[token];!ok{
         i.postings[token]=make(map[string]*Posting)
	 }
	 if _,ok:=i.postings[token][docID]; !ok {
		i.postings[token][docID]=&Posting{}
	 }
	 p:=i.postings[token][docID]
	 p.TF++;
	 p.Positions=append(p.Positions, pos)
   } 
   i.docCount++
}

func (i *InvertedIndex) Get(term string) map[string]*Posting {
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
func (i *InvertedIndex) All() map[string]map[string]*Posting {
	return i.postings
}
func (i *InvertedIndex) Vocabulary() []string {
	terms := make([]string, 0, len(i.postings))
	for t := range i.postings {
		terms = append(terms, t)
	}
	return terms
}