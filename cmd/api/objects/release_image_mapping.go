package objects

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Release_Image_Mapping struct {
	Id        int32
	ReleaseId *int32
	ImageId   *int32
	CreatedAt time.Time
	UpdatedAt time.Time
}

func GetReleaseImageMappingbyId(db *sql.DB, release_image_mapping_id int32) (Release_Image_Mapping, error) {
	var release_image_mapping Release_Image_Mapping
	err := db.QueryRow(`SELECT * FROM release_image_mapping WHERE id = ?`, release_image_mapping_id).Scan(
		&release_image_mapping.Id,
		&release_image_mapping.ReleaseId,
		&release_image_mapping.ImageId,
		&release_image_mapping.CreatedAt,
		&release_image_mapping.UpdatedAt,
	)
	if err != nil {
		return release_image_mapping, err
	}
	return release_image_mapping, nil
}

func GetReleaseImageMapping(db *sql.DB, release_image_mapping Release_Image_Mapping) (Release_Image_Mapping, error) {
	const required_field string = "Missing field \"%s\" required when retrieving a release_image_mapping"
	const sql_error string = "Error retrieving release_image_mapping: %w"
	var retrieved_release_image_mapping Release_Image_Mapping
	if release_image_mapping.ReleaseId == nil {
		err_msg := fmt.Sprintf(required_field, "ReleaseId")
		return retrieved_release_image_mapping, errors.New(err_msg)
	}
	if release_image_mapping.ImageId == nil {
		err_msg := fmt.Sprintf(required_field, "ImageId")
		return retrieved_release_image_mapping, errors.New(err_msg)
	}
	err := db.QueryRow(`SELECT * FROM release_image_mapping WHERE release_id = ? AND image_id = ?`, *release_image_mapping.ReleaseId, *release_image_mapping.ImageId).Scan(
		&retrieved_release_image_mapping.Id,
		&retrieved_release_image_mapping.ReleaseId,
		&retrieved_release_image_mapping.ImageId,
		&retrieved_release_image_mapping.CreatedAt,
		&retrieved_release_image_mapping.UpdatedAt,
	)
	if err != nil {
		return retrieved_release_image_mapping, err
	}
	return retrieved_release_image_mapping, nil
}

func AddReleaseImgMapping(db *sql.DB, new_release_img_mapping Release_Image_Mapping) error {
	const required_field string = "Missing field \"%s\" required when creating a release_image_mapping"
	const sql_error string = "Error creating release_image_mapping: %w"
	if new_release_img_mapping.ReleaseId == nil {
		err_msg := fmt.Sprintf(required_field, "ReleaseId")
		return errors.New(err_msg)
	}
	if new_release_img_mapping.ImageId == nil {
		err_msg := fmt.Sprintf(required_field, "ImageId")
		return errors.New(err_msg)
	}
	_, err := db.Exec("INSERT INTO release_image_mapping (release_id, image_id) VALUES (?, ?)",
		*new_release_img_mapping.ReleaseId, *new_release_img_mapping.ImageId)
	if err != nil {
		return err
	}
	return nil
}

func GetImgMappings(db *sql.DB, release_id int32) ([]Release_Image_Mapping, error) {
	var release_img_mappings []Release_Image_Mapping
	rows, err := db.Query(`SELECT * FROM release_image_mapping WHERE release_id = ?`, release_id)
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

func GetReleaseMappings(db *sql.DB, image_id int32) ([]Release_Image_Mapping, error) {
	var release_img_mappings []Release_Image_Mapping
	rows, err := db.Query(`SELECT * FROM release_image_mapping WHERE image_id = ?`, image_id)
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
	rows, err := db.Query(`SELECT * FROM release_image_mapping`)
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

func DeleteReleaseImgMappingbyId(db *sql.DB, release_img_mapping_id int32) error {
	_, err := db.Exec(`DELETE FROM release_image_mapping WHERE id = ?`, release_img_mapping_id)
	if err != nil {
		return err
	}
	return nil
}

func DeleteReleaseImgMapping(db *sql.DB, release_image_mapping_to_delete Release_Image_Mapping) error {
	const required_field string = "Missing field \"%s\" required when deleting a release_image_mapping"
	const sql_error string = "Error deleting release_image_mapping: %w"
	if release_image_mapping_to_delete.ReleaseId == nil {
		err_msg := fmt.Sprintf(required_field, "ReleaseId")
		return errors.New(err_msg)
	}
	if release_image_mapping_to_delete.ImageId == nil {
		err_msg := fmt.Sprintf(required_field, "ImageId")
		return errors.New(err_msg)
	}
	_, err := db.Exec(`DELETE FROM release_image_mapping WHERE release_id = ? AND image_id = ?`, release_image_mapping_to_delete.ReleaseId, release_image_mapping_to_delete.ImageId)
	if err != nil {
		return err
	}
	return nil
}
