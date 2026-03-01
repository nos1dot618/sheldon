package query

import "strings"

func Tokenize(input string) []string {
	var tokens []string
	var current strings.Builder
	insideString := false

	for i := 0; i < len(input); i++ {
		ch := input[i]

		switch ch {
		case '"':
			insideString = !insideString
			current.WriteByte(ch)
		case ' ':
			if insideString {
				current.WriteByte(ch)
			} else {
				if current.Len() > 0 {
					tokens = append(tokens, current.String())
					current.Reset()
				}
			}
		default:
			current.WriteByte(ch)
		}
	}

	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}

	return tokens
}
