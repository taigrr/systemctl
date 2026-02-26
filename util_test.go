package systemctl

import (
	"reflect"
	"testing"
)

func TestPrepareArgs(t *testing.T) {
	tests := []struct {
		name     string
		base     string
		opts     Options
		extra    []string
		expected []string
	}{
		{
			name:     "system mode no extra",
			base:     "start",
			opts:     Options{},
			extra:    nil,
			expected: []string{"start", "--system"},
		},
		{
			name:     "user mode no extra",
			base:     "start",
			opts:     Options{UserMode: true},
			extra:    nil,
			expected: []string{"start", "--user"},
		},
		{
			name:     "system mode with unit",
			base:     "start",
			opts:     Options{},
			extra:    []string{"nginx.service"},
			expected: []string{"start", "--system", "nginx.service"},
		},
		{
			name:     "user mode with unit and extra args",
			base:     "restart",
			opts:     Options{UserMode: true},
			extra:    []string{"foo.service", "--no-block"},
			expected: []string{"restart", "--user", "foo.service", "--no-block"},
		},
		{
			name:     "daemon-reload no extra",
			base:     "daemon-reload",
			opts:     Options{},
			extra:    nil,
			expected: []string{"daemon-reload", "--system"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := prepareArgs(tt.base, tt.opts, tt.extra...)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("prepareArgs(%q, %+v, %v) = %v, want %v",
					tt.base, tt.opts, tt.extra, got, tt.expected)
			}
		})
	}
}
