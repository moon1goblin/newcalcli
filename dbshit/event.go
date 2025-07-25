package dbshit

import(
	"database/sql"
	"errors"
	"time"

	_ "modernc.org/sqlite"
)

type Event struct {
	// sometimes i use this sometimes i dont idk
	Id int
	Name string
	Begin_time time.Time
}

// insert event into db
func (event Event) Push(db_ptr *sql.DB) error {
	_, err := db_ptr.Exec(
		`
		INSERT INTO main 
		(event_name, begin_datetime) 
		VALUES 
		(?, ?);
		`,
		event.Name,
		event.Begin_time.Format(time.DateTime),
	)
	return err
}

// true if event exists in db, false otherwise duh
// id is not needed
func (event Event) Find(db_ptr *sql.DB) (bool, error) {
	// couldnt find how to return just the count but whatever
	rows, err := db_ptr.Query(
		`
		SELECT * FROM main 
		WHERE event_name=? 
		AND begin_datetime=?
		`,
		event.Name,
		event.Begin_time.Format(time.DateTime),
	)
	if err != nil {
		return false, errors.New("(Event) Find error: " + err.Error())
	}

	found := false
	for rows.Next() {
		found = true
	}
	return found, nil
}
