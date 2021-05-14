package systemctl

import (
	"context"
	"errors"
	"strconv"
)

// TODO
func IsFailed(ctx context.Context, unit string, usermode bool) (bool, error) {

	return false, nil
}

// TODO
func IsActive(ctx context.Context, unit string, usermode bool) (bool, error) {
	return false, nil
}

// TODO
func Status(ctx context.Context, unit string, usermode bool) (bool, error) {
	return false, nil
}

// TODO
func Restart(ctx context.Context, unit string, usermode bool) error {
	return nil
}

// TODO
func Start(ctx context.Context, unit string, usermode bool) error {
	return nil
}

// TODO
func Stop(ctx context.Context, unit string, usermode bool) error {
	return nil
}

// TODO
func Enable(ctx context.Context, unit string, usermode bool) error {
	var args = []string{"enable", "--system", unit}
	if usermode {
		args[1] = "--user"
	}
	_, stderr, code, err := execute(ctx, args)

	if err != nil {
		return err
	}
	err = filterErr(stderr)
	if err != nil {
		return err
	}
	if code != 0 {
		return errors.New("received error code " + strconv.Itoa(code))
	}
	return nil
}

// TODO
func Disable(ctx context.Context, unit string, usermode bool) error {
	return nil
}

// TODO
func IsEnabled(ctx context.Context, unit string, usermode bool) (bool, error) {
	return false, nil
}

// TODO
func DaemonReload(ctx context.Context, unit string, usermode bool) error {
	return nil
}

//TODO
func Show(ctx context.Context, unit string, property string, usermode bool) (string, error) {
	return "", nil
}

//TODO
func Mask(ctx context.Context, unit string, usermode bool) error {
	return nil
}

func Unmask(ctx context.Context, unit string, usermode bool) error {
	return nil
}
