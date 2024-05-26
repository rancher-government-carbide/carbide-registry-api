package middleware

import (
	"carbide-registry-api/pkg/api/utils"
	license "carbide-registry-api/pkg/license"
	"crypto/rsa"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

func SessionAuth(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if !Authorized(r) {
			utils.HttpJSONError(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func Login(w http.ResponseWriter, license license.CarbideLicense, licensePubkeys []*rsa.PublicKey) error {
	userID, err := Authenticate(license, licensePubkeys)
	if err != nil {
		utils.HttpJSONError(w, err.Error(), http.StatusUnauthorized)
		return err
	}
	err = Authorize(w, userID)
	if err != nil {
		return err
	}
	log.Info().Str("customerID", userID).Msg("login succeeded")
	utils.RespondWithJSON(w, "login succeeded")
	return nil
}

func Logout(w http.ResponseWriter) error {
	if err := terminateSessionToken(w); err != nil {
		return err
	}
	log.Info().Msg("logout succeeded")
	utils.RespondWithJSON(w, "logout succeeded")
	return nil
}

// parse customerID from a valid license
func Authenticate(carbideLicense license.CarbideLicense, licensePubkeys []*rsa.PublicKey) (string, error) {
	customerID, err := license.ParseCarbideLicense(*carbideLicense.License, licensePubkeys)
	if err != nil {
		log.Error().Err(err)
		return "", err
	}
	log.Info().Str("customerID", *customerID).Msg("license parsed successfully")
	return *customerID, nil
}

// provide user with session token
func Authorize(w http.ResponseWriter, userID string) error {
	expiration := time.Now().Add(time.Hour * 8)
	if err := setSessionToken(w, userID, expiration); err != nil {
		utils.HttpJSONError(w, err.Error(), http.StatusInternalServerError)
		log.Error().Err(err)
		return err
	}
	return nil
}

// true if request has a valid session token
func Authorized(r *http.Request) bool {
	token, err := parseSessionToken(r)
	if err != nil {
		log.Error().Msg("failed to parse session token")
		return false
	}
	err = verifySessionToken(token)
	if err != nil {
		log.Info().Msg("user is unauthorized")
		return false
	}
	return true
}

var tokenSigningSecret = os.Getenv("JWTSECRET")

func setSessionToken(w http.ResponseWriter, userID string, expiry time.Time) error {
	token, err := generateSessionToken(userID, expiry)
	if err != nil {
		return err
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
	return nil
}

func generateSessionToken(userID string, expiry time.Time) (string, error) {
	tokenSigningSecret := os.Getenv("JWTSECRET")
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = expiry.Unix()
	claims["userid"] = userID
	tokenString, err := token.SignedString([]byte(tokenSigningSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// check for proper signing method and return key to be used for signature verification
func sessionTokenKeyFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return "", errors.New("invalid JWT signing method")
	}
	return []byte(tokenSigningSecret), nil
}

func parseSessionToken(r *http.Request) (*jwt.Token, error) {
	c, err := r.Cookie("token")
	if err != nil {
		return nil, err
	}
	tokenString := c.Value
	token, err := jwt.Parse(tokenString, sessionTokenKeyFunc)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func verifySessionToken(token *jwt.Token) error {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("failed to parse JWT claims")
	}
	exp := claims["exp"].(float64)
	if int64(exp) < time.Now().Local().Unix() {
		return errors.New("token expired")
	}
	return nil
}

func terminateSessionToken(w http.ResponseWriter) error {
	if err := setSessionToken(w, "", time.Now()); err != nil {
		return err
	}
	return nil
}
