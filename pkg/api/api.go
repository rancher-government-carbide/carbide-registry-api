package api

import (
	"carbide-images-api/pkg/api/middleware"
	"database/sql"
	"net/http"
)

func NewRouter(db *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()
	global := middleware.CORS
	withAuth := middleware.ChainHandlers(global, middleware.JWTAuth)
	mux.Handle("POST /user", global(createUserHandler(db)))
	mux.Handle("DELETE /user", global(deleteUserHandler(db)))
	mux.Handle("POST /login", global(loginHandler(db)))
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
	return mux
}
