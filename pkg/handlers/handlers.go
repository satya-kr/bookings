package handlers

import (
	"fmt"
	"htmlweb/pkg/config"
	"htmlweb/pkg/models"
	"htmlweb/pkg/render"
	"net"
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

	render.Template(w, "home", &models.TemplateData{})
}

func (m *Repository) AboutPage(w http.ResponseWriter, r *http.Request) {

	sm := make(map[string]string)
	sm["test"] = "Hello World"

	// remoreIP := m.App.Session.Get(r.Context(), "remote_id") or
	remoreIP := m.App.Session.GetString(r.Context(), "remote_id")
	sm["remote_ip"] = remoreIP

	render.Template(w, "about", &models.TemplateData{
		StringMap: sm,
	})
}

func getClientIPv4(r *http.Request) (string, error) {
	remoteAddr := r.RemoteAddr
	ip, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return "", err
	}

	// Check if the IP address is in IPv4 format
	clientIP := net.ParseIP(ip)
	if clientIP == nil || clientIP.To4() == nil {
		return "", fmt.Errorf("invalid IPv4 address")
	}

	return clientIP.String(), nil
}
