package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func Init(db_user string, db_pass string, db_host string, db_port string, db_name string) (*sql.DB, error) {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", db_user, db_pass, db_host, db_port, db_name)
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return db, err
	} else if err = db.Ping(); err != nil {
		return db, err
	}
	return db, err
}

func SchemaInit(db *sql.DB) error {
	if err := ensureProductTable(db); err != nil {
		return err
	}
	if err := ensureReleasesTable(db); err != nil {
		return err
	}
	if err := ensureImagesTable(db); err != nil {
		return err
	}
	if err := ensureReleaseImageMappingTable(db); err != nil {
		return err
	}
	if err := ensureCreateUpdateImageProcedure(db); err != nil {
		return err
	}
	if err := ensureUpdateSbomFlagsProcedure(db); err != nil {
		return err
	}
	if err := ensureUpdateTrivyFlagsProcedure(db); err != nil {
		return err
	}
	if err := ensureUpdateImageSignedProcedure(db); err != nil {
		return err
	}
	return nil
}

func ensureProductTable(db *sql.DB) error {
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS product (
		  id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		  name VARCHAR(255) NOT NULL,
		  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		)`); err != nil {
		return err
	}
	return nil
}

func ensureReleasesTable(db *sql.DB) error {
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS releases (
		  id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		  product_id INT NOT NULL,
		  name VARCHAR(255) NOT NULL,
		  tarball_link VARCHAR(255),
		  CONSTRAINT fk_release_product FOREIGN KEY (product_id) REFERENCES product(id) ON DELETE CASCADE,
		  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		)`); err != nil {
		return err
	}
	return nil
}

