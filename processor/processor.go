package processor

import (
	"strings"
	"unicode"

)

type Processor struct {
	stopwords map[string]struct{}

}

func New() *Processor {
	return &Processor{
		stopwords:LoadStopwords(),
	}
}

func (p *Processor) Process(text string) ([]string, map[string]int) {

	text = strings.ToLower(text)
     var b strings.Builder
	 b.Grow(len(text))
	for _,ch:=range text{

		if unicode.IsLetter(ch) || unicode.IsDigit(ch) || unicode.IsSpace(ch){
        b.WriteRune(ch)
		}
	}

	words := strings.Fields(b.String())


	tokens := make([]string, 0)
	freq := make(map[string]int)

	for _, w := range words {
		if _, isStop := p.stopwords[w]; isStop {
			continue
		}

		tokens = append(tokens, w)
		freq[w]++
	}

	return tokens, freq
}