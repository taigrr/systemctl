package systemctl

import "strings"

type Options struct {
	UserMode bool
}

type Unit struct {
	Name        string
	Load        string
	Active      string
	Sub         string
	Description string
}

// UnitTypes contains all valid systemd unit type suffixes.
var UnitTypes = []string{
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

// HasValidUnitSuffix checks whether the given unit name ends with a valid
// systemd unit type suffix (e.g. ".service", ".timer").
func HasValidUnitSuffix(unit string) bool {
	for _, t := range UnitTypes {
		if strings.HasSuffix(unit, "."+t) {
			return true
		}
	}
	return false
}
