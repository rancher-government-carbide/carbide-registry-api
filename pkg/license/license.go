package license

import (
	"crypto/rsa"
	"time"

	golicense "github.com/ebauman/golicense/pkg/license"
	"github.com/ebauman/golicense/pkg/types"
	"github.com/google/uuid"
)

func CreateCarbideLicense(privateKey *rsa.PrivateKey, nodeCount int, customerID string, expirationDate time.Time) (*string, error) {
	grants := map[string]int{
		"compliance.cattle.io/stigatron": nodeCount,
	}
	metadata := map[string]string{}
	var notBeforeTime time.Time
	license := types.License{
		Id:        uuid.NewString(),
		Grants:    grants,
		Metadata:  metadata,
		NotAfter:  expirationDate,
		NotBefore: notBeforeTime,
		Licensee:  customerID,
	}
	keystring, err := golicense.GenerateLicenseKey(privateKey, license)
	if err != nil {
		return nil, err
	}
	return &keystring, nil
}

// returns customer ID if valid, err if not
func ParseCarbideLicense(licenseString string, licensePubkeys []*rsa.PublicKey) (*string, error) {
	license, err := golicense.ParseLicenseKey([]byte(licenseString), licensePubkeys)
	if err != nil {
		return nil, err
	}
	return &license.Licensee, nil
}
