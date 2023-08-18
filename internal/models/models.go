package models

import "strings"

type Reservation struct {
	FirstName string
	LastName  string
	Email     string
	Phone     string
}

func (r Reservation) Name() string {
	return strings.ToUpper(r.FirstName + " " + r.LastName)
}
