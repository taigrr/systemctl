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
	err := Enable(ctx, unit, true)
	if err != ErrDoesNotExist {
		t.Errorf("error is %v, but should have been %v", err, ErrDoesNotExist)
	}

}
func TestEnableNoPermissions(t *testing.T) {
	unit := "nginx"
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err := Enable(ctx, unit, false)
	if err != ErrInsufficientPermissions {
		t.Errorf("error is %v, but should have been %v", err, ErrInsufficientPermissions)
	}

}
