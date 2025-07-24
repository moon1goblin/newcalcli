package dbshit

import (
	"database/sql"
	"errors"
	"time"

	_ "modernc.org/sqlite"
)

// returns sorted list of events
// includes begin and doesnt include end, so like a ray [)
// begin and end are optional
func GetEventsInRange(begin *time.Time, end *time.Time, db_ptr *sql.DB) (*[]Event, error) {
	// TODO: count how many rows it returned and allocate the events array accordingly
	var (
		rows *sql.Rows
		err error
	)

	// checking if we even have a begin and an end
	// there must be a better way... but im too dumb to see it
	if begin == nil && end == nil {
		if rows, err = db_ptr.Query(`
			SELECT * 
			FROM sorted_view
			`, 
		); err != nil {
			return nil, nil
		}
	} else if begin != nil {
		if rows, err = db_ptr.Query(`
			SELECT * 
			FROM sorted_view 
			WHERE datetime(begin_datetime) >= ?
			`, 
			begin.Format(time.DateTime),
		); err != nil {
			return nil, nil
		}
	} else if end != nil {
		if rows, err = db_ptr.Query(`
			SELECT * 
			FROM sorted_view 
			WHERE datetime(begin_datetime) < ?
			`, 
			end.Format(time.DateTime),
		); err != nil {
			return nil, nil
		}
	} else {
		if rows, err = db_ptr.Query(`
			SELECT * 
			FROM sorted_view 
			WHERE datetime(begin_datetime) >= ? 
			AND datetime(begin_datetime) < ?
			`, 
			begin.Format(time.DateTime),
			end.Format(time.DateTime),
		); err != nil {
			return nil, nil
		}
	}

	var events []Event
	var begin_dummy_str string

	for rows.Next() {
		new_event := Event{}

		if err := rows.Scan(
			&new_event.Id,
			&new_event.Name,
			&begin_dummy_str,
		); err != nil {
			return &events, err
		}

		new_datetime, err := time.Parse(time.DateTime, begin_dummy_str)
		if err != nil {
			return nil, errors.New("GetEventsInRange error: failed to convert sqlite string to time.Time")
		}
		new_event.Begin_time = new_datetime

		events = append(events, new_event)
	}

	return &events, nil
}

type Event struct {
	Id int
	Name string
	Begin_time time.Time
}


func (event Event) Push(db_ptr *sql.DB) error {
	_, err := db_ptr.Exec(
		`INSERT INTO main 
		(event_name, begin_datetime) 
		VALUES 
		(?, ?)`,
		event.Name,
		event.Begin_time.Format(time.DateTime),
	)
	return err
}

// // FIXME: abomination of a function but whatever
//
// func GetFuckingEverything(db_ptr *sql.DB) (string, error) {
// 	rows, err := db_ptr.Query(
// 		`SELECT * FROM main`,
// 	)
// 	if err != nil {
// 		return "", err
// 	}
//
// 	var builder strings.Builder
//
// 	data := [2]string{}
// 	for rows.Next() {
// 		// how do i pass by address while unpacking? i tried &data...
// 		if err := rows.Scan(&data[0], &data[1]); err != nil {
// 			return builder.String(), err
// 		}
// 		for i := range 2 {
// 			builder.WriteString(data[i])
// 			builder.WriteString(" ")
// 		}
// 		builder.WriteString("\n")
// 	}
//
// 	return builder.String(), nil
// }
