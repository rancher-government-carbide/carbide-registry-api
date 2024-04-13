package api

import (
	"carbide-images-api/pkg/api/utils"
	"carbide-images-api/pkg/azureToken"
	"carbide-images-api/pkg/license"
	"carbide-images-api/pkg/objects"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerregistry/armcontainerregistry"
	log "github.com/sirupsen/logrus"
)

func createCarbideAccountHandler(clientFactory *armcontainerregistry.ClientFactory) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var newLicense objects.CarbideLicense
		err := utils.DecodeJSONObject(w, r, &newLicense)
		if err != nil {
			log.Error(err)
			return
		}
		if newLicense.CustomerID == nil || newLicense.DaysTillExpiry == nil || newLicense.NodeCount == nil {
			log.Errorf("missing customerID, daysTillExpiry, or nodeCount")
			utils.HttpJSONError(w, "missing customerID, daysTillExpiry, or nodeCount", http.StatusBadRequest)
			return
		}
		expiry := time.Now().Add(time.Hour * 24 * time.Duration(*newLicense.DaysTillExpiry))
		*newLicense.License = license.CreateCarbideLicense(*newLicense.NodeCount, license.GOLICENSE_KEY, *newLicense.CustomerID, expiry)
		*newLicense.Token, *newLicense.Password, err = azureToken.CreateCarbideAccount(clientFactory, *newLicense.CustomerID, expiry)
		if err != nil {
			log.Error(err)
			utils.HttpJSONError(w, err.Error(), http.StatusInternalServerError)
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
