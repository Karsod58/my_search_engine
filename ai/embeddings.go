package ai

import (
	"context"
	"math"

	"github.com/ollama/ollama/api"
)

type EmbeddingService struct {
	client *api.Client
	model  string
}

func NewEmbeddingService() (*EmbeddingService, error) {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		return nil,err
	}
	return &EmbeddingService{
		client: client,
		model:  "llama3.2",
	}, nil
}
func (e *EmbeddingService) GetEmbedding(text string) ([]float64, error) {
	ctx := context.Background()
	req := &api.EmbeddingRequest{
		Model:  e.model,
		Prompt: text,
	}
	resp, err := e.client.Embeddings(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Embedding,nil
}
func  CosineSimilarity(a,b []float64) float64{
	if len(a)!=len(b) {
		return 0.0
	}
	var dotProduct,normalA, normalB float64
	for i:=range a {
		dotProduct+=a[i]*b[i]
		normalA+=a[i]*a[i]
		normalB+=b[i]*b[i]
	}
	if normalA==0.0 || normalB==0.0 {
		return   0.0
	}
	return dotProduct/(math.Sqrt(normalA)*math.Sqrt(normalB))
}
