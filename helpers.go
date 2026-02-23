package systemctl

import (
	"context"
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/taigrr/systemctl/properties"
)

const dateFormat = "Mon 2006-01-02 15:04:05 MST"

// Get start time of a service (`systemctl show [unit] --property ExecMainStartTimestamp`) as a `Time` type
func GetStartTime(ctx context.Context, unit string, opts Options) (time.Time, error) {
	value, err := Show(ctx, unit, properties.ExecMainStartTimestamp, opts)
	if err != nil {
		return time.Time{}, err
	}
	// ExecMainStartTimestamp returns an empty string if the unit is not running
	if value == "" {
		return time.Time{}, ErrUnitNotActive
	}
	return time.Parse(dateFormat, value)
}

// Get the number of times a process restarted (`systemctl show [unit] --property NRestarts`) as an int
func GetNumRestarts(ctx context.Context, unit string, opts Options) (int, error) {
	value, err := Show(ctx, unit, properties.NRestarts, opts)
	if err != nil {
		return -1, err
	}
	return strconv.Atoi(value)
}

// Get current memory in bytes (`systemctl show [unit] --property MemoryCurrent`) as an int
func GetMemoryUsage(ctx context.Context, unit string, opts Options) (int, error) {
	value, err := Show(ctx, unit, properties.MemoryCurrent, opts)
	if err != nil {
		return -1, err
	}
	if value == "[not set]" {
		return -1, ErrValueNotSet
	}
	return strconv.Atoi(value)
}

// Get the PID of the main process (`systemctl show [unit] --property MainPID`) as an int
func GetPID(ctx context.Context, unit string, opts Options) (int, error) {
	value, err := Show(ctx, unit, properties.MainPID, opts)
	if err != nil {
		return -1, err
	}
	return strconv.Atoi(value)
}

// GetSocketsForServiceUnit returns the socket units associated with a given service unit.
func GetSocketsForServiceUnit(ctx context.Context, unit string, opts Options) ([]string, error) {
	args := []string{"list-sockets", "--all", "--no-legend", "--no-pager"}
	if opts.UserMode {
		args = append(args, "--user")
	}
	stdout, _, _, err := execute(ctx, args)
	if err != nil {
		return []string{}, err
	}
	lines := strings.Split(stdout, "\n")
	sockets := []string{}
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}
		socketUnit := fields[1]
		serviceUnit := fields[2]
		if serviceUnit == unit+".service" {
			sockets = append(sockets, socketUnit)
		}
	}
	return sockets, nil
}

// GetUnits returns a list of all loaded units and their states.
func GetUnits(ctx context.Context, opts Options) ([]Unit, error) {
	args := []string{"list-units", "--all", "--no-legend", "--full", "--no-pager"}
	if opts.UserMode {
		args = append(args, "--user")
	}
	stdout, stderr, _, err := execute(ctx, args)
	if err != nil {
		return []Unit{}, errors.Join(err, filterErr(stderr))
	}
	lines := strings.Split(stdout, "\n")
	units := []Unit{}
	for _, line := range lines {
		entry := strings.Fields(line)
		if len(entry) < 4 {
			continue
		}
		unit := Unit{
			Name:        entry[0],
			Load:        entry[1],
			Active:      entry[2],
			Sub:         entry[3],
			Description: strings.Join(entry[4:], " "),
		}
		units = append(units, unit)
	}
	return units, nil
}

// GetMaskedUnits returns a list of all masked unit names.
func GetMaskedUnits(ctx context.Context, opts Options) ([]string, error) {
	args := []string{"list-unit-files", "--state=masked"}
	if opts.UserMode {
		args = append(args, "--user")
	}
	stdout, stderr, _, err := execute(ctx, args)
	if err != nil {
		return []string{}, errors.Join(err, filterErr(stderr))
	}
	lines := strings.Split(stdout, "\n")
	units := []string{}
	for _, line := range lines {
		if !strings.Contains(line, "masked") {
			continue
		}
		entry := strings.Split(line, " ")
		if len(entry) < 3 {
			continue
		}
		if entry[1] == "masked" {
			unit := entry[0]
			uName := strings.Split(unit, ".")
			unit = uName[0]
			units = append(units, unit)
		}
	}
	return units, nil
}

// IsSystemd checks if systemd is the current init system by reading /proc/1/comm.
func IsSystemd() (bool, error) {
	b, err := os.ReadFile("/proc/1/comm")
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(string(b)) == "systemd", nil
}

// IsMasked checks if a unit is masked.
func IsMasked(ctx context.Context, unit string, opts Options) (bool, error) {
	units, err := GetMaskedUnits(ctx, opts)
	if err != nil {
		return false, err
	}
	for _, u := range units {
		if u == unit {
			return true, nil
		}
	}
	return false, nil
}

// IsRunning checks if a unit's sub-state is "running".
// See https://unix.stackexchange.com/a/396633 for details.
func IsRunning(ctx context.Context, unit string, opts Options) (bool, error) {
	status, err := Show(ctx, unit, properties.SubState, opts)
	return status == "running", err
}
