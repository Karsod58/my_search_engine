package main

import (
	// "encoding/json"
	"fmt"
	"log"
	

	tea "github.com/charmbracelet/bubbletea"

	"github.com/Karsod58/search_engine/ai"
	"github.com/Karsod58/search_engine/crawler"
	"github.com/Karsod58/search_engine/documents"
	inverted_index "github.com/Karsod58/search_engine/index"
	"github.com/Karsod58/search_engine/processor"
	"github.com/Karsod58/search_engine/search"
	"github.com/Karsod58/search_engine/tui"
)

func main() {
	
  	fmt.Print("Enter URL to crawl: ")
		var seed string
		fmt.Scanln(&seed)
	c := crawler.New(seed, 1) 

	idx := inverted_index.New()
	p := processor.New()
	expander, err := ai.NewQueryExpander()
if err != nil {
	log.Printf("Warning: Could not initialize query expander: %v", err)
	expander = nil
} else {
	log.Println("✓ Query expander initialized")
}
      embedder, err := ai.NewEmbeddingService()
	   if err != nil {
        log.Printf("Warning: Could not initialize embeddings: %v", err)
    }
	summarizer, err := ai.NewSummarizer()
if err != nil {
	log.Printf("Warning: Could not initialize summarizer: %v", err)
	summarizer = nil
} else {
	log.Println("✓ Summarizer initialized")
}
	var docs []documents.Document
	docID := 0


	c.Start(seed, func(url string, text string) {
		tokens, _ := p.Process(text)

		id := fmt.Sprintf("doc-%d", docID)
		docID++

		idx.AddDocument(id, tokens)

		docs = append(docs, documents.Document{
			ID:   id,
			Text: text,
		})
		    if embedder != nil {
            embedding, err := embedder.GetEmbedding(text)
            if err == nil {
                idx.AddEmbedding(id, embedding)
            }
        }
	})

    
	searcher := search.New(idx, p, docs,embedder,expander,summarizer)
	m := tui.New(searcher)
	prog := tea.NewProgram(m)


	_, _ = prog.Run()
}

// // #region agent log
// func logMainEvent(location, message string, data map[string]interface{}, runId, hypothesisId string) {
// 	f, err := os.OpenFile("debug-815281.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
// 	if err != nil {
// 		return
// 	}
// 	defer f.Close()

// 	entry := map[string]interface{}{
// 		"sessionId":    "815281",
// 		"id":           fmt.Sprintf("log_%d", time.Now().UnixNano()),
// 		"timestamp":    time.Now().UnixMilli(),
// 		"location":     location,
// 		"message":      message,
// 		"data":         data,
// 		"runId":        runId,
// 		"hypothesisId": hypothesisId,
// 	}

// 	b, err := json.Marshal(entry)
// 	if err != nil {
// 		return
// 	}

// 	_, _ = f.Write(append(b, '\n'))
// }
// // #endregion