//go:build linux

package systemctl

type Options struct {
	UserMode bool
	// When true, sudo is used to run the systemctl commands.
	Sudo bool
	// Specify the sudo password.
	// Only applies if Sudo is true.
	// Optional if passwordless sudo is enabled or SudoAskPass is used.
	// Mutually exclusive with SudoAskPass.
	SudoPass string
	// When true, the --askpass sudo flag is used. See man sudo.
	// Requires Sudo: true.
	// Mutually exclusive with SudoPass.
	SudoAskPass bool
}

type Unit struct {
	Name        string
	Load        string
	Active      string
	Sub         string
	Description string
}

var unitTypes = []string{
	"automount",
	"device",
	"mount",
	"path",
	"scope",
	"service",
	"slice",
	"snapshot",
	"socket",
	"swap",
	"target",
	"timer",
}
