package api

import (
	"carbide-registry-api/pkg/api/middleware"
	"carbide-registry-api/pkg/api/utils"
	DB "carbide-registry-api/pkg/database"
	"carbide-registry-api/pkg/objects"
	"database/sql"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func authCheckHandler() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		middleware.Authorized(w, r)
		return
	}
	return http.HandlerFunc(fn)
}

func createUserHandler(db *sql.DB) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var newUser objects.User
		err := utils.DecodeJSONObject(w, r, &newUser)
		if err != nil {
			log.Error(err)
			return
		}
		if err := DB.AddUser(db, newUser); err != nil {
			log.Error(err)
			utils.RespondError(w, err.Error(), http.StatusBadRequest)
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
		err = middleware.Authenticate(w, newUser)
		if err != nil {
			log.Error(err)
		}
		utils.Respond(w, "user has been created")
		return
	}
	return http.HandlerFunc(fn)
}

func deleteUserHandler(db *sql.DB) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var userToDelete objects.User
		err := utils.DecodeJSONObject(w, r, &userToDelete)
		if err != nil {
			log.Error(err)
			return
		}
		if err := DB.VerifyUser(db, userToDelete); err != nil {
			log.WithFields(log.Fields{
				"username": *userToDelete.Username,
				"error":    err,
			}).Error("invalid username or password")
			utils.RespondError(w, "invalid username or password", http.StatusBadRequest)
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
		utils.Respond(w, "user has been deleted")
		return
	}
	return http.HandlerFunc(fn)
}

func loginHandler(db *sql.DB) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var userLoggingIn objects.User
		err := utils.DecodeJSONObject(w, r, &userLoggingIn)
		if err != nil {
			log.Error(err)
			return
		}
		if err := DB.VerifyUser(db, userLoggingIn); err != nil {
			log.Error(err)
			utils.RespondError(w, "invalid username or password", http.StatusUnauthorized)
			return
		}
		userLoggingIn, err = DB.GetUser(db, *userLoggingIn.Username)
		if err != nil {
			log.Error(err)
			return
		}
		err = middleware.Authenticate(w, userLoggingIn)
		if err != nil {
			log.Error(err)
		}
		log.WithFields(log.Fields{
			"user": *userLoggingIn.Username,
		}).Info("user logged in successfully")
		utils.Respond(w, "login successfull")
	}
	return http.HandlerFunc(fn)
}
