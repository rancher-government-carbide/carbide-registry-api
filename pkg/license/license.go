package license

import (
	"crypto/rsa"
	"os"
	"time"

	"github.com/ebauman/golicense/pkg/certificate"
	golicense "github.com/ebauman/golicense/pkg/license"
	"github.com/ebauman/golicense/pkg/types"
	"github.com/google/uuid"
)

var PRIVATEKEY *rsa.PrivateKey

func init() {
	var GOLICENSE_KEY = os.Getenv("GOLICENSE_KEY")
	var err error
	PRIVATEKEY, err = certificate.PEMToPrivateKey([]byte(GOLICENSE_KEY))
	if err != nil {
		panic(err)
	}
}

func CreateCarbideLicense(nodeCount int, customerID string, expirationDate time.Time) (*string, error) {
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
	keystring, err := golicense.GenerateLicenseKey(PRIVATEKEY, license)
	if err != nil {
		return nil, err
	}
	return &keystring, nil
}
