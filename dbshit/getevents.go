package dbshit

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

// returns sorted list of events
// includes begin and doesnt include end, so like a ray [)
// begin and end are optional
func GetEventsInRange(begin *TimeStr, end *TimeStr, db_ptr *sql.DB) (*[]Event, error) {
	// TODO: count how many rows it returned and allocate the events array accordingly
	var (
		rows *sql.Rows
		err error
	)
	// checking if we even have a begin and an end
	// there must be a better way... but im too dumb to see it
	if begin == nil && end == nil {
		if rows, err = db_ptr.Query(
			`
			SELECT * 
			FROM sorted_view;
			`, 
		); err != nil {
			return nil, fmt.Errorf("GetEventsInRange error: %w: %w", ErrSqlite, err)
		}
	} else if begin != nil {
		if rows, err = db_ptr.Query(
			`
			SELECT * 
			FROM sorted_view 
			WHERE datetime(begin_datetime) >= ?;
			`, 
			begin.String(),
		); err != nil {
			return nil, fmt.Errorf("GetEventsInRange error: %w: %w", ErrSqlite, err)
		}
	} else if end != nil {
		if rows, err = db_ptr.Query(
			`
			SELECT * 
			FROM sorted_view 
			WHERE datetime(begin_datetime) < ?;
			`, 
			end.String(),
		); err != nil {
			return nil, fmt.Errorf("GetEventsInRange error: %w: %w", ErrSqlite, err)
		}
	} else {
		if rows, err = db_ptr.Query(
			`
			SELECT * 
			FROM sorted_view 
			WHERE datetime(begin_datetime) >= ? 
			AND datetime(begin_datetime) < ?;
			`, 
			begin.String(),
			end.String(),
		); err != nil {
			return nil, fmt.Errorf("GetEventsInRange error: %w: %w", ErrSqlite, err)
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

		timestr, err := TimeStrFromStr(begin_dummy_str)
		if err != nil {
			return nil, fmt.Errorf("GetEventsInRange error while scanning rows: %w", err)
		}
		new_event.Begin_time = *timestr

		// i couldnt figure out how to get the row count
		// just allocate enough events right away
		// so append it is
		events = append(events, new_event)
	}

	return &events, nil
}
