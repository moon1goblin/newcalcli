package event

import (
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

// returns sorted list of events
// includes begin and doesnt include end, so like a ray [)
// begin and end are optional
func GetEventsInRange(begin sql.NullTime, end sql.NullTime, db_ptr *sql.DB) (*[]Event, error) {
	// TODO: count how many rows it returned and allocate the events array accordingly
	var (
		rows *sql.Rows
		err error
	)
	// checking if we even have a begin and an end
	// there must be a better way... but im too dumb to see it
	if !begin.Valid && !end.Valid {
		if rows, err = db_ptr.Query(
			`
			SELECT * 
			FROM sorted_view;
			`, 
		); err != nil {
			return nil, fmt.Errorf("GetEventsInRange error: %w: %w", ErrSqlite, err)
		}
	} else if !begin.Valid {
		if rows, err = db_ptr.Query(
			`
			SELECT * 
			FROM sorted_view 
			WHERE begin_datetime >= ?;
			`, 
			begin.Time.Unix(),
		); err != nil {
			return nil, fmt.Errorf("GetEventsInRange error: %w: %w", ErrSqlite, err)
		}
	} else if !end.Valid {
		if rows, err = db_ptr.Query(
			`
			SELECT * 
			FROM sorted_view 
			WHERE begin_datetime < ?;
			`, 
			end.Time.Unix(),
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
			begin.Time.Unix(),
			end.Time.Unix(),
		); err != nil {
			return nil, fmt.Errorf("GetEventsInRange error: %w: %w", ErrSqlite, err)
		}
	}

	var(
		events []Event
		begin_dummy int64
		end_dummy sql.NullInt64
	)

	for rows.Next() {
		new_event := Event{}
		if err := rows.Scan(
			&new_event.Id,
			&new_event.Name,
			&begin_dummy,
			&end_dummy,
			&new_event.Type,
		); err != nil {
			return &events, fmt.Errorf("GetEventsInRange error while scanning rows: %w: %w", ErrSqlite, err)
		}
		new_event.Begin_time = time.Unix(begin_dummy, 0)
		new_event.End_time = func() sql.NullTime {
			if end_dummy.Valid {
				return sql.NullTime{Time: time.Unix(end_dummy.Int64, 0), Valid: true}
			}
			return sql.NullTime{}
		}()

		// i couldnt figure out how to get the row count
		// just allocate enough events right away
		// so append it is
		events = append(events, new_event)
	}

	return &events, nil
}
