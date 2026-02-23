package systemctl

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

var systemctl string

// killed is the exit code returned when a process is terminated by SIGINT.
const killed = 130

func init() {
	path, _ := exec.LookPath("systemctl")
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

	if systemctl == "" {
		return "", "", 1, ErrNotInstalled
	}
	cmd := exec.CommandContext(ctx, systemctl, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	output = stdout.String()
	warnings = stderr.String()
	code = cmd.ProcessState.ExitCode()

	if code == killed {
		return output, warnings, code, ErrExecTimeout
	}

	customErr := filterErr(warnings)
	if customErr != nil {
		err = customErr
	}
	if code != 0 && err == nil {
		err = fmt.Errorf("received error code %d for stderr `%s`: %w", code, warnings, ErrUnspecified)
	}

	return output, warnings, code, err
}

func filterErr(stderr string) error {
	switch {
	case strings.Contains(stderr, `does not exist`):
		return errors.Join(ErrDoesNotExist, fmt.Errorf("stderr: %s", stderr))
	case strings.Contains(stderr, `not found.`):
		return errors.Join(ErrDoesNotExist, fmt.Errorf("stderr: %s", stderr))
	case strings.Contains(stderr, `not loaded.`):
		return errors.Join(ErrUnitNotLoaded, fmt.Errorf("stderr: %s", stderr))
	case strings.Contains(stderr, `No such file or directory`):
		return errors.Join(ErrDoesNotExist, fmt.Errorf("stderr: %s", stderr))
	case strings.Contains(stderr, `Interactive authentication required`):
		return errors.Join(ErrInsufficientPermissions, fmt.Errorf("stderr: %s", stderr))
	case strings.Contains(stderr, `Access denied`):
		return errors.Join(ErrInsufficientPermissions, fmt.Errorf("stderr: %s", stderr))
	case strings.Contains(stderr, `DBUS_SESSION_BUS_ADDRESS`):
		return errors.Join(ErrBusFailure, fmt.Errorf("stderr: %s", stderr))
	case strings.Contains(stderr, `is masked`):
		return errors.Join(ErrMasked, fmt.Errorf("stderr: %s", stderr))
	case strings.Contains(stderr, `Failed`):
		return errors.Join(ErrUnspecified, fmt.Errorf("stderr: %s", stderr))
	default:
		return nil
	}
}
