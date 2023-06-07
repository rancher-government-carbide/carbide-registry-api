package api

import (
	"carbide-api/cmd/api/objects"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func middleware(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
}

func login_middleware(w http.ResponseWriter, r *http.Request) {
	_, err := verifyJWT(r)
	if err == nil {
		log.Print("User is already logged in\n")
		w.Write([]byte(fmt.Sprintf("User is already logged in")))
		return
	}
}

func enableCors(w http.ResponseWriter) {
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Accept-Charset, Accept-Language, Access-Control-Request-Method, Access-Control-Request-Headers, Authorization, Content-Length, Content-Type, Cookie, Date, Forwarded, Origin, User-Agent")
	(w).Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	(w).Header().Set("Vary", "Origin")
	(w).Header().Set("Vary", "Access-Control-Request-Method")
	(w).Header().Set("Vary", "Access-Control-Request-Headers")
}

// generate JWT from given user - returns err and token
func generateJWT(user objects.User) (string, error) {

	// pull secret from environment
	secret := os.Getenv("JWTSECRET")

	// generate new jwt
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	// add claims payload
	claims["exp"] = time.Now().Add(time.Minute * 60).Unix()
	claims["userid"] = fmt.Sprint(user.Id)

	// stringify token
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// checks if http request is authorized/logged in - returns error and username string; empty if err
func verifyJWT(r *http.Request) (int64, error) {

	secret := os.Getenv("JWTSECRET")

	// get token from cookie
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			return 0, err
		}
		// For any other type of error, return a bad request status
		return 0, err
	}

	// Get the JWT string from the cookie
	tokenstring := c.Value

	// parse and check token validity
	token, err := jwt.Parse(tokenstring, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", errors.New("Invalid JWT Token")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return 0, err
	}

	// parse claims from token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("Failed to parse JWT claims")
	}

	// check if token is expired
	exp := claims["exp"].(float64)
	if int64(exp) < time.Now().Local().Unix() {
		return 0, errors.New("Token Expired")
	}

	s_userid := claims["userid"].(string)
	userid, err := strconv.ParseInt(s_userid, 10, 64)

	return userid, nil
}

func terminateJWT() {
	// replace jwt with another that expires immediately
}
