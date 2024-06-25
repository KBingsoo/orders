package cards

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
	h.router.Get("/cards", h.getCards)
	h.router.Get("/cards/{id}", h.getCard)
	h.router.Post("/cards", h.createCard)
}

func (h *handler) Routes() *chi.Mux {
	return h.router
}

func (h *handler) getCards(w http.ResponseWriter, r *http.Request) {
	cards, err := h.manager.GetAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, cards)
}

func (h *handler) getCard(w http.ResponseWriter, r *http.Request) {
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

func (h *handler) createCard(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	card := new(models.Card)
	if err := json.Unmarshal(b, card); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.manager.Create(r.Context(), card); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, card)
}
