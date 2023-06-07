package main

import (
	"fmt"
	"log"
	"net/http"

	// "net/http"
	"carbide-api/cmd/api"
	"os"
)

func main() {
	dbuser := os.Getenv("DBUSER")
	if dbuser == "" {
		dbuser = "root"
	}
	dbpass := os.Getenv("DBPASS")
	if dbpass == "" {
		dbpass = ""
	}
	dburl := os.Getenv("DBURL")
	if dburl == "" {
		dburl = "0.0.0.0:26257"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	db, err := api.DatabaseInit(dbuser, dbpass, dburl)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Database Connected!")
	defer db.Close()

	fmt.Printf("Starting Server...\n")

	a := &api.Api{
		LoginHandler: &api.LoginHandler{DB: db},
		UserHandler:  &api.UserHandler{DB: db},
	}

	http.ListenAndServe("0.0.0.0:"+port, a)
	fmt.Printf("Now listening on port " + port + ".")
}
