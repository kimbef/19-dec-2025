package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// DB wraps the sql.DB connection
type DB struct {
	*sql.DB
}

// New creates a new database connection
func New(connectionString string) (*DB, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{db}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.DB.Close()
}

// Health checks the database health
func (db *DB) Health() error {
	return db.Ping()
}

// Migrate runs database migrations
func (db *DB) Migrate() error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS pipelines (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			status VARCHAR(50) DEFAULT 'active',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS builds (
			id SERIAL PRIMARY KEY,
			pipeline_id INTEGER NOT NULL REFERENCES pipelines(id),
			build_number INTEGER NOT NULL,
			status VARCHAR(50) DEFAULT 'pending',
			branch VARCHAR(255) NOT NULL,
			commit_hash VARCHAR(255) NOT NULL,
			started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			finished_at TIMESTAMP,
			duration INTEGER DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS deployments (
			id SERIAL PRIMARY KEY,
			build_id INTEGER NOT NULL REFERENCES builds(id),
			environment VARCHAR(50) NOT NULL,
			status VARCHAR(50) DEFAULT 'pending',
			deployed_by VARCHAR(255) NOT NULL,
			deployed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_pipelines_name ON pipelines(name)`,
		`CREATE INDEX IF NOT EXISTS idx_pipelines_status ON pipelines(status)`,
		`CREATE INDEX IF NOT EXISTS idx_builds_pipeline_id ON builds(pipeline_id)`,
		`CREATE INDEX IF NOT EXISTS idx_builds_status ON builds(status)`,
		`CREATE INDEX IF NOT EXISTS idx_deployments_build_id ON deployments(build_id)`,
		`CREATE INDEX IF NOT EXISTS idx_deployments_environment ON deployments(environment)`,
	}

	for _, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return fmt.Errorf("failed to run migration: %w", err)
		}
	}

	return nil
}
