package main

import (
	"context"
	"github.com/moon1goblin/newcalcli/cal"
	"github.com/urfave/cli/v3"
	"log"
	"os"
)

func main() {
	if err := cal.InitDB(); err != nil {
		log.Fatal(err)
	}

	cmds := &cli.Command{
		Commands: []*cli.Command{
			Command_new,
			Command_ls,
			Command_rm,
			Command_find,
		},
	}

	if err := cmds.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
