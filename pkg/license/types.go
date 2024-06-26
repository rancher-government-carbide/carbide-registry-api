package license

import (
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerregistry/armcontainerregistry"
)

type CarbideLicense struct {
	CustomerID     *string
	DaysTillExpiry *int
	NodeCount      *int

	LicenseString *string                             `json:"license,omitempty"`
	Token         *armcontainerregistry.Token         `json:"token,omitempty"`
	Password      *armcontainerregistry.TokenPassword `json:"password,omitempty"`
}
