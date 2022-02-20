package api

import (
	"NFTracker/pkg/db"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func NewAPI() *chi.Mux {
	// setup router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Route("/homes", func(r chi.Router) {
		r.Post("/", createHomes)              // GET /articles/search
		r.Get("/{homeID}", getHomeByID)       // GET /articles/search
		r.Get("/", getHomes)                  // GET /articles/search
		r.Put("/{homeID}", updateHomeById)    // GET /articles/search
		r.Delete("/{homeID}", deleteHomeById) // GET /articles/search
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})

	return r
}

func createHomes(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("create home"))
}

func getHomeByID(w http.ResponseWriter, r *http.Request) {
	homeID := chi.URLParam(r, "homeID")

	w.Write([]byte(fmt.Sprintf("get home: %s", homeID)))
}

type GetHomesResponse struct {
	Homes []db.Home `json:"homes"`
}

func getHomes(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("get home"))
}

func updateHomeById(w http.ResponseWriter, r *http.Request) {
	homeID := chi.URLParam(r, "homeID")

	w.Write([]byte(fmt.Sprintf("update home: %s", homeID)))
}
func deleteHomeById(w http.ResponseWriter, r *http.Request) {
	homeID := chi.URLParam(r, "homeID")

	w.Write([]byte(fmt.Sprintf("delete home: %s", homeID)))
}
