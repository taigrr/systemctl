package systemctl

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestErrorFuncs(t *testing.T) {
	portableUserUnit := userTestUnit(t)
	testCases := []struct {
		name      string
		fn        func(ctx context.Context, unit string, opts Options) error
		lifecycle bool
		errCases  []struct {
			unit      string
			err       error
			opts      Options
			runAsUser bool
		}
	}{
		{
			name:      "Enable",
			fn:        func(ctx context.Context, unit string, opts Options) error { return Enable(ctx, unit, opts) },
			lifecycle: false,
			errCases: []struct {
				unit      string
				err       error
				opts      Options
				runAsUser bool
			}{
				{"nonexistant", ErrDoesNotExist, Options{UserMode: true}, true},
				{portableUserUnit, nil, Options{UserMode: true}, true},
				{"nonexistant", ErrInsufficientPermissions, Options{UserMode: false}, true},
				{"nginx", ErrInsufficientPermissions, Options{UserMode: false}, true},
				{"nonexistant", ErrDoesNotExist, Options{UserMode: false}, false},
				{"nginx", ErrBusFailure, Options{UserMode: true}, false},
				{"nginx", nil, Options{UserMode: false}, false},
			},
		},
		{
			name:      "Disable",
			fn:        func(ctx context.Context, unit string, opts Options) error { return Disable(ctx, unit, opts) },
			lifecycle: false,
			errCases: []struct {
				unit      string
				err       error
				opts      Options
				runAsUser bool
			}{
				{"nonexistant", ErrDoesNotExist, Options{UserMode: true}, true},
				{portableUserUnit, nil, Options{UserMode: true}, true},
				{"nonexistant", ErrInsufficientPermissions, Options{UserMode: false}, true},
				{"nginx", ErrInsufficientPermissions, Options{UserMode: false}, true},
				{"nonexistant", ErrDoesNotExist, Options{UserMode: false}, false},
				{"nginx", ErrBusFailure, Options{UserMode: true}, false},
				{"nginx", nil, Options{UserMode: false}, false},
			},
		},
		{
			name:      "Restart",
			fn:        func(ctx context.Context, unit string, opts Options) error { return Restart(ctx, unit, opts) },
			lifecycle: true,
			errCases: []struct {
				unit      string
				err       error
				opts      Options
				runAsUser bool
			}{
				{"nonexistant", ErrDoesNotExist, Options{UserMode: true}, true},
				{portableUserUnit, nil, Options{UserMode: true}, true},
				{"nonexistant", ErrInsufficientPermissions, Options{UserMode: false}, true},
				{"nginx", ErrInsufficientPermissions, Options{UserMode: false}, true},
				{"nonexistant", ErrDoesNotExist, Options{UserMode: false}, false},
				{"nginx", ErrBusFailure, Options{UserMode: true}, false},
				{"nginx", nil, Options{UserMode: false}, false},
			},
		},
		{
			name:      "Start",
			fn:        func(ctx context.Context, unit string, opts Options) error { return Start(ctx, unit, opts) },
			lifecycle: true,
			errCases: []struct {
				unit      string
				err       error
				opts      Options
				runAsUser bool
			}{
				{"nonexistant", ErrDoesNotExist, Options{UserMode: true}, true},
				{portableUserUnit, nil, Options{UserMode: true}, true},
				{"nonexistant", ErrInsufficientPermissions, Options{UserMode: false}, true},
				{"nginx", ErrInsufficientPermissions, Options{UserMode: false}, true},
				{"nonexistant", ErrDoesNotExist, Options{UserMode: false}, false},
				{"nginx", ErrBusFailure, Options{UserMode: true}, false},
				{"nginx", nil, Options{UserMode: false}, false},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Errorcheck %s", tc.name), func(t *testing.T) {
			for _, errCase := range tc.errCases {
				t.Run(fmt.Sprintf("%s as %s", errCase.unit, userString), func(t *testing.T) {
					if (userString == "root" || userString == "system") && errCase.runAsUser {
						t.Skip("skipping user test while running as superuser")
					} else if (userString != "root" && userString != "system") && !errCase.runAsUser {
						t.Skip("skipping superuser test while running as user")
					}
					if tc.lifecycle && errCase.runAsUser && errCase.opts.UserMode {
						requireUserBus(t)
					}
					ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
					defer cancel()
					err := tc.fn(ctx, errCase.unit, errCase.opts)
					if !errors.Is(err, errCase.err) {
						t.Errorf("error is %v, but should have been %v", err, errCase.err)
					}
				})
			}
		})
	}
}
