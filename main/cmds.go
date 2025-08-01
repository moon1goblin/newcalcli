package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/moon1goblin/newcalcli/cal"

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
	Command_find *cli.Command = &cli.Command{
		Name: "find",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "name",
				Aliases:  []string{"n"},
				Required: true,
			},
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
			// Need to output an id of an event to stdout
			// First SELECT by name. If >1 results - ask for additional flags
			var e, b bool
			if cmd.String("begin") != "" {
				b = true
			}
			if cmd.String("end") != "" {
				e = true // absolute retarded shit, idk how to make it better
			}

			begin_time, _, err := cal.TimeFromStr(cmd.String("begin"))
			if err != nil && !errors.Is(err, cal.ErrEmptyString) {
				return fmt.Errorf("ListEvents error: %w", err)
			}
			// end_time, _, err := cal.TimeFromStr(cmd.String("end"))
			// if err != nil && !errors.Is(err, cal.ErrEmptyString) {
			// 	return fmt.Errorf("ListEvents error: %w", err)
			// }

			var count int

			if !e && !b {
				if err := cal.Db_ptr_g.QueryRow(`
				SELECT COUNT(*) 
				FROM sorted_view 
				WHERE event_name=?;
				`, cmd.String("name")).Scan(&count); err != nil {
					return fmt.Errorf("Erorr occured while searching: %s", err)
				}
			}

			if count > 1 {
				fmt.Printf("Found %d events with specified flags, specify more flags.\n", count)
				return nil
			}
			fmt.Printf("Found 1 shit %s %s", cmd.String("begin_datetime"), begin_time)

			// rows, err := cal.Db_ptr_g.Query(`
			// 	SELECT FROM sorted_view
			// 	WHERE event_name=?;
			// 	`, cmd.String("name"))
			// if err != nil {
			// 	return fmt.Errorf("Error occured: %s", err)
			// }
			return nil
		},
	}
	// TODO: better rm, not just rm by id lol
	Command_rm *cli.Command = &cli.Command{
		Name: "rm",
		Flags: []cli.Flag{
			&cli.Int64Flag{
				Name:     "id",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "name",
				Required: false,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.String("name") != "" {
				cal.Db_ptr_g.Exec(`
					SELECT FROM sorted_view 
					WHERE event_name=?;
				`, cmd.String("name"))
				// Select all of same name, then ask which to delete, like 1,2,3,4...
				fmt.Print("Select... idk\n")
				num, err := NumberSelectPrompt()
				if num != 0 || err != nil {
					return fmt.Errorf("rmAction: failed to delete by name: %w", err)
				}

				if _, err := cal.Db_ptr_g.Exec(
					`
				DELETE FROM main 
				WHERE event_name=?;
				`,
					cmd.String("name"),
				); err != nil {
					return fmt.Errorf("rmAction: failed to delete by name: %w: %w", cal.ErrSqlite, err)
				}
			} else if cmd.Int64("id") != 0 {
				if _, err := cal.Db_ptr_g.Exec(
					`
				DELETE FROM main 
				WHERE event_id=?;
				`,
					cmd.Int64("id"),
				); err != nil {
					return fmt.Errorf("rmAction: failed to delete by id: %w: %w", cal.ErrSqlite, err)
				}
			}

			return nil
		},
	}
)
