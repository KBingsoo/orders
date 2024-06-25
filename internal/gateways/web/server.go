package web

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler interface {
	Routes() *chi.Mux
}
type server struct {
	router  *chi.Mux
	handler Handler
}

func NewServer(cardsHadler Handler) *server {
	router := chi.NewRouter()
	s := &server{
		router:  router,
		handler: cardsHadler,
	}

	s.Init()

	return s
}

func (s *server) Init() {
	s.router.Mount("/cards", s.handler.Routes())
}

func (s *server) Close() {
}

func (s *server) Run(port int) error {
	if port == 0 {
		port = 8080
	}

	return http.ListenAndServe(fmt.Sprintf(":%d", port), s.router)
}
