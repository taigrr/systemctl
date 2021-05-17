package systemctl

import (
	"context"
	"regexp"
	"strings"

	"github.com/taigrr/systemctl/properties"
)

// Reload systemd manager configuration.
//
// This will rerun all generators (see systemd. generator(7)), reload all unit
// files, and recreate the entire dependency tree. While the daemon is being
// reloaded, all sockets systemd listens on behalf of user configuration will
// stay accessible.
func DaemonReload(ctx context.Context, opts Options) error {
	var args = []string{"daemon-reload", "--system"}
	if opts.UserMode {
		args[1] = "--user"
	}
	_, _, _, err := execute(ctx, args)
	return err
}

// Disables one or more units.
//
// This removes all symlinks to the unit files backing the specified units from
// the unit configuration directory, and hence undoes any changes made by
// enable or link.
func Disable(ctx context.Context, unit string, opts Options) error {
	var args = []string{"disable", "--system", unit}
	if opts.UserMode {
		args[1] = "--user"
	}
	_, _, _, err := execute(ctx, args)
	return err
}

// Enable one or more units or unit instances.
//
// This will create a set of symlinks, as encoded in the [Install] sections of
// the indicated unit files. After the symlinks have been created, the system
// manager configuration is reloaded (in a way equivalent to daemon-reload),
// in order to ensure the changes are taken into account immediately.
func Enable(ctx context.Context, unit string, opts Options) error {
	var args = []string{"enable", "--system", unit}
	if opts.UserMode {
		args[1] = "--user"
	}
	_, _, _, err := execute(ctx, args)
	return err
}

// Check whether any of the specified units are active (i.e. running).
//
// Returns true if the unit is active, false if inactive or failed.
// Also returns false in an error case.
func IsActive(ctx context.Context, unit string, opts Options) (bool, error) {
	var args = []string{"is-active", "--system", unit}
	if opts.UserMode {
		args[1] = "--user"
	}
	stdout, _, _, err := execute(ctx, args)
	if matched, _ := regexp.MatchString(`inactive`, stdout); matched {
		return false, nil
	} else if matched, _ := regexp.MatchString(`active`, stdout); matched {
		return true, nil
	} else if matched, _ := regexp.MatchString(`failed`, stdout); matched {
		return false, nil
	}

	return false, err
}

// Checks whether any of the specified unit files are enabled (as with enable).
//
// Returns true if the unit is enabled, aliased, static, indirect, generated
// or transient.
//
// Returns false if disabled. Also returns an error if linked, masked, or bad.
//
// See https://www.freedesktop.org/software/systemd/man/systemctl.html#is-enabled%20UNIT%E2%80%A6
// for more information
func IsEnabled(ctx context.Context, unit string, opts Options) (bool, error) {
	var args = []string{"is-enabled", "--system", unit}
	if opts.UserMode {
		args[1] = "--user"
	}
	stdout, _, _, err := execute(ctx, args)
	switch stdout {
	case "enabled":
		return true, nil
	case "enabled-runtime":
		return true, nil
	case "linked":
		return false, ErrLinked
	case "linked-runtime":
		return false, ErrLinked
	case "alias":
		return true, nil
	case "masked":
		return false, ErrMasked
	case "masked-runtime":
		return false, ErrMasked
	case "static":
		return true, nil
	case "indirect":
		return true, nil
	case "disabled":
		return false, nil
	case "generated":
		return true, nil
	case "transient":
		return true, nil
	}
	if err != nil {
		return false, err
	}
	return false, ErrUnspecified
}

// Check whether any of the specified units are in a "failed" state.
func IsFailed(ctx context.Context, unit string, opts Options) (bool, error) {
	var args = []string{"is-failed", "--system", unit}
	if opts.UserMode {
		args[1] = "--user"
	}
	stdout, _, _, err := execute(ctx, args)
	if matched, _ := regexp.MatchString(`inactive`, stdout); matched {
		return false, nil
	} else if matched, _ := regexp.MatchString(`active`, stdout); matched {
		return false, nil
	} else if matched, _ := regexp.MatchString(`failed`, stdout); matched {
		return true, nil
	}
	return false, err
}

func Mask(ctx context.Context, unit string, opts Options) error {
	var args = []string{"mask", "--system", unit}
	if opts.UserMode {
		args[1] = "--user"
	}
	_, _, _, err := execute(ctx, args)
	return err
}

func Restart(ctx context.Context, unit string, opts Options) error {
	var args = []string{"restart", "--system", unit}
	if opts.UserMode {
		args[1] = "--user"
	}
	_, _, _, err := execute(ctx, args)
	return err
}

func Show(ctx context.Context, unit string, property properties.Property, opts Options) (string, error) {
	var args = []string{"show", "--system", unit, "--property", string(property)}
	if opts.UserMode {
		args[1] = "--user"
	}
	stdout, _, _, err := execute(ctx, args)
	stdout = strings.TrimPrefix(stdout, string(property)+"=")
	stdout = strings.TrimSuffix(stdout, "\n")
	return stdout, err
}

func Start(ctx context.Context, unit string, opts Options) error {
	var args = []string{"start", "--system", unit}
	if opts.UserMode {
		args[1] = "--user"
	}
	_, _, _, err := execute(ctx, args)
	return err
}

func Status(ctx context.Context, unit string, opts Options) (string, error) {
	var args = []string{"status", "--system", unit}
	if opts.UserMode {
		args[1] = "--user"
	}
	stdout, _, _, err := execute(ctx, args)
	return stdout, err
}

func Stop(ctx context.Context, unit string, opts Options) error {
	var args = []string{"stop", "--system", unit}
	if opts.UserMode {
		args[1] = "--user"
	}
	_, _, _, err := execute(ctx, args)
	return err
}

func Unmask(ctx context.Context, unit string, opts Options) error {
	var args = []string{"unmask", "--system", unit}
	if opts.UserMode {
		args[1] = "--user"
	}
	_, _, _, err := execute(ctx, args)
	return err
}
