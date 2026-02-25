package main

import (
	"log"

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
		_, freq := proc.Process(doc.Text)
		idx.AddDocument(doc.ID, freq)
	}

	idx.Finalize()

	searcher := search.New(idx, proc)

	p := tea.NewProgram(tui.New(searcher))

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}