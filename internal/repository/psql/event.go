package psql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/stazoloto/calendar/internal/domain"
)

type Events struct {
	db *sql.DB
}

func NewEvents(db *sql.DB) *Events {
	return &Events{db}
}

// Create создает событие
func (e *Events) Create(ctx context.Context, event *domain.Event) error {
	_, err := e.db.Exec("INSERT INTO events (title, description, start_at, end_at, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)",
		event.Title, event.Description, event.StartAt, event.EndAt, event.CreatedAt, event.UpdatedAt)

	return err
}

// GetByID получает событие
func (e *Events) GetByID(ctx context.Context, id int64) (*domain.Event, error) {
	var event domain.Event
	err := e.db.QueryRow("SELECT id, title, description, start_at, end_at, created_at, updated_at FROM events WHERE id = $1", id).
		Scan(&event.ID, &event.Title, &event.Description, &event.StartAt, &event.EndAt, &event.CreatedAt, &event.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, domain.ErrEventNotFound
	}

	return &event, err
}

func (e *Events) GetAll(ctx context.Context, from, to time.Time) ([]domain.Event, error) {
	rows, err := e.db.QueryContext(ctx, 
		"SELECT id, title, description, start_at, end_at, created_at, updated_at FROM events WHERE start_at >= $1 AND start_at <= $2",
		from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := make([]domain.Event, 0)
	for rows.Next() {
		var event domain.Event
		if err = rows.Scan(&event.ID, &event.Title, &event.Description, &event.StartAt, &event.EndAt, &event.CreatedAt, &event.UpdatedAt); err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	return events, rows.Err()
}

// Update обновляет событие
func (e *Events) Update(ctx context.Context, id int64, inp *domain.UpdateEventInput) error {
	setValues := make([]string, 0)
	args := make([]any, 0)
	argId := 1

	if inp.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *inp.Title)
		argId++
	}

	if inp.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *inp.Description)
		argId++
	}

	if inp.StartAt != nil {
		setValues = append(setValues, fmt.Sprintf("start_at=$%d", argId))
		args = append(args, *inp.StartAt)
		argId++
	}

	if inp.EndAt != nil {
		setValues = append(setValues, fmt.Sprintf("end_at=$%d", argId))
		args = append(args, *inp.EndAt)
		argId++
	}

	if inp.UpdatedAt != nil {
		setValues = append(setValues, fmt.Sprintf("updated_at=$%d", argId))
		args = append(args, *inp.UpdatedAt)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE events SET %s WHERE id=$%d", setQuery, argId)
	args = append(args, id)

	_, err := e.db.Exec(query, args...)
	return err
}

// Delete удаляет событие
func (e *Events) Delete(ctx context.Context, id int64) error {
	_, err := e.db.Exec("DELETE FROM events WHERE id=$1", id)
	return err
}
