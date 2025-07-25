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

	if err := dbshit.CreateDB(db_ptr); err != nil {
		return fmt.Errorf("initAction: failed to create db: %w", err)
	}
	// create a sorted view for the table
	// its used in dbshit.GetEventsInRange or something
	if err := dbshit.CreateSortedView(db_ptr); err != nil {
		return fmt.Errorf("initAction: failed to create a sorted view in db: %w", err)
	}

	return nil
}
