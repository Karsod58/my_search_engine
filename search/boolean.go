package search

import (
	"strings"

	inverted_index "github.com/Karsod58/search_engine/index"

)

func toDocSet(postings map[string]*inverted_index.Posting) map[string]bool{
	set:=make(map[string]bool)
	for docId:=range postings {
		set[docId]=true
	}
	return  set
}
func intersect(a,b map[string]bool)  map[string]bool{
	out:=make(map[string]bool)
	for k:=range a {
		if b[k] {
		out[k]=true
		}
	}
	return  out
}
func union(a,b map[string]bool) map[string] bool{
finalRes:=make(map[string]bool)
for term := range a {
	finalRes[term]=true
}
for term := range b{
	finalRes[term]=true
}
return  finalRes
}
func difference(a,b map[string]bool) map[string]bool{
	output:=make(map[string]bool)
	for  term:=range a {
		if !b[term] {
			output[term]=true
		}
	}
	return output
}
func(s *Searcher) evaluateBoolean(rawQuery string) map[string]bool {
	parts:=strings.Fields(rawQuery)
	var currSet map[string]bool
	var opr string
	 for _,part:=range parts {
		upper:=strings.ToUpper(part)
		if upper=="AND" || upper=="OR" || upper=="NOT" {
			opr=upper
			continue
		}
		postings:=s.idx.Get(strings.ToLower(part))
		docSet:=toDocSet(postings)
		if currSet==nil {
             currSet=docSet
			 continue
		}
		switch opr {
		case "AND":
			currSet=intersect(currSet,docSet)
		case "OR":
			currSet=union(currSet,docSet)
		case "NOT":
			currSet=difference(currSet,docSet)
		default:
			currSet=union(currSet,docSet)
		}
	 }
	 return  currSet
}