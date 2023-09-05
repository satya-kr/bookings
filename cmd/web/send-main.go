package main

import (
	"github.com/satya-kr/bookings/internal/models"
	mail "github.com/xhit/go-simple-mail/v2"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const smtp_server = "mail.astergo.in"
const smtp_usr = "gotest@astergo.in"
const smtp_pass = "Mail@test0321"
const smtp_port = "587"
const recipientEmail = "satyajit.kr.prajapati@gmail.com"

var pathToTemplate, _ = filepath.Abs("../../email-templates")

func listenForMail() {
	//run mail process on background using go routine
	go func() {
		for {
			msg := <-app.MailChan
			sendMsg(msg)
		}
	}()
}

func sendMsg(m models.EmailData) {
	server := mail.NewSMTPClient()
	server.Host = smtp_server
	server.Port, _ = strconv.Atoi(smtp_port)
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second
	server.Username = smtp_usr
	server.Password = smtp_pass
	server.Encryption = mail.EncryptionTLS

	client, err := server.Connect()
	if err != nil {
		errorLog.Println(err)
	}

	email := mail.NewMSG()
	email.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject)
	if m.Template == "" {
		email.SetBody(mail.TextHTML, m.Content)
	} else {

		file, err := os.ReadFile(path.Join(pathToTemplate, m.Template))
		if err != nil {
			app.ErrorLog.Println(err)
		}

		mailTemplate := string(file)
		msgToSend := strings.Replace(mailTemplate, "[%body%]", m.Content, 1)
		email.SetBody(mail.TextHTML, msgToSend)
	}

	err = email.Send(client)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Email Sent!")
	}
}
