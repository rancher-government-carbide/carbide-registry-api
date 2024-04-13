package api

import (
	"carbide-images-api/pkg/api/utils"
	"carbide-images-api/pkg/objects"
	"context"
	"net/http"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"

	// "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerregistry/armcontainerregistry"
	log "github.com/sirupsen/logrus"
)

// var AZURE_TENANT_ID = os.Getenv("AZURE_TENANT_ID")
// var AZURE_CLIENT_ID = os.Getenv("AZURE_CLIENT_ID")
// var AZURE_CLIENT_SECRET = os.Getenv("AZURE_CLIENT_SECRET")
var AZURE_SUBSCRIPTION_ID = os.Getenv("AZURE_SUBSCRIPTION_ID")

const REGISTRY_NAME = "rgcrprod"
const SCOPE_MAP = "_repositories_pull"
const SCOPE_MAP_ID = "/subscriptions/b82373d0-f87a-45fd-b466-f1f97e68fcd1/resourceGroups/rg-rgcr-prod-usgovvirginia/providers/Microsoft.ContainerRegistry/registries/rgcrprod/scopeMaps/" + SCOPE_MAP
const RESOURCE_GROUP = "rg-rgcr-prod-usgovvirginia"

func NewAzureClients() (*armcontainerregistry.ClientFactory, error) {
	options := azidentity.DefaultAzureCredentialOptions{
		ClientOptions: azcore.ClientOptions{
			Cloud: cloud.AzureGovernment,
		},
	}
	cred, err := azidentity.NewDefaultAzureCredential(&options)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	clientFactory, err := armcontainerregistry.NewClientFactory(AZURE_SUBSCRIPTION_ID, cred, nil)
	if err != nil {
		log.Fatalf("failed to create clientFactory: %v", err)
	}
	return clientFactory, err
}

func createCarbideAccountHandler(clientFactory *armcontainerregistry.ClientFactory) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var newLicense objects.CarbideLicense
		err := utils.DecodeJSONObject(w, r, &newLicense)
		if err != nil {
			log.Error(err)
			return
		}
		if newLicense.CustomerID == nil || newLicense.DaysTillExpiry == nil || newLicense.NodeCount == nil {
			utils.HttpJSONError(w, "missing customerID, daysTillExpiry, or nodeCount", http.StatusBadRequest)
			log.Errorf("missing customerID, daysTillExpiry, or nodeCount")
			return
		}
		*newLicense.License = createCarbideLicense()
		*newLicense.Token, *newLicense.Password, err = createCarbideAccount(clientFactory, *newLicense.CustomerID, *newLicense.DaysTillExpiry)
		if err != nil {
			log.Error(err)
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

func createCarbideLicense() string {
	license := "xyz.xyz.xyz=="
	// TODO: use from rancherfederal/golicense as dep
	// golicense create license --grant "compliance.cattle.io/stigatron=$NODE_COUNT" --key "$GOLICENSE_KEY" --licensee "$CUSTOMER_ID" --not-after "$EXPIRATION_DATE"
	return license
}

func createCarbideAccount(clientFactory *armcontainerregistry.ClientFactory, customerID string, daysTillExpiry int) (armcontainerregistry.Token, armcontainerregistry.TokenPassword, error) {
	tokensclient := clientFactory.NewTokensClient()
	token, err := createNewReadToken(tokensclient, customerID)
	if err != nil {
		return armcontainerregistry.Token{}, armcontainerregistry.TokenPassword{}, err
	}
	registriesClient := clientFactory.NewRegistriesClient()
	password, err := createNewPassword(registriesClient, customerID, daysTillExpiry)
	if err != nil {
		return armcontainerregistry.Token{}, armcontainerregistry.TokenPassword{}, err
	}
	return token, password, nil
}

// returns token ID
func createNewReadToken(tokensClient *armcontainerregistry.TokensClient, customerID string) (armcontainerregistry.Token, error) {
	tokenName := customerID + "-read-token"
	ctx := context.Background()
	poller, err := tokensClient.BeginCreate(ctx, RESOURCE_GROUP, REGISTRY_NAME, tokenName, armcontainerregistry.Token{
		Properties: &armcontainerregistry.TokenProperties{
			ScopeMapID: to.Ptr(SCOPE_MAP_ID),
			Status:     to.Ptr(armcontainerregistry.TokenStatusEnabled),
		},
	}, nil)
	if err != nil {
		log.Errorf("failed to finish the request: %v", err)
	}
	res, err := poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Errorf("failed to pull the result: %v", err)
	}
	return res.Token, nil
}

// assumes customer token already exists in registry
func createNewPassword(registriesClient *armcontainerregistry.RegistriesClient, tokenID string, daysTillExpiry int) (armcontainerregistry.TokenPassword, error) {
	ctx := context.Background()
	expiry := time.Now().Add(time.Hour * 24 * time.Duration(daysTillExpiry))
	poller, err := registriesClient.BeginGenerateCredentials(ctx, RESOURCE_GROUP, REGISTRY_NAME, armcontainerregistry.GenerateCredentialsParameters{
		Expiry:  &expiry,
		Name:    &armcontainerregistry.PossibleTokenPasswordNameValues()[1],
		TokenID: &tokenID,
	}, nil)
	if err != nil {
		log.Errorf("failed to finish the request: %v", err)
	}
	res, err := poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Errorf("failed to pull the result: %v", err)
	}
	return *res.Passwords[0], nil
}
