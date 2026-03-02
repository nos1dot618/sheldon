package database

import (
	"database/sql"
	q "sheldon/internal/query"
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
func SearchCommand(filters []string) ([]HistoricCommand, error) {
	db, err := openDb()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	err = initDb(db)
	if err != nil {
		return nil, err
	}

	nodes, err := q.Parse(filters)
	if err != nil {
		return nil, err
	}

	schema, err := q.GetColumns(db)

	query, args, err := q.BuildQuery(nodes, schema)
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(query, args...)
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
