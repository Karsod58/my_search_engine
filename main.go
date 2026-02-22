package main

import (
	"fmt"

	"github.com/Karsod58/search_engine/documents"
	inverted_index "github.com/Karsod58/search_engine/index"
	"github.com/Karsod58/search_engine/processor"
	"github.com/Karsod58/search_engine/search"
)

func main() {
	docs:=documents.Sample()
	proc:=processor.New()
	idx:=inverted_index.New()
	searcher:=search.New(idx)
    for _,doc:=range docs {
      _,freq:=proc.Process(doc.Text)
	  idx.AddDocument(doc.ID,freq)
	} 
	idx.Finalize()	
	searcher.Search()
fmt.Println("=== TF Index ===")
for term, posting := range idx.All() {
	fmt.Println(term, "=>", posting)
}
fmt.Println(" IDF scores")
for term,_ := range idx.All() {
	fmt.Printf("%s => %f \n",term,idx.GetIdf(term))
}
fmt.Println("Scores for each document")
for _,doc:= range docs{
	fmt.Printf("%s : score is %f \n",doc.ID,searcher.Get(doc.ID))
}


}