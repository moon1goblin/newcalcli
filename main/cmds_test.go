package main

import (
	"bytes"
	"context"
	"github.com/moon1goblin/newcalcli/cal"
	"os"
	"testing"
)

func TestFindCommand(t *testing.T) {
	cal.InitDB("test_db")
	// test_db file in /main is a test db. Add new events to it with these lines ↓
	// ev, _ := cal.EventCreate("a", "01.01.12.00", "")
	// ev.Push()
	// ev, _ = cal.EventCreate("a", "01.01.12.00", "01.01.13.00")
	// ev.Push()
	// ev, _ = cal.EventCreate("b", "01.01.12.00", "01.01.13.00")
	// ev.Push()
	// ev, _ = cal.EventCreate("a", "02.01.12.00", "")
	// ev.Push()
	// ev, _ = cal.EventCreate("a", "01.02.12.00", "")
	// ev.Push()
	// ev, _ = cal.EventCreate("a", "01.03.12.00", "")
	// ev.Push()
	// ev, _ = cal.EventCreate("c", "12.12", "")
	// ev.Push()

	tests := []struct {
		name           string
		args           []string
		expectedOutput string
		expectedError  bool
	}{
		{
			name:           "./run find",
			args:           []string{"find"},
			expectedOutput: "7\n",
			expectedError:  false,
		},
		{
			name:           "./run find -n a",
			args:           []string{"find", "-n", "a"},
			expectedOutput: "5\n",
			expectedError:  false,
		},
		{
			name:           "./run find -e 01.01.13.00",
			args:           []string{"find", "-e", "01.01.13.00"},
			expectedOutput: "2\n",
			expectedError:  false,
		},
		{
			name:           "./run find -b 01.01.12.00",
			args:           []string{"find", "-b", "01.01.12.00"},
			expectedOutput: "3\n",
			expectedError:  false,
		},
		{
			name:           "./run find -b 01.01.12.00 -e 01.01.13.00",
			args:           []string{"find", "-b", "01.01.12.00", "-e", "01.01.13.00"},
			expectedOutput: "2\n",
			expectedError:  false,
		},
		{
			name:           "./run find -n a -b 01.01.12.00 -e 01.01.13.00",
			args:           []string{"find", "-n", "a", "-b", "01.01.12.00", "-e", "01.01.13.00"},
			expectedOutput: "1\n",
			expectedError:  false,
		},
		{
			name:           "./run find -b 12.12",
			args:           []string{"find", "-b", "12.12"},
			expectedOutput: "1\n",
			expectedError:  false,
		},
		{
			name:           "./run find -e 12.12",
			args:           []string{"find", "-e", "12.12"},
			expectedOutput: "1\n",
			expectedError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Redirect stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Run the command
			err := Command_find.Run(context.Background(), append([]string{"run"}, tt.args...))

			// Restore stdout
			w.Close()
			os.Stdout = oldStdout

			// Capture output
			var buf bytes.Buffer
			buf.ReadFrom(r)
			output := buf.String()

			// Check for expected error
			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected an error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}

			// Check output
			if output != tt.expectedOutput {
				t.Errorf("Expected output:\n%q\nGot:\n%q", tt.expectedOutput, output)
			}
		})
	}
}
