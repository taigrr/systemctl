package systemctl

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// Testing assumptions
// - there's no unit installed named `nonexistant`
// - the syncthing unit is available but the binary is not on the tester's system
//   this is just what was available on mine, should you want to change it,
//   either to something in this repo or more common, feel free to submit a PR.
// - your 'user' isn't root
// - your user doesn't have a PolKit rule allowing access to configure nginx

func TestGetStartTime(t *testing.T) {
	testCases := []struct {
		unit      string
		err       error
		opts      Options
		runAsUser bool
	}{
		// Run these tests only as a user
		//try nonexistant unit in user mode as user
		{"nonexistant", ErrUnitNotRunning, Options{usermode: false}, true},
		// try existing unit in user mode as user
		{"syncthing", ErrUnitNotRunning, Options{usermode: true}, true},
		// try existing unit in system mode as user
		{"nginx", nil, Options{usermode: false}, true},

		// Run these tests only as a superuser

		// try nonexistant unit in system mode as system
		{"nonexistant", ErrUnitNotRunning, Options{usermode: false}, false},
		// try existing unit in system mode as system
		{"nginx", ErrBusFailure, Options{usermode: true}, false},
		// try existing unit in system mode as system
		{"nginx", nil, Options{usermode: false}, false},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s as %s", tc.unit, userString), func(t *testing.T) {
			if (userString == "root" || userString == "system") && tc.runAsUser {
				t.Skip("skipping user test while running as superuser")
			} else if (userString != "root" && userString != "system") && !tc.runAsUser {
				t.Skip("skipping superuser test while running as user")
			}
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()
			_, err := GetStartTime(ctx, tc.unit, tc.opts)
			if err != tc.err {
				t.Errorf("error is %v, but should have been %v", err, tc.err)
			}
		})
	}
	// Prove start time changes after a restart
	t.Run(fmt.Sprintf("prove start time changes"), func(t *testing.T) {
		if userString != "root" && userString != "system" {
			t.Skip("skipping superuser test while running as user")
		}

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		startTime, err := GetStartTime(ctx, "nginx", Options{usermode: false})
		if err != nil {
			t.Errorf("issue getting start time of nginx: %v", err)
		}
		err = Restart(ctx, "nginx", Options{usermode: false})
		if err != nil {
			t.Errorf("issue restarting nginx as %s: %v", userString, err)
		}
		newStartTime, err := GetStartTime(ctx, "nginx", Options{usermode: false})
		if err != nil {
			t.Errorf("issue getting second start time of nginx: %v", err)
		}
		diff := newStartTime.Sub(startTime).Seconds()
		if diff <= 0 {
			t.Errorf("Expected start diff to be positive, but got: %d", int(diff))
		}
	})

}

func TestGetMemoryUsage(t *testing.T) {
	testCases := []struct {
		unit      string
		err       error
		opts      Options
		runAsUser bool
	}{
		// Run these tests only as a user

		//try nonexistant unit in user mode as user
		{"nonexistant", ErrValueNotSet, Options{usermode: false}, true},
		// try existing unit in user mode as user
		{"syncthing", ErrValueNotSet, Options{usermode: true}, true},
		// try existing unit in system mode as user
		{"nginx", nil, Options{usermode: false}, true},

		// Run these tests only as a superuser

		// try nonexistant unit in system mode as system
		{"nonexistant", ErrValueNotSet, Options{usermode: false}, false},
		// try existing unit in system mode as system
		{"nginx", ErrBusFailure, Options{usermode: true}, false},
		// try existing unit in system mode as system
		{"nginx", nil, Options{usermode: false}, false},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s as %s", tc.unit, userString), func(t *testing.T) {
			if (userString == "root" || userString == "system") && tc.runAsUser {
				t.Skip("skipping user test while running as superuser")
			} else if (userString != "root" && userString != "system") && !tc.runAsUser {
				t.Skip("skipping superuser test while running as user")
			}
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()
			_, err := GetMemoryUsage(ctx, tc.unit, tc.opts)
			if err != tc.err {
				t.Errorf("error is %v, but should have been %v", err, tc.err)
			}
		})
	}
	// Prove start time changes after a restart
	t.Run(fmt.Sprintf("prove start time changes"), func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		bytes, err := GetMemoryUsage(ctx, "nginx", Options{usermode: false})
		if err != nil {
			t.Errorf("issue getting memory usage of nginx: %v", err)
		}
		secondBytes, err := GetMemoryUsage(ctx, "user.slice", Options{usermode: false})
		if err != nil {
			t.Errorf("issue getting second memort usage reading of nginx: %v", err)
		}
		if bytes == secondBytes {
			t.Errorf("Expected memory usage between nginx and user.slice to differ, but both were: %d", bytes)
		}
	})

}
func TestGetPID(t *testing.T) {
	testCases := []struct {
		unit      string
		err       error
		opts      Options
		runAsUser bool
	}{
		// Run these tests only as a user

		//try nonexistant unit in user mode as user
		{"nonexistant", nil, Options{usermode: false}, true},
		// try existing unit in user mode as user
		{"syncthing", nil, Options{usermode: true}, true},
		// try existing unit in system mode as user
		{"nginx", nil, Options{usermode: false}, true},

		// Run these tests only as a superuser

		// try nonexistant unit in system mode as system
		{"nonexistant", nil, Options{usermode: false}, false},
		// try existing unit in system mode as system
		{"nginx", ErrBusFailure, Options{usermode: true}, false},
		// try existing unit in system mode as system
		{"nginx", nil, Options{usermode: false}, false},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s as %s", tc.unit, userString), func(t *testing.T) {
			if (userString == "root" || userString == "system") && tc.runAsUser {
				t.Skip("skipping user test while running as superuser")
			} else if (userString != "root" && userString != "system") && !tc.runAsUser {
				t.Skip("skipping superuser test while running as user")
			}
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()
			_, err := GetPID(ctx, tc.unit, tc.opts)
			if err != tc.err {
				t.Errorf("error is %v, but should have been %v", err, tc.err)
			}
		})
	}
}

// GetMemoryUsage(ctx context.Context, unit string, opts Options) (int, error) {
// GetPID(ctx context.Context, unit string, opts Options) (int, error) {
