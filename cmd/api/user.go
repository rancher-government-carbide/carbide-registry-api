package api

import (
	"carbide-api/cmd/api/objects"
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

type UserHandler struct {
	DB *sql.DB
}

type LoginHandler struct {
	DB *sql.DB
}

func userPost(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	var newuser objects.User

	// read json payload into new user object
	if err := parseJSON(w, r, &newuser); err != nil {
		log.Print(err)
		return
	}

	// add user to database
	err := objects.AddUser(db, newuser)
	if err != nil {
		log.Print(err)
		return
	}

	// reassign stored values in user object
	newuser, err = objects.GetUser(db, newuser.Username)
	if err != nil {
		log.Print(err)
		return
	}
	log.Printf("User %s has been successfully created", newuser.Username)

	// respond with jwt
	token, err := generateJWT(newuser)
	if err != nil {
		log.Print(err)
		return
	}
	fmt.Print(token, "\n")

	// provide token as cookie to frontend
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

	// parse username from jwt token
	userid, err := verifyJWT(r)
	if err != nil {
		log.Print(err)
		http.Error(w, "Missing or expired cookie", 401)
		return
	}

	var user objects.User
	// parse username and password into JSON object
	if err := parseJSON(w, r, &user); err != nil {
		log.Print(err)
		return
	}
	// check credentials
	if err := objects.VerifyUser(db, user); err != nil {
		log.Print(err)
		return
	}

	// remove row with corresponding userid
	if err := objects.DeleteUser(db, userid); err != nil {
		log.Print(err)
		log.Printf("Failed to delete user with id %d", userid)
		w.Write([]byte(fmt.Sprintf("Failed to delete user with id %d", userid)))
		return
	}
	respondSuccess(w)

	return
}

func loginPost(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	var login objects.User

	// payload should contain username and password of user (other fields are ignored)
	if err := parseJSON(w, r, &login); err != nil {
		log.Print(err)
		return
	}
	// checks user+pass against stored database values
	if err := objects.VerifyUser(db, login); err != nil {
		log.Print(err)
		return
	}
	// fill user object with necessary information (userid)
	login, err := objects.GetUser(db, login.Username)
	if err != nil {
		log.Print(err)
		return
	}
	// creates jwt for user
	token, err := generateJWT(login)
	if err != nil {
		log.Print(err)
		return
	}
	// provide token as cookie to frontend
	ck := http.Cookie{
		Name:     "token",
		Value:    token,
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}
	// responds with message to store cookie in browser for future requests
	http.SetCookie(w, &ck)

	log.Printf("%s logged in successfully", login.Username)
	respondSuccess(w)

	return
}
