package http

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

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

func (h *httpHandler) GetHierarchyHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Println("error while reading url parameter:", err)
		http.Error(w, "error while reading url parameter", http.StatusBadRequest)
		return
	}
	hierarchy, err := h.Storage.GetHierarchy(context.Background(), sql.NullInt32{Int32: int32(id), Valid: true})
	if errors.Is(err, sql.ErrNoRows) {
		log.Println("error while getting hierarchy from database:", err)
		http.Error(w, "no such node", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Println("error while getting tree from database:", err)
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	jsonHierarchy, err := json.MarshalIndent(hierarchy, "", "  ")
	if err != nil {
		log.Println("error while marshaling json:", err)
		http.Error(w, "encoding json", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(jsonHierarchy)
	if err != nil {
		log.Println("error while writing response body:", err)
		http.Error(w, "error while writing response body", http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", contentTypeJSON)
	w.WriteHeader(http.StatusOK)
}

func (h *httpHandler) GetNodeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Println("error while reading url parameter:", err)
		http.Error(w, "error while reading url parameter", http.StatusBadRequest)
		return
	}
	node, err := h.Storage.GetOneNode(context.Background(), sql.NullInt32{Int32: int32(id), Valid: true})
	if errors.Is(err, sql.ErrNoRows) {
		log.Println("error while getting node from database:", err)
		http.Error(w, "no such node", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Println("error while getting tree from database:", err)
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	jsonNode, err := json.MarshalIndent(node, "", "  ")
	if err != nil {
		log.Println("error while marshaling json:", err)
		http.Error(w, "encoding json", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(jsonNode)
	if err != nil {
		log.Println("error while writing response body:", err)
		http.Error(w, "error while writing response body", http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", contentTypeJSON)
	w.WriteHeader(http.StatusOK)
}

func (h *httpHandler) Route() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.GetTreeHandler)
	r.Get("/hierarchy/{id}", h.GetHierarchyHandler)
	r.Get("/node/{id}/", h.GetNodeHandler)
	return r
}
