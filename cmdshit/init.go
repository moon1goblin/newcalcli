package cmdshit

import (
	"database/sql"

	"context"

	"github.com/urfave/cli/v3"
	_ "modernc.org/sqlite"
)

var Cmd_init *cli.Command = &cli.Command{
	Name: "init",
	Action: initAction,
}

func initAction(ctx context.Context, cmd *cli.Command) error {
	// TODO: dont create a db if it exists already

	// take the db_ptr out of the context (again idk wtf that is)
	db_ptr := ctx.Value("db_ptr").(*sql.DB)

	if _, err := db_ptr.Exec(
		// make columns not null idk?
		`CREATE TABLE main(
			event_id INTEGER PRIMARY KEY
			, event_name TEXT
			, begin_datetime TEXT
		);`,
	); err != nil {
		return err
	}

	// create a sorted view for the table
	// its used in dbshit.GetEventsInRange or something
	if _, err := db_ptr.Exec(`
		CREATE VIEW sorted_view AS 
		SELECT * 
		FROM main 
		ORDER BY 
		datetime(begin_datetime) 
		ASC;
	`); err != nil {
		return err
	}

	return nil
}
