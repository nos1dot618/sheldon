package query

import "strings"

func Parse(input string) ([]Node, error) {
	tokens := Tokenize(input)
	var nodes []Node

	for _, token := range tokens {
		if strings.ContainsAny(token, "~=<>") {
			node, err := parseFilter(token)
			if err != nil {
				return nil, err
			}
			nodes = append(nodes, node)
		} else {
			nodes = append(nodes, TextNode{Value: token})
		}
	}

	return nodes, nil
}

// TODO: Support >=, <=
func parseFilter(token string) (Node, error) {
	var op FilterOp

	switch {
	case strings.Contains(token, "~"):
		op = OpApprox
	case strings.Contains(token, "="):
		op = OpEquals
	case strings.Contains(token, "<"):
		op = OpLessThan
	case strings.Contains(token, ">"):
		op = OpGreaterThan
	default:
		return nil, nil
	}

	parts := strings.SplitN(token, op.String(), 2)
	field := parts[0]
	value := strings.Trim(parts[1], `"`)

	return FilterNode{Field: field, Op: op, Value: value}, nil
}
