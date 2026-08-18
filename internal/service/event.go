package service

import (
	"context"
	"time"

	"github.com/stazoloto/calendar/internal/domain"
)

type EventsRepository interface {
	Create(ctx context.Context, event *domain.Event) error
	GetByID(ctx context.Context, id int64) (*domain.Event, error)
	Update(ctx context.Context, id int64, inp *domain.UpdateEventInput) error
	Delete(ctx context.Context, id int64) error
	GetAll(ctx context.Context, from, to time.Time) ([]domain.Event, error)
}

type Events struct {
	repo EventsRepository
}

func NewEvents(repo EventsRepository) *Events {
	return &Events{repo: repo}
}

func (e *Events) Create(ctx context.Context, event *domain.Event) error {
	if event.CreatedAt.IsZero() {
		event.CreatedAt = time.Now()
	}

	if event.Title == "" {
		return domain.ErrEmptyTitle
	}

	if event.EndAt.Before(event.StartAt) {
		return domain.ErrInvalidDateRange
	}

	return e.repo.Create(ctx, event)
}

func (e *Events) GetByID(ctx context.Context, id int64) (*domain.Event, error) {
	return e.repo.GetByID(ctx, id)
}

func (e *Events) GetAll(ctx context.Context, from, to time.Time) ([]domain.Event, error) {
	return e.repo.GetAll(ctx, from, to)
}

func (e *Events) Update(ctx context.Context, id int64, inp *domain.UpdateEventInput) error {
	return e.repo.Update(ctx, id, inp)
}

func (e *Events) Delete(ctx context.Context, id int64) error {
	return e.repo.Delete(ctx, id)
}
