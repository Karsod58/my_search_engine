package ai

import (
	"context"
	"fmt"
	"strings"

	"github.com/ollama/ollama/api"
)

type QueryExpander struct {
	client *api.Client
	model string 
}
func NewQueryExpander() (*QueryExpander,error){
	client,err:=api.ClientFromEnvironment()
	if err!=nil{
		return nil,err
	}
	return  &QueryExpander{
		client: client,
		model: "llama3.2",
	},nil
}
func (q *QueryExpander) ExpandQuery(query string) ([]string,error){
	ctx:=context.Background()
	prompt:=fmt.Sprintf(`Given the search query: "%s"

Generate 3-5 related search terms or synonyms that would help find relevant documents.
Return only the terms, one per line, without explanations.

Terms:`, query)
req:=&api.GenerateRequest{
	Model: q.model,
	Prompt: prompt,
	Stream: new(bool),
}
var response strings.Builder
respFunc:=func(resp api.GenerateResponse) error {
	response.WriteString(resp.Response)
	return nil
}
err:=q.client.Generate(ctx,req,respFunc)
if err!=nil{
	return nil,err
}
terms:=[]string{query}
lines:=strings.Split(response.String(),"\n")
for _,line:=range lines {
	line=strings.TrimSpace(line)
	line=strings.TrimPrefix(line,"-")
	line=strings.TrimPrefix(line,"*")
	line=strings.TrimSpace(line)
	if line!="" && line!=query {
		terms=append(terms, line)
	}
	if len(terms)>=6 {
		break
	}
}
return  terms,nil
}
func (q *QueryExpander) DetectIntent(query string) (string,error){
	ctx:=context.Background()
	prompt := fmt.Sprintf(`Classify this search query into ONE category:
Query: "%s"

Categories:
- informational (seeking knowledge: "what is", "how to", "why")
- navigational (looking for specific page/site)
- transactional (wanting to do something: "download", "buy", "create")
- comparison (comparing options: "vs", "better than", "difference")

Return ONLY the category name, nothing else.`, query)
req:=&api.GenerateRequest{
	Model: q.model,
	Prompt: prompt,
	Stream: new(bool),
}
var response strings.Builder 
respFunc:=func(resp api.GenerateResponse) error {
	response.WriteString(resp.Response)
	return  nil
}
err:=q.client.Generate(ctx,req,respFunc)
if err!=nil{
	return "informational",err
}
intent:=strings.ToLower(strings.TrimSpace(response.String()))
validIntent:=[]string{"informational","navigational","transactional","comparison"}
for _,valid:=range validIntent {
	if strings.Contains(intent,valid) {
		return valid,err
	}
}
return "informational",nil
}
func (q *QueryExpander) ReformulateQuery(query string) (string,error){
	ctx:=context.Background()
	prompt := fmt.Sprintf(`Improve this search query to be more effective:
Original: "%s"

Return ONLY the improved query, nothing else. Keep it concise.`, query)
req:=&api.GenerateRequest{
	Model: q.model,
	Prompt: prompt,
	Stream: new(bool),
}
var response strings.Builder
respFunc:=func(resp api.GenerateResponse) error {
	response.WriteString(resp.Response)
	return nil
}
err:=q.client.Generate(ctx,req,respFunc)
if err!=nil{
	return query,err
}
improved:=strings.TrimSpace(response.String())
improved=strings.Trim(improved,`"'`)
if improved==""{
	return query,nil
}
return improved,nil
}