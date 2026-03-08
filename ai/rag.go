package ai

import (
	"context"
	"fmt"
	"strings"

	"github.com/ollama/ollama/api"
)

type ChatMessage struct {
	Role    string
	Content string
}
type RAGChat struct {
	client *api.Client
	model string
	history []ChatMessage
}
func NewRAGChat() (*RAGChat,error){
	client,err:=api.ClientFromEnvironment()
	if err!=nil{
		return nil,err
	}
	return &RAGChat{
		client: client,
		model: "llama3.2",
		history: []ChatMessage{},
	},nil
}
func(r *RAGChat) Chat(question string,relevantDocs []string) (string,error){
	ctx:=context.Background()

	docsContext:=""
	if len(relevantDocs)>0 {
		docsContext="Relevant documents:\n\n"
		for i,doc:=range relevantDocs {
			docsContext+=fmt.Sprintf("Document %d:\n%s\n\n", i+1, doc)
		}
	}
	conversationContext:=""
	if len(r.history)>0 {
		conversationContext="Previous Conversation \n"
		for _,msg:=range r.history{
		conversationContext+=fmt.Sprintf("%s: %s\n",msg.Role,msg.Content)
		}
	}
	prompt := fmt.Sprintf(`You are a helpful assistant that answers questions based on the provided documents.

%s%s
User question: %s

Instructions:
- Answer based ONLY on the provided documents
- If the answer is not in the documents, say "I don't have enough information to answer that"
- Be concise and direct
- Reference specific documents when relevant

Answer:`, conversationContext, docsContext, question)
req:=&api.GenerateRequest{
	Model: r.model,
	Prompt: prompt,
	Stream: new(bool),
}
var response strings.Builder
respFunc:=func(resp api.GenerateResponse) error{
	response.WriteString(resp.Response)
	return nil
}
err:=r.client.Generate(ctx,req,respFunc)
if err!=nil{
	return "",err
}
answer:=strings.TrimSpace(response.String())
r.history=append(r.history, ChatMessage{Role: "user",Content: question})
r.history=append(r.history, ChatMessage{Role: "assistant",Content: answer})
if len(r.history)>6 {
	r.history=r.history[len(r.history)-6:]
}
return answer,nil
}
func(r *RAGChat) ChatStreaming(question string,relevantDocs []string,callback func(string)) error {
	ctx:=context.Background()
	docsContext:=""
	if len(relevantDocs)>0 {
		docsContext="Relevant documents:\n\n"
		for i,doc:=range relevantDocs{
			docsContext+=fmt.Sprintf("Document %d:\n%s\n\n",i+1,doc)
		}
	}
	conversationContext:=""
	if len(r.history)>0{
		conversationContext="Previous conversation:\n"
		for _,msg:=range r.history {
			conversationContext+=fmt.Sprintf("%s: %s\n",msg.Role,msg.Content)
		}
		conversationContext+="\n"
	}
	
	prompt := fmt.Sprintf(`You are a helpful assistant that answers questions based on the provided documents.

%s%s
User question: %s

Instructions:
- Answer based ONLY on the provided documents
- If the answer is not in the documents, say "I don't have enough information to answer that"
- Be concise and direct
- Reference specific documents when relevant

Answer:`, conversationContext, docsContext, question)


	req := &api.GenerateRequest{
		Model:  r.model,
		Prompt: prompt,
		Stream: func() *bool { b := true; return &b }(), // true
	}

	var fullResponse strings.Builder
	respFunc := func(resp api.GenerateResponse) error {
		if resp.Response != "" {
			fullResponse.WriteString(resp.Response)
			callback(resp.Response)
		}
		return nil
	}

	err := r.client.Generate(ctx, req, respFunc)
	if err != nil {
		return err
	}

	r.history = append(r.history, ChatMessage{Role: "user", Content: question})
	r.history = append(r.history, ChatMessage{Role: "assistant", Content: fullResponse.String()})

	if len(r.history) > 6 {
		r.history = r.history[len(r.history)-6:]
	}

	return nil
}


func (r *RAGChat) ClearHistory() {
	r.history = []ChatMessage{}
}


func (r *RAGChat) GetHistory() []ChatMessage {
	return r.history
}