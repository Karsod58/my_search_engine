package search

import (
	"strconv"
	"strings"
)

func tokenize(input string) []string {
	return strings.Fields(input)
}

func parse(query string) *Node {

	tokens := tokenize(query)

	var output []*Node
	var ops []string

	for i := 0; i < len(tokens); i++ {
		token := tokens[i]
		upper := strings.ToUpper(token)

		// Parentheses
		if token == "(" {
			ops = append(ops, token)
			continue
		}

		if token == ")" {
			for len(ops) > 0 && ops[len(ops)-1] != "(" {
				popOperator(&output, &ops)
			}
			ops = ops[:len(ops)-1] // remove "("
			continue
		}

		// NEAR/n
		if strings.HasPrefix(upper, "NEAR/") {
			distStr := strings.TrimPrefix(upper, "NEAR/")
			dist, _ := strconv.Atoi(distStr)
			ops = append(ops, "NEAR/"+strconv.Itoa(dist))
			continue
		}

		// Boolean operators
		if upper == "AND" || upper == "OR" || upper == "NOT" {

			for len(ops) > 0 && precedence(ops[len(ops)-1]) >= precedence(upper) {
				popOperator(&output, &ops)
			}

			ops = append(ops, upper)
			continue
		}

		// TERM
		output = append(output, &Node{
			Type:  TERM,
			Value: strings.ToLower(token),
		})
	}

	for len(ops) > 0 {
		popOperator(&output, &ops)
	}

	return output[0]
}
func popOperator(output *[]*Node, ops *[]string) {

	op := (*ops)[len(*ops)-1]
	*ops = (*ops)[:len(*ops)-1]

	var node Node

	switch {
	case op == "AND":
		right := (*output)[len(*output)-1]
		left := (*output)[len(*output)-2]
		*output = (*output)[:len(*output)-2]

		node = Node{Type: AND, Left: left, Right: right}

	case op == "OR":
		right := (*output)[len(*output)-1]
		left := (*output)[len(*output)-2]
		*output = (*output)[:len(*output)-2]

		node = Node{Type: OR, Left: left, Right: right}

	case op == "NOT":
		child := (*output)[len(*output)-1]
		*output = (*output)[:len(*output)-1]

		node = Node{Type: NOT, Left: child}

	case strings.HasPrefix(op, "NEAR/"):
		right := (*output)[len(*output)-1]
		left := (*output)[len(*output)-2]
		*output = (*output)[:len(*output)-2]

		distStr := strings.TrimPrefix(op, "NEAR/")
		dist, _ := strconv.Atoi(distStr)

		node = Node{
			Type:     NEAR,
			Left:     left,
			Right:    right,
			Distance: dist,
		}
	}

	*output = append(*output, &node)
}