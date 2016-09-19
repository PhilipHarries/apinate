package exec

import (
	"os/exec"
	"strings"
)

func SplitStringToSlice(s string) []string {
	rStr := strings.Split(s, "\n")
	if rStr[len(rStr)-1] == "" {
		rStr = rStr[:len(rStr)-1]
	}
	return rStr
}

func Exec(commandString string) ([]string, error) {

	splitString := strings.Fields(commandString)
	command := splitString[0]
	args := splitString[1:]

	out, err := exec.Command(command, args...).CombinedOutput()
	arrayOut := SplitStringToSlice(string(out))

	return arrayOut, err
}
