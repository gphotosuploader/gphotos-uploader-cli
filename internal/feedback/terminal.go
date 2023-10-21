package feedback

import (
	"bufio"
	"errors"
	"fmt"
	"golang.org/x/term"
	"os"
	"strings"
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

func YesNoPrompt(prompt string, def bool) (bool, error) {
	choices := "Y/n"
	if !def {
		choices = "y/N"
	}

	prompt = fmt.Sprintf("%s (%s)", prompt, choices)

	for {
		s, err := InputUserField(prompt, false)
		if err != nil {
			return def, err
		}
		if s == "" {
			return def, nil
		}
		s = strings.ToLower(s)
		if s == "y" || s == "yes" {
			return true, nil
		}
		if s == "n" || s == "no" {
			return false, nil
		}
	}
}
