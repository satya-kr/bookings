package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/satya-kr/bookings/internal/config"
	"github.com/satya-kr/bookings/internal/forms"
	"github.com/satya-kr/bookings/internal/models"
	"github.com/satya-kr/bookings/internal/render"
	"net/http"
)

type Repository struct {
	App *config.AppConfig
}

var Repo *Repository

func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) HomePage(w http.ResponseWriter, r *http.Request) {

	// remoteIP, err := getClientIPv4(r)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	remoteIP := r.RemoteAddr

	fmt.Println(remoteIP)

	m.App.Session.Put(r.Context(), "remote_id", remoteIP)

	render.Template(w, r, "home", &models.TemplateData{})
}

func (m *Repository) AboutPage(w http.ResponseWriter, r *http.Request) {

	sm := make(map[string]string)
	sm["test"] = "Hello World"

	// remoreIP := m.App.Session.Get(r.Context(), "remote_id") or
	remoreIP := m.App.Session.GetString(r.Context(), "remote_id")
	sm["remote_ip"] = remoreIP

	render.Template(w, r, "about", &models.TemplateData{
		StringMap: sm,
	})
}

func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "generals", &models.TemplateData{})
}

func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "majors", &models.TemplateData{})
}

func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {

	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	render.Template(w, r, "make-reservation", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
		return
	}

	//fmt.Printf("%v", r.PostForm)
	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
	}

	form := forms.New(r.PostForm)

	//form.Has("first_name", r)
	//form.Has("last_name", r)
	//form.Has("email", r)
	//form.Has("phone", r)

	form.Required("first_name", "last_name", "email", "phone")
	form.MinLength("first_name", 5, r)
	form.MinLength("last_name", 5, r)
	form.MinLength("email", 5, r)
	form.MinLength("phone", 10, r)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.Template(w, r, "make-reservation", &models.TemplateData{
			Form: form,
			Data: data,
		})

		return
	}
	m.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		fmt.Println("Faild to get item from session")
		return
	}
	data := make(map[string]interface{})
	data["reservation"] = reservation
	render.Template(w, r, "reservation-summary", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) Availablity(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availablity", &models.TemplateData{})
}
func (m *Repository) PostAvailablity(w http.ResponseWriter, r *http.Request) {

	start := r.Form.Get("start")
	end := r.Form.Get("end")

	w.Write([]byte(fmt.Sprintf("start date is %s and end date is %s", start, end)))
}

type JsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

func (m *Repository) PostAvailablityAjax(w http.ResponseWriter, r *http.Request) {

	res := JsonResponse{
		OK:      true,
		Message: "Available!",
	}

	out, err := json.MarshalIndent(res, "", " ")
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (m *Repository) Login(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "login", &models.TemplateData{})
}

// func getClientIPv4(r *http.Request) (string, error) {
// 	remoteAddr := r.RemoteAddr
// 	ip, _, err := net.SplitHostPort(remoteAddr)
// 	if err != nil {
// 		return "", err
// 	}

// 	// Check if the IP address is in IPv4 format
// 	clientIP := net.ParseIP(ip)
// 	if clientIP == nil || clientIP.To4() == nil {
// 		return "", fmt.Errorf("invalid IPv4 address")
// 	}

// 	return clientIP.String(), nil
// }
