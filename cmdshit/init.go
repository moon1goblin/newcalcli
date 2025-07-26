package cmdshit

import (
	"calcli/event"
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
		return fmt.Errorf("initAction: failed to create db: %w: %w", event.ErrSqlite, err)
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
		return fmt.Errorf("initAction: failed to create a sorted view in db: %w: %w", event.ErrSqlite, err)
	}

	return nil
}
