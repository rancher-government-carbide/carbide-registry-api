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
	mux.Handle("POST /login", authMiddleware.Then(loginHandler(db)))
	mux.Handle("GET /product", authMiddleware.Then(getAllProductsHandler(db)))
	mux.Handle("POST /product", authMiddleware.Then(createProductHandler(db)))
	mux.Handle("GET /product/{name}", authMiddleware.Then(getProductHandler(db)))
	mux.Handle("PUT /product/{name}", authMiddleware.Then(updateProductHandler(db)))
	mux.Handle("DELETE /product/{name}", authMiddleware.Then(deleteProductHandler(db)))
	mux.Handle("GET /release", authMiddleware.Then(getAllReleasesHandler(db)))
	mux.Handle("POST /release", authMiddleware.Then(createReleaseHandler(db)))
	mux.Handle("GET /release/{name}", authMiddleware.Then(getReleaseHandler(db)))
	mux.Handle("DELETE /release/{name}", authMiddleware.Then(deleteReleaseHandler(db)))
	mux.Handle("GET /image", authMiddleware.Then(getAllImagesHandler(db)))
	mux.Handle("POST /image", authMiddleware.Then(createImageHandler(db)))
	mux.Handle("GET /image/{id}", authMiddleware.Then(getImageHandler(db)))
	mux.Handle("PUT /image/{id}", authMiddleware.Then(updateImageHandler(db)))
	mux.Handle("DELETE /image/{id}", authMiddleware.Then(deleteImageHandler(db)))
	mux.Handle("GET /releaseImageMapping", authMiddleware.Then(getAllReleaseImageMappingsHandler(db)))
	mux.Handle("POST /releaseImageMapping", authMiddleware.Then(createReleaseImageMappingHandler(db)))
	mux.Handle("DELETE /releaseImageMapping", authMiddleware.Then(deleteReleaseImageMappingHandler(db)))
	return mux
}
