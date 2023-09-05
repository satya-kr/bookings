package models

import (
	"strings"
	"time"
)

//type Reservation struct {
//	FirstName string
//	LastName  string
//	Email     string
//	Phone     string
//}

func (r Reservation) Name() string {
	return strings.ToUpper(r.FirstName + " " + r.LastName)
}

type User struct {
	ID          int
	FirstName   string
	LastName    string
	Email       string
	Password    string
	AccessLevel int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Room struct {
	ID        int
	RoomName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Restriction struct {
	ID              int
	RestrictionName string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type Reservation struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	Phone     string
	StartDate time.Time
	EndDate   time.Time
	RoomID    int
	CreatedAt time.Time
	UpdatedAt time.Time
	Room      Room
}

type RoomRestriction struct {
	ID            int
	StartDate     time.Time
	EndDate       time.Time
	ReservationID int
	RestrictionID int
	RoomID        int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Room          Room
	Reservation   Reservation
	Restriction   Restriction
}

type EmailData struct {
	To      string
	From    string
	Subject string
	//Content template.HTML
	Content  string
	Template string
}
