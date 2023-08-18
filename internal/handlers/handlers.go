package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/satya-kr/bookings/pkg/config"
	"github.com/satya-kr/bookings/pkg/models"
	"github.com/satya-kr/bookings/pkg/render"
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
	render.Template(w, r, "make-reservation", &models.TemplateData{})
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
