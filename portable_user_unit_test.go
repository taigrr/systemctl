package systemctl

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func userBusAvailable(t *testing.T) bool {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := GetUnits(ctx, Options{UserMode: true})
	return err == nil
}

func requireUserBus(t *testing.T) {
	t.Helper()
	if !userBusAvailable(t) {
		t.Skip("skipping user-mode lifecycle test without a reachable user systemd bus")
	}
}

func userTestUnit(t *testing.T) string {
	t.Helper()
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("get user home dir: %v", err)
	}
	unitDir := filepath.Join(homeDir, ".local", "share", "systemd", "user")
	if err := os.MkdirAll(unitDir, 0o755); err != nil {
		t.Fatalf("create user systemd dir: %v", err)
	}

	unitName := "openclaw-test"
	unitPath := filepath.Join(unitDir, unitName+".service")
	unitFile := `[Unit]
Description=OpenClaw portable test service

[Service]
Type=simple
ExecStart=/bin/sh -c 'sleep 300'
Restart=on-failure

[Install]
WantedBy=default.target
`
	if err := os.WriteFile(unitPath, []byte(unitFile), 0o644); err != nil {
		t.Fatalf("write user test unit: %v", err)
	}
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = Stop(ctx, unitName, Options{UserMode: true})
		_ = Disable(ctx, unitName, Options{UserMode: true})
		_ = Unmask(ctx, unitName, Options{UserMode: true})
		_ = os.Remove(filepath.Join(unitDir, "default.target.wants", unitName+".service"))
		_ = os.Remove(unitPath)
		_ = DaemonReload(ctx, Options{UserMode: true})
	})

	if userBusAvailable(t) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := DaemonReload(ctx, Options{UserMode: true}); err != nil {
			t.Fatalf("reload user daemon: %v", err)
		}
	}
	return unitName
}
