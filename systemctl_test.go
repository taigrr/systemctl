package systemctl

import (
	"context"
	"testing"
	"time"
)

func TestEnableNonexistant(t *testing.T) {
	unit := "nonexistant"
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	opts := Options{
		usermode: true,
	}
	err := Enable(ctx, unit, opts)
	if err != ErrDoesNotExist {
		t.Errorf("error is %v, but should have been %v", err, ErrDoesNotExist)
	}

}

// Note: test assumes your user isn't root and doesn't have a PolKit rule allowing access
//       to configure nginx. Whether it's installed should be irrelevant.
func TestEnableNoPermissions(t *testing.T) {
	unit := "nginx"
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	opts := Options{
		usermode: false,
	}
	err := Enable(ctx, unit, opts)
	if err != ErrInsufficientPermissions {
		t.Errorf("error is %v, but should have been %v", err, ErrInsufficientPermissions)
	}

}

// Note: requires the syncthing unit to be available on the tester's system.
//       this is just what was available on mine, should you want to change it,
//       either to something in this repo or more common, feel free to submit a PR.
func TestEnableSuccess(t *testing.T) {
	unit := "syncthing"
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	opts := Options{
		usermode: true,
	}
	err := Enable(ctx, unit, opts)
	if err != nil {
		t.Errorf("error is %v, but should have been %v", err, nil)
	}
}
