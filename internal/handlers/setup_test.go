package handlers

import (
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/justinas/nosurf"
	"github.com/satya-kr/bookings/internal/config"
	"github.com/satya-kr/bookings/internal/models"
	"github.com/satya-kr/bookings/internal/render"
	"log"
	"net/http"
	"path/filepath"
	"text/template"
	"time"
)

var app config.AppConfig
var session *scs.SessionManager
var pathToTemplate, _ = filepath.Abs("../../templates")
var functions = template.FuncMap{}

func getRoutes() http.Handler {
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

	tc, err := GetTestTemplatesCache()
	if err != nil {
		log.Fatal("Cannot create template cache!")
	}
	app.TemplateCache = tc
	app.UseCache = true

	repo := NewRepo(&app)
	NewHandlers(repo)

	render.NewRenderer(&app)

	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	//if run post url test make comment NoSurf
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", Repo.HomePage)
	mux.Get("/about", Repo.AboutPage)
	mux.Get("/contact", Repo.AboutPage)
	mux.Get("/generals-quarters", Repo.Generals)
	mux.Get("/majors-suites", Repo.Majors)
	mux.Get("/search-Availability", Repo.Availability)
	mux.Post("/search-availability", Repo.PostAvailability)
	mux.Post("/search-availability-ajax", Repo.PostAvailabilityAjax)
	mux.Get("/make-reservation", Repo.Reservation)
	mux.Post("/make-reservation", Repo.PostReservation)
	mux.Get("/reservation-summary", Repo.ReservationSummary)
	mux.Get("/login", Repo.Login)

	fileServer := http.FileServer(http.Dir("../../static"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}

// NoSurf adds CSRF protection to all POST requests
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

// SessionLoad loads and saves the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

// GetTestTemplatesCache collect all templates then merge them with layout
func GetTestTemplatesCache() (map[string]*template.Template, error) {

	//get the Template Cache from app congig
	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplate))
	if err != nil {
		return myCache, err
	}

	matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplate))
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplate))
			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = ts
	}

	return myCache, nil
}
