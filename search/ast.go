package search

type NodeType int 

const (
	TERM NodeType=iota
  AND
  OR
  NOT
  NEAR
)
type Node struct{
	Type NodeType
	Value string
	Distance int
	Left *Node
	Right *Node
}

func precedence(op string) int {
	switch op {
	case "NOT":
		return 3
	case "AND":
		return 2
	case "OR":
		return  1
	}
	return  0
}