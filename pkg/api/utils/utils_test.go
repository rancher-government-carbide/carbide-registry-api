package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRespondError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := RespondError(w, "Test Error Message", http.StatusBadRequest)
		if err != nil {
			t.Fatalf("httpJSONError failed: %v", err)
		}
	}))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		t.Fatalf("HTTP request failed: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, res.StatusCode)
	}
	var jsonResponse ErrorResponse
	err = json.NewDecoder(res.Body).Decode(&jsonResponse)
	if err != nil {
		t.Fatalf("Failed to decode JSON response: %v", err)
	}
	expectedErrorMessage := "Test Error Message"
	if jsonResponse.ErrorMessage != expectedErrorMessage {
		t.Errorf("Expected error message %s, got %s", expectedErrorMessage, jsonResponse.ErrorMessage)
	}
}
