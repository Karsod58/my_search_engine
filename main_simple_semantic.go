package main

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/Karsod58/search_engine/ai"
	"github.com/Karsod58/search_engine/documents"
	inverted_index "github.com/Karsod58/search_engine/index"
	"github.com/Karsod58/search_engine/processor"
	"github.com/Karsod58/search_engine/search"
	"github.com/Karsod58/search_engine/tui"
)

func main() {
	fmt.Println("Initializing search engine with sample documents...")

	// Initialize embedding service
	embedder, err := ai.NewEmbeddingService()
	if err != nil {
		log.Printf("Warning: Could not initialize embeddings: %v", err)
		embedder = nil
	} else {
		fmt.Println("✓ Embedding service initialized")
	}

	// Use sample documents instead of crawler
	docs := documents.Sample()
	idx := inverted_index.New()
	p := processor.New()

	fmt.Println("Indexing documents...")
	for i, doc := range docs {
		tokens, _ := p.Process(doc.Text)
		idx.AddDocument(doc.ID, tokens)

		// Generate embeddings (this is slow)
		if embedder != nil {
			fmt.Printf("  Generating embedding for doc %d/%d...\n", i+1, len(docs))
			embedding, err := embedder.GetEmbedding(doc.Text)
			if err == nil {
				idx.AddEmbedding(doc.ID, embedding)
			}
		}
	}

	fmt.Println("✓ Indexing complete!")
	fmt.Println("\nStarting TUI... (Press 'S' to toggle semantic search)")

	searcher := search.New(idx, p, docs, embedder)
	m := tui.New(searcher)
	prog := tea.NewProgram(m)

	_, _ = prog.Run()
}
