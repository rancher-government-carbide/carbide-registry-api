package api

import (
	DB "carbide-images-api/pkg/database"
	"carbide-images-api/pkg/objects"
	"database/sql"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func createUserHandler(db *sql.DB) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var newUser objects.User
		err := decodeJSONObject(w, r, &newUser)
		if err != nil {
			log.Error(err)
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
	return http.HandlerFunc(fn)
}

func deleteUserHandler(db *sql.DB) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var userToDelete objects.User
		err := decodeJSONObject(w, r, &userToDelete)
		if err != nil {
			log.Error(err)
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
	return http.HandlerFunc(fn)
}

func loginHandler(db *sql.DB) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var userLoggingIn objects.User
		err := decodeJSONObject(w, r, &userLoggingIn)
		if err != nil {
			log.Error(err)
			return
		}
		if err := DB.VerifyUser(db, userLoggingIn); err != nil {
			log.Error(err)
			httpJSONError(w, "invalid username or password", http.StatusUnauthorized)
			return
		}
		userLoggingIn, err = DB.GetUser(db, *userLoggingIn.Username)
		if err != nil {
			log.Error(err)
			return
		}
		err = setAuthCookie(w, userLoggingIn)
		if err != nil {
			log.Error(err)
		}
		log.WithFields(log.Fields{
			"user": *userLoggingIn.Username,
		}).Info("user logged in successfully")
		respondWithJSON(w, "login successfull")
	}
	return http.HandlerFunc(fn)
}
