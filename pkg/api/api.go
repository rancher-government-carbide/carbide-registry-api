package api

import (
	"carbide-registry-api/pkg/api/middleware"
	"crypto/rsa"
	"database/sql"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerregistry/armcontainerregistry"
)

// api schema for public usage
//
// enables authentication;
// disables carbide account creation endpoint
func PublicRouter(db *sql.DB, licensePubkeys []*rsa.PublicKey) http.Handler {
	mux := http.NewServeMux()
	withAuth := middleware.SessionAuth
	mux.Handle("GET /healthcheck", middleware.Healthcheck())
	mux.Handle("GET /auth", authCheckHandler())
	mux.Handle("POST /login", loginHandler(licensePubkeys))
	mux.Handle("POST /logout", logoutHandler())
	mux.Handle("GET /product", withAuth(getAllProductsHandler(db)))
	mux.Handle("POST /product", withAuth(createProductHandler(db)))
	mux.Handle("GET /product/{productName}", withAuth(getProductHandler(db)))
	mux.Handle("PUT /product/{productName}", withAuth(updateProductHandler(db)))
	mux.Handle("DELETE /product/{productName}", withAuth(deleteProductHandler(db)))
	mux.Handle("GET /product/{productName}/release", withAuth(getAllReleasesHandler(db)))
	mux.Handle("POST /product/{productName}/release", withAuth(createReleaseHandler(db)))
	mux.Handle("GET /product/{productName}/release/{releaseName}", withAuth(getReleaseHandler(db)))
	mux.Handle("PUT /product/{productName}/release/{releaseName}", withAuth(updateReleaseHandler(db)))
	mux.Handle("DELETE /product/{proudctName}/release/{releaseName}", withAuth(deleteReleaseHandler(db)))
	mux.Handle("GET /image", withAuth(getAllImagesHandler(db)))
	mux.Handle("POST /image", withAuth(createImageHandler(db)))
	mux.Handle("GET /image/{imageID}", withAuth(getImageHandler(db)))
	mux.Handle("PUT /image/{imageID}", withAuth(updateImageHandler(db)))
	mux.Handle("DELETE /image/{imageID}", withAuth(deleteImageHandler(db)))
	mux.Handle("GET /releaseImageMapping", withAuth(getAllReleaseImageMappingsHandler(db)))
	mux.Handle("POST /releaseImageMapping", withAuth(createReleaseImageMappingHandler(db)))
	mux.Handle("DELETE /releaseImageMapping", withAuth(deleteReleaseImageMappingHandler(db)))
	withCors := middleware.CORS
	newMux := withCors(mux)
	return newMux
}

// api schema for internal ssf usage
//
// disables authentication;
// enables carbide account creation endpoint
func PrivateRouter(db *sql.DB, clientFactory *armcontainerregistry.ClientFactory, licensePrivkey *rsa.PrivateKey) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET /healthcheck", middleware.Healthcheck())
	mux.Handle("POST /carbide/license", createCarbideAccountHandler(clientFactory, licensePrivkey))
	mux.Handle("GET /product", getAllProductsHandler(db))
	mux.Handle("POST /product", createProductHandler(db))
	mux.Handle("GET /product/{productName}", getProductHandler(db))
	mux.Handle("PUT /product/{productName}", updateProductHandler(db))
	mux.Handle("DELETE /product/{productName}", deleteProductHandler(db))
	mux.Handle("GET /product/{productName}/release", getAllReleasesHandler(db))
	mux.Handle("POST /product/{productName}/release", createReleaseHandler(db))
	mux.Handle("GET /product/{productName}/release/{releaseName}", getReleaseHandler(db))
	mux.Handle("PUT /product/{productName}/release/{releaseName}", updateReleaseHandler(db))
	mux.Handle("DELETE /product/{proudctName}/release/{releaseName}", deleteReleaseHandler(db))
	mux.Handle("GET /image", getAllImagesHandler(db))
	mux.Handle("POST /image", createImageHandler(db))
	mux.Handle("GET /image/{imageID}", getImageHandler(db))
	mux.Handle("PUT /image/{imageID}", updateImageHandler(db))
	mux.Handle("DELETE /image/{imageID}", deleteImageHandler(db))
	mux.Handle("GET /releaseImageMapping", getAllReleaseImageMappingsHandler(db))
	mux.Handle("POST /releaseImageMapping", createReleaseImageMappingHandler(db))
	mux.Handle("DELETE /releaseImageMapping", deleteReleaseImageMappingHandler(db))
	withCors := middleware.CORS
	newMux := withCors(mux)
	return newMux
}
