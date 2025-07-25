package dbshit

import(
	"database/sql"

	_ "modernc.org/sqlite"

)

// i wanted to have all sql be in this folder because i got confused
func CreateDB(db_ptr *sql.DB) error {
	_, err := db_ptr.Exec(
		`
		CREATE TABLE IF NOT EXISTS main(
			event_id INTEGER PRIMARY KEY
			, event_name TEXT NOT NULL
			, begin_datetime TEXT NOT NULL
		);
		`,
	)
	return err
}

// create a sorted view for the table
// its used in dbshit.GetEventsInRange or something
func CreateSortedView(db_ptr *sql.DB) error {
	_, err := db_ptr.Exec(
		`
		CREATE VIEW IF NOT EXISTS sorted_view AS 
		SELECT * 
		FROM main 
		ORDER BY 
		datetime(begin_datetime) 
		ASC;
		`,
	)
	return err
}

func DeleteById(id int, db_ptr *sql.DB) error {
	_, err := db_ptr.Exec(
		`
		DELETE FROM main 
		WHERE event_id=?;
		`,
		id,
	)
	return err
}
