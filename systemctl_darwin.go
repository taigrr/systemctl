//go:build !linux

package systemctl

import (
	"context"

	"github.com/taigrr/systemctl/properties"
)

func daemonReload(ctx context.Context, opts Options) error {
	return nil
}

func reenable(ctx context.Context, unit string, opts Options) error {
	return nil
}

func disable(ctx context.Context, unit string, opts Options) error {
	return nil
}

func enable(ctx context.Context, unit string, opts Options) error {
	return nil
}

func isActive(ctx context.Context, unit string, opts Options) (bool, error) {
	return false, nil
}

func isEnabled(ctx context.Context, unit string, opts Options) (bool, error) {
	return false, nil
}

func isFailed(ctx context.Context, unit string, opts Options) (bool, error) {
	return false, nil
}

func mask(ctx context.Context, unit string, opts Options) error {
	return nil
}

func restart(ctx context.Context, unit string, opts Options) error {
	return nil
}

func show(ctx context.Context, unit string, property properties.Property, opts Options) (string, error) {
	return "", nil
}

func start(ctx context.Context, unit string, opts Options) error {
	return nil
}

func status(ctx context.Context, unit string, opts Options) (string, error) {
	return "", nil
}

func stop(ctx context.Context, unit string, opts Options) error {
	return nil
}

func unmask(ctx context.Context, unit string, opts Options) error {
	return nil
}
