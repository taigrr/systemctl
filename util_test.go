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

func TestParseMaskedUnits(t *testing.T) {
	stdout := `UNIT FILE                         STATE  PRESET
foo.service                       masked enabled
bar-baz.timer                     disabled enabled
foo.bar.service                   masked enabled
quux@one.service                  masked -
	tabbed.service	masked	enabled

4 unit files listed.`

	got := parseMaskedUnits(stdout)
	expected := []string{"foo", "foo.bar", "quux@one", "tabbed"}
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("parseMaskedUnits() = %v, want %v", got, expected)
	}
}

func TestServiceUnitName(t *testing.T) {
	tests := []struct {
		name     string
		unit     string
		expected string
	}{
		{name: "bare service name", unit: "nginx", expected: "nginx.service"},
		{name: "service suffix", unit: "nginx.service", expected: "nginx.service"},
		{name: "timer suffix", unit: "backup.timer", expected: "backup.timer"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := serviceUnitName(tt.unit)
			if got != tt.expected {
				t.Fatalf("serviceUnitName(%q) = %q, want %q", tt.unit, got, tt.expected)
			}
		})
	}
}

func TestUnitNameWithoutSuffix(t *testing.T) {
	tests := []struct {
		name     string
		unit     string
		expected string
	}{
		{name: "bare name", unit: "nginx", expected: "nginx"},
		{name: "service suffix", unit: "nginx.service", expected: "nginx"},
		{name: "timer suffix", unit: "backup.timer", expected: "backup"},
		{name: "preserves dotted prefix", unit: "foo.bar.service", expected: "foo.bar"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := unitNameWithoutSuffix(tt.unit)
			if got != tt.expected {
				t.Fatalf("unitNameWithoutSuffix(%q) = %q, want %q", tt.unit, got, tt.expected)
			}
		})
	}
}
