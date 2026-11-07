package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"

	"github.com/oleksandr-vlasov/task-manager/internal/model"
)

type DB struct {
	Conn *sql.DB
}

// NewDB intialization new connection to the PostgreSQL
func NewDB(dataSourceName string) (*DB, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("error opening database %w", err)

	}

	// Test connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	log.Println("Successfully connected to PostgreSQL database!")

	instance := &DB{Conn: db}
	if err := instance.createTable(); err != nil {
		return nil, fmt.Errorf("error crating table: %w", err)
	}

	return instance, nil
}

// createTable creates table tasks, if this table is existing
func (db *DB) createTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS tasks (
		id SERIAL PRIMARY KEY,
		title TEXT NOT NULL,
		description TEXT,
		completed BOOLEAN NOT NULL DEFAULT FALSE,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL
	);`

	_, err := db.Conn.Exec(query)
	return err
}

// Create - adding new task in Postgresql
func (db *DB) Create(title, description string) (model.Task, error) {
	query := `
		INSERT INTO tasks (title, description, completed, created_at)
		VALUES ($1, $2, false, $3)
		RETURNING id, title, description, completed, created_at
	`

	task := model.Task{
		CreatedAt: time.Now(),
	}

	err := db.Conn.QueryRow(query, title, description, task.CreatedAt).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Completed,
		&task.CreatedAt,
	)
	return task, err
}

// GetAll retrievs all tasks from the database
func (db *DB) GetAll() ([]model.Task, error) {
	query := `SELECT id, title, description, completed, created_at FROM tasks ORDER BY id ASC`

	rows, err := db.Conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []model.Task
	for rows.Next() {
		var t model.Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Completed, &t.CreatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if tasks == nil {
		tasks = []model.Task{}
	}
	return tasks, nil
}
