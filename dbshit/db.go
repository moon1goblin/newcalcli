package dbshit

import(
	// "time"
	"database/sql"
	"strings"

	_ "modernc.org/sqlite"
)

func CreateDb(db_ptr *sql.DB) error {
	_, err := db_ptr.Exec(
		// make columns not null idk?

		`CREATE TABLE main(
			event_name TEXT
			, begin_datetime TEXT
		);`,
	)
	return err
}

// FIXME: abomination of a function but whatever
func GetFuckingEverything(db_ptr *sql.DB) (string, error) {
	rows, err := db_ptr.Query(
		`SELECT * FROM main`,
	)
	if err != nil {
		return "", err
	}

	var builder strings.Builder

	data := [2]string{}
	for rows.Next() {
		// how do i pass by address while unpacking? i tried &data...
		if err := rows.Scan(&data[0], &data[1]); err != nil {
			return builder.String(), err
		}
		for i := range 2 {
			builder.WriteString(data[i])
			builder.WriteString(" ")
		}
		builder.WriteString("\n")
	}
	
	return builder.String(), nil
}

type Event struct {
	Name string
	Begin_time string
	// Begin_time time.Time
	// end_time time.Time
	// TODO: write a time parser and store time instead of just string
}

func (event Event) Push(db_ptr *sql.DB) error {
	_, err := db_ptr.Exec(
		`INSERT INTO main
		(event_name, begin_datetime)
		VALUES 
		(?, ?)`, 
		event.Name,
		event.Begin_time,
		// event.Begin_time.Format(time.DateTime),
	)
	return err
}
