// +build integration

package azure

import (
	"testing"
	"time"
)

const CUSTOMER_ID_PREFIX = "DELETE-ME-TEST-"

func TestNewAzureClients(t *testing.T) {
	clientFactory, err := NewAzureClients()
	if err != nil {
		t.Error(err)
	}
	_ = clientFactory
}

func TestCreateCarbideAccount(t *testing.T) {
	clientFactory, err := NewAzureClients()
	if err != nil {
		t.Error(err)
	}
	daysTillExpiry := 10
	expiry := time.Now().Add(time.Hour*24*time.Duration(daysTillExpiry))
	customerID := CUSTOMER_ID_PREFIX+time.Now().Format("01-02-04-05")
	token, password, err := CreateCarbideAccount(clientFactory, customerID, expiry)
	if err != nil {
		t.Error(err)
	}
	t.Logf("Token: %v", token)
	t.Logf("Password: %v", *password)

}
