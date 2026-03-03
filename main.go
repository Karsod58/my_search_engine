package main

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/Karsod58/search_engine/crawler"
	"github.com/Karsod58/search_engine/documents"
	inverted_index "github.com/Karsod58/search_engine/index"
	"github.com/Karsod58/search_engine/processor"
	"github.com/Karsod58/search_engine/search"
	"github.com/Karsod58/search_engine/tui"
)

func main() {

	docs := documents.Sample()
	proc := processor.New()
    seed := "https://pkg.go.dev/"
	crawler := crawler.New(seed, 1)
    idxFile:="index.json"
	var idx *inverted_index.InvertedIndex

if _, err := os.Stat(idxFile); err == nil {

    
    idx, err = inverted_index.Load(idxFile)
    if err != nil {
        panic(err)
    }

    println("Index loaded from disk")

} else {
   
    idx = inverted_index.New()

  crawler.Start(seed, func(url string, text string) {
	tokens, _ := proc.Process(text)
	idx.AddDocument(url, tokens)
})

    if err := idx.Save(idxFile); err != nil {
        panic(err)
    }

    println("Index built and saved")
}
	for _, doc := range docs {
		token, _ := proc.Process(doc.Text)
		idx.AddDocument(doc.ID,token )
	}

searcher := search.New(idx, proc, docs)

	m := tui.New(searcher)
p := tea.NewProgram(m)
_, _ = p.Run()
}