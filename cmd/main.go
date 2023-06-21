package main

import (
	"carbide-api/cmd/api"
	"log"
	"net/http"
	"os"
)

func main() {
	db_user := os.Getenv("DBUSER")
	if db_user == "" {
		db_user = "clayton"
	}
	db_pass := os.Getenv("DBPASS")
	if db_pass == "" {
		db_pass = "applevisioncurescancer"
	}
	db_host := os.Getenv("DBHOST")
	if db_host == "" {
		db_host = "127.0.0.1"
	}
	db_port := os.Getenv("DBPORT")
	if db_port == "" {
		db_port = "3306"
	}
	db_name := os.Getenv("DBNAME")
	if db_name == "" {
		db_name = "carbide"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	db, err := api.DatabaseInit(db_user, db_pass, db_host, db_port, db_name)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Database connected!")
	defer db.Close()

	log.Printf("Starting server on port " + port + "...")
	http.ListenAndServe("0.0.0.0:"+port, &api.Serve{DB: db})
}
