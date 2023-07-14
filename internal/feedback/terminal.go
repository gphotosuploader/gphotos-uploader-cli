package feedback

import (
	"bufio"
	"errors"
	"fmt"
	"golang.org/x/term"
	"os"
)

func isTerminal() bool {
	return term.IsTerminal(int(os.Stdin.Fd()))
}

// InputUserField prompts the user to input the provided user field.
func InputUserField(prompt string, secret bool) (string, error) {
	if !isTerminal() {
		return "", errors.New("user input not supported in non interactive mode")
	}

	fmt.Fprintf(stdOut, "%s: ", prompt)

	if secret {
		// Read and return a password (no character echoed on terminal)
		value, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Fprintln(stdOut)
		return string(value), err
	}

	// Read and return an input line
	sc := bufio.NewScanner(os.Stdin)
	sc.Scan()
	return sc.Text(), sc.Err()
}
