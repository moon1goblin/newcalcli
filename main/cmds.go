package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/moon1goblin/newcalcli/cal"
	"github.com/urfave/cli/v3"
	"os"
	"strconv"
	"strings"
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
				Required: false,
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
			begin_time, _, err := cal.TimeFromStr(cmd.String("begin"))
			if err != nil && !errors.Is(err, cal.ErrEmptyString) {
				return fmt.Errorf("ListEvents error: %w", err)
			}
			end_time, _, err := cal.TimeFromStr(cmd.String("end"))
			if err != nil && !errors.Is(err, cal.ErrEmptyString) {
				return fmt.Errorf("ListEvents error: %w", err)
			}

			var bld strings.Builder
			args := []any{}
			bld.WriteString("SELECT COUNT(*) FROM sorted_view WHERE 1=1")
			if cmd.String("name") != "" {
				bld.WriteString(" AND event_name=?")
				args = append(args, cmd.String("name"))
			}
			if cmd.String("begin") != "" {
				bld.WriteString(" AND begin_datetime=?")
				args = append(args, begin_time.Time.Unix())
			}
			if cmd.String("end") != "" {
				bld.WriteString(" AND end_datetime=?")
				args = append(args, end_time.Time.Unix())
			}
			bld.WriteString(";")

			var count int
			if err := cal.Db_ptr_g.QueryRow(bld.String(), args...).Scan(&count); err != nil {
				return fmt.Errorf("Error while copying from row to values pointed by dest: %s", err)
			}

			fmt.Println(count)
			return nil
		},
	}
	Command_rm *cli.Command = &cli.Command{
		Name: "rm",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "name",
				Aliases:  []string{"n"},
				Required: false,
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
			&cli.BoolFlag{
				Name:     "yes",
				Aliases:  []string{"y"},
				Required: false,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			begin_time, _, err := cal.TimeFromStr(cmd.String("begin"))
			if err != nil && !errors.Is(err, cal.ErrEmptyString) {
				return fmt.Errorf("ListEvents error: %w", err)
			}
			end_time, _, err := cal.TimeFromStr(cmd.String("end"))
			if err != nil && !errors.Is(err, cal.ErrEmptyString) {
				return fmt.Errorf("ListEvents error: %w", err)
			}

			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Run the command
			if err := Command_find.Run(context.Background(), []string{"run", "find", "-n", cmd.String("name"), "-b", cmd.String("begin"), "-e", cmd.String("end")}); err != nil {
				return fmt.Errorf("Erorr in find command: %s", err)
			}

			// Restore stdout
			w.Close()
			os.Stdout = oldStdout

			// Capture output
			var buf bytes.Buffer
			buf.ReadFrom(r)
			output := buf.String()
			count, err := strconv.Atoi(strings.TrimRight(output, "\n"))
			if err != nil {
				return fmt.Errorf("Error while converting string to integer: %s", err)
			}
			switch count {
			default:
				fmt.Printf("Found %d events with given flags. Specify more flags or change existing ones.\n", count)
			case 1:
				var bld strings.Builder
				args := []any{}
				bld.WriteString("DELETE FROM main WHERE 1=1")
				if cmd.String("name") != "" {
					bld.WriteString(" AND event_name=?")
					args = append(args, cmd.String("name"))
				}
				if cmd.String("begin") != "" {
					bld.WriteString(" AND begin_datetime=?")
					args = append(args, begin_time.Time.Unix())
				}
				if cmd.String("end") != "" {
					bld.WriteString(" AND end_datetime=?")
					args = append(args, end_time.Time.Unix())
				}
				bld.WriteString(";")
				if !cmd.Bool("yes") {
					fmt.Print("Confirm? [Y/n]: ")
					defer fmt.Print("\n")
					if confirmed, err := ConfirmYNPrompt(); err != nil {
						return fmt.Errorf("error on command new: %w", err)
					} else if !confirmed {
						return nil
					}
				}
				if _, err := cal.Db_ptr_g.Exec(bld.String(), args...); err != nil {
					return fmt.Errorf("Error while execution query: %s", err)
				}
			}
			return nil
		},
	}
)
