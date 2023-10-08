package domain

import (
	"31arthur/drive-editor/helper"
	typos "31arthur/drive-editor/models"
	"encoding/json"
	"net/http"
)

// basic random function to check
func (s *APIServer) HandleRequests(w http.ResponseWriter, r *http.Request) error {
	list := []string{"samples"}
	return helper.WriteJSON(w, http.StatusOK, list)
}

func (s *APIServer) HandleEditRequest(w http.ResponseWriter, r *http.Request) error {

	/*
		This is the sample post request for the HandleEditRequest method
		{
		"id": "1WS6xpxfcW1dOSnmS70QEfSP-njs-BWJJ",
		"caseNumber": "CATU143",
		"letterType": "hello",
		"summary": "",
		"touched":true,
		"deliveryMode": "yolo2",
		"deliveryID": "tello2"
		}
	*/

	if r.Method == "POST" {
		eGFile := new(typos.EGFile)
		if err := json.NewDecoder(r.Body).Decode(eGFile); err != nil {
			return err
		}
		// fmt.Println("Gfile", eGFile)
		if err := s.Store.UpdateFileRequest(*eGFile); err != nil {
			return err
		}
		return helper.WriteJSON(w, http.StatusOK, map[string]string{"payload": "Success"})
	}

	//returns this for every method other than POST
	return helper.WriteJSON(w, http.StatusBadRequest, map[string]string{"payload": "You can't get anything with this request :P"})

}

// this is used to fork out all the data from the db.
func (s *APIServer) HandleAccessAllRequests(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		list := s.Store.AccessAll()
		return helper.WriteJSON(w, http.StatusOK, list)
	}
	return helper.WriteJSON(w, http.StatusBadRequest, map[string]string{"payload": "You can't GET a POST with this request :P"})
}
