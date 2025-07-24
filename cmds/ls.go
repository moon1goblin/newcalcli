package cmds

import (
	"fmt"
	"context"

	"github.com/urfave/cli/v3"
)

var Cmd_ls *cli.Command = &cli.Command{
	Name: "ls",
	Action: func(context.Context, *cli.Command) error {
		fmt.Println("list lol kek cheburek")
		return nil
	},
}
