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
	_, err := db.Exec("INSERT INTO releases (id, release_id, name, tarball_link, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
		new_release.Id, new_release.ProductId, new_release.Name, new_release.TarballLink, new_release.CreatedAt.Format("2006-01-02 15:04:05"), new_release.UpdatedAt.Format("2006-01-02 15:04:05"))
	if err != nil {
		return err
	}
	return nil
}

func GetRelease(db *sql.DB, name string) (Release, error) {
	var release Release
	err := db.QueryRow(`SELECT * FROM releases WHERE name = ?`, name).Scan(&release.Id, &release.ProductId, &release.Name, &release.TarballLink, &release.CreatedAt, &release.UpdatedAt)
	if err != nil {
		return release, err
	}
	return release, nil
}

func GetReleases(db *sql.DB) ([]Release, error) {
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

func DeleteRelease(db *sql.DB, name string) error {
	_, err := db.Exec(`DELETE FROM releases WHERE name = ?`, name)
	if err != nil {
		return err
	}
	return nil
}
