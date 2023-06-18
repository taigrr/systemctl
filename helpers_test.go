package systemctl

import (
	"context"
	"errors"
	"fmt"
	"syscall"
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
	if testing.Short() {
		t.Skip("skipping in short mode")
	}
	testCases := []struct {
		unit      string
		err       error
		opts      Options
		runAsUser bool
	}{
		// Run these tests only as a user
		// try nonexistant unit in user mode as user
		{"nonexistant", ErrUnitNotActive, Options{UserMode: false}, true},
		// try existing unit in user mode as user
		{"syncthing", ErrUnitNotActive, Options{UserMode: true}, true},
		// try existing unit in system mode as user
		{"nginx", nil, Options{UserMode: false}, true},

		// Run these tests only as a superuser

		// try nonexistant unit in system mode as system
		{"nonexistant", ErrUnitNotActive, Options{UserMode: false}, false},
		// try existing unit in system mode as system
		{"nginx", ErrBusFailure, Options{UserMode: true}, false},
		// try existing unit in system mode as system
		{"nginx", nil, Options{UserMode: false}, false},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	Restart(ctx, "syncthing", Options{UserMode: true})
	Stop(ctx, "syncthing", Options{UserMode: true})
	time.Sleep(1 * time.Second)
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s as %s, UserMode=%v", tc.unit, userString, tc.opts.UserMode), func(t *testing.T) {
			if (userString == "root" || userString == "system") && tc.runAsUser {
				t.Skip("skipping user test while running as superuser")
			} else if (userString != "root" && userString != "system") && !tc.runAsUser {
				t.Skip("skipping superuser test while running as user")
			}
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()
			_, err := GetStartTime(ctx, tc.unit, tc.opts)
			if !errors.Is(err, tc.err) {
				t.Errorf("error is %v, but should have been %v", err, tc.err)
			}
		})
	}
	// Prove start time changes after a restart
	t.Run("prove start time changes", func(t *testing.T) {
		if userString != "root" && userString != "system" {
			t.Skip("skipping superuser test while running as user")
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		startTime, err := GetStartTime(ctx, "nginx", Options{UserMode: false})
		if err != nil {
			t.Errorf("issue getting start time of nginx: %v", err)
		}
		time.Sleep(1 * time.Second)
		err = Restart(ctx, "nginx", Options{UserMode: false})
		if err != nil {
			t.Errorf("issue restarting nginx as %s: %v", userString, err)
		}
		time.Sleep(100 * time.Millisecond)
		newStartTime, err := GetStartTime(ctx, "nginx", Options{UserMode: false})
		if err != nil {
			t.Errorf("issue getting second start time of nginx: %v", err)
		}
		diff := newStartTime.Sub(startTime).Seconds()
		if diff <= 0 {
			t.Errorf("Expected start diff to be positive, but got: %d", int(diff))
		}
	})
}

func TestGetNumRestarts(t *testing.T) {
	type testCase struct {
		unit      string
		err       error
		opts      Options
		runAsUser bool
	}
	testCases := []testCase{
		// Run these tests only as a user

		// try nonexistant unit in user mode as user
		{"nonexistant", ErrValueNotSet, Options{UserMode: false}, true},
		// try existing unit in user mode as user
		{"syncthing", ErrValueNotSet, Options{UserMode: true}, true},
		// try existing unit in system mode as user
		{"nginx", nil, Options{UserMode: false}, true},

		// Run these tests only as a superuser

		// try nonexistant unit in system mode as system
		{"nonexistant", ErrValueNotSet, Options{UserMode: false}, false},
		// try existing unit in system mode as system
		{"nginx", ErrBusFailure, Options{UserMode: true}, false},
		// try existing unit in system mode as system
		{"nginx", nil, Options{UserMode: false}, false},
	}
	for _, tc := range testCases {
		func(tc testCase) {
			t.Run(fmt.Sprintf("%s as %s", tc.unit, userString), func(t *testing.T) {
				t.Parallel()
				if (userString == "root" || userString == "system") && tc.runAsUser {
					t.Skip("skipping user test while running as superuser")
				} else if (userString != "root" && userString != "system") && !tc.runAsUser {
					t.Skip("skipping superuser test while running as user")
				}
				ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
				defer cancel()
				_, err := GetNumRestarts(ctx, tc.unit, tc.opts)
				if !errors.Is(err, tc.err) {
					t.Errorf("error is %v, but should have been %v", err, tc.err)
				}
			})
		}(tc)
	}
	// Prove restart count increases by one after a restart
	t.Run("prove restart count increases by one after a restart", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping in short mode")
		}
		if userString != "root" && userString != "system" {
			t.Skip("skipping superuser test while running as user")
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		restarts, err := GetNumRestarts(ctx, "nginx", Options{UserMode: false})
		if err != nil {
			t.Errorf("issue getting number of restarts for nginx: %v", err)
		}
		pid, err := GetPID(ctx, "nginx", Options{UserMode: false})
		if err != nil {
			t.Errorf("issue getting MainPID for nginx as %s: %v", userString, err)
		}
		syscall.Kill(pid, syscall.SIGKILL)
		for {
			running, errIsActive := IsActive(ctx, "nginx", Options{UserMode: false})
			if errIsActive != nil {
				t.Errorf("error asserting nginx is up: %v", errIsActive)
				break
			} else if running {
				break
			}
		}
		secondRestarts, err := GetNumRestarts(ctx, "nginx", Options{UserMode: false})
		if err != nil {
			t.Errorf("issue getting second reading on number of restarts for nginx: %v", err)
		}
		if restarts+1 != secondRestarts {
			t.Errorf("Expected restart count to differ by one, but difference was: %d", secondRestarts-restarts)
		}
	})
}

