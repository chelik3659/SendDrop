package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	
	"SendDrop/internal/storage"
)

type Handler struct {
	storage *storage.Manager
	server  *Server
}

type FileResponse struct {
	Files []storage.FileInfo `json:"files"`
	Total int                `json:"total"`
}

func NewHandler(shareDir string) *Handler {
	manager, _ := storage.NewManager(shareDir)
	return &Handler{
		storage: manager,
	}
}

func (h *Handler) SetServer(s *Server) {
	h.server = s
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	
	switch {
	case r.Method == "GET" && r.URL.Path == "/api/files":
		h.handleListFiles(w, r)
	case r.Method == "POST" && r.URL.Path == "/api/upload":
		h.handleUpload(w, r)
	case r.Method == "GET" && strings.HasPrefix(r.URL.Path, "/api/download/"):
		h.handleDownload(w, r)
	case r.Method == "DELETE" && strings.HasPrefix(r.URL.Path, "/api/delete/"):
		h.handleDelete(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (h *Handler) handleListFiles(w http.ResponseWriter, r *http.Request) {
	files, err := h.storage.ListFiles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	response := FileResponse{
		Files: files,
		Total: len(files),
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) handleUpload(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(100 << 30); err != nil { // 100GB max
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()
	
	filename := header.Filename
	if err := h.storage.SaveFile(filename, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "File uploaded successfully",
		"file":    filename,
	})
}

func (h *Handler) handleDownload(w http.ResponseWriter, r *http.Request) {
	filename := strings.TrimPrefix(r.URL.Path, "/api/download/")
	
	filePath, err := h.storage.GetFile(filename)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	http.ServeFile(w, r, filePath)
}

func (h *Handler) handleDelete(w http.ResponseWriter, r *http.Request) {
	filename := strings.TrimPrefix(r.URL.Path, "/api/delete/")
	
	if err := h.storage.DeleteFile(filename); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "File deleted successfully",
	})
}

func (h *Handler) GetStorage() *storage.Manager {
	return h.storage
}