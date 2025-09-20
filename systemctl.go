package systemctl

import (
	"context"

	"github.com/taigrr/systemctl/properties"
)

// Reload systemd manager configuration.
//
// This will rerun all generators (see systemd. generator(7)), reload all unit
// files, and recreate the entire dependency tree. While the daemon is being
// reloaded, all sockets systemd listens on behalf of user configuration will
// stay accessible.
func DaemonReload(ctx context.Context, opts Options) error {
	return daemonReload(ctx, opts)
}

// Reenables one or more units.
//
// This removes all symlinks to the unit files backing the specified units from
// the unit configuration directory, then recreates the symlink to the unit again,
// atomically. Can be used to change the symlink target.
func Reenable(ctx context.Context, unit string, opts Options) error {
	return reenable(ctx, unit, opts)
}

// Disables one or more units.
//
// This removes all symlinks to the unit files backing the specified units from
// the unit configuration directory, and hence undoes any changes made by
// enable or link.
func Disable(ctx context.Context, unit string, opts Options) error {
	return disable(ctx, unit, opts)
}

// Enable one or more units or unit instances.
//
// This will create a set of symlinks, as encoded in the [Install] sections of
// the indicated unit files. After the symlinks have been created, the system
// manager configuration is reloaded (in a way equivalent to daemon-reload),
// in order to ensure the changes are taken into account immediately.
func Enable(ctx context.Context, unit string, opts Options) error {
	return enable(ctx, unit, opts)
}

// Check whether any of the specified units are active (i.e. running).
//
// Returns true if the unit is active, false if inactive or failed.
// Also returns false in an error case.
func IsActive(ctx context.Context, unit string, opts Options) (bool, error) {
	result, err := isActive(ctx, unit, opts)
	return result, err
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
	result, err := isEnabled(ctx, unit, opts)
	return result, err
}

// Check whether any of the specified units are in a "failed" state.
func IsFailed(ctx context.Context, unit string, opts Options) (bool, error) {
	result, err := isFailed(ctx, unit, opts)
	return result, err
}

// Mask one or more units, as specified on the command line. This will link
// these unit files to /dev/null, making it impossible to start them.
//
// Notably, Mask may return ErrDoesNotExist if a unit doesn't exist, but it will
// continue masking anyway. Calling Mask on a non-existing masked unit does not
// return an error. Similarly, see Unmask.
func Mask(ctx context.Context, unit string, opts Options) error {
	return mask(ctx, unit, opts)
}

// Stop and then start one or more units specified on the command line.
// If the units are not running yet, they will be started.
func Restart(ctx context.Context, unit string, opts Options) error {
	return restart(ctx, unit, opts)
}

// Show a selected property of a unit. Accepted properties are predefined in the
// properties subpackage to guarantee properties are valid and assist code-completion.
func Show(ctx context.Context, unit string, property properties.Property, opts Options) (string, error) {
	str, err := show(ctx, unit, property, opts)
	return str, err
}

// Start (activate) a given unit
func Start(ctx context.Context, unit string, opts Options) error {
	return start(ctx, unit, opts)
}

// Get back the status string which would be returned by running
// `systemctl status [unit]`.
//
// Generally, it makes more sense to programatically retrieve the properties
// using Show, but this command is provided for the sake of completeness
func Status(ctx context.Context, unit string, opts Options) (string, error) {
	stat, err := status(ctx, unit, opts)
	return stat, err
}

// Stop (deactivate) a given unit
func Stop(ctx context.Context, unit string, opts Options) error {
	return stop(ctx, unit, opts)
}

// Unmask one or more unit files, as specified on the command line.
// This will undo the effect of Mask.
//
// In line with systemd, Unmask will return ErrDoesNotExist if the unit
// doesn't exist, but only if it's not already masked.
// If the unit doesn't exist but it's masked anyway, no error will be
// returned. Gross, I know. Take it up with Poettering.
func Unmask(ctx context.Context, unit string, opts Options) error {
	return unmask(ctx, unit, opts)
}