func TestGetMemoryUsage(t *testing.T) {
	type testCase struct {
		unit      string
		err       error
		opts      Options
		runAsUser bool
	}
	testCases := []testCase{
		// Run these tests only as a user

		// try nonexistant unit in user mode as user
		{"nonexistant", ErrValueNotSet, Options{UserMode: false}, true},
		// try existing unit in user mode as user
		{"syncthing", ErrValueNotSet, Options{UserMode: true}, true},
		// try existing unit in system mode as user
		{"nginx", nil, Options{UserMode: false}, true},

		// Run these tests only as a superuser

		// try nonexistant unit in system mode as system
		{"nonexistant", ErrValueNotSet, Options{UserMode: false}, false},
		// try existing unit in system mode as system
		{"nginx", ErrBusFailure, Options{UserMode: true}, false},
		// try existing unit in system mode as system
		{"nginx", nil, Options{UserMode: false}, false},
	}
	for _, tc := range testCases {
		func(tc testCase) {
			t.Run(fmt.Sprintf("%s as %s", tc.unit, userString), func(t *testing.T) {
				t.Parallel()
				if (userString == "root" || userString == "system") && tc.runAsUser {
					t.Skip("skipping user test while running as superuser")
				} else if (userString != "root" && userString != "system") && !tc.runAsUser {
					t.Skip("skipping superuser test while running as user")
				}
				ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
				defer cancel()
				_, err := GetMemoryUsage(ctx, tc.unit, tc.opts)
				if !errors.Is(err, tc.err) {
					t.Errorf("error is %v, but should have been %v", err, tc.err)
				}
			})
		}(tc)
	}
	// Prove memory usage values change across services
	t.Run("prove memory usage values change across services", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		bytes, err := GetMemoryUsage(ctx, "nginx", Options{UserMode: false})
		if err != nil {
			t.Errorf("issue getting memory usage of nginx: %v", err)
		}
		secondBytes, err := GetMemoryUsage(ctx, "user.slice", Options{UserMode: false})
		if err != nil {
			t.Errorf("issue getting second memory usage reading of nginx: %v", err)
		}
		if bytes == secondBytes {
			t.Errorf("Expected memory usage between nginx and user.slice to differ, but both were: %d", bytes)
		}
	})
}

func TestGetPID(t *testing.T) {
	type testCase struct {
		unit      string
		err       error
		opts      Options
		runAsUser bool
	}

	testCases := []testCase{
		// Run these tests only as a user

		// try nonexistant unit in user mode as user
		{"nonexistant", nil, Options{UserMode: false}, true},
		// try existing unit in user mode as user
		{"syncthing", nil, Options{UserMode: true}, true},
		// try existing unit in system mode as user
		{"nginx", nil, Options{UserMode: false}, true},

		// Run these tests only as a superuser

		// try nonexistant unit in system mode as system
		{"nonexistant", nil, Options{UserMode: false}, false},
		// try existing unit in system mode as system
		{"nginx", ErrBusFailure, Options{UserMode: true}, false},
		// try existing unit in system mode as system
		{"nginx", nil, Options{UserMode: false}, false},
	}
	for _, tc := range testCases {
		func(tc testCase) {
			t.Run(fmt.Sprintf("%s as %s", tc.unit, userString), func(t *testing.T) {
				t.Parallel()
				if (userString == "root" || userString == "system") && tc.runAsUser {
					t.Skip("skipping user test while running as superuser")
				} else if (userString != "root" && userString != "system") && !tc.runAsUser {
					t.Skip("skipping superuser test while running as user")
				}
				ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
				defer cancel()
				_, err := GetPID(ctx, tc.unit, tc.opts)
				if !errors.Is(err, tc.err) {
					t.Errorf("error is %v, but should have been %v", err, tc.err)
				}
			})
		}(tc)
	}
	t.Run("prove pid changes", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping in short mode")
		}
		if userString != "root" && userString != "system" {
			t.Skip("skipping superuser test while running as user")
		}
		unit := "nginx"
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		Restart(ctx, unit, Options{UserMode: true})
		pid, err := GetPID(ctx, unit, Options{UserMode: false})
		if err != nil {
			t.Errorf("issue getting MainPID for nginx as %s: %v", userString, err)
		}
		syscall.Kill(pid, syscall.SIGKILL)
		secondPid, err := GetPID(ctx, unit, Options{UserMode: false})
		if err != nil {
			t.Errorf("issue getting second MainPID for nginx as %s: %v", userString, err)
		}
		if pid == secondPid {
			t.Errorf("Expected pid != secondPid, but both were: %d", pid)
		}
	})
}
