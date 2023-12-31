package api

import (
	DB "carbide-images-api/pkg/database"
	"carbide-images-api/pkg/objects"
	"database/sql"
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func decodeUser(w http.ResponseWriter, r *http.Request) (objects.User, error) {
	var user objects.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusBadRequest)
		log.Error(err)
		return user, err
	}
	return user, nil
}

func userPost(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	newUser, err := decodeUser(w, r)
	if err != nil {
		return
	}
	if err := DB.AddUser(db, newUser); err != nil {
		log.Error(err)
		httpJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	newUser, err = DB.GetUser(db, *newUser.Username)
	if err != nil {
		log.Error(err)
		return
	}
	log.WithFields(log.Fields{
		"user": *newUser.Username,
	}).Info("user has been successfully created")
	err = setAuthCookie(w, newUser)
	if err != nil {
		log.Error(err)
	}
	respondWithJSON(w, "user has been created")
	return
}

func userDelete(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	userToDelete, err := decodeUser(w, r)
	if err != nil {
		return
	}
	if err := DB.VerifyUser(db, userToDelete); err != nil {
		log.WithFields(log.Fields{
			"username": *userToDelete.Username,
			"error":    err,
		}).Error("invalid username or password")
		httpJSONError(w, "invalid username or password", http.StatusBadRequest)
		return
	}
	if err := DB.DeleteUserByUsername(db, *userToDelete.Username); err != nil {
		log.WithFields(log.Fields{
			"username": *userToDelete.Username,
			"error":    err,
		}).Error("failed to delete user")
		return
	}
	log.WithFields(log.Fields{
		"user": *userToDelete.Username,
	}).Info("user has been successfully deleted")
	respondWithJSON(w, "user has been deleted")
	return
}

func loginPost(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	login, err := decodeUser(w, r)
	if err != nil {
		return
	}
	if err := DB.VerifyUser(db, login); err != nil {
		log.Error(err)
		return
	}
	login, err = DB.GetUser(db, *login.Username)
	if err != nil {
		log.Error(err)
		return
	}
	err = setAuthCookie(w, login)
	if err != nil {
		log.Error(err)
	}
	log.WithFields(log.Fields{
		"user": *login.Username,
	}).Info("user logged in successfully")
	respondWithJSON(w, "login successfull")
}
