package database

import (
	"carbide-registry-api/pkg/objects"
	"database/sql"
	"errors"
)

func AddProduct(db *sql.DB, newProduct objects.Product) (int64, error) {
	if err := newProduct.Validate(); err != nil {
		return -1, err
	}
	result, err := db.Exec("INSERT INTO product (name, logo_url) VALUES (?, ?)", *newProduct.Name, *newProduct.LogoUrl)
	if err != nil {
		return -1, err
	}
	return result.LastInsertId()
}

func GetProduct(db *sql.DB, id int32) (objects.Product, error) {
	var retrievedProduct objects.Product
	err := db.QueryRow(`SELECT * FROM product WHERE id = ?`, id).Scan(&retrievedProduct.Id, &retrievedProduct.Name, &retrievedProduct.LogoUrl, &retrievedProduct.CreatedAt, &retrievedProduct.UpdatedAt)
	if err != nil {
		return retrievedProduct, err
	}
	return retrievedProduct, nil
}

func GetProductByName(db *sql.DB, name string) (objects.Product, error) {
	var retrievedProduct objects.Product
	err := db.QueryRow(`SELECT * FROM product WHERE name = ?`, name).Scan(&retrievedProduct.Id, &retrievedProduct.Name, &retrievedProduct.LogoUrl, &retrievedProduct.CreatedAt, &retrievedProduct.UpdatedAt)
	if err != nil {
		return retrievedProduct, err
	}
	return retrievedProduct, nil
}

func GetAllProducts(db *sql.DB, limit int, offset int) ([]objects.Product, error) {
	var products []objects.Product
	rows, err := db.Query(`SELECT * FROM product LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		products = nil
		return products, err
	}
	defer rows.Close()
	for rows.Next() {
		var product objects.Product
		err = rows.Scan(&product.Id, &product.Name, &product.LogoUrl, &product.CreatedAt, &product.UpdatedAt)
		if err != nil {
			products = nil
			return products, err
		}
		products = append(products, product)
	}
	if err = rows.Err(); err != nil {
		products = nil
		return products, err
	}
	return products, nil
}

func UpdateProduct(db *sql.DB, newProduct objects.Product, name string) (int64, error) {
	if newProduct.Name == nil && newProduct.LogoUrl == nil {
		return -1, errors.New("invalid product")
	}
	if newProduct.Name == nil {
		result, err := db.Exec(`UPDATE product SET logo_url = ? WHERE name = ?`, *newProduct.LogoUrl, name)
		if err != nil {
			return -1, err
		}
		return result.LastInsertId()
	}
	if newProduct.LogoUrl == nil {
		result, err := db.Exec(`UPDATE product SET name = ? WHERE name = ?`, *newProduct.Name, name)
		if err != nil {
			return -1, err
		}
		return result.LastInsertId()
	}
	result, err := db.Exec(`UPDATE product SET name = ?, logo_url = ? WHERE name = ?`, *newProduct.Name, *newProduct.LogoUrl, name)
	if err != nil {
		return -1, err
	}
	return result.LastInsertId()
}

func DeleteProduct(db *sql.DB, name string) error {
	_, err := db.Exec(`DELETE FROM product WHERE name = ?`, name)
	if err != nil {
		return err
	}
	return nil
}
