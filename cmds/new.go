package cmds

import (
	"fmt"
	"context"

	"github.com/urfave/cli/v3"
)

var Cmd_new *cli.Command = &cli.Command{
	Name: "new",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "begin",
			Aliases: []string{"b"},
			Required: true,
		},
		&cli.StringFlag{
			Name: "name",
			Aliases: []string{"n"},
			Required: true,
		},
	},
	Action: newFunc,
}

func newFunc(_ context.Context, cmd *cli.Command) error {
	fmt.Printf("hi im new :), %s\n", cmd.String("name"))

	return nil
}
