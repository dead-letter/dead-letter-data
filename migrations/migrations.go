package migrations

import (
	"database/sql"
	"embed"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var embedMigrations embed.FS

func Run(dsn string, dev bool) error {
	// Run migrations
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("failed to open db: %v", err)
	}
	defer db.Close()

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %v", err)
	}

	if dev {
		if err := goose.Reset(db, "migrations"); err != nil {
			return fmt.Errorf("failed to reset migrations: %v", err)
		}
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return fmt.Errorf("failed to apply migrations: %v", err)
	}

	return nil
}
