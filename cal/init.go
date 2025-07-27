package cal

import (
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
)

// TODO: just fucking make db_ptr a global variable
func InitDB(db_ptr *sql.DB) error {
	// TODO: because were storing time in seconds since Epoch
	// store timezone too so when switching itd stay the same?
	if _, err := db_ptr.Exec(
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
		return fmt.Errorf("initAction: failed to create db: %w: %w", ErrSqlite, err)
	}

	// create a sorted view for the table
	// its used in dbshit.GetEventsInRange or something
		// datetime(begin_datetime)
	if _, err := db_ptr.Exec(
		`
		CREATE VIEW IF NOT EXISTS sorted_view AS 
		SELECT * FROM main ORDER BY begin_datetime ASC;
		`,
	); err != nil {
		return fmt.Errorf("initAction: failed to create a sorted view in db: %w: %w", ErrSqlite, err)
	}
	return nil
}
