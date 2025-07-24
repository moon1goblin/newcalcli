package cmdshit

import (
	"fmt"
	"context"
	"calcli/dbshit"
	"database/sql"

	"github.com/urfave/cli/v3"
	_ "modernc.org/sqlite"
)

var Cmd_ls *cli.Command = &cli.Command{
	Name: "ls",
	Action: func(ctx context.Context, _ *cli.Command) error {
		db_ptr := ctx.Value("db_ptr").(*sql.DB)
		all_events_str, err := dbshit.GetFuckingEverything(db_ptr)
		if err != nil {
			return err
		}
		// FIXME: flags think their values are only the first word :(
		// so like "-n im batman"'s value is im
		fmt.Println(all_events_str)
		return nil
	},
}
