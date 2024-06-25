package main

import (
	"carbide-registry-api/pkg/api"
	"carbide-registry-api/pkg/azure"
	"carbide-registry-api/pkg/database"
	"crypto/rsa"
	"net/http"
	"os"

	"github.com/ebauman/golicense/pkg/certificate"
	log "github.com/sirupsen/logrus"
)

func main() {
	DBUSER := os.Getenv("DBUSER")
	if DBUSER == "" {
		DBUSER = "clayton"
	}
	DBPASS := os.Getenv("DBPASS")
	if DBPASS == "" {
		DBPASS = "applevisioncurescancer"
	}
	DBHOST := os.Getenv("DBHOST")
	if DBHOST == "" {
		DBHOST = "127.0.0.1"
	}
	DBPORT := os.Getenv("DBPORT")
	if DBPORT == "" {
		DBPORT = "3306"
	}
	DBNAME := os.Getenv("DBNAME")
	if DBNAME == "" {
		DBNAME = "carbide"
	}
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "5000"
	}

	GOLICENSE_KEY := os.Getenv("GOLICENSE_KEY")
	if GOLICENSE_KEY == "" {
		log.Fatal("Missing GOLICENSE_KEY env variable (carbide license private key), exiting...")
	}
	licensePrivkey, err := certificate.PEMToPrivateKey([]byte(GOLICENSE_KEY))
	if err != nil {
		panic(err)
	}
	licensePubkeys := []*rsa.PublicKey{&licensePrivkey.PublicKey}

	db, err := database.Init(DBUSER, DBPASS, DBHOST, DBPORT, DBNAME)
	if err != nil {
		log.Error("Database connection failed, exiting...")
		log.Fatal(err)
	}
	log.Info("Database connected!")
	defer db.Close()
	log.Info("Initializing db schema...")
	err = database.SchemaInit(db)
	if err != nil {
		log.Error("Database schema init failed, exiting...")
		log.Fatal(err)
	}

	clientFactory, err := azure.NewAzureClients()
	log.Info("Starting server on PORT " + PORT + "...")
	log.Fatal(http.ListenAndServe(":"+PORT, api.NewRouter(db, clientFactory, licensePrivkey, licensePubkeys)))
}
