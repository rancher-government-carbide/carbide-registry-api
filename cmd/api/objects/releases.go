package objects

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Release struct {
	Id          int32
	ProductId   *int32
	Name        *string
	TarballLink *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func AddRelease(db *sql.DB, new_release Release) error {
	const required_field string = "Missing field \"%s\" required when creating a new release"
	const sql_error string = "Error creating new release: %w"
	if new_release.ProductId == nil {
		err_msg := fmt.Sprintf(required_field, "Product Id")
		return errors.New(err_msg)
	}
	if new_release.Name == nil {
		err_msg := fmt.Sprintf(required_field, "Name")
		return errors.New(err_msg)
	}
	if new_release.TarballLink == nil {
		_, err := db.Exec(
			"INSERT INTO releases (product_id, name) VALUES (?, ?)",
			*new_release.ProductId, *new_release.Name)
		if err != nil {
			return fmt.Errorf(sql_error, err)
		}
	} else {
		_, err := db.Exec(
			"INSERT INTO releases (product_id, name, tarball_link) VALUES (?, ?, ?)",
			*new_release.ProductId, *new_release.Name, *new_release.TarballLink)
		if err != nil {
			return fmt.Errorf(sql_error, err)
		}
	}
	return nil
}

func GetRelease(db *sql.DB, release Release) (Release, error) {
	const required_field string = "Missing field \"%s\" required when retrieving a release"
	const sql_error string = "Error finding release: %w"
	var retrieved_release Release
	if release.ProductId == nil {
		err_msg := fmt.Sprintf(required_field, "Product Id")
		return retrieved_release, errors.New(err_msg)
	}
	if release.Name == nil {
		err_msg := fmt.Sprintf(required_field, "Name")
		return retrieved_release, errors.New(err_msg)
	}
	err := db.QueryRow(
		`SELECT * FROM releases WHERE name = ? AND product_id = ?`, *release.Name, *release.ProductId).Scan(
		&retrieved_release.Id, retrieved_release.ProductId, retrieved_release.Name, retrieved_release.TarballLink, &retrieved_release.CreatedAt, &retrieved_release.UpdatedAt)
	if err != nil {
		return retrieved_release, fmt.Errorf(sql_error, err)
	}
	return retrieved_release, nil
}

func GetAllReleasesforProduct(db *sql.DB, product_name string) ([]Release, error) {

	product, err := GetProduct(db, product_name)
	product_id := product.Id

	var releases []Release
	rows, err := db.Query(`SELECT * FROM releases WHERE product_id = ?`, product_id)
	if err != nil {
		releases = nil
		return releases, err
	}
	defer rows.Close()

	for rows.Next() {
		var release Release
		err = rows.Scan(&release.Id, release.ProductId, release.Name, release.TarballLink, &release.CreatedAt, &release.UpdatedAt)
		if err != nil {
			releases = nil
			return releases, err
		}
		releases = append(releases, release)
	}
	if err = rows.Err(); err != nil {
		releases = nil
		return releases, err
	}

	return releases, nil
}

func GetAllReleases(db *sql.DB) ([]Release, error) {
	var releases []Release
	rows, err := db.Query(`SELECT * FROM releases`)
	if err != nil {
		releases = nil
		return releases, err
	}
	defer rows.Close()

	for rows.Next() {
		var release Release
		err = rows.Scan(&release.Id, &release.ProductId, &release.Name, &release.TarballLink, &release.CreatedAt, &release.UpdatedAt)
		if err != nil {
			releases = nil
			return releases, err
		}
		releases = append(releases, release)
	}
	if err = rows.Err(); err != nil {
		releases = nil
		return releases, err
	}

	return releases, nil
}

func UpdateRelease(db *sql.DB, updated_release Release) error {
	const missing_field string = "Missing field %s (needed to locate release in DB)"
	const sql_error string = "Error updating new release: %w"
	if updated_release.ProductId == nil {
		err_msg := fmt.Sprintf(missing_field, "Product Id")
		return errors.New(err_msg)
	}
	if updated_release.Name == nil {
		err_msg := fmt.Sprintf(missing_field, "Name")
		return errors.New(err_msg)
	}
	if updated_release.TarballLink == nil {
		return errors.New("No new data to update release with")
	} else {
		_, err := db.Exec(
			`UPDATE releases SET tarball_link = ? WHERE name = ? AND product_id = ?`,
			*updated_release.TarballLink, *updated_release.Name, *updated_release.ProductId)
		if err != nil {
			return fmt.Errorf(sql_error, err)
		}
	}
	return nil
}

func DeleteRelease(db *sql.DB, release_to_delete Release) error {
	const missing_field string = "Missing field %s (needed to locate release in DB)"
	const sql_error string = "Error updating new release: %w"
	if release_to_delete.ProductId == nil {
		err_msg := fmt.Sprintf(missing_field, "Product Id")
		return errors.New(err_msg)
	}
	if release_to_delete.Name == nil {
		err_msg := fmt.Sprintf(missing_field, "Name")
		return errors.New(err_msg)
	}
	_, err := db.Exec(
		`DELETE FROM releases WHERE name = ? AND product_id = ?`,
		*release_to_delete.Name, *release_to_delete.ProductId)
	if err != nil {
		return err
	}
	return nil
}
