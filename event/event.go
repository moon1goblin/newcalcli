package event

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

type EventType int
const (
	// for testing
	NilEvent EventType = iota
	FullDayEvent
	InstantEvent
	WithDurationEvent
)

type Event struct {
	// sometimes i use this sometimes i dont idk
	Id int
	Name string
	Begin_time time.Time
	// TODO: find an optinal<T> package or something
	End_time sql.NullTime
	Type EventType
}

var ErrSqlite = errors.New("Sqlite error")

// insert event into db
func (event Event) Push(db_ptr *sql.DB) error {
	_, err := db_ptr.Exec(
		`
		INSERT INTO main(
			event_name
			, begin_datetime
			, end_datetime
			, event_type
		) VALUES 
		(?, ?, ?, ?);
		`,
		event.Name,
		event.Begin_time.Unix(),
		func() *int64 { 
			if event.End_time.Valid {
				// cant take addres of return value
				// long live the garbage collector
				hi := event.End_time.Time.Unix()
				return &hi
			}
			return nil
		}(),
		event.Type,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to push event with name %s and begin_time %s: %w: %w",
			event.Name,
			event.Begin_time.String(),
			ErrSqlite,
			err,
		)
	}
	return err
}

// true if event exists in db, false otherwise duh
// id is not needed
func (event Event) Find(db_ptr *sql.DB) (bool, error) {
	count := sql.NullInt32{}

	if err := db_ptr.QueryRow(
		`
		SELECT COUNT(*) FROM main 
		WHERE event_name=? 
		AND begin_datetime=?
		AND end_datetime=?
		AND event_type=?
		`,
		event.Name,
		event.Begin_time.Unix(),
		func() *int64 {
			if event.End_time.Valid {
				str := event.End_time.Time.Unix()
				return &str
			}
			return nil
		}(),
		event.Type,
	).Scan(&count); err != nil || !count.Valid {
		return false, fmt.Errorf("Event Find error: %w: %w", ErrSqlite, err)
	}
	return count.Int32 > 0, nil
}
