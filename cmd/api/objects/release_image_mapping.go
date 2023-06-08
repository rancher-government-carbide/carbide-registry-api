package objects

import (
	"database/sql"
	"time"
)

type Release_Image_Mapping struct {
	Id        int32
	ReleaseId int32
	ImageId   int32
	CreatedAt time.Time
	UpdatedAt time.Time
}

func AddReleaseImgMapping(db *sql.DB, new_release_img_mapping Release_Image_Mapping) error {
	_, err := db.Exec("INSERT INTO release_image_mapping (id, release_id, image_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		new_release_img_mapping.Id, new_release_img_mapping.ReleaseId, new_release_img_mapping.ImageId, new_release_img_mapping.CreatedAt.Format("2006-01-02 15:04:05"), new_release_img_mapping.UpdatedAt.Format("2006-01-02 15:04:05"))
	if err != nil {
		return err
	}
	return nil
}

func GetReleaseImgMappings(db *sql.DB, image_id int32) ([]Release_Image_Mapping, error) {
	var release_img_mappings []Release_Image_Mapping
	rows, err := db.Query(`SELECT * FROM release_img_mapping WHERE image_id = ?`, image_id)
	if err != nil {
		release_img_mappings = nil
		return release_img_mappings, err
	}
	defer rows.Close()

	for rows.Next() {
		var release_img_mapping Release_Image_Mapping
		err = rows.Scan(&release_img_mapping.Id, &release_img_mapping.ReleaseId, &release_img_mapping.ImageId, &release_img_mapping.CreatedAt, &release_img_mapping.UpdatedAt)
		if err != nil {
			release_img_mappings = nil
			return release_img_mappings, err
		}
		release_img_mappings = append(release_img_mappings, release_img_mapping)
	}
	if err = rows.Err(); err != nil {
		release_img_mappings = nil
		return release_img_mappings, err
	}

	return release_img_mappings, nil
}

func GetAllReleaseImgMappings(db *sql.DB) ([]Release_Image_Mapping, error) {
	var release_img_mappings []Release_Image_Mapping
	rows, err := db.Query(`SELECT * FROM release_img_mapping`)
	if err != nil {
		release_img_mappings = nil
		return release_img_mappings, err
	}
	defer rows.Close()

	for rows.Next() {
		var release_img_mapping Release_Image_Mapping
		err = rows.Scan(&release_img_mapping.Id, &release_img_mapping.ReleaseId, &release_img_mapping.ImageId, &release_img_mapping.CreatedAt, &release_img_mapping.UpdatedAt)
		if err != nil {
			release_img_mappings = nil
			return release_img_mappings, err
		}
		release_img_mappings = append(release_img_mappings, release_img_mapping)
	}
	if err = rows.Err(); err != nil {
		release_img_mappings = nil
		return release_img_mappings, err
	}

	return release_img_mappings, nil
}

func DeleteReleaseImgMapping(db *sql.DB, name string) error {
	_, err := db.Exec(`DELETE FROM release_img_mapping WHERE release_img_mapping_name = ?`, name)
	if err != nil {
		return err
	}
	return nil
}
