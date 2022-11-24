package dbrepo

import (
	"context"
	"time"

	"github.com/DaniilShd/booking/internal/models"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

// InsertReservation insert reservation into DB
func (m *postgresDBRepo) InsertReservation(res models.Reservation) (int, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	stmt := `insert into reservation 
	(first_name, last_name, email, phone, start_date, end_date, room_id, created_at, updated_at) 
	values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`

	var newID int

	err := m.DB.QueryRowContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.Startdate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		return 0, err
	}
	return newID, nil
}

func (m *postgresDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	stmt := `insert into room_restriction 
	(start_date, end_date, room_id,reservation_id, restriction_id, created_at, updated_at) 
	values ($1, $2, $3, $4, $5, $6, $7)`

	_, err := m.DB.ExecContext(ctx, stmt,

		r.Startdate,
		r.EndDate,
		r.RoomID,
		r.ReservationID,
		r.RestrictionID,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return err
	}

	return nil
}

// SearchAvailailityByDates return true if availability exists for room_id
func (m *postgresDBRepo) SearchAvailailityByDatesBuRoomID(start, end time.Time, roomID int) (bool, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	var numRows int
	query := `select
	count(id)
	from 
	room_restriction
	where 
	room_id= $1
	$2 < end_date and $3 > start_date;
	`
	row := m.DB.QueryRowContext(ctx, query, roomID, start, end)
	err := row.Scan(&numRows)
	if err != nil {
		return false, err
	}

	if numRows == 0 {
		return true, nil
	}
	return false, nil
}

// AddNotNullToReservationIDForRestriction return a slice rooms, if any givan date range
func (m *postgresDBRepo) AddNotNullToReservationIDForRestriction(start, end time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var rooms []models.Room

	query := `
	select
	r,id, r.room_name
		from 
	rooms r
		where 
	r.id not in (select rr.room_id from room_restriction rr where $1 < rr.end_date and $2 > rr.start_date)`

	rows, err := m.DB.QueryContext(ctx, query, start, end)
	if err != nil {
		return nil, err
	}

	var room models.Room
	for rows.Next() {
		err := rows.Scan(&room.ID, &room.RoomName)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return rooms, nil
}
