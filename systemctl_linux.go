//go:build linux

package systemctl

import (
	"context"
	"strings"

	"github.com/taigrr/systemctl/properties"
)

func daemonReload(ctx context.Context, opts Options, args ...string) error {
	a := prepareArgs("daemon-reload", opts, args...)
	_, _, _, err := execute(ctx, a)
	return err
}

func reenable(ctx context.Context, unit string, opts Options, args ...string) error {
	a := prepareArgs("reenable", opts, append([]string{unit}, args...)...)
	_, _, _, err := execute(ctx, a)
	return err
}

func disable(ctx context.Context, unit string, opts Options, args ...string) error {
	a := prepareArgs("disable", opts, append([]string{unit}, args...)...)
	_, _, _, err := execute(ctx, a)
	return err
}

func enable(ctx context.Context, unit string, opts Options, args ...string) error {
	a := prepareArgs("enable", opts, append([]string{unit}, args...)...)
	_, _, _, err := execute(ctx, a)
	return err
}

func isActive(ctx context.Context, unit string, opts Options, args ...string) (bool, error) {
	a := prepareArgs("is-active", opts, append([]string{unit}, args...)...)
	stdout, _, _, err := execute(ctx, a)
	stdout = strings.TrimSuffix(stdout, "\n")
	switch stdout {
	case "inactive":
		return false, nil
	case "active":
		return true, nil
	case "failed":
		return false, nil
	case "activating":
		return false, nil
	default:
		return false, err
	}
}

func isEnabled(ctx context.Context, unit string, opts Options, args ...string) (bool, error) {
	a := prepareArgs("is-enabled", opts, append([]string{unit}, args...)...)
	stdout, _, _, err := execute(ctx, a)
	stdout = strings.TrimSuffix(stdout, "\n")
	switch stdout {
	case "enabled":
		return true, nil
	case "enabled-runtime":
		return true, nil
	case "linked":
		return false, ErrLinked
	case "linked-runtime":
		return false, ErrLinked
	case "alias":
		return true, nil
	case "masked":
		return false, ErrMasked
	case "masked-runtime":
		return false, ErrMasked
	case "static":
		return true, nil
	case "indirect":
		return true, nil
	case "disabled":
		return false, nil
	case "generated":
		return true, nil
	case "transient":
		return true, nil
	}
	if err != nil {
		return false, err
	}
	return false, ErrUnspecified
}

func isFailed(ctx context.Context, unit string, opts Options, args ...string) (bool, error) {
	a := prepareArgs("is-failed", opts, append([]string{unit}, args...)...)
	stdout, _, _, err := execute(ctx, a)
	stdout = strings.TrimSuffix(stdout, "\n")
	switch stdout {
	case "inactive":
		return false, nil
	case "active":
		return false, nil
	case "failed":
		return true, nil
	default:
		return false, err
	}
}

func mask(ctx context.Context, unit string, opts Options, args ...string) error {
	a := prepareArgs("mask", opts, append([]string{unit}, args...)...)
	_, _, _, err := execute(ctx, a)
	return err
}

func restart(ctx context.Context, unit string, opts Options, args ...string) error {
	a := prepareArgs("restart", opts, append([]string{unit}, args...)...)
	_, _, _, err := execute(ctx, a)
	return err
}

func show(ctx context.Context, unit string, property properties.Property, opts Options, args ...string) (string, error) {
	extra := append([]string{unit, "--property", string(property)}, args...)
	a := prepareArgs("show", opts, extra...)
	stdout, _, _, err := execute(ctx, a)
	stdout = strings.TrimPrefix(stdout, string(property)+"=")
	stdout = strings.TrimSuffix(stdout, "\n")
	return stdout, err
}

func start(ctx context.Context, unit string, opts Options, args ...string) error {
	a := prepareArgs("start", opts, append([]string{unit}, args...)...)
	_, _, _, err := execute(ctx, a)
	return err
}

func status(ctx context.Context, unit string, opts Options, args ...string) (string, error) {
	a := prepareArgs("status", opts, append([]string{unit}, args...)...)
	stdout, _, _, err := execute(ctx, a)
	return stdout, err
}

func stop(ctx context.Context, unit string, opts Options, args ...string) error {
	a := prepareArgs("stop", opts, append([]string{unit}, args...)...)
	_, _, _, err := execute(ctx, a)
	return err
}

func unmask(ctx context.Context, unit string, opts Options, args ...string) error {
	a := prepareArgs("unmask", opts, append([]string{unit}, args...)...)
	_, _, _, err := execute(ctx, a)
	return err
}
