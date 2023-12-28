package api

import (
	DB "carbide-images-api/pkg/database"
	"carbide-images-api/pkg/objects"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func userPost(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var newUser objects.User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusBadRequest)
		log.Error(err)
		return
	}
	if err := DB.AddUser(db, newUser); err != nil {
		log.Error(err)
		return
	}
	newUser, err = DB.GetUser(db, *newUser.Username)
	if err != nil {
		log.Error(err)
		return
	}
	log.WithFields(log.Fields{
		"user": *newUser.Username,
	}).Info("User has been successfully created")
	token, err := generateJWT(newUser)
	if err != nil {
		log.Error(err)
		return
	}
	ck := http.Cookie{
		Name:     "token",
		Value:    token,
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}
	http.SetCookie(w, &ck)
	respondSuccess(w)
	return
}

func userDelete(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	userid, err := verifyJWT(r)
	if err != nil {
		log.Error(err)
		http.Error(w, "Missing or expired cookie", 401)
		return
	}
	var userToDelete objects.User
	err = json.NewDecoder(r.Body).Decode(&userToDelete)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusBadRequest)
		log.Error(err)
		return
	}
	if err := DB.VerifyUser(db, userToDelete); err != nil {
		log.Error(err)
		return
	}
	if err := DB.DeleteUser(db, userid); err != nil {
		log.Error(err)
		log.Errorf("Failed to delete user with id %d", userid)
		w.Write([]byte(fmt.Sprintf("Failed to delete user with id %d", userid)))
		return
	}
	log.WithFields(log.Fields{
		"user": *userToDelete.Username,
	}).Info("User has been successfully deleted or didn't exist in the first place")
	respondSuccess(w)
	return
}

func loginPost(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var login objects.User
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusBadRequest)
		log.Error(err)
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
	token, err := generateJWT(login)
	if err != nil {
		log.Error(err)
		return
	}
	ck := http.Cookie{
		Name:     "token",
		Value:    token,
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}
	http.SetCookie(w, &ck)
	log.WithFields(log.Fields{
		"user": *login.Username,
	}).Info("User logged in successfully")
	respondSuccess(w)
}