func ensureImagesTable(db *sql.DB) error {
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS images (
		  id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		  image_name VARCHAR(255) NOT NULL,
		  CONSTRAINT unique_image_name UNIQUE (image_name),
		  image_signed BOOLEAN DEFAULT false,
		  trivy_signed BOOLEAN DEFAULT false,
		  trivy_valid BOOLEAN DEFAULT false,
		  sbom_signed BOOLEAN DEFAULT false,
		  sbom_valid BOOLEAN DEFAULT false,
		  last_scanned_at DATETIME DEFAULT NULL,
		  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		)`); err != nil {
		return err
	}
	return nil
}

func ensureReleaseImageMappingTable(db *sql.DB) error {
	if _, err := db.Exec(`
   	CREATE TABLE IF NOT EXISTS release_image_mapping (
   	  id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
   	  release_id INT NOT NULL,
   	  image_id INT NOT NULL,
   	  CONSTRAINT fk_releases_map FOREIGN KEY (release_id) REFERENCES releases(id) ON DELETE CASCADE,
   	  CONSTRAINT fk_images_map FOREIGN KEY (image_id) REFERENCES images(id) ON DELETE CASCADE,
   	  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   	  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
   	)`); err != nil {
		return err
	}
	return nil
}

func ensureCreateUpdateImageProcedure(db *sql.DB) error {
	if _, err := db.Exec(`
		CREATE PROCEDURE IF NOT EXISTS create_update_image(
		    IN p_product_name VARCHAR(255),
		    IN p_release_name VARCHAR(255),
		    IN p_image_name VARCHAR(255)
		)
		BEGIN
		    DECLARE p_product_id INT;
		    DECLARE p_release_id INT;
		    DECLARE p_image_id INT;
		    DECLARE p_mapping_id INT;
		    -- Check if the product exists
		    SELECT id INTO p_product_id
		    FROM product
		    WHERE name = p_product_name;
		    -- If the product doesn't exist, insert a new record
		    IF p_product_id IS NULL THEN
		        INSERT INTO product (name)
		        VALUES (p_product_name);
		        SET p_product_id = LAST_INSERT_ID();
		    END IF;
		    -- Check if the release exists
		    SELECT id INTO p_release_id
		    FROM releases
		    WHERE product_id = p_product_id AND name = p_release_name;
		    -- If the release doesn't exist, insert a new record
		    IF p_release_id IS NULL THEN
		        INSERT INTO releases (product_id, name)
		        VALUES (p_product_id, p_release_name);
		        SET p_release_id = LAST_INSERT_ID();
		    END IF;
		    -- Check if the image exists
		    SELECT id INTO p_image_id
		    FROM images
		    WHERE image_name = TRIM(p_image_name);
		    -- If the image doesn't exist, insert a new record
		    IF p_image_id IS NULL THEN
		        INSERT INTO images (image_name)
		        VALUES (TRIM(p_image_name));
		        SET p_image_id = LAST_INSERT_ID();
		    END IF;
		    -- Check if the release image mapping record exists
		    SELECT id INTO p_mapping_id
		    FROM release_image_mapping
		    WHERE release_id = p_release_id and image_id = p_image_id;
		    -- If the mapping doesn't exist, insert a new record
		    IF p_mapping_id IS NULL THEN
		        INSERT INTO release_image_mapping (release_id, image_id)
		        VALUES (p_release_id, p_image_id);
		        SET p_mapping_id = LAST_INSERT_ID();
		    END IF;
		END;
	`); err != nil {
		return err
	}
	return nil
}

func ensureUpdateImageSignedProcedure(db *sql.DB) error {
	if _, err := db.Exec(`
		CREATE PROCEDURE IF NOT EXISTS update_image_signed(
		    IN p_image_name VARCHAR(255),
		    IN p_image_signed BOOLEAN
		)
		BEGIN
		    DECLARE p_image_id INT;
		    -- Check if the image exists
		    SELECT id INTO p_image_id
		    FROM images
		    WHERE image_name = TRIM(p_image_name);
		    -- If the image doesn't exist, insert a new record
		    IF p_image_id IS NULL THEN
		        INSERT INTO images (image_name, image_signed)
		        VALUES (TRIM(p_image_name), p_image_signed);
		    ELSE
		        -- Update the record
		        UPDATE images
		        SET image_signed = p_image_signed
		        WHERE image_name = TRIM(p_image_name);
		    END IF;
		END;
	`); err != nil {
		return err
	}
	return nil
}

func ensureUpdateTrivyFlagsProcedure(db *sql.DB) error {
	if _, err := db.Exec(`
		CREATE PROCEDURE IF NOT EXISTS update_trivy_flags(
		    IN p_image_name VARCHAR(255),
		    IN p_trivy_signed BOOLEAN,
		    IN p_trivy_valid BOOLEAN
		)
		BEGIN
		    DECLARE p_image_id INT;
		    -- Check if the image exists
		    SELECT id INTO p_image_id
		    FROM images
		    WHERE image_name = TRIM(p_image_name);
		    -- If the image doesn't exist, insert a new record
		    IF p_image_id IS NULL THEN
		        INSERT INTO images (image_name, trivy_signed, trivy_valid, last_scanned_at)
		        VALUES (TRIM(p_image_name), p_trivy_signed, p_trivy_valid, CURRENT_TIMESTAMP);
		    ELSE
		        -- Update the record
		        UPDATE images
		        SET trivy_signed = p_trivy_signed,
		            trivy_valid = p_trivy_valid,
		            last_scanned_at = CURRENT_TIMESTAMP
		        WHERE image_name = TRIM(p_image_name);
		    END IF;
		END;
	`); err != nil {
		return err
	}
	return nil
}

func ensureUpdateSbomFlagsProcedure(db *sql.DB) error {
	if _, err := db.Exec(`
		CREATE PROCEDURE IF NOT EXISTS update_sbom_flags(
		    IN p_image_name VARCHAR(255),
		    IN p_sbom_signed BOOLEAN,
		    IN p_sbom_valid BOOLEAN
		)
		BEGIN
		    DECLARE p_image_id INT;
		    -- Check if the image exists
		    SELECT id INTO p_image_id
		    FROM images
		    WHERE image_name = TRIM(p_image_name);
		    -- If the image doesn't exist, insert a new record
		    IF p_image_id IS NULL THEN
		        INSERT INTO images (image_name, sbom_signed, sbom_valid)
		        VALUES (TRIM(p_image_name), p_sbom_signed, p_sbom_valid);
		    ELSE
		        -- Update the record
		        UPDATE images
		        SET sbom_signed = p_sbom_signed,
		            sbom_valid = p_sbom_valid
		        WHERE image_name = TRIM(p_image_name);
		    END IF;
		END;
	`); err != nil {
		return err
	}
	return nil
}
