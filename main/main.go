package main

import (
	"github.com/moon1goblin/newcalcli/cal"
	"context"
	"log"
	"os"
	"github.com/urfave/cli/v3"
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
		},
	}

	if err := cmds.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
