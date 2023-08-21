// package http describes HTTP server's work
package http

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"

	db "github.com/AbramovArseniy/Companies/internal/storage/postgres/db"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const contentTypeJSON = "application/json"

// httpHandler describes http server
type httpHandler struct {
	Storage db.Querier
}

// New creates a new httpHandler with config
func New(dbPool *pgxpool.Pool) (*httpHandler, error) {
	dbConn, err := dbPool.Acquire(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error while acquiring database connection: %w", err)
	}
	storage := db.New(dbConn)

	return &httpHandler{
		Storage: storage,
	}, nil
}

// GetTree returns information about all the nodes in the tree
func (h *httpHandler) GetTreeHandler(w http.ResponseWriter, r *http.Request) {
	tree, err := h.Storage.GetAllTree(context.Background())
	if err != nil {
		log.Println("error while getting tree from database:", err)
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	jsonTree, err := json.Marshal(tree)
	if err != nil {
		log.Println("error while marshaling json:", err)
		http.Error(w, "encoding json", http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", contentTypeJSON)
	_, err = w.Write(jsonTree)
	if err != nil {
		log.Println("error while writing response body:", err)
		http.Error(w, "error while writing response body", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// GetHierarchy returns information about hierarchy of a node by the node id
func (h *httpHandler) GetHierarchyHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Println("error while reading url parameter:", err)
		http.Error(w, "error while reading url parameter", http.StatusBadRequest)
		return
	}
	hierarchy, err := h.Storage.GetHierarchy(context.Background(), int32(id))
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
	jsonHierarchy, err := json.Marshal(hierarchy)
	if err != nil {
		log.Println("error while marshaling json:", err)
		http.Error(w, "encoding json", http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", contentTypeJSON)
	_, err = w.Write(jsonHierarchy)
	if err != nil {
		log.Println("error while writing response body:", err)
		http.Error(w, "error while writing response body", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// GetNode returns information about one node by the node id
func (h *httpHandler) GetNodeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Println("error while reading url parameter:", err)
		http.Error(w, "error while reading url parameter", http.StatusBadRequest)
		return
	}
	node, err := h.Storage.GetOneNode(context.Background(), int32(id))
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
	jsonNode, err := json.Marshal(node)
	if err != nil {
		log.Println("error while marshaling json:", err)
		http.Error(w, "encoding json", http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", contentTypeJSON)
	_, err = w.Write(jsonNode)
	if err != nil {
		log.Println("error while writing response body:", err)
		http.Error(w, "error while writing response body", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Route creates an http router
func (h *httpHandler) Route() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.GetTreeHandler)
	r.Get("/hierarchy/{id}", h.GetHierarchyHandler)
	r.Get("/node/{id}", h.GetNodeHandler)
	return r
}
