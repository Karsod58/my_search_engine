package main

import (
	"fmt"

	"github.com/Karsod58/search_engine/documents"
	inverted_index "github.com/Karsod58/search_engine/index"
	"github.com/Karsod58/search_engine/processor"
	"github.com/Karsod58/search_engine/search"
)

func main() {

	docs := documents.Sample()

	proc := processor.New()
	idx := inverted_index.New()


	for _, doc := range docs {
		_, freq := proc.Process(doc.Text)
		idx.AddDocument(doc.ID, freq)
	}

	idx.Finalize()

	searcher := search.New(idx, proc)


	query := "go easy"
	results := searcher.Search(query, 2)


	fmt.Println("=== TF Index ===")
	for term, posting := range idx.All() {
		fmt.Println(term, "=>", posting)
	}

	fmt.Println("\n=== IDF Scores ===")
	for term := range idx.All() {
		fmt.Printf("%s => %.6f\n", term, idx.GetIdf(term))
	}

	fmt.Println("\nTop results for query:", query)
	for _, r := range results {
		fmt.Printf("%s => %.6f\n", r.DocID, r.Score)
	}


	fmt.Println("\nFetched Documents:")
	for _, r := range results {
		doc := documents.GetByID(docs, r.DocID)
		fmt.Printf("%s => %s\n", r.DocID, doc.Text)
	}
}