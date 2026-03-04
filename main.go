package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/Karsod58/search_engine/crawler"
	"github.com/Karsod58/search_engine/documents"
	inverted_index "github.com/Karsod58/search_engine/index"
	"github.com/Karsod58/search_engine/processor"
	"github.com/Karsod58/search_engine/search"
	"github.com/Karsod58/search_engine/tui"
)

func main() {
	logMainEvent("main.go:main", "start", map[string]interface{}{}, "run2", "HM")

	seed := "https://pkg.go.dev/"
	c := crawler.New(seed, 1) // Reduced depth from 2 to 1

	idx := inverted_index.New()
	p := processor.New()

	var docs []documents.Document
	docID := 0

	logMainEvent("main.go:main", "before_crawler_start", map[string]interface{}{
		"seed": seed,
	}, "run2", "HM")

	c.Start(seed, func(url string, text string) {
		tokens, _ := p.Process(text)

		id := fmt.Sprintf("doc-%d", docID)
		docID++

		idx.AddDocument(id, tokens)

		docs = append(docs, documents.Document{
			ID:   id,
			Text: text,
		})
	})

	logMainEvent("main.go:main", "after_crawler_return", map[string]interface{}{
		"docCount": len(docs),
	}, "run2", "HM")

	searcher := search.New(idx, p, docs)
	m := tui.New(searcher)
	prog := tea.NewProgram(m)

	logMainEvent("main.go:main", "before_tui_run", map[string]interface{}{}, "run2", "HM")
	_, _ = prog.Run()
}

// #region agent log
func logMainEvent(location, message string, data map[string]interface{}, runId, hypothesisId string) {
	f, err := os.OpenFile("debug-815281.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()

	entry := map[string]interface{}{
		"sessionId":    "815281",
		"id":           fmt.Sprintf("log_%d", time.Now().UnixNano()),
		"timestamp":    time.Now().UnixMilli(),
		"location":     location,
		"message":      message,
		"data":         data,
		"runId":        runId,
		"hypothesisId": hypothesisId,
	}

	b, err := json.Marshal(entry)
	if err != nil {
		return
	}

	_, _ = f.Write(append(b, '\n'))
}
// #endregion