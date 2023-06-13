package objects

import (
	"database/sql"
	"time"
)

type Image struct {
	Id            int32
	ImageName     string
	ImageSigned   bool
	TrivySigned   bool
	TrivyValid    bool
	SbomSigned    bool
	SbomValid     bool
	LastScannedAt time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time

	ReleaseMappings []Release_Image_Mapping
}

func AddImage(db *sql.DB, new_image Image) error {
	_, err := db.Exec("INSERT INTO images (id, image_name, image_signed, trivy_signed, trivy_valid, sbom_signed, sbom_valid, last_scanned_at, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		new_image.Id, new_image.ImageName, new_image.ImageSigned, new_image.TrivySigned, new_image.TrivyValid, new_image.SbomSigned, new_image.SbomValid, new_image.LastScannedAt.Format("2006-01-02 15:04:05"))
	if err != nil {
		return err
	}
	return nil
}

func GetImage(db *sql.DB, name string) (Image, error) {
	var image Image
	err := db.QueryRow(`SELECT * FROM images WHERE image_name = ?`, name).Scan(&image.Id, &image.ImageName, &image.ImageSigned, &image.TrivySigned, &image.TrivyValid, &image.SbomSigned, &image.SbomValid, &image.LastScannedAt, &image.CreatedAt, &image.UpdatedAt)
	if err != nil {
		return image, err
	}
	image.ReleaseMappings, err = GetReleaseImgMappings(db, image.Id)
	if err != nil {
		return image, err
	}
	return image, nil
}

func GetAllImages(db *sql.DB) ([]Image, error) {
	var images []Image
	rows, err := db.Query(`SELECT * FROM images`)
	if err != nil {
		images = nil
		return images, err
	}
	defer rows.Close()

	for rows.Next() {
		var image Image
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

func GetAllImagesforProduct(db *sql.DB, product_name string, release_name string) ([]Product, error) {
	var products []Product
	rows, err := db.Query(`SELECT * FROM images WHERE product_name`)
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

func UpdateImage(db *sql.DB, updated_image Image) error {
	if updated_image.ImageSigned {
		stmt, err := db.Prepare("CALL update_image_signed(?, ?)")
		if err != nil {
			return err
		}
		defer stmt.Close()
		_, err = stmt.Exec(updated_image.ImageName, updated_image.ImageSigned)
		if err != nil {
			return err
		}
	}
	if updated_image.TrivySigned || updated_image.TrivyValid {
		stmt, err := db.Prepare("CALL update_trivy_flags(?, ?, ?)")
		if err != nil {
			return err
		}
		defer stmt.Close()
		_, err = stmt.Exec(updated_image.ImageName, updated_image.TrivySigned, updated_image.TrivyValid)
		if err != nil {
			return err
		}
	}
	if updated_image.SbomSigned || updated_image.SbomValid {
		stmt, err := db.Prepare("CALL update_sbom_flags(?, ?, ?)")
		if err != nil {
			return err
		}
		defer stmt.Close()
		_, err = stmt.Exec(updated_image.ImageName, updated_image.SbomSigned, updated_image.SbomValid)
		if err != nil {
			return err
		}
	}
	return nil
}

func DeleteImage(db *sql.DB, name string) error {
	_, err := db.Exec(`DELETE FROM images WHERE image_name = ?`, name)
	if err != nil {
		return err
	}
	return nil
}

func get_image_metadata(db *sql.DB, image_name string) (release string, product string) {
	// tbd
	return "", ""
}
