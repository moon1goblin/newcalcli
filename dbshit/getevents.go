package dbshit

import (
	"database/sql"
	"errors"
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

	var(
		events []Event
		begin_dummy_str *string
		end_dummy_str *string
	)

	for rows.Next() {
		new_event := Event{}

		// FIXME: i read the src code comments and apparently
		// fucking apparently you can scan directly into fucking *time.Time????
		// TODO: fucking simplify everything with that knowledge???
		if err := rows.Scan(
			&new_event.Id,
			&new_event.Name,
			&begin_dummy_str,
			&end_dummy_str,
			&new_event.Type,
		); err != nil {
			return &events, fmt.Errorf("GetEventsInRange error while scanning rows: %w: %w", ErrSqlite, err)
		}

		if begin_timestr, err := TimeStrFromStr(begin_dummy_str); err != nil {
			return nil, fmt.Errorf("GetEventsInRange error while scanning rows: %w", err)
		} else {
			new_event.Begin_time = *begin_timestr
		}
		if end_timestr, err := TimeStrFromStr(end_dummy_str); !errors.Is(err, ErrNullString) && err != nil {
			return nil, fmt.Errorf("GetEventsInRange error while scanning rows: %w", err)
		} else {
			new_event.End_time = end_timestr
		}

		// i couldnt figure out how to get the row count
		// just allocate enough events right away
		// so append it is
		events = append(events, new_event)
	}

	return &events, nil
}
