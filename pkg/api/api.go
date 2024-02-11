package api

import (
	"carbide-images-api/pkg/api/middleware"
	"database/sql"
	"net/http"

	"github.com/justinas/alice"
)

func NewRouter(db *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()
	globalMiddleware := alice.New(middleware.CORS)
	authMiddleware := alice.New(middleware.CORS, middleware.JWTAuth)
	mux.Handle("POST /user", globalMiddleware.Then(createUserHandler(db)))
	mux.Handle("DELETE /user", globalMiddleware.Then(deleteUserHandler(db)))
	mux.Handle("POST /login", globalMiddleware.Then(loginHandler(db)))
	mux.Handle("GET /product", authMiddleware.Then(getAllProductsHandler(db)))
	mux.Handle("POST /product", authMiddleware.Then(createProductHandler(db)))
	mux.Handle("GET /product/{productName}", authMiddleware.Then(getProductHandler(db)))
	mux.Handle("PUT /product/{productName}", authMiddleware.Then(updateProductHandler(db)))
	mux.Handle("DELETE /product/{productName}", authMiddleware.Then(deleteProductHandler(db)))
	mux.Handle("GET /product/{productName}/release", authMiddleware.Then(getAllReleasesHandler(db)))
	mux.Handle("POST /product/{productName}/release", authMiddleware.Then(createReleaseHandler(db)))
	mux.Handle("GET /product/{productName}/release/{releaseName}", authMiddleware.Then(getReleaseHandler(db)))
	mux.Handle("DELETE /product/{proudctName}/release/{releaseName}", authMiddleware.Then(deleteReleaseHandler(db)))
	mux.Handle("GET /image", authMiddleware.Then(getAllImagesHandler(db)))
	mux.Handle("POST /image", authMiddleware.Then(createImageHandler(db)))
	mux.Handle("GET /image/{imageID}", authMiddleware.Then(getImageHandler(db)))
	mux.Handle("PUT /image/{imageID}", authMiddleware.Then(updateImageHandler(db)))
	mux.Handle("DELETE /image/{imageID}", authMiddleware.Then(deleteImageHandler(db)))
	mux.Handle("GET /releaseImageMapping", authMiddleware.Then(getAllReleaseImageMappingsHandler(db)))
	mux.Handle("POST /releaseImageMapping", authMiddleware.Then(createReleaseImageMappingHandler(db)))
	mux.Handle("DELETE /releaseImageMapping", authMiddleware.Then(deleteReleaseImageMappingHandler(db)))
	return mux
}
