package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	sqliteDriver "github.com/golang-migrate/migrate/v4/database/sqlite"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Migrator struct {
	m *migrate.Migrate
}

//todo: close connection closed all connections

func NewMigrator(db *sql.DB) (*Migrator, error) {
	driver, err := sqliteDriver.WithInstance(db, &sqliteDriver.Config{})
	if err != nil {
		return nil, fmt.Errorf("create migration driver: %w", err)
	}
	m, err := migrate.NewWithDatabaseInstance("file://internal/database/migrations", "sqlite", driver)
	if err != nil {
		return nil, fmt.Errorf("create migration: %w", err)
	}
	return &Migrator{m: m}, nil
}

func (m *Migrator) Close() error {
	sourceErr, databaseErr := m.m.Close()

	if sourceErr != nil {
		return fmt.Errorf("close migration source: %w", sourceErr)
	}

	if databaseErr != nil {
		return fmt.Errorf("close migration database: %w", databaseErr)
	}
	return nil
}

func (m *Migrator) Up() error {
	if err := m.m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("run migrations: %w", err)
	} else if errors.Is(err, migrate.ErrNoChange) {
		log.Printf("No migrations found")
	}
	return nil
}

func (m *Migrator) Down() error {
	if err := m.m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("run migrations: %w", err)
	} else if errors.Is(err, migrate.ErrNoChange) {
		log.Printf("No migrations found")
	}
	return nil
}

func (m *Migrator) Force(version int) error {
	if err := m.m.Force(version); err != nil {
		return fmt.Errorf("force migration: %w", err)
	}

	return nil
}
