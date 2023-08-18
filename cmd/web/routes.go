package main

import (
	"github.com/satya-kr/bookings/internal/config"
	"github.com/satya-kr/bookings/internal/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func routes(app *config.AppConfig) http.Handler {

	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)

	// mux.Use(WriteToConsole)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.HomePage)
	mux.Get("/about", handlers.Repo.AboutPage)
	mux.Get("/contact", handlers.Repo.AboutPage)
	mux.Get("/generals-quarters", handlers.Repo.Generals)
	mux.Get("/majors-suites", handlers.Repo.Majors)
	mux.Get("/search-availablity", handlers.Repo.Availablity)
	mux.Post("/search-availability", handlers.Repo.PostAvailablity)
	mux.Post("/search-availability-ajax", handlers.Repo.PostAvailablityAjax)
	mux.Get("/make-reservation", handlers.Repo.Reservation)
	mux.Post("/make-reservation", handlers.Repo.PostReservation)
	mux.Get("/reservation-summary", handlers.Repo.ReservationSummary)
	mux.Get("/login", handlers.Repo.Login)

	fileServer := http.FileServer(http.Dir("../../static"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
