package orders

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/KBingsoo/entities/pkg/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type handler struct {
	manager *manager
	router  *chi.Mux
}

func NewHandler(manager *manager) *handler {
	r := chi.NewRouter()
	h := handler{
		manager: manager,
		router:  r,
	}

	h.init()

	return &h
}

func (h *handler) init() {
	h.router.Get("/{id}", h.getOrder)
	h.router.Post("/", h.createOrder)
}

func (h *handler) Routes() *chi.Mux {
	return h.router
}

func (h *handler) getOrder(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	card, err := h.manager.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, card)
}

func (h *handler) createOrder(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	order := new(models.Order)
	if err := json.Unmarshal(b, order); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.manager.Create(r.Context(), order); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, order)
}
