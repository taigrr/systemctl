//go:build linux

package systemctl

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestExecuteContextErrors(t *testing.T) {
	fakeSystemctl := filepath.Join(t.TempDir(), "systemctl")
	// Replace the shell so cancellation cannot leave a child holding the pipes.
	if err := os.WriteFile(fakeSystemctl, []byte("#!/bin/sh\nprintf 'Failed to connect to bus' >&2\nexec sleep 30\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	original := systemctl
	systemctl = fakeSystemctl
	t.Cleanup(func() { systemctl = original })

	tests := []struct {
		name    string
		context func() (context.Context, context.CancelFunc)
		want    error
	}{
		{
			name: "already canceled",
			context: func() (context.Context, context.CancelFunc) {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx, cancel
			},
			want: context.Canceled,
		},
		{
			name: "already expired",
			context: func() (context.Context, context.CancelFunc) {
				return context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
			},
			want: context.DeadlineExceeded,
		},
		{
			name: "deadline during execution",
			context: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), 200*time.Millisecond)
			},
			want: context.DeadlineExceeded,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := test.context()
			defer cancel()
			_, _, code, err := execute(ctx, []string{"start", "example.service"})
			if !errors.Is(err, test.want) {
				t.Errorf("execute() error = %v, want %v", err, test.want)
			}
			if !errors.Is(err, ErrExecTimeout) {
				t.Errorf("execute() error = %v, want ErrExecTimeout", err)
			}
			if code == 0 {
				t.Error("execute() exit code = 0, want failure")
			}
		})
	}
}
