package query

import (
	"database/sql"
	"strings"
)

type ColumnType int
type Schema map[string]ColumnInfo

type ColumnInfo struct {
	Name    string
	RawType string // Raw SQL type.
	Type    ColumnType
}

const (
	TypeText ColumnType = iota
	TypeInt
	TypeReal
	TypeTime
	TypeUnknown
)

func (columnType ColumnType) String() string {
	switch columnType {
	case TypeText:
		return "TEXT"
	case TypeInt:
		return "INTEGER"
	case TypeReal:
		return "REAL"
	case TypeTime:
		return "DATETIME"
	default:
		return "UNKNOWN"
	}
}

func GetColumns(db *sql.DB) (Schema, error) {
	rows, err := db.Query("PRAGMA table_info(History);")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	schema := make(Schema)
	for rows.Next() {
		var columnId int
		var name string
		var rawType string
		var isNotNull int
		var defaultValue sql.NullString
		var isPrimaryKey bool
		if err := rows.Scan(&columnId, &name, &rawType, &isNotNull,
			&defaultValue, &isPrimaryKey); err != nil {
			return nil, err
		}

		schema[strings.ToLower(name)] = ColumnInfo{Name: name, RawType: rawType, Type: normalizedType(rawType)}
	}

	return schema, nil
}

func normalizedType(rawType string) ColumnType {
	rawType = strings.ToUpper(rawType)

	switch {
	case strings.Contains(rawType, "INT"):
		return TypeInt
	case strings.Contains(rawType, "CHAR"), strings.Contains(rawType, "TEXT"), strings.Contains(rawType, "CLOB"):
		return TypeText
	case strings.Contains(rawType, "REAL"), strings.Contains(rawType, "FLOA"), strings.Contains(rawType, "DOUB"):
		return TypeReal
	case strings.Contains(rawType, "DATE"), strings.Contains(rawType, "TIME"):
		return TypeTime
	default:
		return TypeUnknown
	}
}
