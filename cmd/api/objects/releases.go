package objects

import (
	"database/sql"
	"time"
)

type Release struct {
	Id          int32
	ProductId   int32
	Name        string
	TarballLink string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func AddRelease(db *sql.DB, new_release Release) error {
	_, err := db.Exec("INSERT INTO releases (product_id, name, tarball_link) VALUES (?, ?, ?)",
		new_release.ProductId, new_release.Name, new_release.TarballLink)
	if err != nil {
		return err
	}
	return nil
}

// TODO: should make sure product id matches too
func GetRelease(db *sql.DB, name string) (Release, error) {
	var release Release
	err := db.QueryRow(`SELECT * FROM releases WHERE name = ?`, name).Scan(&release.Id, &release.ProductId, &release.Name, &release.TarballLink, &release.CreatedAt, &release.UpdatedAt)
	if err != nil {
		return release, err
	}
	return release, nil
}

func GetAllReleasesforProduct(db *sql.DB, product_name string) ([]Release, error) {

	product, err := GetProduct(db, product_name)
	product_id := product.Name

	var releases []Release
	rows, err := db.Query(`SELECT * FROM releases WHERE product_id = ?`, product_id)
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

// TODO: should make sure product id matches too
func UpdateRelease(db *sql.DB, updated_release Release) error {
	if _, err := db.Exec(
		`UPDATE releases SET tarball_link = ?, WHERE name = ?;`,
		updated_release.TarballLink, updated_release.Name); err != nil {
		return err
	}
	return nil
}

// TODO: should make sure product id matches too
func DeleteRelease(db *sql.DB, name string) error {
	_, err := db.Exec(`DELETE FROM releases WHERE name = ?`, name)
	if err != nil {
		return err
	}
	return nil
}
