package cmdshit

import (
	"context"
	"database/sql"

	"github.com/urfave/cli/v3"
	_ "modernc.org/sqlite"
)

var Cmd_rm *cli.Command = &cli.Command{
	Name: "rm",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "id",
			Required: true,
		},
		// &cli.StringFlag{
		// 	Name: "name",
		// 	Aliases: []string{"n"},
		// 	Required: false,
		// },
		// &cli.StringFlag{
		// 	Name: "begin",
		// 	Aliases: []string{"b"},
		// 	Required: false,
		// },
		// &cli.StringFlag{
		// 	Name: "name",
		// 	Aliases: []string{"n"},
		// 	Required: false,
		// },
	},
	Action: rmAction,
}

func rmAction(ctx context.Context, cmd *cli.Command) error {
	db_ptr := ctx.Value("db_ptr").(*sql.DB)

	if _, err := db_ptr.Exec(
		`DELETE FROM main 
		WHERE event_id=?;
		`,
		cmd.String("id"),
	); err != nil {
		return err
	}

	return nil
}
