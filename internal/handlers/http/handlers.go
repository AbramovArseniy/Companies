package http

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/AbramovArseniy/Companies/internal/cfg"
	db "github.com/AbramovArseniy/Companies/internal/storage/postgres/db"
	"github.com/go-chi/chi"
	_ "github.com/jackc/pgx/v5/stdlib"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const contentTypeJSON = "application/json"

type httpHandler struct {
	Storage db.Querier
}

func New(cfg *cfg.Config) *httpHandler {
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

func (h *httpHandler) GetTreeHandler(w http.ResponseWriter, r *http.Request) {
	tree, err := h.Storage.GetAllTree(context.Background())
	if err != nil {
		log.Println("error while getting tree from database:", err)
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	jsonTree, err := json.MarshalIndent(tree, "", "  ")
	if err != nil {
		log.Println("error while marshaling json:", err)
		http.Error(w, "encoding json", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(jsonTree)
	if err != nil {
		log.Println("error while writing response body:", err)
		http.Error(w, "error while writing response body", http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", contentTypeJSON)
	w.WriteHeader(http.StatusOK)
}

func (h *httpHandler) GetHierarchyHandler(w http.ResponseWriter, r *http.Request) {}

func (h *httpHandler) GetNodeHandler(w http.ResponseWriter, r *http.Request) {}

func (h *httpHandler) Route() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.GetTreeHandler)
	r.Get("/hierarchy/", h.GetHierarchyHandler)
	r.Get("/node/{id}/", h.GetNodeHandler)
	return r
}
