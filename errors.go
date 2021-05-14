package systemctl

import (
	"errors"
)

var (
	// ErrSystemctlNotInstalled means that upon trying to manually locate systemctl in the user's path,
	// it was not found. If this error occurs, the library isn't entirely useful.
	ErrNotInstalled = errors.New("systemd binary was not found")

	// ErrExecTimeout means that the provided context was done before the command finished execution.
	ErrExecTimeout = errors.New("command timed out")

	// ErrInsufficientPermissions means the calling executable was invoked without enough permissions to run the selected command
	ErrInsufficientPermissions = errors.New("insufficient permissions for action")

	// ErrDoesNotExist means the unit specified doesn't exist or can't be found
	ErrDoesNotExist = errors.New("Unit does not exist")
)
