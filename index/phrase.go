package inverted_index

func HasPhrase(postings []map[string]*Posting) map[string]bool{
	res:=make(map[string]bool)
	if len(postings)==0 {
		return  res
	}
	for docId:=range postings[0] {
		matched:=true
		basePositions:=postings[0][docId].Positions
		for i:=1;i<len(postings);i++ {
			nextPosting,ok:=postings[i][docId]
			if !ok{
				matched=false
				break
			}
			found:=false

			for _,pos:=range basePositions{
				for _,nextPos:=range nextPosting.Positions {
					if nextPos==pos+i {
						found=true
						break
					}
				}
				if found {
					break
				}
			}
			if !found {
				matched=false
				break
			}
		}
		if matched {
			res[docId]=true
		}
	}
	return  res
}