package main

import (
	"calcli/cmdshit"

	"database/sql"
	"context"
	"log"
	"os"

	_ "modernc.org/sqlite"
	"github.com/urfave/cli/v3"
)

func main() {
	cmds := &cli.Command{
		// keep in mind idk wtf context is
		// i just know i can use it to pass my data to subcommands
		Before: func(ctx context.Context, _ *cli.Command) (context.Context, error) {
			// connect to sqlite instance
			db_ptr, err := sql.Open("sqlite", "db")
			ctx = context.WithValue(ctx, "db_ptr", db_ptr)
			return ctx, err
		},

		After: func(ctx context.Context, _ *cli.Command) error {
			// take the db_ptr out of the context (again idk wtf that is)
			db_ptr := ctx.Value("db_ptr").(*sql.DB)

			// must close db for changes to occur, but im not even sure about that
			err := db_ptr.Close()
			return err
		},

		Commands: []*cli.Command{
			cmdshit.Cmd_new,
			cmdshit.Cmd_ls,
			cmdshit.Cmd_init,
			cmdshit.Cmd_rm,
		},
	}

	// TODO: prettify error messages, they look mad ugly

	if err := cmds.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
