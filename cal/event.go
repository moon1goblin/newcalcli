package cal

import (
	"database/sql"
	"strings"
	"errors"
	"time"
	"fmt"

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
	Id         int
	Name       string
	Begin_time time.Time
	End_time sql.NullTime
	Type     EventType
}

var ErrEventAlreadyExists = errors.New("event already exists")

func EventCreate(name_str, begin_datetime_str, end_datetime_str string) (*Event, error) {
	my_event, err := ProcessDates(name_str, begin_datetime_str, end_datetime_str)
	if err != nil {
		return nil, fmt.Errorf("EventCreate error: %w", err)
	}
	if found, err := my_event.Find(); err != nil {
		return nil, fmt.Errorf("EventCreate error: %w", err)
	} else if found {
		return nil, ErrEventAlreadyExists
	}
	return my_event, nil
}

func (event Event) String() string {
	return event.string(false)
}

func (event Event) StringWithDate() string {
	return event.string(true)
}

func (event Event) string(withdate bool) string {
	var (
		builder      strings.Builder
		begin_format string
		end_format   string
	)

	if withdate {
		begin_format = "01.02.2006 15:04"
	} else {
		begin_format = "15:04"
	}

	isSameDate := func(lhs, rhs time.Time) bool {
		ly, lm, ld := lhs.Date()
		ry, rm, rd := rhs.Date()
		return ld == rd && lm == rm && ly == ry
	}
	if event.End_time.Valid && isSameDate(event.Begin_time, event.End_time.Time) {
		end_format = "15:04"
	} else {
		end_format = begin_format
	}

	switch event.Type {
	case FullDayEvent:
		if withdate {
			builder.WriteString(event.Begin_time.Format("01.02.2006"))
			builder.WriteString(" ")
		}
	case InstantEvent:
		builder.WriteString(event.Begin_time.Format(begin_format))
		builder.WriteString(" ")
	case WithDurationEvent:
		builder.WriteString(event.Begin_time.Format(begin_format))
		builder.WriteString("-")
		builder.WriteString(event.End_time.Time.Format(end_format))
		builder.WriteString(" ")
	}

	builder.WriteString(event.Name)

	return builder.String()
}

var ErrSqlite = errors.New("Sqlite error")

// insert event into db
func (event Event) Push() error {
	_, err := Db_ptr_g.Exec(
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
// FIXME: this broke for some reason
func (event Event) Find() (bool, error) {
	var count int
	if err := Db_ptr_g.QueryRow(
		`
		SELECT COUNT(*) FROM main
		WHERE event_name=?
		AND begin_datetime=?
		AND end_datetime=?
		AND event_type=?;
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
	).Scan(&count); err != nil {
		return false, fmt.Errorf("Event Find error: %w: %w", ErrSqlite, err)
	}
	return count > 0, nil
}

// returns sorted list of events
// includes begin and doesnt include end, so like a ray [)
// begin and end are optional
func GetEventsInRange(begin sql.NullTime, end sql.NullTime) (*[]Event, error) {
	// TODO: count how many rows it returned and allocate the events array accordingly
	var (
		rows *sql.Rows
		err  error
	)
	// checking if we even have a begin and an end
	// there must be a better way... but im too dumb to see it
	if !begin.Valid && !end.Valid {
		if rows, err = Db_ptr_g.Query(
			`
			SELECT * 
			FROM sorted_view;
			`,
		); err != nil {
			return nil, fmt.Errorf("GetEventsInRange error: %w: %w", ErrSqlite, err)
		}
	} else if !begin.Valid {
		if rows, err = Db_ptr_g.Query(
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
		if rows, err = Db_ptr_g.Query(
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
		if rows, err = Db_ptr_g.Query(
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

	var (
		events      []Event
		begin_dummy int64
		end_dummy   sql.NullInt64
	)

	new_event := Event{}
	for rows.Next() {
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
		if end_dummy.Valid {
			new_event.End_time = sql.NullTime{Time: time.Unix(end_dummy.Int64, 0), Valid: true}
		} else {
			new_event.End_time = sql.NullTime{}
		}

		// i couldnt figure out how to get the row count to allocate enough events right away
		// so append it is
		events = append(events, new_event)
	}

	return &events, nil
}
