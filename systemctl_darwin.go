//go:build !linux

package systemctl

import (
	"context"

	"github.com/taigrr/systemctl/properties"
)

func daemonReload(_ context.Context, _ Options, _ ...string) error {
	return nil
}

func reenable(_ context.Context, _ string, _ Options, _ ...string) error {
	return nil
}

func disable(_ context.Context, _ string, _ Options, _ ...string) error {
	return nil
}

func enable(_ context.Context, _ string, _ Options, _ ...string) error {
	return nil
}

func isActive(_ context.Context, _ string, _ Options, _ ...string) (bool, error) {
	return false, nil
}

func isEnabled(_ context.Context, _ string, _ Options, _ ...string) (bool, error) {
	return false, nil
}

func isFailed(_ context.Context, _ string, _ Options, _ ...string) (bool, error) {
	return false, nil
}

func mask(_ context.Context, _ string, _ Options, _ ...string) error {
	return nil
}

func restart(_ context.Context, _ string, _ Options, _ ...string) error {
	return nil
}

func show(_ context.Context, _ string, _ properties.Property, _ Options, _ ...string) (string, error) {
	return "", nil
}

func start(_ context.Context, _ string, _ Options, _ ...string) error {
	return nil
}

func status(_ context.Context, _ string, _ Options, _ ...string) (string, error) {
	return "", nil
}

func stop(_ context.Context, _ string, _ Options, _ ...string) error {
	return nil
}

func unmask(_ context.Context, _ string, _ Options, _ ...string) error {
	return nil
}
