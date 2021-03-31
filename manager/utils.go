package manager

import (
	"log"
	"time"

	mail "github.com/xhit/go-simple-mail"
)

func sendConfirmEmail(addressTo map[string]bool, addressCC map[string]bool, subject string, body string) {
	log.Println("sendConfirmEmail start")

	server := mail.NewSMTPClient()

	//smtpUser := os.Getenv("SMTP_USR")
	//smtpPassword := os.Getenv("SMTP_PWD")
	smtpUser := "mailuser"
	smtpPassword := "mailpass"

	server.Host = "mail.sorintdev.it"
	server.Port = 25
	server.Username = smtpUser
	server.Password = smtpPassword
	// server.Encryption = mail.EncryptionSSL
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
	email = email.SetFrom("Papagaio <no-reply@sorint.it>")

	if addressTo != nil {
		for s, _ := range addressTo {
			log.Printf("Aggiungo address to " + s)
			email = email.AddTo(s)
		}
	}

	if addressCC != nil {
		for s, _ := range addressCC {
			log.Printf("Aggiungo address cc " + s)
			email = email.AddCc(s)
		}
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
