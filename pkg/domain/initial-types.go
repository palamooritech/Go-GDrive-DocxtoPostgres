package domain

import (
	storage "31arthur/drive-editor/pkg/adapter/postgres"
	"net/http"
)

type APIServer struct {
	ListenAddr string
	Store      storage.PGXStorage
}

type APIFunc func(w http.ResponseWriter, r *http.Request) error

type APIError struct {
	Error string `json:"error"`
}
