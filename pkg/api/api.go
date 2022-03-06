package api

import (
	"NFTracker/pkg/db"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-pg/pg/v10"
	"log"
	"net/http"
	"strconv"
)

func NewAPI(pgdb *pg.DB) *chi.Mux {
	// setup router
	r := chi.NewRouter()
	r.Use(middleware.Logger, middleware.WithValue("DB", pgdb))
	r.Route("/homes", func(r chi.Router) {
		r.Post("/", createHome)               // GET /articles/search
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
	AgentID     int64  `json:"agent_id"`
}

type HomeResponse struct {
	Success bool     `json:"success"`
	Error   string   `json:"error"`
	Home    *db.Home `json:"home"`
}

type CreateHomeResponse struct {
	Success bool     `json:"success"`
	Error   string   `json:"error"`
	Home    *db.Home `json:"home"`
}

func createHome(w http.ResponseWriter, r *http.Request) {
	// parse in the request body
	req := &CreateHomeRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		res := &HomeResponse{
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

	// get the db from context
	pgdb, ok := r.Context().Value("DB").(*pg.DB)
	if !ok {
		res := &HomeResponse{
			Success: false,
			Error:   "could not get databse from context",
			Home:    nil,
		}
		err = json.NewEncoder(w).Encode(res)
		if err != nil {
			log.Printf("error sending resopnse: %v\n", err)
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// insert our home
	home, err := db.CreateHome(pgdb, &db.Home{
		Price:       req.Price,
		Description: req.Description,
		Address:     req.Address,
		AgentID:     req.AgentID,
	})
	if err != nil {
		res := &HomeResponse{
			Success: false,
			Error:   err.Error(),
			Home:    nil,
		}
		err = json.NewEncoder(w).Encode(res)
		if err != nil {
			log.Printf("error sending response: %v\n", err)
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// return a response
	res := &HomeResponse{
		Success: true,
		Error:   "",
		Home:    home,
	}
	_ = json.NewEncoder(w).Encode(res)
	w.WriteHeader(http.StatusOK)
}

func getHomeByID(w http.ResponseWriter, r *http.Request) {
	homeID := chi.URLParam(r, "homeID")

	// get the database from context
	pgdb, ok := r.Context().Value("DB").(*pg.DB)
	if !ok {
		res := &HomeResponse{
			Success: false,
			Error:   "could not get databse from context",
			Home:    nil,
		}
		err := json.NewEncoder(w).Encode(res)
		if err != nil {
			log.Printf("error sending resopnse: %v\n", err)
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// query for the home
	home, err := db.GetHome(pgdb, homeID)
	if err != nil {
		res := &HomeResponse{
			Success: false,
			Error:   err.Error(),
			Home:    nil,
		}
		err = json.NewEncoder(w).Encode(res)
		if err != nil {
			log.Printf("error sending response: %v\n", err)
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// return a response
	res := &HomeResponse{
		Success: true,
		Error:   "",
		Home:    home,
	}
	_ = json.NewEncoder(w).Encode(res)
	w.WriteHeader(http.StatusOK)
}

type HomesResponse struct {
	Success bool       `json:"success"`
	Error   string     `json:"error"`
	Home    []*db.Home `json:"home"`
}

func getHomes(w http.ResponseWriter, r *http.Request) {
	// get the database from context
	pgdb, ok := r.Context().Value("DB").(*pg.DB)
	if !ok {
		res := &HomesResponse{
			Success: false,
			Error:   "could not get databse from context",
			Home:    nil,
		}
		err := json.NewEncoder(w).Encode(res)
		if err != nil {
			log.Printf("error sending resopnse: %v\n", err)
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	homes, err := db.GetHomes(pgdb)
	if err != nil {
		res := &HomesResponse{
			Success: false,
			Error:   err.Error(),
			Home:    nil,
		}
		err = json.NewEncoder(w).Encode(res)
		if err != nil {
			log.Printf("error sending response: %v\n", err)
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// send response
	res := &HomesResponse{
		Success: true,
		Error:   "",
		Home:    homes,
	}
	_ = json.NewEncoder(w).Encode(res)
	w.WriteHeader(http.StatusOK)
}

type UpdateHomeByIDRequest struct {
	Price       int64  `json:"price"`
	Description string `json:"description"`
	Address     string `json:"address"`
	AgentID     int64  `json:"agent_id"`
}

func updateHomeById(w http.ResponseWriter, r *http.Request) {
	homeID := chi.URLParam(r, "homeID")
	intHomeID, err := strconv.ParseInt(homeID, 10, 64)
	if err != nil {
		res := &HomeResponse{
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

	// parse in the request body
	req := &UpdateHomeByIDRequest{}
	err = json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		res := &HomeResponse{
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

	// get the db from context
	pgdb, ok := r.Context().Value("DB").(*pg.DB)
	if !ok {
		res := &HomeResponse{
			Success: false,
			Error:   "could not get databse from context",
			Home:    nil,
		}
		err = json.NewEncoder(w).Encode(res)
		if err != nil {
			log.Printf("error sending resopnse: %v\n", err)
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// update home
	home, err := db.UpdateHome(pgdb, &db.Home{
		ID:          intHomeID,
		Price:       req.Price,
		Description: req.Description,
		Address:     req.Address,
		AgentID:     req.AgentID,
	})
	if err != nil {
		res := &HomeResponse{
			Success: false,
			Error:   err.Error(),
			Home:    nil,
		}
		err = json.NewEncoder(w).Encode(res)
		if err != nil {
			log.Printf("error sending response: %v\n", err)
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// return a response
	res := &HomeResponse{
		Success: true,
		Error:   "",
		Home:    home,
	}
	_ = json.NewEncoder(w).Encode(res)
	w.WriteHeader(http.StatusOK)
}

func deleteHomeById(w http.ResponseWriter, r *http.Request) {
	homeID := chi.URLParam(r, "homeID")
	intHomeID, err := strconv.ParseInt(homeID, 10, 64)
	if err != nil {
		res := &HomeResponse{
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

	// get the db from context
	pgdb, ok := r.Context().Value("DB").(*pg.DB)
	if !ok {
		res := &HomeResponse{
			Success: false,
			Error:   "could not get databse from context",
			Home:    nil,
		}
		err = json.NewEncoder(w).Encode(res)
		if err != nil {
			log.Printf("error sending resopnse: %v\n", err)
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// delete home
	err = db.DeleteHome(pgdb, intHomeID)
	if err != nil {
		res := &HomeResponse{
			Success: false,
			Error:   err.Error(),
			Home:    nil,
		}
		err = json.NewEncoder(w).Encode(res)
		if err != nil {
			log.Printf("error sending response: %v\n", err)
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// return a response
	res := &HomeResponse{
		Success: true,
		Error:   "",
		Home:    nil,
	}
	_ = json.NewEncoder(w).Encode(res)
	w.WriteHeader(http.StatusOK)
}
