package systemctl

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestErrorFuncs(t *testing.T) {
	errFuncs := []func(ctx context.Context, unit string, opts Options) error{
		Enable,
		Disable,
		Restart,
		Start,
	}
	errCases := []struct {
		unit      string
		err       error
		opts      Options
		runAsUser bool
	}{
		/* Run these tests only as an unpriviledged user */

		// try nonexistant unit in user mode as user
		{"nonexistant", ErrDoesNotExist, Options{UserMode: true}, true},
		// try existing unit in user mode as user
		{"syncthing", nil, Options{UserMode: true}, true},
		// try nonexisting unit in system mode as user
		{"nonexistant", ErrInsufficientPermissions, Options{UserMode: false}, true},
		// try existing unit in system mode as user
		{"nginx", ErrInsufficientPermissions, Options{UserMode: false}, true},

		/* End user tests*/

		/* Run these tests only as a superuser */

		// try nonexistant unit in system mode as system
		{"nonexistant", ErrDoesNotExist, Options{UserMode: false}, false},
		// try existing unit in system mode as system
		{"nginx", ErrBusFailure, Options{UserMode: true}, false},
		// try existing unit in system mode as system
		{"nginx", nil, Options{UserMode: false}, false},

		/* End superuser tests*/

	}

	for _, f := range errFuncs {
		fName := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
		fName = strings.TrimPrefix(fName, "github.com/taigrr/")
		t.Run(fmt.Sprintf("Errorcheck %s", fName), func(t *testing.T) {
			for _, tc := range errCases {
				t.Run(fmt.Sprintf("%s as %s", tc.unit, userString), func(t *testing.T) {
					if (userString == "root" || userString == "system") && tc.runAsUser {
						t.Skip("skipping user test while running as superuser")
					} else if (userString != "root" && userString != "system") && !tc.runAsUser {
						t.Skip("skipping superuser test while running as user")
					}
					ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
					defer cancel()
					err := f(ctx, tc.unit, tc.opts)
					if !errors.Is(err, tc.err) {
						t.Errorf("error is %v, but should have been %v", err, tc.err)
					}
				})
			}
		})
	}
}
