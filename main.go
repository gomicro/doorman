package main

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/gomicro/ledger"
	"github.com/gomicro/steward"
	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
)

const (
	app = "doorman"
)

var (
	version string
	config  configuration
	log     *ledger.Ledger
)

type configuration struct {
	Host string `default:"0.0.0.0"`
	Port string `default:"4567"`

	LogLevel string `default:"debug"`
}

func main() {
	configure()

	err := startService()
	if err != nil {
		log.Errorf("Something went wrong: %v", err.Error())
		return
	}

	log.Info("Server stopping")
}

func configure() {
	err := envconfig.Process(app, &config)
	if err != nil {
		fmt.Printf("failed to process env vars: %v\n", err.Error())
		os.Exit(1)
	}

	log = ledger.New(os.Stdout, ledger.ParseLevel(config.LogLevel))
	log.Debug("Logger configured")

	if version == "" {
		version = "local-dev"
	}

	steward.SetStatusEndpoint("/v1/status")

	log.Info("Configuration complete")
}

func startService() error {
	log.Infof("Starting service (%v)", version)
	log.Infof("Listening on %v:%v", config.Host, config.Port)

	http.Handle("/", registerEndpoints())
	return http.ListenAndServe(net.JoinHostPort(config.Host, config.Port), nil)
}

func registerEndpoints() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/o/oauth2/auth", handleGetGoogleAuth).Methods("GET").Queries("client_id", "{client_id}", "redirect_uri", "{redirect_uri}", "response_type", "{response_type}")
	r.HandleFunc("/o/oauth2/auth", handlePostGoogleAuth).Methods("POST")
	r.HandleFunc("/oauth2/v3/userinfo", handleUserInfo)

	return r
}
