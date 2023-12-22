package database

import (
	"carbide-images-api/pkg/objects"
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

func AddUser(db *sql.DB, newuser objects.User) error {
	// hash user password
	if bytes, err := bcrypt.GenerateFromPassword([]byte(newuser.Password), 4); err != nil {
		return err
	} else {
		newuser.Password = string(bytes)
	}
	// insert new user object in database
	if _, err := db.Exec(
		"INSERT INTO users (username, password) VALUES (?, ?)",
		newuser.Username, newuser.Password); err != nil {
		return err
	}
	return nil
}

// returns error and full user object for given username
func GetUser(db *sql.DB, username string) (objects.User, error) {
	var user objects.User
	err := db.QueryRow(`SELECT * FROM users WHERE username = $1`, username).Scan(&user.Id, &user.Username, &user.Password)
	if err != nil {
		return user, err
	}
	return user, nil
}

// queries database for all users - returns err and a slice of all users; empty if err
func GetUsers(db *sql.DB) ([]objects.User, error) {
	var users []objects.User
	rows, err := db.Query(`SELECT * FROM users`)
	if err != nil {
		users = nil
		return users, err
	}
	defer rows.Close()
	for rows.Next() {
		var usr objects.User
		if err := rows.Scan(&usr.Id, &usr.Username, &usr.Password); err != nil {
			users = nil
			return users, err
		}
		users = append(users, usr)
	}
	if err = rows.Err(); err != nil {
		users = nil
		return users, err
	}
	return users, nil
}

// checks if user credentials are valid - returns nil on success; err otherwise
func VerifyUser(db *sql.DB, user objects.User) error {
	verified, err := GetUser(db, user.Username)
	if err != nil {
		return err
	}
	// compare the user's stored password with the one provided
	if err := bcrypt.CompareHashAndPassword([]byte(verified.Password), []byte(user.Password)); err != nil {
		return err
	}
	return nil
}

// delete corresponding row in users table
func DeleteUser(db *sql.DB, userid int64) error {
	if _, err := db.Exec(
		`DELETE FROM users WHERE id = $1;`, userid); err != nil {
		return err
	}
	return nil
}
