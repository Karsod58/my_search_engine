package main

import (
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
	idx := inverted_index.New()
    seed := "https://google.com"

crawler := crawler.New(seed, 1)

crawler.Crawl(seed, 0, func(url string, text string) {

	tokens, _ := proc.Process(text)

	idx.AddDocument(url, tokens)
})
	for _, doc := range docs {
		token, _ := proc.Process(doc.Text)
		idx.AddDocument(doc.ID,token )
	}

searcher := search.New(idx, proc, docs)

	m := tui.New(searcher)
p := tea.NewProgram(m)
_, _ = p.Run()
}