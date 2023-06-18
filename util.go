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
	case strings.Contains(`does not exist`, stderr):
		return errors.Join(ErrDoesNotExist, fmt.Errorf("stderr: %s", stderr))
	case strings.Contains(`not found.`, stderr):
		return errors.Join(ErrDoesNotExist, fmt.Errorf("stderr: %s", stderr))
	case strings.Contains(`not loaded.`, stderr):
		return errors.Join(ErrUnitNotLoaded, fmt.Errorf("stderr: %s", stderr))
	case strings.Contains(`No such file or directory`, stderr):
		return errors.Join(ErrDoesNotExist, fmt.Errorf("stderr: %s", stderr))
	case strings.Contains(`Interactive authentication required`, stderr):
		return errors.Join(ErrInsufficientPermissions, fmt.Errorf("stderr: %s", stderr))
	case strings.Contains(`Access denied`, stderr):
		return errors.Join(ErrInsufficientPermissions, fmt.Errorf("stderr: %s", stderr))
	case strings.Contains(`DBUS_SESSION_BUS_ADDRESS`, stderr):
		return errors.Join(ErrBusFailure, fmt.Errorf("stderr: %s", stderr))
	case strings.Contains(`is masked`, stderr):
		return errors.Join(ErrMasked, fmt.Errorf("stderr: %s", stderr))
	case strings.Contains(`Failed`, stderr):
		return errors.Join(ErrUnspecified, fmt.Errorf("stderr: %s", stderr))
	default:
		return nil
	}
}
