package main

import (
	"encoding/gob"
	"fmt"
	"github.com/satya-kr/bookings/internal/config"
	"github.com/satya-kr/bookings/internal/handlers"
	"github.com/satya-kr/bookings/internal/models"
	"github.com/satya-kr/bookings/internal/render"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
)

const port = ":8080"

var app config.AppConfig
var session *scs.SessionManager

func main() {
	//put in the session
	gob.Register(models.Reservation{})

	//change this to true when we are in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction //set false because we use http insted of https
	app.Session = session

	tc, err := render.GetTemplatesCache()
	if err != nil {
		log.Fatal("Cannot create template cache!")
	}
	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	// http.HandleFunc("/", handlers.Repo.HomePage)
	// http.HandleFunc("/about", handlers.Repo.AboutPage)

	fmt.Printf("Server is running on port %s ...", port)

	// err = http.ListenAndServe(port, nil)
	// if err != nil {
	// 	panic(err)
	// }

	svr := &http.Server{
		Addr:    port,
		Handler: routes(&app),
	}

	err = svr.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
