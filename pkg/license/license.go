package license

import (
	"os"
	"time"
)

var GOLICENSE_KEY = os.Getenv("GOLICENSE_KEY")

func CreateCarbideLicense(nodeCount int, golicenseKey string, customerID string, expirationDate time.Time) string {
	license := "xyz.xyz.xyz=="
	// TODO: use from rancherfederal/golicense as dep
	// golicense create license --grant "compliance.cattle.io/stigatron=$NODE_COUNT" --key "$GOLICENSE_KEY" --licensee "$CUSTOMER_ID" --not-after "$EXPIRATION_DATE"
	return license
}
