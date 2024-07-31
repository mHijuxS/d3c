package helpers

import "strings"

func SeparateCommand(completeCommand string) (separatedCommand []string) {
	separatedCommand = strings.Split(strings.TrimSuffix(completeCommand, "\n"), " ")
	return separatedCommand
}
