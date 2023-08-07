package http

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/AbramovArseniy/Companies/internal/cfg"
	db "github.com/AbramovArseniy/Companies/internal/storage/postgres/db"
	"github.com/go-chi/chi"
)

type httpHandler struct {
	Storage db.Querier
}

func New(cfg cfg.Config) *httpHandler {
	database, err := sql.Open("pgx", cfg.DBAddress)
	if err != nil {
		log.Println("error while opening database:", err)
		return nil
	}
	querier := db.New(database)
	return &httpHandler{
		Storage: querier,
	}
}

func (h *httpHandler) GetTreeHandler(w http.ResponseWriter, r *http.Request) {}

func (h *httpHandler) GetHierarchyHandler(w http.ResponseWriter, r *http.Request) {}

func (h *httpHandler) GetNodeHandler(w http.ResponseWriter, r *http.Request) {}

func (h *httpHandler) Route() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.GetTreeHandler)
	r.Get("/hierarchy/", h.GetHierarchyHandler)
	r.Get("/node/{id}/", h.GetNodeHandler)
	return r
}
