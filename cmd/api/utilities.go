package api

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Message string
}

func respondFailure(w http.ResponseWriter) {
	var success Response
	success.Message = "FAILURE"
	renderJSON(w, success)
}

func respondSuccess(w http.ResponseWriter) {
	var success Response
	success.Message = "SUCCESS"
	renderJSON(w, success)
}

func renderJSON(w http.ResponseWriter, v interface{}) error {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	return nil
}

func parseJSON(w http.ResponseWriter, r *http.Request, v interface{}) error {
	err := json.NewDecoder(r.Body).Decode(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}
	return nil
}
