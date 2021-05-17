package systemctl

import (
	"context"
	"strconv"
	"time"

	"github.com/taigrr/systemctl/properties"
)

const dateFormat = "Mon 2006-01-02 15:04:05 MST"

// Get start time of a service (`systemctl show [unit] --property ExecMainStartTimestamp`) as a `Time` type
func GetStartTime(ctx context.Context, unit string, opts Options) (time.Time, error) {
	value, err := Show(ctx, unit, properties.ExecMainStartTimestamp, opts)

	if err != nil {
		return time.Unix(0, 0), err
	}
	// ExecMainStartTimestamp returns an empty string if the unit is not running
	if value == "" {
		return time.Unix(0, 0), ErrUnitNotActive
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
