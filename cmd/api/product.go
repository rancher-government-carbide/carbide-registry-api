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
		case http.MethodPost:
			productPost(w, r, db)
			return
		case http.MethodGet:
			productGet(w, r, db)
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

func productGet(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	products, err := objects.GetAllProducts(db)
	if err != nil {
		log.Print(err)
		return
	}
	products_json, err := json.Marshal(products)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(products_json)
	if err != nil {
		log.Print(err)
	}
	return

}

func productPost(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	var product objects.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = objects.AddProduct(db, product)
	if err != nil {
		log.Print(err)
		return
	}
	product, err = objects.GetProduct(db, *product.Name)
	if err != nil {
		log.Print(err)
		return
	}
	log.Printf("New product %s has been successfully created", product.Name)

	return
}

func productGet1(w http.ResponseWriter, r *http.Request, db *sql.DB, product_name string) {

	var product objects.Product
	product, err := objects.GetProduct(db, product_name)
	if err != nil {
		log.Print(err)
		return
	}
	product_json, err := json.Marshal(product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(product_json)
	if err != nil {
		log.Print(err)
	}
	return

}

func productPut1(w http.ResponseWriter, r *http.Request, db *sql.DB, product_name string) {

	var product objects.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	*product.Name = product_name
	err = objects.UpdateProduct(db, *product.Name, product_name)
	if err != nil {
		log.Print(err)
		return
	}
	product, err = objects.GetProduct(db, *product.Name)
	if err != nil {
		log.Print(err)
		return
	}
	log.Printf("Product %s has been successfully updated", product.Name)

	return
}

func productDelete1(w http.ResponseWriter, r *http.Request, db *sql.DB, product_name string) {

	err := objects.DeleteProduct(db, product_name)
	if err != nil {
		log.Print(err)
		return
	}
	log.Printf("Product %s has been successfully deleted", product_name)
	return
}
