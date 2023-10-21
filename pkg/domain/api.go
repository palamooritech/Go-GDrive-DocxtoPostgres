package domain

import (
	"31arthur/drive-editor/helper"
	storage "31arthur/drive-editor/pkg/adapter/postgres"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// the makeHTTPHandleFunc is like a wrapper function which takes http handlers
// with the return type of error. This helps in achieving uniform error
// addressal throughout the application
func makeHTTPHandleFunc(f APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			helper.WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error() + " in MakeHTTPHandleFunc"})
		}
	}
}

// returns a APIServer immutable instance by using pointer reference, which
// helps in maintaining a reference to the PGSXStore throughout the
// application
func NewAPIServer(listenAddr string, store storage.PGXStorage) *APIServer {
	return &APIServer{
		ListenAddr: listenAddr,
		Store:      store,
	}
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Replace "*" with specific origins as needed
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *APIServer) Run() {
	//router initialization and the routes initialization
	router := mux.NewRouter()
	router.Use(enableCORS)
	router.HandleFunc("/api", makeHTTPHandleFunc(s.HandleRequests))
	router.HandleFunc("/api/edit", makeHTTPHandleFunc(s.HandleEditRequest))
	router.HandleFunc("/api/accessall", makeHTTPHandleFunc(s.HandleAccessAllRequests))
	router.HandleFunc("/api/searchall", makeHTTPHandleFunc(s.HandleSearchAllRequests))
	router.HandleFunc("/api/uploadresponse", makeHTTPHandleFunc(s.HandleResponseUpload))

	log.Println("JSON API server running on", s.ListenAddr)
	http.ListenAndServe(s.ListenAddr, router)
}
