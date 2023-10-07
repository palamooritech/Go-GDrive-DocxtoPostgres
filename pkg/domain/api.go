package domain

import (
	"31arthur/drive-editor/helper"
	storage "31arthur/drive-editor/pkg/adapter/postgres"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func makeHTTPHandleFunc(f APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			helper.WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error() + " in MakeHTTPHandleFunc"})
		}
	}
}

func NewAPIServer(listenAddr string, store storage.PGXStorage) *APIServer {
	return &APIServer{
		ListenAddr: listenAddr,
		Store:      store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/api", makeHTTPHandleFunc(s.HandleRequests))

	log.Println("JSON API server running on", s.ListenAddr)
	http.ListenAndServe(s.ListenAddr, router)
}
