package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"io"
	"log"
	"net/http"
)

type ServerImpl struct {
	service Service
}

type Service interface {
	GetShortUrl(link []byte) (string, error)
	GetLongUrl(shortLink string) (string, error)
}

func NewServer(service Service) *ServerImpl {
	return &ServerImpl{
		service: service,
	}
}

func (s *ServerImpl) Handlers() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		shortLink, err := s.service.GetShortUrl(body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
		}
		_, err = w.Write([]byte(shortLink + "\n"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
		}
	})
	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		longLink, err := s.service.GetLongUrl(r.URL.Path)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
		}

		http.Redirect(w, r, longLink, http.StatusFound)
	})
	return r
}
