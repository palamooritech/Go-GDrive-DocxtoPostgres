package domain

import (
	"31arthur/drive-editor/helper"
	"net/http"
)

func (s *APIServer) HandleRequests(w http.ResponseWriter, r *http.Request) error {
	list := []string{"samples"}
	return helper.WriteJSON(w, http.StatusOK, list)
}
