package database

import (
	"database/sql"
	"strings"
	"time"
)

// TODO: User may able to configure fields he wants to store, such as Shell, User, TimeElapsed
// Most likely every record has the same value, but could be useful for some users.
type HistoricCommand struct {
	Id        int
	Command   string
	Cwd       string
	ExitCode  int
	Timestamp time.Time
}

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

		schema[name] = ColumnInfo{Name: name, RawType: rawType, Type: normalizedType(rawType)}
	}

	return schema, nil
}

// TODO: Add encryption.
func RecordCommand(command string, cwd string, exitcode int) error {
	db, err := openDb()
	if err != nil {
		return err
	}
	defer db.Close()

	err = initDb(db)
	if err != nil {
		return err
	}

	_, err = db.Exec(
		"INSERT INTO History (Command, Cwd, ExitCode) VALUES (?, ?, ?)",
		command,
		cwd,
		exitcode,
	)
	return err
}

// TODO: Only return a range of results.
// TODO: Support more expressive queries.
func SearchCommand(query string) ([]HistoricCommand, error) {
	db, err := openDb()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	err = initDb(db)
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(
		"SELECT * FROM History WHERE Command LIKE ? ORDER BY Id DESC",
		"%"+query+"%",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []HistoricCommand
	for rows.Next() {
		var command HistoricCommand
		if err := rows.Scan(&command.Id, &command.Command, &command.Cwd,
			&command.ExitCode, &command.Timestamp); err != nil {
			return nil, err
		}
		results = append(results, command)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

func initDb(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS History (
			Id INTEGER PRIMARY KEY AUTOINCREMENT,
			Command TEXT,
			Cwd TEXT,
			ExitCode INTEGER,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`

	_, err := db.Exec(query)
	return err
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
