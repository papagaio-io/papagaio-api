package config

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"io/ioutil"
	"log"
)

// Configuration contains all informations required to run stroodle
type Configuration struct {
	// Server configuration
	Server Server
	// Enable the logging of the http requests
	LogHTTPRequest bool
	// Database specific informations
	Database DbConfig
	// Disable SSL certificate validation of keycloak
	DisableSSLCertificateValidation bool
	// Keycloak configuration
	Keycloak KeycloakConfig
}

// DbConfig contains all informations required to connect to repository
type DbConfig struct {
	// Type of the repository to use. Actually, only postgres is supported
	DbType string
	// Address(or hostname) of the repository
	DbAddr string
	// Port of the repository
	DbPort string
	// Name of the repository
	DbName string
	// User of the repository
	DbUser string
	// Password of the repository
	DbPass string
	// SSL of the
	Dbsslmode string
}

// Server contains all informations required to setup our config
type Server struct {
	// Port on which our config must listen and serve
	Port string
}

type KeycloakConfig struct {
	Realm         string
	AuthURL       string
	Resource      string
	PubKey        string
	TokenValidity int `json:"Token-validity"`
}

// Config contains global configuration read with config.ReadConfig()
var Config Configuration

var KeycloakPubKey interface{}

func readConfig() {
	var raw []byte
	var err error

	if raw, err = ioutil.ReadFile("/app/config.json"); err != nil {
		if raw, err = ioutil.ReadFile("config.json"); err != nil {
			log.Fatal("Unable to read configuration file: ", err)
		}
	}

	if err = json.Unmarshal(raw, &Config); err != nil {
		log.Fatal("Unable to parse configuration file: ", err)
	}
}

// SetupConfig load the configuration from config.json and set config.Config to it
func SetupConfig() {
	readConfig()
	readKeycloakConfig()
}

// parsePubKey parsing the public key generated from keycloak
func parsePubKey(raw []byte) interface{} {
	block, _ := pem.Decode(raw)
	if block == nil {
		log.Fatal("Unable to parse the public key")
	}

	pubkey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Fatal("Unable to parse the public key: ", err)
	}

	return pubkey
}

// readKeycloakConfig read the keycloak configuration file from "keycloak.json"
func readKeycloakConfig() {
	var err error
	var raw []byte

	if raw, err = ioutil.ReadFile("/app/keycloak.json"); err != nil {
		if raw, err = ioutil.ReadFile("keycloak.json"); err != nil {
			log.Fatal("Unable to read keycloak.json: ", err)
		}
	}
	if err = json.Unmarshal(raw, &Config.Keycloak); err != nil {
		log.Fatal("Unable to parse keycloak.json", err)
	}
	if raw, err = ioutil.ReadFile("/app/pubkey.pub"); err != nil {
		if raw, err = ioutil.ReadFile("keys/pubkey.pub"); err != nil {
			log.Fatal("Public key (pubkey.pub) was not found: ", err)
		}
	}

	KeycloakPubKey = parsePubKey(raw)
	Config.Keycloak.PubKey = string(raw)
}