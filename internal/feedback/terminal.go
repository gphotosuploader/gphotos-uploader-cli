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

	printPrompt(prompt)

	if secret {
		return readPassword()
	}

	return readInputLine()
}

// printPrompt prints the prompt to the user.
func printPrompt(prompt string) {
	fmt.Fprintf(stdOut, "%s: ", prompt)
}

// readPassword reads and returns a password from the user.
func readPassword() (string, error) {
	value, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Fprintln(stdOut)
	return string(value), err
}

// readInputLine reads and returns an input line from the user.
func readInputLine() (string, error) {
	sc := bufio.NewScanner(os.Stdin)
	sc.Scan()
	return sc.Text(), sc.Err()
}

func YesNoPrompt(prompt string, def bool) (bool, error) {
	choices := getChoices(def)
	prompt = fmt.Sprintf("%s (%s)", prompt, choices)

	for {
		s, err := InputUserField(prompt, false)
		if err != nil || s == "" {
			return def, err
		}

		s = strings.ToLower(s)
		if isYes(s) {
			return true, nil
		}
		if isNo(s) {
			return false, nil
		}
	}
}

func getChoices(def bool) string {
	if !def {
		return "y/N"
	}
	return "Y/n"
}

func isYes(s string) bool {
	return s == "y" || s == "yes"
}

func isNo(s string) bool {
	return s == "n" || s == "no"
}
