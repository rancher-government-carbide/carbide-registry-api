package objects

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Product struct {
	Id        int32
	Name      *string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func AddProduct(db *sql.DB, new_product Product) error {
	const required_field string = "Missing field \"%s\" required when creating a new product"
	const sql_error string = "Error creating new product: %w"
	if new_product.Name == nil {
		err_msg := fmt.Sprintf(required_field, "Name")
		return errors.New(err_msg)
	} else {
		_, err := db.Exec("INSERT INTO product (name) VALUES (?)", *new_product.Name)
		if err != nil {
			return fmt.Errorf(sql_error, err)
		}
	}
	return nil
}

func GetProduct(db *sql.DB, name string) (Product, error) {
	var retrieved_product Product
	err := db.QueryRow(`SELECT * FROM product WHERE name = ?`, name).Scan(&retrieved_product.Id, retrieved_product.Name, &retrieved_product.CreatedAt, &retrieved_product.UpdatedAt)
	if err != nil {
		return retrieved_product, err
	}
	return retrieved_product, nil
}

func GetAllProducts(db *sql.DB) ([]Product, error) {
	var products []Product
	rows, err := db.Query(`SELECT * FROM product`)
	if err != nil {
		products = nil
		return products, err
	}
	defer rows.Close()

	for rows.Next() {
		var product Product
		err = rows.Scan(&product.Id, &product.Name, &product.CreatedAt, &product.UpdatedAt)
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

func UpdateProduct(db *sql.DB, new_name string, name string) error {
	if _, err := db.Exec(
		`UPDATE product SET name = ?, WHERE name = ?;`, new_name, name); err != nil {
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
