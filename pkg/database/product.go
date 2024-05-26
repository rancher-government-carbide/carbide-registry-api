package database

import (
	"carbide-registry-api/pkg/objects"
	"database/sql"
)

func AddProduct(db *sql.DB, newProduct objects.Product) error {
	if err := newProduct.Validate(); err != nil {
		return err
	}
	_, err := db.Exec("INSERT INTO product (name, logo_url) VALUES (?, ?)", *newProduct.Name, *newProduct.LogoUrl)
	if err != nil {
		return err
	}
	return nil
}

func GetProduct(db *sql.DB, name string) (objects.Product, error) {
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

func UpdateProduct(db *sql.DB, newName string, name string) error {
	if _, err := db.Exec(
		`UPDATE product SET name = ? WHERE name = ?`, newName, name); err != nil {
		return err
	}
	return nil
}

func DeleteProduct(db *sql.DB, name string) error {
	_, err := db.Exec(`DELETE FROM product WHERE name = ?`, name)
	if err != nil {
		return err
	}
	return nil
}
