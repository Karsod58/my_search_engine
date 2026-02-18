package main

import (
	"github.com/Karsod58/search_engine/processor"
	"fmt"
)

func main() {
	text := "Go is an open-source database engine written in Go."

	p := processor.New()

	tokens, freq := p.Process(text)

	fmt.Println("Tokens:", tokens)
	fmt.Println("Freq:", freq)

}