package systemctl

import "context"

// TODO
func IsFailed(ctx context.Context, unit string) (bool, error) {

	return false, nil
}

// TODO
func IsActive(ctx context.Context, unit string) (bool, error) {
	return false, nil
}

// TODO
func Status(ctx context.Context, unit string) (bool, error) {
	return false, nil
}

// TODO
func Restart(ctx context.Context, unit string) error {
	return nil
}

// TODO
func Start(ctx context.Context, unit string) error {
	return nil
}

// TODO
func Stop(ctx context.Context, unit string) error {
	return nil
}

// TODO
func Enable(ctx context.Context, unit string) error {
	return nil
}

// TODO
func Disable(ctx context.Context, unit string) error {
	return nil
}

// TODO
func IsEnabled(ctx context.Context, unit string) (bool, error) {
	return false, nil
}

// TODO
func DaemonReload(ctx context.Context, unit string) error {
	return nil
}

//TODO
func Show(ctx context.Context, unit string, property string) (string, error) {
	return "", nil
}

//TODO
func Mask(ctx context.Context, unit string) error {
	return nil
}

func Unmask(ctx context.Context, unit string) error {
	return nil
}
