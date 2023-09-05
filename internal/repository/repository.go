package repository

import (
	"github.com/satya-kr/bookings/internal/models"
	"time"
)

type DatabaseRepo interface {
	AllUsers() bool

	StoreReservation(res models.Reservation) (int, error)
	StoreRoomRestriction(r models.RoomRestriction) error
	SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error)
	SearchAvailabilityByDates(start, end time.Time, roomID int) (bool, error)
	GetRoomByID(id int) (models.Room, error)
	GetUserByID(id int) (models.User, error)
	UpdateUser(u models.User) error
	Authenticate(email, testPassword string) (int, string, error)
	AllReservations() ([]models.Reservation, error)
	AllNewReservations() ([]models.Reservation, error)
	GetReservationByID(id int) (models.Reservation, error)
	UpdateReservation(u models.Reservation) error
	DeleteReservation(id int) error
	UpdateProcessedForReservation(id, processed int) error
	AllRooms() ([]models.Room, error)
}
