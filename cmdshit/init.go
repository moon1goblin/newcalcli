package cmdshit

import (
	"calcli/dbshit"
	"database/sql"
	"fmt"

	"context"

	"github.com/urfave/cli/v3"
	_ "modernc.org/sqlite"
)

var Cmd_init *cli.Command = &cli.Command{
	Name: "init",
	Action: initAction,
}

func initAction(ctx context.Context, cmd *cli.Command) error {
	// take the db_ptr out of the context (again idk wtf that is)
	db_ptr := ctx.Value("db_ptr").(*sql.DB)

	if _, err := db_ptr.Exec(
		`
		CREATE TABLE IF NOT EXISTS main(
			event_id INTEGER PRIMARY KEY
			, event_name TEXT NOT NULL
			, begin_datetime TEXT NOT NULL
			, end_datetime TEXT
			, event_type INTEGER
		);
		`,
	); err != nil {
		return fmt.Errorf("initAction: failed to create db: %w: %w", dbshit.ErrSqlite, err)
	}

	// create a sorted view for the table
	// its used in dbshit.GetEventsInRange or something
	if _, err := db_ptr.Exec(
		`
		CREATE VIEW IF NOT EXISTS sorted_view AS 
		SELECT * 
		FROM main 
		ORDER BY 
		datetime(begin_datetime) 
		ASC;
		`,
	); err != nil {
		return fmt.Errorf("initAction: failed to create a sorted view in db: %w: %w", dbshit.ErrSqlite, err)
	}

	return nil
}
