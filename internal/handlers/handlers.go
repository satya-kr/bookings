package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/satya-kr/bookings/internal/config"
	"github.com/satya-kr/bookings/internal/driver"
	"github.com/satya-kr/bookings/internal/forms"
	"github.com/satya-kr/bookings/internal/helpers"
	"github.com/satya-kr/bookings/internal/models"
	"github.com/satya-kr/bookings/internal/render"
	"github.com/satya-kr/bookings/internal/repository"
	"github.com/satya-kr/bookings/internal/repository/dbrepo"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

var Repo *Repository

func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
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

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("can not get reservation from session"))
		return
	}

	room, err := m.DB.GetRoomByID(res.RoomID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res.Room.RoomName = room.RoomName
	m.App.Session.Put(r.Context(), "reservation", res)

	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	strMap := make(map[string]string)
	strMap["start_date"] = sd
	strMap["end_date"] = ed

	data := make(map[string]interface{})
	data["reservation"] = res

	render.Template(w, r, "make-reservation", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: strMap,
	})
}

func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("can not get reservation from session"))
		return
	}

	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	fmt.Printf("%v", r.PostForm)

	//sd := r.Form.Get("start_date")
	//ed := r.Form.Get("end_date")

	//2023-08-29 -- 01/02 03:04:05PM '06 -0700
	layout := "2006-01-02"
	//
	//startDate, err := time.Parse(layout, sd)
	//if err != nil {
	//	helpers.ServerError(w, err)
	//	return
	//}
	//endDate, err := time.Parse(layout, ed)
	//if err != nil {
	//	helpers.ServerError(w, err)
	//	return
	//}
	//
	//roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	//if err != nil {
	//	helpers.ServerError(w, err)
	//	return
	//}

	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Email = r.Form.Get("email")
	reservation.Phone = r.Form.Get("phone")

	//reservation := models.Reservation{
	//	FirstName: r.Form.Get("first_name"),
	//	LastName:  r.Form.Get("last_name"),
	//	Email:     r.Form.Get("email"),
	//	Phone:     r.Form.Get("phone"),
	//	StartDate: startDate,
	//	EndDate:   endDate,
	//	RoomID:    roomID,
	//}

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

	var newReservationID int
	newReservationID, err = m.DB.StoreReservation(reservation)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	restriction := models.RoomRestriction{
		StartDate:     reservation.StartDate,
		EndDate:       reservation.EndDate,
		RoomID:        reservation.RoomID,
		ReservationID: newReservationID,
		RestrictionID: 1,
	}

	err = m.DB.StoreRoomRestriction(restriction)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	message := fmt.Sprintf(`<strong>Reservation Confermation</strong><br><br>
	Dear: %s,<br>
	This is confirm your reservation from %s to %s!
	<br><br>
	<i>Thank you for business with us.</i>
	`, reservation.FirstName, reservation.StartDate.Format(layout), reservation.EndDate.Format(layout))

	msg := models.EmailData{
		To:       reservation.Email,
		From:     "gotest@astergo.in",
		Subject:  "Reservation is confirmed!",
		Content:  message,
		Template: "basic.html",
	}

	m.App.MailChan <- msg

	message = fmt.Sprintf(`<strong>
	Reservation Notification</strong><br><br>
	A reservation has been made for %s from %s to %s!`, reservation.FirstName, reservation.StartDate.Format(layout), reservation.EndDate.Format(layout))

	msg = models.EmailData{
		To:      "satyajit.kr.prajapati@gmail.com",
		From:    "gotest@astergo.in",
		Subject: "Reservation Notification",
		Content: message,
	}

	m.App.MailChan <- msg

	m.App.Session.Put(r.Context(), "reservation", reservation)
	m.App.Session.Put(r.Context(), "flash", "Your reservation is successful.")
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		//fmt.Println("Failed to get item from session")
		m.App.ErrorLog.Println("Can't get item from session")
		m.App.Session.Put(r.Context(), "error", "can't get reservation from session!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	m.App.Session.Remove(r.Context(), "reservation")
	data := make(map[string]interface{})
	data["reservation"] = reservation

	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")
	strMap := make(map[string]string)
	strMap["start_date"] = sd
	strMap["end_date"] = ed

	render.Template(w, r, "reservation-summary", &models.TemplateData{
		Data:      data,
		StringMap: strMap,
	})
}

