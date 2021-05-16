package systemctl

import (
	"context"
	"regexp"
	"strings"

	"github.com/taigrr/systemctl/properties"
)

func IsFailed(ctx context.Context, unit string, opts Options) (bool, error) {
	var args = []string{"is-failed", "--system", unit}
	if opts.UserMode {
		args[1] = "--user"
	}
	stdout, _, _, err := execute(ctx, args)
	if matched, _ := regexp.MatchString(`inactive`, stdout); matched {
		return false, nil
	} else if matched, _ := regexp.MatchString(`active`, stdout); matched {
		return false, nil
	} else if matched, _ := regexp.MatchString(`failed`, stdout); matched {
		return true, nil
	}
	return false, err
}

func IsActive(ctx context.Context, unit string, opts Options) (bool, error) {
	var args = []string{"is-active", "--system", unit}
	if opts.UserMode {
		args[1] = "--user"
	}
	stdout, _, _, err := execute(ctx, args)
	if matched, _ := regexp.MatchString(`inactive`, stdout); matched {
		return false, nil
	} else if matched, _ := regexp.MatchString(`active`, stdout); matched {
		return true, nil
	} else if matched, _ := regexp.MatchString(`failed`, stdout); matched {
		return false, nil
	}

	return false, err
}

func IsEnabled(ctx context.Context, unit string, opts Options) (bool, error) {
	var args = []string{"is-enabled", "--system", unit}
	if opts.UserMode {
		args[1] = "--user"
	}
	stdout, _, _, err := execute(ctx, args)
	if matched, _ := regexp.MatchString(`enabled`, stdout); matched {
		return true, nil
	} else if matched, _ := regexp.MatchString(`disabled`, stdout); matched {
		return false, nil
	} else if matched, _ := regexp.MatchString(`masked`, stdout); matched {
		return false, nil
	}

	return false, err
}

func Status(ctx context.Context, unit string, opts Options) (string, error) {
	var args = []string{"status", "--system", unit}
	if opts.UserMode {
		args[1] = "--user"
	}
	stdout, _, _, err := execute(ctx, args)
	return stdout, err
}

func Restart(ctx context.Context, unit string, opts Options) error {
	var args = []string{"restart", "--system", unit}
	if opts.UserMode {
		args[1] = "--user"
	}
	_, _, _, err := execute(ctx, args)
	return err
}

func Start(ctx context.Context, unit string, opts Options) error {
	var args = []string{"start", "--system", unit}
	if opts.UserMode {
		args[1] = "--user"
	}
	_, _, _, err := execute(ctx, args)
	return err
}

func Stop(ctx context.Context, unit string, opts Options) error {
	var args = []string{"stop", "--system", unit}
	if opts.UserMode {
		args[1] = "--user"
	}
	_, _, _, err := execute(ctx, args)
	return err
}

func Enable(ctx context.Context, unit string, opts Options) error {
	var args = []string{"enable", "--system", unit}
	if opts.UserMode {
		args[1] = "--user"
	}
	_, _, _, err := execute(ctx, args)
	return err
}

func Disable(ctx context.Context, unit string, opts Options) error {
	var args = []string{"disable", "--system", unit}
	if opts.UserMode {
		args[1] = "--user"
	}
	_, _, _, err := execute(ctx, args)
	return err
}

func DaemonReload(ctx context.Context, opts Options) error {
	var args = []string{"daemon-reload", "--system"}
	if opts.UserMode {
		args[1] = "--user"
	}
	_, _, _, err := execute(ctx, args)
	return err
}

func Show(ctx context.Context, unit string, property properties.Property, opts Options) (string, error) {
	var args = []string{"show", "--system", unit, "--property", string(property)}
	if opts.UserMode {
		args[1] = "--user"
	}
	stdout, _, _, err := execute(ctx, args)
	stdout = strings.TrimPrefix(stdout, string(property)+"=")
	stdout = strings.TrimSuffix(stdout, "\n")
	return stdout, err
}

func Mask(ctx context.Context, unit string, opts Options) error {
	var args = []string{"mask", "--system", unit}
	if opts.UserMode {
		args[1] = "--user"
	}
	_, _, _, err := execute(ctx, args)
	return err
}

func Unmask(ctx context.Context, unit string, opts Options) error {
	var args = []string{"unmask", "--system", unit}
	if opts.UserMode {
		args[1] = "--user"
	}
	_, _, _, err := execute(ctx, args)
	return err
}
