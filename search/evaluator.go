package search


func (s *Searcher) evaluate(node *Node) map[string]bool {

	switch node.Type {

	case TERM:
		postings := s.idx.Get(node.Value)
		set := make(map[string]bool)
		for docID := range postings {
			set[docID] = true
		}
		return set

	case AND:
		left := s.evaluate(node.Left)
		right := s.evaluate(node.Right)
		return intersect(left, right)

	case OR:
		left := s.evaluate(node.Left)
		right := s.evaluate(node.Right)
		return union(left, right)

	case NOT:
		all := s.allDocs()
		child := s.evaluate(node.Left)
		return difference(all, child)

	case NEAR:
		return s.evaluateNear(node)

	}

	return nil
}
func (s *Searcher) evaluateNear(node *Node) map[string]bool {

	leftTerm := node.Left.Value
	rightTerm := node.Right.Value

	post1 := s.idx.Get(leftTerm)
	post2 := s.idx.Get(rightTerm)

	result := make(map[string]bool)

	for docID, p1 := range post1 {
		p2, ok := post2[docID]
		if !ok {
			continue
		}

		for _, pos1 := range p1.Positions {
			for _, pos2 := range p2.Positions {
				if abs(pos1-pos2) <= node.Distance {
					result[docID] = true
					break
				}
			}
			if result[docID] {
				break
			}
		}
	}

	return result
}
func (s *Searcher) allDocs() map[string]bool {
	set := make(map[string]bool)
	for _, doc := range s.docs {
		set[doc.ID] = true
	}
	return set
}