package search

import (
	"fmt"
	"time"

	"github.com/Karsod58/search_engine/ai"
	"github.com/Karsod58/search_engine/documents"
	inverted_index "github.com/Karsod58/search_engine/index"
	"github.com/Karsod58/search_engine/processor"
)

type Searcher struct {
	idx       *inverted_index.InvertedIndex
	processor *processor.Processor
	docs      []documents.Document
	auto      *AutoComplete
	embedder  *ai.EmbeddingService
}

func New(
	idx *inverted_index.InvertedIndex,
	p *processor.Processor,
	docs []documents.Document,
	embedder *ai.EmbeddingService,
) *Searcher {

	ac := NewAutoComplete(idx.Vocabulary())

	return &Searcher{
		idx:       idx,
		processor: p,
		docs:      docs,
		auto:      ac,
		embedder: embedder,
	}
}

func (s *Searcher) getBM25Scores(query string) map[string]float64 {
	if s == nil || s.idx == nil || s.processor == nil {
		return nil
	}

	terms, _ := s.processor.Process(query)
	scores := make(map[string]float64)

	for _, term := range terms {
		postings := s.idx.Get(term)
		idf := s.idx.GetIdf(term)

		for docID, posting := range postings {
			k1 := 1.5
			b := 0.75
			avgDocLen := s.idx.AvgDocLength()

			tf := posting.TF
			docLen := float64(s.idx.DocLength(docID))

			numerator := tf * (k1 + 1)
			denominator := tf + k1*(1-b+b*(docLen/avgDocLen))

			bm25 := idf * (numerator / denominator)

			scores[docID] += bm25
		}
	}

	return scores
}

func (s *Searcher) Suggest(prefix string) []string {
	return s.auto.Suggest(prefix, 5)
}
func (s *Searcher) SemanticSearch(query string,k int,alpha float64) []Result{
	start:=time.Now()
	queryEmbed,err:=s.embedder.GetEmbedding(query)
	if err!=nil {
		return s.Search(query,k)
	}
	bm25Score:=s.getBM25Scores(query)
	semanticScores:=make(map[string]float64)
	for docId:=range s.idx.Postings {
		docEmbed:=s.idx.GetEmbedding(docId)
		if docEmbed!=nil {
			similarity:=ai.CosineSimilarity(queryEmbed,docEmbed)
			semanticScores[docId]=similarity
		}
	}
	bm25Normalized:=normalizeScores(bm25Score)
	semanticNormalized:=normalizeScores(semanticScores)
    finalScore:=make(map[string]float64)
	allDocs:=make(map[string]bool)
	for docId := range bm25Normalized {
		allDocs[docId]=true
	}
	for docId:= range semanticNormalized {
		allDocs[docId]=true
	}
	for docId:=range  allDocs {
		bm25:=bm25Normalized[docId]
		semantic:=semanticNormalized[docId]
		finalScore[docId]=alpha*bm25+(1-alpha)*semantic
	}

	terms, _ := s.processor.Process(query)
	
	totalResults := len(finalScore)
	elapsed := time.Since(start)
	stats := &SearchStats{
		QueryTime:    fmt.Sprintf("%.2fms", float64(elapsed.Microseconds())/1000.0),
		DocsSearched: len(s.docs),
		TermsMatched: len(terms),
		TotalResults: totalResults,
	}


	top := TopK(finalScore, k)
	
	results := make([]Result, 0, len(top))
	for _, t := range top {
		doc := documents.GetByID(s.docs, t.DocID)
		if doc == nil {
			continue
		}

		snippet := ""
		if doc.Text != "" {
			snippet = s.bestSnippet(doc.Text, terms, 25)
			snippet = Highlight(snippet, terms)
		}

		title := doc.Title
		if title == "" && len(doc.Text) > 0 {
			if len(doc.Text) > 50 {
				title = doc.Text[:50] + "..."
			} else {
				title = doc.Text
			}
		}

		results = append(results, Result{
			DocID:   t.DocID,
			Score:   t.Score,
			Snippet: snippet,
			Title:   title,
			URL:     doc.URL,
			Stats:   stats,
		})
	}
	return results
}

