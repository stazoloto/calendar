package service

import (
	"context"
	"testing"
	"time"

	"github.com/stazoloto/calendar/internal/domain"
)

type mockRepo struct {
	events []domain.Event
	err    error
}

func (m *mockRepo) Create(ctx context.Context, event *domain.Event) error {
	if m.err != nil {
		return m.err
	}

	event.ID = int64(len(m.events) + 1)
	m.events = append(m.events, *event)
	return nil
}

func (m *mockRepo) GetByID(ctx context.Context, id int64) (*domain.Event, error) {
	for _, e := range m.events {
		if e.ID == id {
			return &e, nil
		}
	}

	return nil, domain.ErrEventNotFound
}

func (m *mockRepo) Update(ctx context.Context, id int64, inp *domain.UpdateEventInput) error {
	return m.err
}

func (m *mockRepo) Delete(ctx context.Context, id int64) error {
	return m.err
}

func (m *mockRepo) GetAll(ctx context.Context, from, to time.Time, limit, offset int) ([]domain.Event, error) {
	if m.err != nil {
		return nil, m.err
	}

	return m.events, nil
}

func TestCreate_EmptyTitle(t *testing.T) {
	svc := NewEvents(&mockRepo{})
	err := svc.Create(context.Background(), &domain.Event{
		Title:   "",
		StartAt: time.Now(),
		EndAt:   time.Now().Add(time.Hour),
	})
	if err != domain.ErrEmptyTitle {
		t.Errorf("expected ErrEmptyTitle, got %v", err)
	}
}

func TestCreate_InvalidDateRange(t *testing.T) {
	svc := NewEvents(&mockRepo{})
	now := time.Now()
	err := svc.Create(context.Background(), &domain.Event{
		Title:   "test",
		StartAt: now,
		EndAt:   now.Add(-time.Hour),
	})
	if err != domain.ErrInvalidDateRange {
		t.Errorf("expected ErrInvalidDateRange, got %v", err)
	}
}

func TestCreate_Success(t *testing.T) {
	svc := NewEvents(&mockRepo{})
	now := time.Now()
	event := &domain.Event{
		Title:   "Test",
		StartAt: now,
		EndAt:   now.Add(time.Hour),
	}
	err := svc.Create(context.Background(), event)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if event.ID == 0 {
		t.Error("expected ID to be set")
	}

	if event.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
}
