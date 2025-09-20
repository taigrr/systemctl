package systemctl

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
