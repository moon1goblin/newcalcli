package main

import (
	"github.com/moon1goblin/newcalcli/cal"
	"context"
	"errors"
	"fmt"

	"github.com/urfave/cli/v3"
)

// TODO: descriptions for commands and flags (later)

var (
	Command_new *cli.Command = &cli.Command{
		Name: "new",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "begin",
				Aliases:  []string{"b"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "end",
				Aliases:  []string{"e"},
				Required: false,
			},
			&cli.StringFlag{
				Name:     "name",
				Aliases:  []string{"n"},
				Required: true,
			},
			&cli.BoolFlag{
				Name:     "yes",
				Aliases:  []string{"y"},
				Required: false,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			my_event, err := cal.EventCreate(cmd.String("name"), cmd.String("begin"), cmd.String("end"))
			if err != nil {
				return fmt.Errorf("error on command new: %w", err)
			}

			if !cmd.Bool("yes") {
				fmt.Printf("New event: %s\nConfirm? [Y/n]: ", my_event.StringWithDate())
				defer fmt.Print("\n")
				if confirmed, err := ConfirmYNPrompt(); err != nil {
					return fmt.Errorf("error on command new: %w", err)
				} else if !confirmed {
					return nil
				}
			}

			if err := my_event.Push(); err != nil {
				return fmt.Errorf("newAction error: %w", err)
			}

			return nil
		},
	}
	Command_ls *cli.Command = &cli.Command{
		Name: "ls",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "begin",
				Aliases:  []string{"b"},
				Required: false,
			},
			&cli.StringFlag{
				Name:     "end",
				Aliases:  []string{"e"},
				Required: false,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			// process dates
			begin_time, _, err := cal.TimeFromStr(cmd.String("begin"))
			if err != nil && !errors.Is(err, cal.ErrEmptyString) {
				return fmt.Errorf("ListEvents error: %w", err)
			}
			end_time, _, err := cal.TimeFromStr(cmd.String("end"))
			if err != nil && !errors.Is(err, cal.ErrEmptyString) {
				return fmt.Errorf("ListEvents error: %w", err)
			}

			// get sorted events in range [begin, end)
			events, err := cal.GetEventsInRange(begin_time, end_time)
			if err != nil {
				return fmt.Errorf("ListEvents error: %w", err)
			}

			fmt.Print(PrintEvents(events))

			return nil
		},
	}
	// TODO: better rm, not just rm by id lol
	Command_rm *cli.Command = &cli.Command{
		Name: "rm",
		Flags: []cli.Flag{
			&cli.Int64Flag{
				Name:     "id",
				Required: true,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if _, err := cal.Db_ptr_g.Exec(
				`
				DELETE FROM main 
				WHERE event_id=?;
				`,
				cmd.Int64("id"),
			); err != nil {
				return fmt.Errorf("rmAction: failed to delete by id: %w: %w", cal.ErrSqlite, err)
			}

			return nil
		},
	}
)
