package search

import (
	"strconv"
	"strings"
)

func (s *Searcher) evaluateProximity(rawQuery string) map[string]bool {
	parts := strings.Fields(rawQuery)
if len(parts)!=3 {
	return nil
}
term1:=strings.ToLower(parts[0])
opr:=strings.ToUpper(parts[1])
term2:=strings.ToLower(parts[2])

if !strings.HasPrefix(opr,"NEAR/") {
	return  nil
}
trimmedStr:=strings.TrimPrefix(opr,"/NEAR")
n,err:=strconv.Atoi(trimmedStr)
if err!=nil{
	return  nil
}
post1:=s.idx.Get(term1)
post2:=s.idx.Get(term2)
res:=make(map[string]bool)
for docId,p1:=range post1 {
	p2,ok:=post2[docId]
	if !ok {
		continue
	}
		// check distance
		for _, pos1 := range p1.Positions {
			for _, pos2 := range p2.Positions {

				if abs(pos1-pos2) <= n {
					res[docId] = true
					break
				}
			}
			if res[docId] {
				break
			}
		}
	
}
return  res
}
func abs(x int) int {
	if x<0 {
		return -1*x
	}else{
		return x
	}
	
}