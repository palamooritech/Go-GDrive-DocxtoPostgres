package domain

import (
	"31arthur/drive-editor/helper"
	typos "31arthur/drive-editor/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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
		fmt.Println("Hello, this has been called")
		if err := json.NewDecoder(r.Body).Decode(eGFile); err != nil {
			fmt.Println(err)
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

/*
	 the data format
		id	"1W0oB8ml_PhQzh1idpG2GoG_wyZ3HLkpu"
		lid	"DCL Court PoW Case Filings"
		file_name	"509efc74-db46-41ac-96e8-5c1c23938575"
		created_time	"2023-10-04T19:41:55.14+05:30"
		modified_time	"2023-10-04T19:15:03.779+05:30"
		touched	true
		case_number	"CATU143"
		letter_type	"hello"
		summary	"Yellow"
		delivery_mode	"yoloadsadasd2"
		delivery_id	"telloaadskalsadkladsklds2"
		file_url	"https://drive.google.com…Z3HLkpu/view?usp=sharing"
*/
func (s *APIServer) HandleAccessAllRequests(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		list := s.Store.AccessAll()
		return helper.WriteJSON(w, http.StatusOK, map[string]any{"payload": list})
	}
	return helper.WriteJSON(w, http.StatusBadRequest, map[string]string{"payload": "You can't GET a POST with this request :P"})
}

func (s *APIServer) HandleSearchAllRequests(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "POST" {
		searchData := new(typos.SearchData)
		if err := json.NewDecoder(r.Body).Decode(searchData); err != nil {
			return err
		}
		list := s.Store.SearchAll(searchData.Keyword)
		if len(list) == 0 {
			return helper.WriteJSON(w, http.StatusOK, map[string]any{"payload": "No Data"})
		}
		return helper.WriteJSON(w, http.StatusOK, map[string]any{"payload": list})
	}
	return helper.WriteJSON(w, http.StatusBadRequest, map[string]string{"payload": "You can't GET a POST with this request :P"})
}

func (s *APIServer) HandleResponseUpload(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "POST" {
		r.ParseMultipartForm(10 << 20) // Set a reasonable maximum file size
		file, handler, err := r.FormFile("file")
		if err != nil {
			fmt.Println("Error retrieving the file:", err)
			return err
		}
		defer file.Close()

		f, err := os.OpenFile(handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println("Error creating the file:", err)
			return err
		}
		defer f.Close()

		// Copy the uploaded file to the new file on the server
		io.Copy(f, file)

		jsonResponse := map[string]interface{}{"payload": "success"}

		// Set the response content type
		w.Header().Set("Content-Type", "application/json")

		// Write the JSON response
		return helper.WriteJSON(w, http.StatusOK, jsonResponse)
	}

	// For other HTTP methods (e.g., GET)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	jsonResponse := map[string]interface{}{"payload": "You can't GET a POST with this request :P"}
	return helper.WriteJSON(w, http.StatusBadRequest, jsonResponse)
}
