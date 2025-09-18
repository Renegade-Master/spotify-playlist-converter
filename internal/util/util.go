package util

import "os/exec"

// CheckIfCommandExists checks if executable 'e' is in PATH
func CheckIfCommandExists(e ...string) bool {
	anyCommandFound := false

	for _, command := range e {
		_, err := exec.LookPath(command)

		if err != nil {
			anyCommandFound = true
			break
		}
	}

	return anyCommandFound
}
