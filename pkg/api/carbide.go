package api

import (
	"carbide-registry-api/pkg/api/utils"
	"carbide-registry-api/pkg/azure"
	"carbide-registry-api/pkg/license"
	"crypto/rsa"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerregistry/armcontainerregistry"
	log "github.com/sirupsen/logrus"
)

func createCarbideAccountHandler(clientFactory *armcontainerregistry.ClientFactory, privateKey *rsa.PrivateKey) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var newLicense license.CarbideLicense
		err := utils.DecodeJSONObject(w, r, &newLicense)
		if err != nil {
			return
		}
		if newLicense.CustomerID == nil || newLicense.DaysTillExpiry == nil || newLicense.NodeCount == nil {
			log.Errorf("missing customerID, daysTillExpiry, or nodeCount")
			utils.HttpJSONError(w, "missing customerID, daysTillExpiry, or nodeCount", http.StatusBadRequest)
			return
		}
		expiry := time.Now().Add(time.Hour * 24 * time.Duration(*newLicense.DaysTillExpiry))
		newLicense.LicenseString, err = license.CreateCarbideLicense(privateKey, *newLicense.NodeCount, *newLicense.CustomerID, expiry)
		if err != nil {
			log.Error(err)
			utils.HttpJSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		newLicense.Token, newLicense.Password, err = azure.CreateCarbideAccount(clientFactory, *newLicense.CustomerID, expiry)
		if err != nil {
			log.Error(err)
			utils.HttpJSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		err = utils.SendAsJSON(w, newLicense)
		if err != nil {
			log.Error(err)
		}
		return
	}
	return http.HandlerFunc(fn)
}
