package repository

import (
	"time"

	"github.com/DaniilShd/booking/internal/models"
)

type DatabaseRepo interface {
	AllUsers() bool
	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(r models.RoomRestriction) error
	SearchAvailailityByDatesBuRoomID(start, end time.Time, roomID int) (bool, error)
	AddNotNullToReservationIDForRestriction(start, end time.Time) ([]models.Room, error)
}
