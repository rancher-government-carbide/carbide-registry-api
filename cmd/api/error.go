package api

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	ErrorMessage string `json:"error_message"`
}

func http_json_error(w http.ResponseWriter, error string, http_status_code int) error {

	response := ErrorResponse{
		ErrorMessage: error,
	}

	json_response, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	w.WriteHeader(http_status_code)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(json_response)
	if err != nil {
		return err
	}

	return nil
}
