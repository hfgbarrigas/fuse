// Package process is a bridge for OS commands
package process

import (
	"io/ioutil"
	"os/exec"

	"github.com/pkg/errors"
)

// ExecuteProcess will run a bash command through os.Exec.
// If no workingDir is provided, its nil or empty, it will run in calling process's current directory.
func ExecuteProcess(command string, workingDir *string) (stdoutStr, stderrStr string, err error) {
	cmd := exec.Command("/bin/sh", "-c", command)

	if workingDir != nil && *workingDir != "" {
		cmd.Dir = *workingDir
	}

	stdout, err := cmd.StdoutPipe()

	if err != nil {
		return "", "", errors.Wrap(err, "Process error")
	}

	stderr, err := cmd.StderrPipe()

	if err != nil {
		return "", "", errors.Wrap(err, "Process error")
	}

	if err = cmd.Start(); err != nil {
		return "", "", errors.Wrap(err, "Process error")
	}

	stderrBytes, err := ioutil.ReadAll(stderr)
	if err != nil {
		return "", "", errors.Wrap(err, "Process error")
	}

	stdoutBytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		return "", "", errors.Wrap(err, "Process error")
	}

	if err := cmd.Wait(); err != nil {
		return string(stdoutBytes), string(stderrBytes), errors.Wrap(err, "Process error")
	}

	return string(stdoutBytes), string(stderrBytes), nil
}
