package api

import (
	"NFTracker/pkg/db"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
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

type CreateHomeRequest struct {
	Price       int64  `json:"price"`
	Description string `json:"description"`
	Address     string `json:"address"`
	AgentID     int64  `pg:"rel:has-one" json:"agent"`
}

type CreateHomeResponse struct {
	Success bool     `json:"success"`
	Error   string   `json:"error"`
	Home    *db.Home `json:"home"`
}

func createHomes(w http.ResponseWriter, r *http.Request) {
	// parse in the request body
	req := &CreateHomeRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		res := &CreateHomeResponse{
			Success: false,
			Error:   err.Error(),
			Home:    nil,
		}
		err = json.NewEncoder(w).Encode(res)
		if err != nil {
			log.Printf("error sending resopnse: %v\n", err)
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// get the db somehow

	// insert our home

	// return a response
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
