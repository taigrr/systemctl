package systemctl

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/taigrr/systemctl/properties"
)

func systemTestUnit(t *testing.T) string {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	units, err := GetUnits(ctx, Options{UserMode: false})
	if err != nil {
		t.Fatalf("get system units: %v", err)
	}
	for _, unit := range units {
		if unit.Active != "active" || unit.Sub != "running" || !strings.HasSuffix(unit.Name, ".service") {
			continue
		}
		name := strings.TrimSuffix(unit.Name, ".service")
		startTime, err := Show(ctx, name, properties.ExecMainStartTimestamp, Options{UserMode: false})
		if err != nil || startTime == "" {
			continue
		}
		restarts, err := Show(ctx, name, properties.NRestarts, Options{UserMode: false})
		if err != nil || restarts == "" || restarts == "[not set]" {
			continue
		}
		memory, err := Show(ctx, name, properties.MemoryCurrent, Options{UserMode: false})
		if err != nil || memory == "" || memory == "[not set]" {
			continue
		}
		return name
	}

	t.Skip("no readable active system service found for read-only tests")
	return ""
}
