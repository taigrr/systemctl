//go:build linux

package systemctl

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestReloadBuildsExpectedCommand(t *testing.T) {
	tests := []struct {
		name string
		opts Options
		want []string
	}{
		{
			name: "system mode",
			opts: Options{},
			want: []string{"reload", "--system", "nginx.service"},
		},
		{
			name: "user mode",
			opts: Options{UserMode: true},
			want: []string{"reload", "--user", "nginx.service"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			logFile := filepath.Join(tempDir, "args.log")
			fakeSystemctl := filepath.Join(tempDir, "systemctl")
			script := "#!/bin/sh\nprintf '%s\\n' \"$@\" > '" + logFile + "'\n"

			if err := os.WriteFile(fakeSystemctl, []byte(script), 0o755); err != nil {
				t.Fatalf("write fake systemctl: %v", err)
			}

			original := systemctl
			systemctl = fakeSystemctl
			t.Cleanup(func() {
				systemctl = original
			})

			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			if err := Reload(ctx, "nginx.service", tt.opts); err != nil {
				t.Fatalf("Reload returned error: %v", err)
			}

			gotBytes, err := os.ReadFile(logFile)
			if err != nil {
				t.Fatalf("read captured args: %v", err)
			}
			got := strings.Fields(strings.TrimSpace(string(gotBytes)))
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("Reload args = %v, want %v", got, tt.want)
			}
		})
	}
}
