package api

import (
	"carbide-api/cmd/api/objects"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func serveProduct(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var product_name string
	product_name, r.URL.Path = ShiftPath(r.URL.Path)
	if r.URL.Path != "/" {
		var head string
		head, r.URL.Path = ShiftPath(r.URL.Path)
		switch head {
		case "release":
			serveRelease(w, r, db, product_name)
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
		}
		return
	}
	if product_name == "" {
		switch r.Method {
		case http.MethodGet:
			productGet(w, r, db)
			return
		case http.MethodPost:
			productPost(w, r, db)
			return
		case http.MethodOptions:
			return
		default:
			http.Error(w, fmt.Sprintf("Expected method POST or OPTIONS, got %v", r.Method), http.StatusMethodNotAllowed)
			return
		}
	} else {
		switch r.Method {
		case http.MethodGet:
			productGet1(w, r, db, product_name)
			return
		case http.MethodPut:
			productPut1(w, r, db, product_name)
			return
		case http.MethodDelete:
			productDelete1(w, r, db, product_name)
			return
		case http.MethodOptions:
			return
		default:
			http.Error(w, fmt.Sprintf("Expected method POST or OPTIONS, got %v", r.Method), http.StatusMethodNotAllowed)
			return
		}
	}
}

// Responds with a JSON array of all products in the database
//
// Success Code: 200 OK
func productGet(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	products, err := objects.GetAllProducts(db)
	if err != nil {
		http_json_error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}
	products_json, err := json.Marshal(products)
	if err != nil {
		http_json_error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(products_json)
	if err != nil {
		log.Print(err)
	}
	return
}

// Accepts a JSON payload of a new product and responds with the new JSON object after it's been successfully created in the database
//
// Success Code: 201 OK
func productPost(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var created_product objects.Product
	err := json.NewDecoder(r.Body).Decode(&created_product)
	if err != nil {
		http_json_error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = objects.AddProduct(db, created_product)
	if err != nil {
		http_json_error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}
	created_product, err = objects.GetProduct(db, *created_product.Name)
	if err != nil {
		http_json_error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}
	log.Printf("New product %s has been successfully created", *created_product.Name)
	created_product_json, err := json.Marshal(created_product)
	if err != nil {
		http_json_error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(created_product_json)
	if err != nil {
		log.Print(err)
	}
	return
}

// Responds with the JSON representation of a product
//
// Success Code: 200 OK
func productGet1(w http.ResponseWriter, r *http.Request, db *sql.DB, product_name string) {
	var retrieved_product objects.Product
	retrieved_product, err := objects.GetProduct(db, product_name)
	if err != nil {
		log.Print(err)
		return
	}
	retrieved_product_json, err := json.Marshal(retrieved_product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(retrieved_product_json)
	if err != nil {
		log.Print(err)
	}
	return
}

// Responds with the JSON representation of a product
//
// Success Code: 200 OK
func productPut1(w http.ResponseWriter, r *http.Request, db *sql.DB, product_name string) {
	var updated_product objects.Product
	err := json.NewDecoder(r.Body).Decode(&updated_product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = objects.UpdateProduct(db, *updated_product.Name, product_name)
	if err != nil {
		log.Print(err)
		return
	}
	updated_product, err = objects.GetProduct(db, *updated_product.Name)
	if err != nil {
		log.Print(err)
		return
	}
	log.Printf("Product %s has been successfully updated", *updated_product.Name)
	updated_product_json, err := json.Marshal(updated_product)
	if err != nil {
		http_json_error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(updated_product_json)
	if err != nil {
		log.Print(err)
	}
	return
}

// Deletes the product and responds with an empty payload
//
// Success Code: 204 No Content
func productDelete1(w http.ResponseWriter, r *http.Request, db *sql.DB, product_name string) {
	err := objects.DeleteProduct(db, product_name)
	if err != nil {
		http_json_error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}
	log.Printf("Product %s has been successfully deleted", product_name)
	w.WriteHeader(http.StatusNoContent)
	return
}
