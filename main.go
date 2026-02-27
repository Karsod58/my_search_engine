package main

import (
	

	tea "github.com/charmbracelet/bubbletea"

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

	for _, doc := range docs {
		token, _ := proc.Process(doc.Text)
		idx.AddDocument(doc.ID,token )
	}

	idx.Finalize()
searcher := search.New(idx, proc, docs)

	m := tui.New(searcher)
p := tea.NewProgram(m)
_, _ = p.Run()
}