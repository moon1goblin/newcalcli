package main

import (
	"context"
	"log"
	"os"
	"calcli/cmds"

	"github.com/urfave/cli/v3"
)


func main() {
	cmds := &cli.Command{
		Commands: []*cli.Command{
			cmds.Cmd_new,
			cmds.Cmd_ls,
		},
	}

	if err := cmds.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
