package systemctl

import (
	"errors"
	"testing"
)

func TestFilterErr(t *testing.T) {
	tests := []struct {
		name   string
		stderr string
		want   error
	}{
		{
			name:   "empty stderr",
			stderr: "",
			want:   nil,
		},
		{
			name:   "unit does not exist",
			stderr: "Unit foo.service does not exist, proceeding anyway.",
			want:   ErrDoesNotExist,
		},
		{
			name:   "unit not found",
			stderr: "Unit foo.service not found.",
			want:   ErrDoesNotExist,
		},
		{
			name:   "unit not loaded",
			stderr: "Unit foo.service not loaded.",
			want:   ErrUnitNotLoaded,
		},
		{
			name:   "no such file or directory",
			stderr: "No such file or directory",
			want:   ErrDoesNotExist,
		},
		{
			name:   "interactive authentication required",
			stderr: "Interactive authentication required.",
			want:   ErrInsufficientPermissions,
		},
		{
			name:   "access denied",
			stderr: "Access denied",
			want:   ErrInsufficientPermissions,
		},
		{
			name:   "dbus session bus address",
			stderr: "Failed to connect to bus: $DBUS_SESSION_BUS_ADDRESS not set",
			want:   ErrBusFailure,
		},
		{
			name:   "unit is masked",
			stderr: "Unit foo.service is masked.",
			want:   ErrMasked,
		},
		{
			name:   "generic failed",
			stderr: "Failed to do something unknown",
			want:   ErrUnspecified,
		},
		{
			name:   "unrecognized warning",
			stderr: "Warning: something benign happened",
			want:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := filterErr(tt.stderr)
			if tt.want == nil {
				if got != nil {
					t.Errorf("filterErr(%q) = %v, want nil", tt.stderr, got)
				}
				return
			}
			if !errors.Is(got, tt.want) {
				t.Errorf("filterErr(%q) = %v, want error wrapping %v", tt.stderr, got, tt.want)
			}
		})
	}
}

func TestHasValidUnitSuffix(t *testing.T) {
	tests := []struct {
		unit string
		want bool
	}{
		{"nginx.service", true},
		{"sshd.socket", true},
		{"backup.timer", true},
		{"dev-sda1.device", true},
		{"home.mount", true},
		{"dev-sda1.swap", true},
		{"user.slice", true},
		{"multi-user.target", true},
		{"session-1.scope", true},
		{"foo.automount", true},
		{"backup.path", true},
		{"foo.snapshot", true},
		{"nginx", false},
		{"", false},
		{"foo.bar", false},
		{"foo.services", false},
		{".service", true},
	}

	for _, tt := range tests {
		t.Run(tt.unit, func(t *testing.T) {
			got := HasValidUnitSuffix(tt.unit)
			if got != tt.want {
				t.Errorf("HasValidUnitSuffix(%q) = %v, want %v", tt.unit, got, tt.want)
			}
		})
	}
}
