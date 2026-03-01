package main

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func openDb() (*sql.DB, error) {
	configDirectory, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	// Create app directory, if not exists.
	appDirectory := filepath.Join(configDirectory, "Sheldon")
	os.MkdirAll(appDirectory, 0755)

	dbPath := filepath.Join(appDirectory, "history.db")
	return sql.Open("sqlite", dbPath)
}
