package cmdshit

import (
	"errors"
	"fmt"
	"os"
	"golang.org/x/term"
) 

var ErrYNPrompt = errors.New("y/n prompt error")

// so basically anything other than y or Y is false, good design hello
func ConfirmYNPrompt() (bool, error) {
	// stty into raw mode so we dont have to press enter
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return false, fmt.Errorf("ConfirmYNPrompt error setting terminal into raw mode: %w", err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	input := make([]byte, 1)
	if bytes_read, err := os.Stdin.Read(input); err != nil || bytes_read != 1 {
		return false, fmt.Errorf("ConfirmYNPrompt error reading char from stdin: %w", err)
	}

	fmt.Printf("%c", input[0])

	if input[0] == 'y' || input[0] == 'Y' {
		return true, nil
	} else {
		return false, nil
	}
}
