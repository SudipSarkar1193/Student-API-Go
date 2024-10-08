package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/SudipSarkar1193/students-API-Go/internal/config"

	"github.com/SudipSarkar1193/students-API-Go/internal/types"

	_ "github.com/mattn/go-sqlite3"
)

func New(cfg *config.Config) (*types.Sqlite, error) {
	sqliteDB, err := sql.Open("sqlite3", cfg.StoragePath)

	if err != nil {
		return nil, err
	}

	sqlResult, err := sqliteDB.Exec(`CREATE TABLE IF NOT EXISTS students(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT,
	email TEXT,
	age INTEGER,
	)`)
	if err != nil {
		return nil, err
	}

	fmt.Println("sqlResult", sqlResult)

	return &types.Sqlite{
		Db: sqliteDB,
	}, nil

}
