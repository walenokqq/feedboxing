package postgres

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type StoragePG struct {
	DB *sql.DB
}

func NewStoragePG(dataSourceName string) (*StoragePG, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	log.Println("Database connection established")
	return &StoragePG{DB: db}, nil
}

func (d *StoragePG) InitTables() error {
	// Создаём тип enum для статуса фидбэка
	if _, err := d.DB.Exec(`
		DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'feedback_status') THEN
				CREATE TYPE feedback_status AS ENUM ('pending', 'reviewed', 'archived');
			END IF;
		END$$;
	`); err != nil {
		return fmt.Errorf("failed to create feedback_status type: %v", err)
	}

	tables := []string{
		`CREATE TABLE IF NOT EXISTS projects (
			id SERIAL PRIMARY KEY,
			title VARCHAR(64) NOT NULL,
			description TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,

		`CREATE TABLE IF NOT EXISTS forms (
			id SERIAL PRIMARY KEY,
			project_id INTEGER NOT NULL,
			title VARCHAR(64) NOT NULL,
			description TEXT NOT NULL,
			schema JSONB NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
		)`,

		`CREATE TABLE IF NOT EXISTS feedback (
			id SERIAL PRIMARY KEY,
			form_id INTEGER NOT NULL,
			data JSONB NOT NULL,
			status feedback_status NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (form_id) REFERENCES forms(id) ON DELETE CASCADE
		)`,
	}

	for _, table := range tables {
		if _, err := d.DB.Exec(table); err != nil {
			return fmt.Errorf("failed to create table: %v", err)
		}
	}

	log.Println("Tables initialized successfully")
	return nil
}

func (d *StoragePG) Close() error {
	return d.DB.Close()
}
