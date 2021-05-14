package systemctl

import (
	"bytes"
	"context"
	"os/exec"
	"regexp"
)

var systemctl string

const killed = 130

func init() {
	path, err := exec.LookPath("systemctl")
	if err != nil {
		panic(ErrNotInstalled)
	}
	systemctl = path
}

func execute(ctx context.Context, args []string) (string, string, int, error) {
	var (
		err      error
		stderr   bytes.Buffer
		stdout   bytes.Buffer
		code     int
		output   string
		warnings string
	)

	cmd := exec.CommandContext(ctx, systemctl, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	output = stdout.String()
	warnings = stderr.String()
	code = cmd.ProcessState.ExitCode()

	return output, warnings, code, err
}

func filterErr(stderr string) error {
	matched, _ := regexp.MatchString(`does not exist`, stderr)
	if matched {
		return ErrDoesNotExist
	}

	return nil
}
