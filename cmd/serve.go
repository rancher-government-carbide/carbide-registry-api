package cmd

import (
	"carbide-registry-api/pkg/api"
	"carbide-registry-api/pkg/azure"
	"carbide-registry-api/pkg/database"
	"crypto/rsa"
	"fmt"
	"net/http"
	"os"

	"github.com/ebauman/golicense/pkg/certificate"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var port int
var private bool

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the carbide REST API",
	// Long: ``,
	Run: func(cmd *cobra.Command, args []string) {

		GOLICENSE_KEY := os.Getenv("GOLICENSE_KEY")
		if GOLICENSE_KEY == "" {
			log.Fatal().Msg("Missing GOLICENSE_KEY env variable (carbide license private key), exiting...")
		}
		licensePrivkey, err := certificate.PEMToPrivateKey([]byte(GOLICENSE_KEY))
		if err != nil {
			panic(err)
		}
		licensePubkeys := []*rsa.PublicKey{&licensePrivkey.PublicKey}

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

		db, err := database.Init(DBUSER, DBPASS, DBHOST, DBPORT, DBNAME)
		if err != nil {
			log.Error().Msg("Database connection failed, exiting...")
			log.Fatal().Err(err)
		}
		log.Info().Msg("Database connected!")
		defer db.Close()
		log.Info().Msg("Initializing db schema...")
		err = database.SchemaInit(db)
		if err != nil {
			log.Error().Msg("Database schema init failed, exiting...")
			log.Fatal().Err(err)
		}

		portString := fmt.Sprint(port)

		if private {
			clientFactory, err := azure.NewAzureClients()
			if err != nil {
				log.Fatal().Err(err)
			}
			log.Info().Msg("Starting server in private mode. This disables authentication and enables the carbide account creation endpoint.")
			log.Info().Msg("Listening on port " + portString + "...")
			log.Fatal().Err(http.ListenAndServe(":"+portString, api.PrivateRouter(db, clientFactory, licensePrivkey)))
		} else {
			log.Info().Msg("Starting server in public mode. This enables authentication and disables the carbide account creation endpoint.")
			log.Info().Msg("Listening on port " + portString + "...")
			log.Fatal().Err(http.ListenAndServe(":"+portString, api.PublicRouter(db, licensePubkeys)))
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command is called directly
	serveCmd.Flags().BoolVarP(&private, "ssf", "s", false, "Disable auth for internal pipeline use")
	serveCmd.Flags().IntVarP(&port, "port", "p", 5000, "Port to listen on")
}
