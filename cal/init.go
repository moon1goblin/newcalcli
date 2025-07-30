package cal

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

var Db_ptr_g *sql.DB

func InitDB() error {
	var err error
	if Db_ptr_g, err = sql.Open("sqlite", "db"); err != nil {
		return fmt.Errorf("InitDB: failed to connect to db: %w: %w", ErrSqlite, err)
	}

	if _, err := Db_ptr_g.Exec(
		// TODO: because were storing time in seconds since Epoch
		// store timezone too so when switching itd stay the same?
		`
		CREATE TABLE IF NOT EXISTS main(
			event_id INTEGER PRIMARY KEY
			, event_name TEXT NOT NULL
			, begin_datetime INTEGER NOT NULL
			, end_datetime INTEGER
			, event_type INTEGER NOT NULL
		);
		`,
	); err != nil {
		return fmt.Errorf("InitDB: failed to create db: %w: %w", ErrSqlite, err)
	}

	// create a sorted view for the table
	// its used in dbshit.GetEventsInRange or something
		// datetime(begin_datetime)
	if _, err := Db_ptr_g.Exec(
		`
		CREATE VIEW IF NOT EXISTS sorted_view AS 
		SELECT * FROM main ORDER BY begin_datetime ASC;
		`,
	); err != nil {
		return fmt.Errorf("InitDB: failed to create a sorted view in db: %w: %w", ErrSqlite, err)
	}
	return nil
}
