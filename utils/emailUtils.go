package utils

import (
	"log"
	"net/url"
	"os"
	"time"

	mail "github.com/xhit/go-simple-mail"
	"wecode.sorint.it/opensource/papagaio-api/config"
)

const defaultSMTPServer string = "mail.sorintdev.it"
const defaultSMTPPort int = 25
const defaultFrom string = "Papagaio <no-reply@sorint.it>"
const defaultEncryption string = "NONE"

func SendConfirmEmail(addressTo map[string]bool, addressCC map[string]bool, subject string, body string) {
	log.Println("sendConfirmEmail start")

	server := mail.NewSMTPClient()

	smtpUser := getUsername()
	smtpPassword := getPassword()

	server.Host = getSMTPServer()
	server.Port = getSMTPPort()
	server.Username = smtpUser
	server.Password = smtpPassword

	switch getEncryption() {
	case "NONE":
		server.Encryption = mail.EncryptionNone
	case "SSL":
		server.Encryption = mail.EncryptionSSL
	case "TLS":
		server.Encryption = mail.EncryptionTLS
	default:
		server.Encryption = mail.EncryptionNone
	}

	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	smtpClient, err := server.Connect()

	if err != nil {
		log.Fatal("Unable connect to sorint smtp: ", err)
	} else {
		log.Printf("connected with user " + smtpUser)
	}

	email := mail.NewMSG()
	email = email.SetFrom(getFrom())

	for s := range addressTo {
		log.Printf("Aggiungo address to " + s)
		email = email.AddTo(s)
	}

	for s := range addressCC {
		log.Printf("Aggiungo address cc " + s)
		email = email.AddCc(s)
	}

	email = email.SetSubject(subject)
	email.SetBody(mail.TextHTML, body)

	err = email.Send(smtpClient)
	if err != nil {
		log.Println("Email send error ", err)
	} else {
		log.Println("Email Sent")
	}

	log.Println("sendConfirmEmail end")
}

func CanSendEmail() bool {
	_, err := url.ParseRequestURI(getSMTPServer())
	return err == nil && getSMTPPort() > 0 && len(getUsername()) > 0 && len(getPassword()) > 0 && len(getFrom()) > 0
}

func getUsername() string {
	if config.Config.Email != nil && config.Config.Email.Username != nil {
		return *config.Config.Email.Username
	}

	return os.Getenv("SMTP_USR")
}

func getPassword() string {
	if config.Config.Email != nil && config.Config.Email.Password != nil {
		return *config.Config.Email.Password
	}

	return os.Getenv("SMTP_PWD")
}

func getSMTPServer() string {
	if config.Config.Email != nil && config.Config.Email.SMTPPort != nil {
		return *config.Config.Email.SMTPServer
	}

	return defaultSMTPServer
}

func getSMTPPort() int {
	if config.Config.Email != nil && config.Config.Email.SMTPPort != nil {
		return *config.Config.Email.SMTPPort
	}

	return defaultSMTPPort
}

func getFrom() string {
	if config.Config.Email != nil && config.Config.Email.From != nil {
		return *config.Config.Email.From
	}

	return defaultFrom
}

func getEncryption() string {
	if config.Config.Email != nil && config.Config.Email.Encryption != nil {
		return *config.Config.Email.Encryption
	}

	return defaultEncryption
}