func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability", &models.TemplateData{})
}
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {

	start := r.Form.Get("start")
	end := r.Form.Get("end")

	//2023-08-29 -- 01/02 03:04:05PM '06 -0700
	layout := "2006-01-02"

	startDate, err := time.Parse(layout, start)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(layout, end)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	for _, i := range rooms {
		m.App.InfoLog.Println("Room:", i.ID, i.RoomName)
	}

	if len(rooms) == 0 {
		//m.App.InfoLog.Println("Not Available!")
		m.App.Session.Put(r.Context(), "error", "Not Available!")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	//w.Write([]byte(fmt.Sprintf("start date is %s and end date is %s", start, end)))

	data := make(map[string]interface{})
	data["rooms"] = rooms

	resv := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}

	m.App.Session.Put(r.Context(), "reservation", resv)

	render.Template(w, r, "choose-room", &models.TemplateData{
		Data: data,
	})
}

type JsonResponse struct {
	OK        bool   `json:"ok"`
	Message   string `json:"message"`
	RoomID    string `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

func (m *Repository) PostAvailabilityAjax(w http.ResponseWriter, r *http.Request) {

	sd := r.Form.Get("start")
	ed := r.Form.Get("end")

	layout := "2006-01-02"
	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)

	roomID, _ := strconv.Atoi(r.Form.Get("room_id"))

	available, _ := m.DB.SearchAvailabilityByDates(startDate, endDate, roomID)

	res := JsonResponse{
		OK:        available,
		Message:   "Room is Available!",
		StartDate: sd,
		EndDate:   ed,
		RoomID:    strconv.Itoa(roomID),
	}

	out, err := json.MarshalIndent(res, "", " ")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (m *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	//m.App.Session.Get(r.Context(), "reservation")
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, err)
		return
	}

	res.RoomID = roomID
	m.App.Session.Put(r.Context(), "reservation", res)
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

func (m *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {
	roomID, _ := strconv.Atoi(r.URL.Query().Get("id"))
	sd, ed := r.URL.Query().Get("s"), r.URL.Query().Get("e")

	layout := "2006-01-02"
	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)

	var res models.Reservation

	room, err := m.DB.GetRoomByID(roomID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res.RoomID = roomID
	res.Room.RoomName = room.RoomName
	res.StartDate = startDate
	res.EndDate = endDate

	m.App.Session.Put(r.Context(), "reservation", res)
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

func (m *Repository) Login(w http.ResponseWriter, r *http.Request) {

	if helpers.IsAuth(r) {
		m.App.Session.Put(r.Context(), "flash", "Logged in successfully!")
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
	}

	render.Template(w, r, "login", &models.TemplateData{
		Form: forms.New(nil),
	})
}

func (m *Repository) PostLogin(w http.ResponseWriter, r *http.Request) {

	//so every time session will renew the token for best possible security
	_ = m.App.Session.RenewToken(r.Context())
	err := r.ParseForm()
	if err != nil {
		log.Print(err)
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	form := forms.New(r.PostForm)
	form.Required("email", "password")
	if !form.Valid() {

		render.Template(w, r, "login", &models.TemplateData{
			Form: form,
		})

		return
	}

	id, _, err := m.DB.Authenticate(email, password)
	if err != nil {
		log.Print(err)
		m.App.Session.Put(r.Context(), "error", "Invalid login credentials")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	m.App.Session.Put(r.Context(), "user_id", id)
	m.App.Session.Put(r.Context(), "flash", "Logged in successfully!")
	http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
	//render.Template(w, r, "login", &models.TemplateData{
	//	Form: forms.New(nil),
	//})
}

func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (m *Repository) Dashboard(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "dashboard", &models.TemplateData{})
}

func (m *Repository) AdminNewReservations(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin-new-reservations", &models.TemplateData{})
}

func (m *Repository) AdminAllReservations(w http.ResponseWriter, r *http.Request) {

	reservations, err := m.DB.AllReservations()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["reservations"] = reservations

	render.Template(w, r, "admin-all-reservations", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) AdminReservationCalendar(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin-reservations-calendar", &models.TemplateData{})
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
