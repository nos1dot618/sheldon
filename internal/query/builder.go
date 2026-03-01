package query

import (
	"fmt"
	db "sheldon/internal/database"
	time "sheldon/internal/timeutil"
	"strconv"
	"strings"
)

func BuildQuery(nodes []Node, schema db.Schema) (string, []interface{}, error) {
	var query strings.Builder
	query.WriteString("SELECT * FROM History WHERE 1=1")

	var args []interface{}

	for _, node := range nodes {
		switch derived := node.(type) {
		case TextNode:
			query.WriteString(" AND Command LIKE ?")
			args = append(args, "%"+derived.Value+"%")
		case FilterNode:
			filter, arg, err := buildFilterNodeQuery(derived, schema)
			if err != nil {
				return "", nil, err
			}
			query.WriteString(" AND ")
			query.WriteString(filter)
			args = append(args, arg)
		}
	}

	query.WriteString(" ORDER BY Timestamp DESC")
	return query.String(), args, nil
}

func buildFilterNodeQuery(node FilterNode, schema db.Schema) (string, interface{}, error) {
	column := schema[node.Field]

	switch column.Type {
	case db.TypeText:
		switch node.Op {
		case OpApprox:
			return column.Name + " LIKE ?", "%" + node.Value + "%", nil
		case OpEquals:
			return column.Name + "=?", node.Value, nil
		default:
			return "", nil, fmt.Errorf("Operator '%s' is not allowed for text fields like '%s'.", node.Op, column.Name)
		}

	case db.TypeInt:
		switch node.Op {
		case OpEquals, OpLessThan, OpGreaterThan:
			value, err := strconv.Atoi(node.Value)
			if err != nil {
				return "", nil, fmt.Errorf("Invalid integer value '%s' for field '%s'.", node.Value, column.Name)
			}
			return column.Name + node.Op.String() + "?", value, nil
		default:
			return "", nil, fmt.Errorf("Operator '%s' is not allowed for integer fields like '%s'.",
				node.Op, column.Name)
		}

	case db.TypeReal:
		switch node.Op {
		case OpEquals, OpLessThan, OpGreaterThan:
			value, err := strconv.ParseFloat(node.Value, 64)
			if err != nil {
				return "", nil, fmt.Errorf("Invalid real value '%s' for field '%s'.", node.Value, column.Name)
			}
			return column.Name + node.Op.String() + "?", value, nil
		default:
			return "", nil, fmt.Errorf("Operator '%s' is not allowed for real fields like '%s'.",
				node.Op, column.Name)
		}

	case db.TypeTime:
		switch node.Op {
		case OpEquals, OpLessThan, OpGreaterThan:
			value, err := time.ParseTime(node.Value)
			if err != nil {
				return "", nil, fmt.Errorf("Invalid datetime value '%s' for field '%s'.", node.Value, column.Name)
			}
			return column.Name + node.Op.String() + "?", value, nil
		default:
			return "", nil, fmt.Errorf("Operator '%s' is not allowed for datetime fields like '%s'.",
				node.Op, column.Name)
		}

	default:
		return "", nil, fmt.Errorf("Unknown column type found '%s'.", column.Type)
	}
}
