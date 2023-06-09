package objects

import (
	"database/sql"
	"time"
)

type Product struct {
	Id          int32
	Name        string
	TarballLink string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func AddProduct(db *sql.DB, new_product Product) error {
	_, err := db.Exec("INSERT INTO product (id, product_name, product_signed, trivy_signed, trivy_valid, sbom_signed, sbom_valid, last_scanned_at, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		new_product.Id, new_product.Name, new_product.TarballLink, new_product.CreatedAt.Format("2006-01-02 15:04:05"), new_product.UpdatedAt.Format("2006-01-02 15:04:05"))
	if err != nil {
		return err
	}
	return nil
}

func GetProduct(db *sql.DB, name string) (Product, error) {
	var product Product
	err := db.QueryRow(`SELECT * FROM product WHERE product_name = ?`, name).Scan(&product.Id, &product.Name, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		return product, err
	}
	return product, nil
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

func UpdateProduct(db *sql.DB, updated_product Product) error {
	if _, err := db.Exec(
		`UPDATE product SET tarball_link = ?, WHERE name = ?;`,
		updated_product.TarballLink, updated_product.Name); err != nil {
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
