package rest

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/stazoloto/calendar/internal/domain"
)

type Events interface {
	Create(ctx context.Context, event *domain.Event) error
	GetByID(ctx context.Context, id int64) (*domain.Event, error)
	Update(ctx context.Context, id int64, inp *domain.UpdateEventInput) error
	Delete(ctx context.Context, id int64) error
	GetAll(ctx context.Context, from, to time.Time) ([]domain.Event, error)
}

type Handler struct {
	eventService Events
}

func NewHandler(events Events) *Handler {
	return &Handler{
		eventService: events,
	}
}

func (h *Handler) InitRouter() *mux.Router {
	r := mux.NewRouter()

	events := r.PathPrefix("/events").Subrouter()
	{
		events.HandleFunc("", h.createEvent).Methods(http.MethodPost)
		events.HandleFunc("", h.getAllEvents).Methods(http.MethodGet)
		events.HandleFunc("/{id:[0-9]+}", h.getEventByID).Methods(http.MethodGet)
		events.HandleFunc("/{id:[0-9]+}", h.updateEvent).Methods(http.MethodPut)
		events.HandleFunc("/{id:[0-9]+}", h.deleteEvent).Methods(http.MethodDelete)
	}
	return r
}

func (h *Handler) getEventByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := getIdFromRequest(r)
	if err != nil {
		log.Println("getEventByID() error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	event, err := h.eventService.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrEventNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		log.Println("getEventByID() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(event)
	if err != nil {
		log.Println("getEventByID() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(response)
}

func (h *Handler) getAllEvents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	query := r.URL.Query()

	from, err := time.Parse(time.RFC3339, query.Get("from"))
	if err != nil {
		http.Error(w, "invalid 'from' parameter", http.StatusBadRequest)
		return
	}

	to, err := time.Parse(time.RFC3339, query.Get("to"))
	if err != nil {
		http.Error(w, "invalid 'to' parameter", http.StatusBadRequest)
		return
	}

	events, err := h.eventService.GetAll(ctx, from, to)
	if err != nil {
		log.Println("getAllEvents() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(events)
	if err != nil {
		log.Println("getAllEvents() marshal error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(response)
}

func (h *Handler) createEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqBytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var event domain.Event
	if err = json.Unmarshal(reqBytes, &event); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.eventService.Create(ctx, &event)
	if err != nil {
		log.Println("createEvent() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) deleteEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := getIdFromRequest(r)
	if err != nil {
		log.Println("deleteEvent() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = h.eventService.Delete(ctx, id); err != nil {
		log.Println("deleteEvent() error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) updateEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := getIdFromRequest(r)
	if err != nil {
		log.Println("updateEvent() error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	reqBytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var inp domain.UpdateEventInput
	if err = json.Unmarshal(reqBytes, &inp); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.eventService.Update(ctx, id, &inp)
	if err != nil {
		log.Println("error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func getIdFromRequest(r *http.Request) (int64, error) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		return 0, err
	}

	if id == 0 {
		return 0, errors.New("id can't be 0")
	}

	return id, nil
}
