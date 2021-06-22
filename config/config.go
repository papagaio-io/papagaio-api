package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"wecode.sorint.it/opensource/papagaio-api/common"
)

// Configuration contains all informations required to run papagaio
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
	//Agola address
	Agola AgolaConfig
	//Papagaio admin token
	AdminToken string

	//Cmd conficuration
	CmdConfig CmdConfig
	//Timers
	TriggersConfig TriggersConfig
	// Email configuration
	Email *EmailConfig

	TokenSigning TokenSigning
}

type TriggersConfig struct {
	OrganizationsDefaultTriggerTime uint
	RunFailedDefaultTriggerTime     uint
}

type AgolaConfig struct {
	AgolaAddr  string
	AdminToken string
}

type DbConfig struct {
	DbPath string
	DbName string
}

type CmdConfig struct {
	DefaultGatewayURL string
	Token             string
}

// Server contains all informations required to setup our config
type Server struct {
	// Port on which our config must listen and serve
	Port             string
	LocalHostAddress string
}

type KeycloakConfig struct {
	Realm         string
	AuthURL       string
	Resource      string
	PubKey        string
	TokenValidity int `json:"Token-validity"`
}

type EmailConfig struct {
	// Hostname/ip of the smtp server
	SMTPServer *string
	// Port of the smtp server
	SMTPPort *int
	// Username of the smtp server
	Username *string
	// Password of the smtp server
	Password *string
	// From
	From *string
	// Encryption
	Encryption *string
}

type TokenSigning struct {
	// token duration (defaults to 12 hours)
	Duration time.Duration `yaml:"duration"`
	// signing method: "hmac" or "rsa"
	Method string `yaml:"method"`
	// signing key. Used only with HMAC signing method
	Key string `yaml:"key"`
	// path to a file containing a pem encoded private key. Used only with RSA signing method
	PrivateKeyPath string `yaml:"privateKeyPath"`
	// path to a file containing a pem encoded public key. Used only with RSA signing method
	PublicKeyPath string `yaml:"publicKeyPath"`
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
}

func InitTokenSigninData(tokenSigning *TokenSigning) (*common.TokenSigningData, error) {
	sd := &common.TokenSigningData{Duration: tokenSigning.Duration}
	switch tokenSigning.Method {
	case "hmac":
		sd.Method = jwt.SigningMethodHS256
		if tokenSigning.Key == "" {
			return nil, errors.Errorf("empty token signing key for hmac method")
		}
		sd.Key = []byte(tokenSigning.Key)
	case "rsa":
		if tokenSigning.PrivateKeyPath == "" {
			return nil, errors.Errorf("token signing private key file for rsa method not defined")
		}
		if tokenSigning.PublicKeyPath == "" {
			return nil, errors.Errorf("token signing public key file for rsa method not defined")
		}

		sd.Method = jwt.SigningMethodRS256
		privateKeyData, err := ioutil.ReadFile(tokenSigning.PrivateKeyPath)
		if err != nil {
			return nil, errors.Errorf("error reading token signing private key: %w", err)
		}
		sd.PrivateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
		if err != nil {
			return nil, errors.Errorf("error parsing token signing private key: %w", err)
		}
		publicKeyData, err := ioutil.ReadFile(tokenSigning.PublicKeyPath)
		if err != nil {
			return nil, errors.Errorf("error reading token signing public key: %w", err)
		}
		sd.PublicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicKeyData)
		if err != nil {
			return nil, errors.Errorf("error parsing token signing public key: %w", err)
		}
	case "":
		return nil, errors.Errorf("missing token signing method")
	default:
		return nil, errors.Errorf("unknown token signing method: %q", tokenSigning.Method)
	}

	return sd, nil
}
