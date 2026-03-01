package inverted_index
 import "math"
type InvertedIndex struct {
	postings map[string]map[string]*Posting
	docCount int
	docLengths map[string]int
	totalTerms int 
}

func New() *InvertedIndex {
	return &InvertedIndex{
		postings: make(map[string]map[string]*Posting),
		docLengths: make(map[string]int),
	}
}

func (i *InvertedIndex) AddDocument(docID string, tokens []string) {
	i.docLengths[docID]=len(tokens)
	i.totalTerms+=len(tokens)
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
func(i *InvertedIndex) DocLength(docId string)int {
	return  i.docLengths[docId]
}
func (i *InvertedIndex) Get(term string) map[string]*Posting {
	return i.postings[term]
}
func(i *InvertedIndex) AvgDocLength() float64 {
	if i.docCount==0{
		return  0
	}
	return float64(i.totalTerms)/float64(i.docCount)
}
func (i *InvertedIndex) GetIdf(term string) float64 {

	df := float64(len(i.postings[term]))
	N := float64(i.docCount)

	if df == 0 {
		return 0
	}

	return math.Log((N - df + 0.5) / (df + 0.5) + 1)
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