package cmdshit

import (
	"calcli/dbshit"
	"context"
	"fmt"
	"database/sql"

	"github.com/urfave/cli/v3"
	_ "modernc.org/sqlite"
)

var Cmd_rm *cli.Command = &cli.Command{
	Name: "rm",
	Flags: []cli.Flag{
		&cli.Int64Flag{
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

	if err := dbshit.DeleteById(cmd.Int("id"), db_ptr); err != nil {
		return fmt.Errorf("rmAction: failed to delete by id: %w", err)
	}

	return nil
}
