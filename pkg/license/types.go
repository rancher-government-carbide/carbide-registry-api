package license

import (
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerregistry/armcontainerregistry"
)

type CarbideLicense struct {
	CustomerID     *string
	DaysTillExpiry *int
	NodeCount      *int

	License  *string
	Token    *armcontainerregistry.Token
	Password *armcontainerregistry.TokenPassword
}
