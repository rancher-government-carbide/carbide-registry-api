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

	product, err = objects.GetProduct(db, product.Name)
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
	product.Name = product_name
	err = objects.UpdateProduct(db, product)
	if err != nil {
		log.Print(err)
		return
	}
	product, err = objects.GetProduct(db, product.Name)
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

// {
// 	if listname == "" {
// 		switch r.Method {
// 		case http.MethodGet:
//
// 			// get array of list objects
// 			lists, err := getLists(h.db, userid)
// 			if err != nil {
// 				log.Print("Error getting lists from database")
// 				log.Print(err)
// 				return
// 			}
// 			// send array of lists as http response
// 			if err := renderJSON(w, lists); err != nil {
// 				log.Print(err)
// 				return
// 			}
//
// 			return
// 		case http.MethodPost:
//
// 			var newlist List
//
// 			// parse payload into newlist object
// 			if err := parseJSON(w, r, &newlist); err != nil {
// 				log.Print(err)
// 				return
// 			}
//
// 			// set userid to jwt result to prevent impersonation
// 			newlist.Userid = userid
//
// 			// add list to database
// 			if err := addList(h.db, newlist); err != nil {
// 				log.Print(err)
// 				return
// 			}
//
// 			log.Printf("Added \"%s\" to lists!", newlist.Name)
// 			respondSuccess(w)
//
// 			return
// 		case http.MethodOptions:
// 			return
// 		default:
// 			http.Error(w, fmt.Sprintf("Expected method GET, POST, or OPTIONS got %v", r.Method), http.StatusMethodNotAllowed)
// 			return
// 		}
// 	}
// 	switch r.Method {
// 	case http.MethodGet:
//
// 		list, err := getList(h.db, listname, userid)
// 		if err != nil {
// 			log.Print("Error getting list from database")
// 			log.Print(err)
// 			return
// 		}
// 		if err := renderJSON(w, list); err != nil {
// 			log.Print(err)
// 			return
// 		}
//
// 		return
// 	case http.MethodPut:
//
// 		var updatedlist List
// 		if err := parseJSON(w, r, &updatedlist); err != nil {
// 			log.Print("Error parsing json payload\n")
// 			log.Print(err)
// 			return
// 		}
// 		if err := updateList(h.db, updatedlist, listname, userid); err != nil {
// 			log.Printf("Error updating list %s\n", listname)
// 			log.Print(err)
// 			return
// 		}
// 		log.Printf("Updated list \"%s\"\n", updatedlist.Name)
// 		respondSuccess(w)
//
// 		return
// 	case http.MethodDelete:
// 		var list List
// 		list.Userid = userid
// 		list.Name = listname
// 		if err := deleteList(h.db, list); err != nil {
// 			log.Print(err)
// 			log.Print("Failed to delete list\n")
// 		}
// 		log.Printf("Deleted list %s\n", list.Name)
// 		respondSuccess(w)
//
// 		return
// 	case http.MethodOptions:
// 		return
// 	default:
// 		http.Error(w, fmt.Sprintf("Expected method GET, POST, PUT, OPTIONS, or DELETE, got %v", r.Method), http.StatusMethodNotAllowed)
// 		return
// 	}
// }
