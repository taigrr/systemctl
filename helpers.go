package systemctl

import (
	"context"
	"errors"
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

// Get current memory in bytes (`systemctl show [unit] --property MemoryCurrent`) an an int
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

// check if a service is masked
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

// check if a service is running
// https://unix.stackexchange.com/a/396633
func IsRunning(ctx context.Context, unit string, opts Options) (bool, error) {
	status, err := Show(ctx, unit, properties.SubState, opts)
	return status == "running", err
}
