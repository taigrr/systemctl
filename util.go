//go:build linux

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
var sudoPath string

const killed = 130

func init() {
	path, _ := exec.LookPath("systemctl")
	sudoPath, _ = exec.LookPath("sudo")
	systemctl = path
}

func execute(ctx context.Context, args []string, opts Options) (string, string, int, error) {
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

	if opts.UserMode {
		args = append(args, "--user")
	}

	var cmd *exec.Cmd
	if opts.Sudo {
		if sudoPath == "" {
			return "", "", 1, ErrNoSudo
		}

		var sudoArgs []string
		if opts.SudoAskPass {
			sudoArgs = append(sudoArgs, "--askpass")
		}
		if opts.SudoPass != "" {
			sudoArgs = append(sudoArgs, "--stdin")
		}
		sudoArgs = append(sudoArgs, systemctl)

		args = append(sudoArgs, args...)
		cmd = exec.CommandContext(ctx, "sudo", args...)
		if opts.SudoPass != "" {
			cmd.Stdin = strings.NewReader(opts.SudoPass)
		}
	} else {
		cmd = exec.CommandContext(ctx, systemctl, args...)
	}
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
	case strings.Contains(stderr, "sudo: a terminal is required to read the password"):
		return ErrSudoPasswordEntryFail
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
