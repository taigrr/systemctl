package systemctl

import (
	"errors"
)

var (
	// $DBUS_SESSION_BUS_ADDRESS and $XDG_RUNTIME_DIR were not defined
	// This usually is the result of running in usermode as root
	ErrBusFailure = errors.New("bus connection failure")
	// The unit specified doesn't exist or can't be found
	ErrDoesNotExist = errors.New("unit does not exist")
	// The provided context was cancelled before the command finished execution
	ErrExecTimeout = errors.New("command timed out")
	// The executable was invoked without enough permissions to run the selected command
	// Running as superuser or adding the correct PolicyKit definitions can fix this
	// See https://wiki.debian.org/PolicyKit for more information
	ErrInsufficientPermissions = errors.New("insufficient permissions")
	// Selected unit file resides outside of the unit file search path
	ErrLinked = errors.New("unit file linked")
	// Masked units can only be unmasked, but something else was attempted
	// Unmask the unit before enabling or disabling it
	ErrMasked = errors.New("unit masked")
	// Make sure systemctl is in the PATH before calling again
	ErrNotInstalled = errors.New("systemctl not in $PATH")
	// A unit was expected to be running but was found inactive
	// This can happen when calling GetStartTime on a dead unit, for example
	ErrUnitNotActive = errors.New("unit not active")
	// A unit was expected to be loaded, but was not.
	// This can happen when trying to Stop a unit which does not exist, for example
	ErrUnitNotLoaded = errors.New("unit not loaded")
	// An expected value is unavailable, but the unit may be running
	// This can happen when calling GetMemoryUsage on systemd itself, for example
	ErrValueNotSet = errors.New("value not set")

	// Something in the stderr output contains the word `Failed`, but it is not a known case
	// This is a catch-all, and if it's ever seen in the wild, please submit a PR
	ErrUnspecified = errors.New("unknown error, please submit an issue at github.com/taigrr/systemctl")
)
