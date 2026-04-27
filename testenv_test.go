package systemctl

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"
)

var (
	testUnitOnce   sync.Once
	testUserUnit   string
	testSystemUnit string
)

func initTestUnits(t *testing.T) {
	t.Helper()
	testUnitOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		testUserUnit = findUserTestUnit(ctx)
		testSystemUnit = findSystemTestUnit(ctx)
	})
}

func requireUserTestUnit(t *testing.T) string {
	t.Helper()
	initTestUnits(t)
	if testUserUnit == "" {
		t.Skip("skipping: no manageable active user service found")
	}
	return testUserUnit
}

func requireSystemTestUnit(t *testing.T) string {
	t.Helper()
	initTestUnits(t)
	if testSystemUnit == "" {
		t.Skip("skipping: no readable active system service found")
	}
	return testSystemUnit
}

func findUserTestUnit(ctx context.Context) string {
	units, err := GetUnits(ctx, Options{UserMode: true})
	if err != nil {
		return ""
	}
	preferred := []string{"ha-to-openclaw.service", "mail-to-openclaw.service", "buxfer-sync.service", "openclaw-gateway.service"}
	for _, unit := range preferred {
		if userUnitUsable(ctx, unit) {
			return trimServiceSuffix(unit)
		}
	}
	for _, unit := range units {
		if unit.Load != "loaded" || unit.Active != "active" {
			continue
		}
		if userUnitUsable(ctx, unit.Name) {
			return trimServiceSuffix(unit.Name)
		}
	}
	return ""
}

func userUnitUsable(ctx context.Context, unit string) bool {
	trimmed := trimServiceSuffix(unit)
	if _, err := GetPID(ctx, trimmed, Options{UserMode: true}); err != nil {
		return false
	}
	if _, err := GetStartTime(ctx, trimmed, Options{UserMode: true}); err != nil {
		return false
	}
	if _, err := GetMemoryUsage(ctx, trimmed, Options{UserMode: true}); err != nil {
		return false
	}
	return true
}

func findSystemTestUnit(ctx context.Context) string {
	units, err := GetUnits(ctx, Options{UserMode: false})
	if err != nil {
		return ""
	}
	preferred := []string{"ssh.service", "cron.service", "dbus.service", "chrony.service", "containerd.service"}
	for _, unit := range preferred {
		if systemUnitUsable(ctx, unit) {
			return trimServiceSuffix(unit)
		}
	}
	for _, unit := range units {
		if unit.Load != "loaded" || unit.Active != "active" || !strings.HasSuffix(unit.Name, ".service") {
			continue
		}
		if systemUnitUsable(ctx, unit.Name) {
			return trimServiceSuffix(unit.Name)
		}
	}
	return ""
}

func systemUnitUsable(ctx context.Context, unit string) bool {
	trimmed := trimServiceSuffix(unit)
	if _, err := GetPID(ctx, trimmed, Options{UserMode: false}); err != nil {
		return false
	}
	if _, err := GetStartTime(ctx, trimmed, Options{UserMode: false}); err != nil {
		return false
	}
	if _, err := GetMemoryUsage(ctx, trimmed, Options{UserMode: false}); err != nil {
		return false
	}
	return true
}

func trimServiceSuffix(unit string) string {
	return strings.TrimSuffix(unit, ".service")
}

func waitForActiveState(ctx context.Context, unit string, opts Options, want bool) error {
	for {
		active, err := IsActive(ctx, unit, opts)
		if err != nil {
			return err
		}
		if active == want {
			return nil
		}
		select {
		case <-ctx.Done():
			return fmt.Errorf("timed out waiting for %s active=%v: %w", unit, want, ctx.Err())
		case <-time.After(100 * time.Millisecond):
		}
	}
}
