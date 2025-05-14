package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	// _ "github.com/lib/pq"
)

type PostgresStorage struct {
	db *sql.DB
}

// New создает новое подключение к PostgreSQL
func New(connString string) (*PostgresStorage, error) {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	return &PostgresStorage{db: db}, nil
}

// Init создает таблицу если ее нет
func (s *PostgresStorage) Init(ctx context.Context) error {
	query := `
        BEGIN;
        
        CREATE TABLE IF NOT EXISTS projects (
            id SERIAL PRIMARY KEY,
            title VARCHAR(64) NOT NULL,
            description text NOT NULL,
            created_at TIMESTAMP NOT NULL DEFAULT NOW()
        );

        CREATE TABLE IF NOT EXISTS forms (         
            id SERIAL PRIMARY KEY,
            title VARCHAR(64) NOT NULL,
            description VARCHAR(64) NOT NULL,
            schema JSONB NOT NULL,
            created_at TIMESTAMP NOT NULL DEFAULT NOW()
        );

        CREATE TABLE IF NOT EXISTS feedback (
            id SERIAL PRIMARY KEY,
            data JSONB NOT NULL,
            status varchar(64) NOT NULL,
            created_at TIMESTAMP NOT NULL DEFAULT NOW()

            
        );

        COMMIT;
    `

	_, err := s.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	return nil
}

func (s *PostgresStorage) SaveFeedback(ctx context.Context)

func (s *PostgresStorage) SaveDeveloper(ctx context.Context, dev *models.Developer) error {
	query := `
        INSERT INTO developers (firstname, last_name) 
        VALUES ($1, $2)
        RETURNING id, created_at, modified_at
    `
	return s.db.QueryRowContext(ctx, query, dev.Firstname, dev.LastName).
		Scan(&dev.ID, &dev.CreatedAt, &dev.ModifiedAt)
}

func (s *PostgresStorage) GetDeveloper(ctx context.Context, id uint) (*models.Developer, error) {
	query := `
        SELECT 
            id, 
            firstname, 
            last_name, 
            created_at, 
            modified_at, 
            deleted_at 
        FROM developers 
        WHERE id = $1
    `

	var dev models.Developer
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&dev.ID,
		&dev.Firstname,
		&dev.LastName,
		&dev.CreatedAt,
		&dev.ModifiedAt,
		&dev.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("разработчик с ID %d не найден", id)
		}
		return nil, fmt.Errorf("ошибка при получении разработчика: %w", err)
	}

	return &dev, nil
}

// Close закрывает соединение с БД
func (s *PostgresStorage) Close() error {
	return s.db.Close()
}
