package util

import (
	"fmt"
	"os/exec"
	"strings"
)

func RunCmd(cmd string, args ...string) error {
	// Split the command into the executable and its arguments
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return fmt.Errorf("error: is there an empty command?")
	}

	execCmd := exec.Command(parts[0], parts[1:]...)

	output, err := execCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error: command %q failed with %v: %s", cmd, err, output)
	}
	return nil
}
