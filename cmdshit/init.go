package cmdshit

import (
	"calcli/dbshit"
	"database/sql"

	"context"

	"github.com/urfave/cli/v3"
	_ "modernc.org/sqlite"
)

var Cmd_init *cli.Command = &cli.Command{
	Name: "init",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		// TODO: dont create db if it exists already

		// take the db_ptr out of the context (again idk wtf that is)
		db_ptr := ctx.Value("db_ptr").(*sql.DB)

		// rvo? never heard of her
		// this is go anyways who cares
		return dbshit.CreateDb(db_ptr)
	},
}

