package azure

import (
	"context"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerregistry/armcontainerregistry"
	log "github.com/sirupsen/logrus"
)

const REGISTRY_NAME = "rgcrprod"
const SCOPE_MAP = "_repositories_pull"
const RESOURCE_GROUP = "rg-rgcr-prod-usgovvirginia"

var AZURE_SUBSCRIPTION_ID = os.Getenv("AZURE_SUBSCRIPTION_ID")
var SCOPE_MAP_ID = "/subscriptions/" + AZURE_SUBSCRIPTION_ID + "/resourceGroups/" + RESOURCE_GROUP + "/providers/Microsoft.ContainerRegistry/registries/rgcrprod/scopeMaps/" + SCOPE_MAP

func NewAzureClients() (*armcontainerregistry.ClientFactory, error) {
	credentialOptions := azidentity.DefaultAzureCredentialOptions{
		ClientOptions: azcore.ClientOptions{
			Cloud: cloud.AzureGovernment,
		},
	}
	clientOptions := policy.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Cloud: cloud.AzureGovernment,
		},
	}
	cred, err := azidentity.NewDefaultAzureCredential(&credentialOptions)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
		return &armcontainerregistry.ClientFactory{}, err
	}
	clientFactory, err := armcontainerregistry.NewClientFactory(AZURE_SUBSCRIPTION_ID, cred, &clientOptions)
	if err != nil {
		log.Fatalf("failed to create clientFactory: %v", err)
		return &armcontainerregistry.ClientFactory{}, err
	}
	return clientFactory, err
}

func CreateCarbideAccount(clientFactory *armcontainerregistry.ClientFactory, customerID string, expiry time.Time) (*armcontainerregistry.Token, *armcontainerregistry.TokenPassword, error) {
	tokensclient := clientFactory.NewTokensClient()
	token, err := createNewReadToken(tokensclient, customerID)
	if err != nil {
		return &armcontainerregistry.Token{}, &armcontainerregistry.TokenPassword{}, err
	}
	registriesClient := clientFactory.NewRegistriesClient()
	password, err := createNewPassword(registriesClient, *token.ID, expiry)
	if err != nil {
		return &armcontainerregistry.Token{}, &armcontainerregistry.TokenPassword{}, err
	}
	return token, password, nil
}

func createNewReadToken(tokensClient *armcontainerregistry.TokensClient, customerID string) (*armcontainerregistry.Token, error) {
	tokenName := customerID + "-read-token"
	scopeMapID := SCOPE_MAP_ID
	tokenStatus := armcontainerregistry.TokenStatusEnabled
	ctx := context.Background()
	poller, err := tokensClient.BeginCreate(ctx, RESOURCE_GROUP, REGISTRY_NAME, tokenName, armcontainerregistry.Token{
		Properties: &armcontainerregistry.TokenProperties{
			ScopeMapID: &scopeMapID,
			Status:     &tokenStatus,
		},
	}, nil)
	if err != nil {
		log.Errorf("failed to finish the request: %v", err)
		return &armcontainerregistry.Token{}, err
	}
	res, err := poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Errorf("failed to pull the result: %v", err)
		return &armcontainerregistry.Token{}, err
	}
	return &res.Token, nil
}

// assumes customer token already exists in registry
func createNewPassword(registriesClient *armcontainerregistry.RegistriesClient, tokenID string, expiry time.Time) (*armcontainerregistry.TokenPassword, error) {
	passwordName := armcontainerregistry.TokenPasswordNamePassword1
	ctx := context.Background()
	poller, err := registriesClient.BeginGenerateCredentials(ctx, RESOURCE_GROUP, REGISTRY_NAME, armcontainerregistry.GenerateCredentialsParameters{
		Expiry:  &expiry,
		Name:    &passwordName,
		TokenID: &tokenID,
	}, nil)
	if err != nil {
		log.Errorf("failed to finish the request: %v", err)
		return &armcontainerregistry.TokenPassword{}, err
	}
	res, err := poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Errorf("failed to pull the result: %v", err)
		return &armcontainerregistry.TokenPassword{}, err
	}
	return res.Passwords[0], nil
}
