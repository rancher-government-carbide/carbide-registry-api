package database

import (
	"carbide-api/pkg/objects"
	"database/sql"
	"errors"
	"fmt"
)

func GetImagebyId(db *sql.DB, image_id int32) (objects.Image, error) {
	var image objects.Image
	err := db.QueryRow(`SELECT * FROM images WHERE id = ?`, image_id).Scan(&image.Id, &image.ImageName, &image.ImageSigned, &image.TrivySigned, &image.TrivyValid, &image.SbomSigned, &image.SbomValid, &image.LastScannedAt, &image.CreatedAt, &image.UpdatedAt)
	if err != nil {
		return image, err
	}
	image.Releases, err = GetAllReleasesforImage(db, image_id)
	if err != nil {
		return image, err
	}
	return image, nil
}

func AddImage(db *sql.DB, new_image objects.Image) error {
	const required_field string = "Missing field \"%s\" required when creating a new image"
	const sql_error string = "Error creating new image: %w"

	var (
		imageName     sql.NullString
		imageSigned   sql.NullBool
		trivySigned   sql.NullBool
		trivyValid    sql.NullBool
		sbomSigned    sql.NullBool
		sbomValid     sql.NullBool
		lastScannedAt sql.NullTime
	)

	stmt, err := db.Prepare(`
		INSERT INTO images (image_name, image_signed, trivy_signed, trivy_valid, sbom_signed, sbom_valid, last_scanned_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if new_image.ImageName == nil {
		err_msg := fmt.Sprintf(required_field, "ImageName")
		return errors.New(err_msg)
	} else {
		imageName.String = *new_image.ImageName
		imageName.Valid = true
	}
	if new_image.ImageSigned != nil {
		imageSigned.Bool = *new_image.ImageSigned
		imageSigned.Valid = true
	}
	if new_image.TrivySigned != nil {
		trivySigned.Bool = *new_image.TrivySigned
		trivySigned.Valid = true
	}
	if new_image.TrivyValid != nil {
		trivyValid.Bool = *new_image.TrivyValid
		trivyValid.Valid = true
	}
	if new_image.SbomSigned != nil {
		sbomSigned.Bool = *new_image.SbomSigned
		sbomSigned.Valid = true
	}
	if new_image.SbomValid != nil {
		sbomValid.Bool = *new_image.SbomValid
		sbomValid.Valid = true
	}
	if new_image.LastScannedAt != nil {
		lastScannedAt.Time = *new_image.LastScannedAt
		lastScannedAt.Valid = true
	}
	_, err = stmt.Exec(
		imageName,
		imageSigned,
		trivySigned,
		trivyValid,
		sbomSigned,
		sbomValid,
		lastScannedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func GetImagebyName(db *sql.DB, image_name string) (objects.Image, error) {
	var image objects.Image
	err := db.QueryRow(`SELECT * FROM images WHERE image_name = ?`, image_name).Scan(&image.Id, &image.ImageName, &image.ImageSigned, &image.TrivySigned, &image.TrivyValid, &image.SbomSigned, &image.SbomValid, &image.LastScannedAt, &image.CreatedAt, &image.UpdatedAt)
	if err != nil {
		return image, err
	}
	image.Releases, err = GetAllReleasesforImage(db, image.Id)
	if err != nil {
		return image, err
	}
	return image, nil
}

func GetAllImages(db *sql.DB) ([]objects.Image, error) {
	var images []objects.Image
	rows, err := db.Query(`SELECT * FROM images`)
	if err != nil {
		images = nil
		return images, err
	}
	defer rows.Close()

	for rows.Next() {
		var image objects.Image
		err = rows.Scan(&image.Id, &image.ImageName, &image.ImageSigned, &image.TrivySigned, &image.TrivyValid, &image.SbomSigned, &image.SbomValid, &image.LastScannedAt, &image.CreatedAt, &image.UpdatedAt)
		if err != nil {
			images = nil
			return images, err
		}
		images = append(images, image)
	}
	if err = rows.Err(); err != nil {
		images = nil
		return images, err
	}

	return images, nil
}

// func GetAllImagesforProduct(db *sql.DB, product_name string, release_name string) ([]Product, error) {
// 	var products []Product
// 	rows, err := db.Query(`SELECT * FROM images WHERE product_name`)
// 	if err != nil {
// 		products = nil
// 		return products, err
// 	}
// 	defer rows.Close()
//
// 	for rows.Next() {
// 		var product Product
// 		err = rows.Scan(&product.Id, &product.Name, &product.CreatedAt, &product.UpdatedAt)
// 		if err != nil {
// 			products = nil
// 			return products, err
// 		}
// 		products = append(products, product)
// 	}
// 	if err = rows.Err(); err != nil {
// 		products = nil
// 		return products, err
// 	}
//
// 	return products, nil
// }

func UpdateImage(db *sql.DB, updated_image objects.Image) error {

	const required_field string = "Missing field \"%s\" required when updating an image"
	const sql_error string = "Error updating image: %w"

	var (
		imageid       sql.NullInt32
		imageSigned   sql.NullBool
		trivySigned   sql.NullBool
		trivyValid    sql.NullBool
		sbomSigned    sql.NullBool
		sbomValid     sql.NullBool
		lastScannedAt sql.NullTime
	)

	stmt, err := db.Prepare(`
		UPDATE images 
		SET image_signed = ?, 
		trivy_signed = ?, 
		trivy_valid = ?, 
		sbom_signed = ?, 
		sbom_valid = ?, 
		last_scanned_at = ?
		WHERE id = ?
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if updated_image.Id == 0 {
		err_msg := fmt.Sprintf(required_field, "Id")
		return errors.New(err_msg)
	} else {
		imageid.Int32 = updated_image.Id
		imageid.Valid = true
	}
	if updated_image.ImageSigned != nil {
		imageSigned.Bool = *updated_image.ImageSigned
		imageSigned.Valid = true
	}
	if updated_image.TrivySigned != nil {
		trivySigned.Bool = *updated_image.TrivySigned
		trivySigned.Valid = true
	}
	if updated_image.TrivyValid != nil {
		trivyValid.Bool = *updated_image.TrivyValid
		trivyValid.Valid = true
	}
	if updated_image.SbomSigned != nil {
		sbomSigned.Bool = *updated_image.SbomSigned
		sbomSigned.Valid = true
	}
	if updated_image.SbomValid != nil {
		sbomValid.Bool = *updated_image.SbomValid
		sbomValid.Valid = true
	}
	if updated_image.LastScannedAt != nil {
		lastScannedAt.Time = *updated_image.LastScannedAt
		lastScannedAt.Valid = true
	}
	_, err = stmt.Exec(
		imageSigned,
		trivySigned,
		trivyValid,
		sbomSigned,
		sbomValid,
		lastScannedAt,
		imageid,
	)
	if err != nil {
		return err
	}

	return nil
}

func DeleteImage(db *sql.DB, id int32) error {
	_, err := db.Exec(`DELETE FROM images WHERE id = ?`, id)
	if err != nil {
		return err
	}
	return nil
}

func GetImageWithoutReleases(db *sql.DB, image_id int32) (objects.Image, error) {
	var retrieved_image objects.Image
	const sql_error string = "Error fetching image: %w"
	err := db.QueryRow(`SELECT * FROM images WHERE id = ?`, image_id).Scan(&retrieved_image.Id, &retrieved_image.ImageName, &retrieved_image.ImageSigned, &retrieved_image.TrivySigned, &retrieved_image.TrivyValid, &retrieved_image.SbomSigned, &retrieved_image.SbomValid, &retrieved_image.LastScannedAt, &retrieved_image.CreatedAt, &retrieved_image.UpdatedAt)
	if err != nil {
		return retrieved_image, fmt.Errorf(sql_error, err)
	}
	return retrieved_image, nil
}

func GetAllImagesforRelease(db *sql.DB, release_id int32) ([]objects.Image, error) {
	var fetched_images []objects.Image

	var release_img_mappings []objects.Release_Image_Mapping
	release_img_mappings, err := GetImgMappings(db, release_id)
	if err != nil {
		return fetched_images, err
	}

	for _, release_image_mapping := range release_img_mappings {
		image, err := GetImageWithoutReleases(db, *release_image_mapping.ImageId)
		if err != nil {
			return fetched_images, err
		}
		fetched_images = append(fetched_images, image)
	}

	return fetched_images, nil
}

// 	if updated_image.ImageSigned != nil {
// 		stmt, err := db.Prepare("CALL update_image_signed(?, ?)")
// 		if err != nil {
// 			return err
// 		}
// 		defer stmt.Close()
// 		_, err = stmt.Exec(updated_image.ImageName, updated_image.ImageSigned)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	if updated_image.TrivySigned != nil && updated_image.TrivyValid != nil {
// 		stmt, err := db.Prepare("CALL update_trivy_flags(?, ?, ?)")
// 		if err != nil {
// 			return err
// 		}
// 		defer stmt.Close()
// 		_, err = stmt.Exec(updated_image.ImageName, updated_image.TrivySigned, updated_image.TrivyValid)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	if updated_image.SbomSigned != nil && updated_image.SbomValid != nil {
// 		stmt, err := db.Prepare("CALL update_sbom_flags(?, ?, ?)")
// 		if err != nil {
// 			return err
// 		}
// 		defer stmt.Close()
// 		_, err = stmt.Exec(updated_image.ImageName, updated_image.SbomSigned, updated_image.SbomValid)
// 		if err != nil {
// 			return err
// 		}
// 	}
