package api

import (
	DB "carbide-api/pkg/database"
	"carbide-api/pkg/objects"
	"database/sql"
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// Responds with a JSON array of all products in the database
//
// Success Code: 200 OK
func productGet(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	products, err := DB.GetAllProducts(db)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	productsJSON, err := json.Marshal(products)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(productsJSON)
	if err != nil {
		log.Error(err)
	}
	return
}

// Accepts a JSON payload of a new product and responds with the new JSON object after it's been successfully created in the database
//
// Success Code: 201 OK
func productPost(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var createdProduct objects.Product
	err := json.NewDecoder(r.Body).Decode(&createdProduct)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = DB.AddProduct(db, createdProduct)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	createdProduct, err = DB.GetProduct(db, *createdProduct.Name)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	log.WithFields(log.Fields{
		"product": *createdProduct.Name,
	}).Info("Product has been successfully created")
	createdProductJSON, err := json.Marshal(createdProduct)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(createdProductJSON)
	if err != nil {
		log.Error(err)
	}
	return
}

// Responds with the JSON representation of a product
//
// Success Code: 200 OK
func productGetByName(w http.ResponseWriter, r *http.Request, db *sql.DB, productName string) {
	var retrievedProduct objects.Product
	retrievedProduct, err := DB.GetProduct(db, productName)
	if err != nil {
		log.Error(err)
		return
	}
	retrievedProductJSON, err := json.Marshal(retrievedProduct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(retrievedProductJSON)
	if err != nil {
		log.Error(err)
	}
	return
}

// Responds with the JSON representation of a product
//
// Success Code: 200 OK
func productPutByName(w http.ResponseWriter, r *http.Request, db *sql.DB, productName string) {
	var updatedProduct objects.Product
	err := json.NewDecoder(r.Body).Decode(&updatedProduct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = DB.UpdateProduct(db, *updatedProduct.Name, productName)
	if err != nil {
		log.Error(err)
		return
	}
	updatedProduct, err = DB.GetProduct(db, *updatedProduct.Name)
	if err != nil {
		log.Error(err)
		return
	}
	log.WithFields(log.Fields{
		"product": *updatedProduct.Name,
	}).Info("Product has been successfully updated")
	updatedProductJSON, err := json.Marshal(updatedProduct)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(updatedProductJSON)
	if err != nil {
		log.Error(err)
	}
	return
}

// Deletes the product and responds with an empty payload
//
// Success Code: 204 No Content
func productDeleteByName(w http.ResponseWriter, r *http.Request, db *sql.DB, productName string) {
	err := DB.DeleteProduct(db, productName)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	log.WithFields(log.Fields{
		"product": productName,
	}).Info("Product has been successfully deleted")
	w.WriteHeader(http.StatusNoContent)
	return
}
