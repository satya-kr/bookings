package main

import (
	"encoding/gob"
	"github.com/satya-kr/bookings/internal/config"
	"github.com/satya-kr/bookings/internal/driver"
	"github.com/satya-kr/bookings/internal/handlers"
	"github.com/satya-kr/bookings/internal/helpers"
	"github.com/satya-kr/bookings/internal/models"
	"github.com/satya-kr/bookings/internal/render"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
)

const port = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

//const smtp_server = "mail.astergo.in"
//const smtp_usr = "gotest@astergo.in"
//const smtp_pass = "Mail@test0321"
//const smtp_port = "587"
//const recipientEmail = "satyajit.kr.prajapati@gmail.com"

func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	defer close(app.MailChan)
	listenForMail()

	//msg := models.EmailData{
	//	To:      "satyajit.kr.prajapati@gmail.com",
	//	From:    "gotest@astergo.in",
	//	Subject: "Test MSG from Go Lang Application",
	//	Content: "",
	//}
	//
	//app.MailChan <- msg
	//
	////Email Serd
	//subject := "Hello, Golang SMTP!"
	//message := "This is a test email sent using Go's net/smtp package."
	//
	//emailMessage := []byte(
	//	"Subject: " + subject + "\r\n" +
	//		"To: " + recipientEmail + "\r\n" +
	//		"From: " + smtp_usr + "\r\n" +
	//		"\r\n" +
	//		message,
	//)
	//
	//auth := smtp.PlainAuth("", smtp_usr, smtp_pass, smtp_server)
	//err = smtp.SendMail(smtp_server+":"+smtp_port, auth, smtp_usr, []string{recipientEmail}, emailMessage)
	//if err != nil {
	//	log.Fatal(err)
	//}

	svr := &http.Server{
		Addr:    port,
		Handler: routes(&app),
	}

	log.Printf("Server is running on http://127.0.0.1%s", port)
	err = svr.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func run() (*driver.DB, error) {
	//put in the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})

	mailChan := make(chan models.EmailData)
	app.MailChan = mailChan

	//change this to true when we are in production
	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "Error:\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction //set false because we use http instead of https
	app.Session = session

	//connect to database
	log.Println("Connecting to database ...")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings_db user=postgres password=000000")
	if err != nil {
		log.Fatal("Cannot connect to database, Dying ...")
	}
	log.Println("Connected to database !")

	tc, err := render.GetTemplatesCache()
	if err != nil {
		log.Fatal("Cannot create template cache!")
		return nil, err
	}
	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