func (s *Searcher) Search(query string, k int) []Result {
	start := time.Now()

	if s == nil || s.idx == nil || s.processor == nil {
		return nil
	}

	terms, _ := s.processor.Process(query)
	corrections := make(map[string]string)

	expandedTerms := make([]string, 0, len(terms))
	for _, term := range terms {
		postings := s.idx.Get(term)
		if postings == nil || len(postings) == 0 {
			matches := processor.FindClosest(term, s.idx.Vocabulary(), 2)
			if len(matches) > 0 {
				corrections[term] = matches[0]
				expandedTerms = append(expandedTerms, matches[0])
			} else {
				expandedTerms = append(expandedTerms, term)
			}
		} else {
			expandedTerms = append(expandedTerms, term)
		}
	}

	terms = expandedTerms

	root := parse(query)
	allowed := s.evaluate(root)
	nearDocs := s.evaluateProximity(query)
	scores := make(map[string]float64)

	for _, term := range terms {
		postings := s.idx.Get(term)
		idf := s.idx.GetIdf(term)

		for docID, posting := range postings {
			if allowed != nil {
				if !allowed[docID] {
					continue
				}
			}
			k1 := 1.5
			b := 0.75
			avgDocLen := s.idx.AvgDocLength()

			tf := posting.TF
			docLen := float64(s.idx.DocLength(docID))

			numerator := tf * (k1 + 1)
			denominator := tf + k1*(1-b+b*(docLen/avgDocLen))

			bm25 := idf * (numerator / denominator)

			scores[docID] += bm25
		}
		if nearDocs != nil {
			for docId := range nearDocs {
				if !nearDocs[docId] {
					delete(scores, docId)
				}
			}
		}
	}

	postingsList := []map[string]*inverted_index.Posting{}

	for _, term := range terms {
		postingsList = append(postingsList, s.idx.Get(term))
	}

	phraseDocs := inverted_index.HasPhrase(postingsList)

	for docID := range phraseDocs {
		scores[docID] += 2.0
	}

	totalResults := len(scores)
	elapsed := time.Since(start)
	stats := &SearchStats{
		QueryTime:    fmt.Sprintf("%.2fms", float64(elapsed.Microseconds())/1000.0),
		DocsSearched: len(s.docs),
		TermsMatched: len(terms),
		TotalResults: totalResults,
	}

	top := TopK(scores, k)
	results := make([]Result, 0, len(top))

	for _, t := range top {
		doc := documents.GetByID(s.docs, t.DocID)
		if doc == nil {
			continue
		}

		snippet := ""
		if doc.Text != "" {
			snippet = s.bestSnippet(doc.Text, terms, 25)
			snippet = Highlight(snippet, terms)
		}

		title := doc.Title
		if title == "" && len(doc.Text) > 0 {
			if len(doc.Text) > 50 {
				title = doc.Text[:50] + "..."
			} else {
				title = doc.Text
			}
		}

		result := Result{
			DocID:   t.DocID,
			Score:   t.Score,
			Snippet: snippet,
			Title:   title,
			URL:     doc.URL,
			Stats:   stats,
		}

		if len(results) == 0 && len(corrections) > 0 {
			result.Corrections = corrections
		}

		results = append(results, result)
	}

	return results
}
func normalizeScores(scores map[string]float64) map[string]float64 {
    if len(scores) == 0 {
        return scores
    }
    
    var maxScore float64
    for _, score := range scores {
        if score > maxScore {
            maxScore = score
        }
    }
    
    if maxScore == 0 {
        return scores
    }
    
    normalized := make(map[string]float64)
    for docID, score := range scores {
        normalized[docID] = score / maxScore
    }
    
    return normalized
}