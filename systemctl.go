//go:build linux

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
	args := []string{"daemon-reload"}
	_, _, _, err := execute(ctx, args, opts)
	return err
}

// Reenables one or more units.
//
// This removes all symlinks to the unit files backing the specified units from
// the unit configuration directory, then recreates the symlink to the unit again,
// atomically. Can be used to change the symlink target.
func Reenable(ctx context.Context, unit string, opts Options) error {
	args := []string{"reenable", unit}
	_, _, _, err := execute(ctx, args, opts)
	return err
}

// Disables one or more units.
//
// This removes all symlinks to the unit files backing the specified units from
// the unit configuration directory, and hence undoes any changes made by
// enable or link.
func Disable(ctx context.Context, unit string, opts Options) error {
	args := []string{"disable", unit}
	_, _, _, err := execute(ctx, args, opts)
	return err
}

// Enable one or more units or unit instances.
//
// This will create a set of symlinks, as encoded in the [Install] sections of
// the indicated unit files. After the symlinks have been created, the system
// manager configuration is reloaded (in a way equivalent to daemon-reload),
// in order to ensure the changes are taken into account immediately.
func Enable(ctx context.Context, unit string, opts Options) error {
	args := []string{"enable", unit}
	_, _, _, err := execute(ctx, args, opts)
	return err
}

// Check whether any of the specified units are active (i.e. running).
//
// Returns true if the unit is active, false if inactive or failed.
// Also returns false in an error case.
func IsActive(ctx context.Context, unit string, opts Options) (bool, error) {
	args := []string{"is-active", unit}
	stdout, _, _, err := execute(ctx, args, opts)
	stdout = strings.TrimSuffix(stdout, "\n")
	switch stdout {
	case "inactive":
		return false, nil
	case "active":
		return true, nil
	case "failed":
		return false, nil
	case "activating":
		return false, nil
	default:
		return false, err
	}
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
	args := []string{"is-enabled", unit}
	stdout, _, _, err := execute(ctx, args, opts)
	stdout = strings.TrimSuffix(stdout, "\n")
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
	args := []string{"is-failed", unit}
	stdout, _, _, err := execute(ctx, args, opts)
	if matched, _ := regexp.MatchString(`inactive`, stdout); matched {
		return false, nil
	} else if matched, _ := regexp.MatchString(`active`, stdout); matched {
		return false, nil
	} else if matched, _ := regexp.MatchString(`failed`, stdout); matched {
		return true, nil
	}
	return false, err
}

// Mask one or more units, as specified on the command line. This will link
// these unit files to /dev/null, making it impossible to start them.
//
// Notably, Mask may return ErrDoesNotExist if a unit doesn't exist, but it will
// continue masking anyway. Calling Mask on a non-existing masked unit does not
// return an error. Similarly, see Unmask.
func Mask(ctx context.Context, unit string, opts Options) error {
	args := []string{"mask", unit}
	_, _, _, err := execute(ctx, args, opts)
	return err
}

// Stop and then start one or more units specified on the command line.
// If the units are not running yet, they will be started.
func Restart(ctx context.Context, unit string, opts Options) error {
	args := []string{"restart", unit}
	_, _, _, err := execute(ctx, args, opts)
	return err
}

// Show a selected property of a unit. Accepted properties are predefined in the
// properties subpackage to guarantee properties are valid and assist code-completion.
func Show(ctx context.Context, unit string, property properties.Property, opts Options) (string, error) {
	args := []string{"show", unit, "--property", string(property)}
	stdout, _, _, err := execute(ctx, args, opts)
	stdout = strings.TrimPrefix(stdout, string(property)+"=")
	stdout = strings.TrimSuffix(stdout, "\n")
	return stdout, err
}

// Start (activate) a given unit
func Start(ctx context.Context, unit string, opts Options) error {
	args := []string{"start", unit}
	_, _, _, err := execute(ctx, args, opts)
	return err
}

// Get back the status string which would be returned by running
// `systemctl status [unit]`.
//
// Generally, it makes more sense to programatically retrieve the properties
// using Show, but this command is provided for the sake of completeness
func Status(ctx context.Context, unit string, opts Options) (string, error) {
	args := []string{"status", unit}
	stdout, _, _, err := execute(ctx, args, opts)
	return stdout, err
}

// Stop (deactivate) a given unit
func Stop(ctx context.Context, unit string, opts Options) error {
	args := []string{"stop", unit}
	_, _, _, err := execute(ctx, args, opts)
	return err
}

// Unmask one or more unit files, as specified on the command line.
// This will undo the effect of Mask.
//
// In line with systemd, Unmask will return ErrDoesNotExist if the unit
// doesn't exist, but only if it's not already masked.
// If the unit doesn't exist but it's masked anyway, no error will be
// returned. Gross, I know. Take it up with Poettering.
func Unmask(ctx context.Context, unit string, opts Options) error {
	args := []string{"unmask", unit}
	_, _, _, err := execute(ctx, args, opts)
	return err
}
