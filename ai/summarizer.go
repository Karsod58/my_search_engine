package ai

import (
	"context"
	"fmt"
	"strings"

	"github.com/ollama/ollama/api"
)

type Summarizer struct {
	client *api.Client
	model  string
}
func NewSummarizer() (*Summarizer,error){
	client,err:=api.ClientFromEnvironment()
	if err!=nil{
		return nil,err
	}
	return  &Summarizer{
		client: client,
		model: "llama3.2",
	},nil
}
func(s *Summarizer) SummarizeResults(query string,results []string) (string,error) {
	ctx:=context.Background()
	context:=strings.Join(results,"\n\n---\n\n")
	
	prompt := fmt.Sprintf(`Based on the following search results for the query "%s", provide a concise summary (2-3 sentences) that answers the query or explains the key findings.

Search Results:
%s

Summary:`, query, context)
req:=&api.GenerateRequest{
	Model: s.model,
	Prompt: prompt,
	Stream: new(bool),
}
var  response strings.Builder
respFunc:=func(resp api.GenerateResponse) error{
	response.WriteString(resp.Response)
	return nil
}
err:=s.client.Generate(ctx,req,respFunc)
if err!=nil {
	return "",err
}
return strings.TrimSpace(response.String()),nil
}
func(s *Summarizer) SummarizeResultsStreaming(query string,results[] string,callback func(string)) error {
	ctx:=context.Background()
	context:=strings.Join(results,"\n\n---\n\n")
	prompt:=fmt.Sprintf(`Based on the following search results for the query "%s", provide a concise summary (2-3 sentences) that answers the query or explains the key findings.

Search Results:
%s

Summary:`, query, context)
req:=&api.GenerateRequest{
	Model: s.model,
	Prompt: prompt,
	Stream: func() *bool{b:=true; return &b} (),
}
respFunc:=func(resp api.GenerateResponse) error {
	if resp.Response!=""{
		callback(resp.Response)
	}
	return nil
}
return s.client.Generate(ctx,req,respFunc)
}
func(s *Summarizer) AnswerQuestion(question string,documents []string) (string,error) {
	ctx:=context.Background()
	context:=strings.Join(documents,"\n\n---\n\n")
	prompt:=fmt.Sprintf(`Answer the following question based ONLY on the provided documents. If the answer is not in the documents, say "I don't have enough information to answer that."

Question: %s

Documents:
%s

Answer:`, question, context)
req:=&api.GenerateRequest{
	Model: s.model,
	Prompt: prompt,
	Stream: new(bool),
}
var response strings.Builder
respFunc:=func (resp api.GenerateResponse) error {
	response.WriteString(resp.Response)
	return nil
}
err:=s.client.Generate(ctx,req,respFunc)
if err!=nil{
	return "",err
}
return strings.TrimSpace(response.String()),nil
}
func (s *Summarizer) ExtractKeyInsights(documents []string) ([]string, error) {
	ctx := context.Background()

	context := strings.Join(documents, "\n\n---\n\n")

	prompt := fmt.Sprintf(`Extract 3-5 key insights from these documents. Return them as a bullet list.

Documents:
%s

Key Insights:`, context)

	req := &api.GenerateRequest{
		Model:  s.model,
		Prompt: prompt,
		Stream: new(bool),
	}

	var response strings.Builder
	respFunc := func(resp api.GenerateResponse) error {
		response.WriteString(resp.Response)
		return nil
	}

	err := s.client.Generate(ctx, req, respFunc)
	if err != nil {
		return nil, err
	}

	
	insights := []string{}
	lines := strings.Split(response.String(), "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		line = strings.TrimPrefix(line, "-")
		line = strings.TrimPrefix(line, "*")
		line = strings.TrimPrefix(line, "•")
		line = strings.TrimSpace(line)
		
		if line != "" {
			insights = append(insights, line)
		}
	}

	return insights, nil
}