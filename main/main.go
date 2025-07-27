package main

import (
	"calcli/cal"
	"context"
	"database/sql"
	"fmt"

	// "log"
	"os"

	"github.com/urfave/cli/v3"

	_ "modernc.org/sqlite"
)

// TODO: just fucking call init here and set db_ptr
// so like, make it global in the packages and just set it here

func main() {
	cmds := &cli.Command{
		// keep in mind idk wtf context is
		// i just know i can use it to pass my data to subcommands
		Before: func(ctx context.Context, _ *cli.Command) (context.Context, error) {
			// connect to sqlite instance
			db_ptr, err := sql.Open("sqlite", "db")
			if err != nil {
				return ctx, err
			}
			ctx = context.WithValue(ctx, "db_ptr", db_ptr)

			if err := cal.InitDB(db_ptr); err != nil {
				return ctx, err
			}
			return ctx, nil
		},

		After: func(ctx context.Context, _ *cli.Command) error {
			// take the db_ptr out of the context (again idk wtf that is)
			db_ptr := ctx.Value("db_ptr").(*sql.DB)

			// must close db for changes to occur, but im not even sure about that
			err := db_ptr.Close()
			return err
		},

		Commands: []*cli.Command{
			Command_new,
			Command_ls,
			Command_rm,
		},
	}

	if err := cmds.Run(context.Background(), os.Args); err != nil {
		// log.Fatal(err)
		fmt.Println(err)
	}
}
