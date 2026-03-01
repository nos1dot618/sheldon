package query

type Node interface {
	isNode()
}

type TextNode struct {
	Value string
}

func (TextNode) isNode() {}

type FilterOp int

const (
	OpApprox FilterOp = iota
	OpEquals
	OpLessThan
	OpGreaterThan
)

func (op FilterOp) String() string {
	switch op {
	case OpApprox:
		return "~"
	case OpEquals:
		return "="
	case OpLessThan:
		return "<"
	case OpGreaterThan:
		return ">"
	default:
		return "unknown"
	}
}

type FilterNode struct {
	Field string
	Op    FilterOp
	Value string
}

func (FilterNode) isNode() {}
