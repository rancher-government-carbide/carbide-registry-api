package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"carbide-registry-api/pkg/azure"
	"carbide-registry-api/pkg/objects"
)

func TestCreateCarbideAccountHandler(t *testing.T) {
	clientFactory, err := azure.NewAzureClients()
	if err != nil {
		t.Fatalf("failed to initialize azure client factory")
	}

	// Define test payloads
	testPayloads := []map[string]interface{}{
		{
			"customerID":     "DELETE-ME-API-HANDLER-TEST" + time.Now().Format("01-02-04-05"),
			"daysTillExpiry": 30,
			"nodeCount":      5,
		},
	}

	// Iterate over test payloads
	for _, testData := range testPayloads {
		inputPayloadJSON, err := json.Marshal(testData)
		if err != nil {
			t.Fatalf("Failed to marshal input payload to JSON: %v", err)
		}
		req, err := http.NewRequest("POST", "/carbide-account", bytes.NewBuffer(inputPayloadJSON))
		if err != nil {
			t.Fatalf("Failed to create HTTP request: %v", err)
		}
		rr := httptest.NewRecorder()
		handler := createCarbideAccountHandler(clientFactory)
		handler.ServeHTTP(rr, req)
		t.Logf("Got here 5")
		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusCreated)
		}

		var createdLicesnse objects.CarbideLicense
		err = json.Unmarshal(rr.Body.Bytes(), &createdLicesnse)
		if err != nil {
			t.Errorf("Failed to unmarshal response body: %v", err)
		}

	}

}
