//go:build linux

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
