package main

import (
	"fmt"

	"github.com/Karsod58/search_engine/documents"
	inverted_index "github.com/Karsod58/search_engine/index"
	"github.com/Karsod58/search_engine/processor"
)

func main() {
	docs:=documents.Sample()
	proc:=processor.New()
	idx:=inverted_index.New()
    for _,doc:=range docs {
      _,freq:=proc.Process(doc.Text)
	  idx.AddDocument(doc.ID,freq)
	} 
	idx.Finalize()	
fmt.Println("=== TF Index ===")
for term, posting := range idx.All() {
	fmt.Println(term, "=>", posting)
}
fmt.Println(" IDF scores")
for term,_ := range idx.All() {
	fmt.Printf("%s => %f \n",term,idx.GetIdf(term))
}


}