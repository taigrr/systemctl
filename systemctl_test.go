package systemctl

import (
	"context"
	"fmt"
	"os"
	"os/user"
	"syscall"
	"testing"
	"time"

	"github.com/taigrr/systemctl/properties"
)

var userString string

// Testing assumptions
// - there's no unit installed named `nonexistant`
// - the syncthing unit to be available on the tester's system.
//   this is just what was available on mine, should you want to change it,
//   either to something in this repo or more common, feel free to submit a PR.
// - your 'user' isn't root
// - your user doesn't have a PolKit rule allowing access to configure nginx

func TestMain(m *testing.M) {
	curUser, err := user.Current()
	if err != nil {
		fmt.Println("Could not determine running user")
	}
	userString = curUser.Username
	fmt.Printf("currently running tests as: %s \n", userString)
	fmt.Println("Don't forget to run both root and user tests.")
	os.Exit(m.Run())
}
func TestDaemonReload(t *testing.T) {
	testCases := []struct {
		unit      string
		err       error
		opts      Options
		runAsUser bool
	}{
		/* Run these tests only as a user */

		// fail to reload system daemon as user
		{"", ErrInsufficientPermissions, Options{UserMode: false}, true},
		// reload user's scope daemon
		{"", nil, Options{UserMode: true}, true},
		/* End user tests*/

		/* Run these tests only as a superuser */

		// succeed to reload daemon
		{"", nil, Options{UserMode: false}, false},
		// fail to connect to user bus as system
		{"", ErrBusFailure, Options{UserMode: true}, false},

		/* End superuser tests*/
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
			err := DaemonReload(ctx, tc.opts)
			if err != tc.err {
				t.Errorf("error is %v, but should have been %v", err, tc.err)
			}
		})
	}
}
func TestDisable(t *testing.T) {
	t.Run(fmt.Sprintf(""), func(t *testing.T) {
		if userString != "root" && userString != "system" {
			t.Skip("skipping superuser test while running as user")
		}
		unit := "nginx"
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		err := Mask(ctx, unit, Options{UserMode: false})
		defer cancel()
		if err != nil {
			Unmask(ctx, unit, Options{UserMode: false})
			t.Errorf("Unable to mask %s", unit)
		}
		err = Disable(ctx, unit, Options{UserMode: false})
		if err != ErrMasked {
			Unmask(ctx, unit, Options{UserMode: false})
			t.Errorf("error is %v, but should have been %v", err, ErrMasked)
		}
		err = Unmask(ctx, unit, Options{UserMode: false})
		if err != nil {
			t.Errorf("Unable to unmask %s", unit)
		}
	})

}
func TestEnable(t *testing.T) {

	t.Run(fmt.Sprintf(""), func(t *testing.T) {
		if userString != "root" && userString != "system" {
			t.Skip("skipping superuser test while running as user")
		}
		unit := "nginx"
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err := Mask(ctx, unit, Options{UserMode: false})
		if err != nil {
			Unmask(ctx, unit, Options{UserMode: false})
			t.Errorf("Unable to mask %s", unit)
		}
		err = Enable(ctx, unit, Options{UserMode: false})
		if err != ErrMasked {
			Unmask(ctx, unit, Options{UserMode: false})
			t.Errorf("error is %v, but should have been %v", err, ErrMasked)
		}
		err = Unmask(ctx, unit, Options{UserMode: false})
		if err != nil {
			t.Errorf("Unable to unmask %s", unit)
		}
	})

}
func ExampleEnable() {
	unit := "syncthing"
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := Enable(ctx, unit, Options{UserMode: true})
	switch err {
	case ErrMasked:
		fmt.Printf("%s is masked, unmask it before enabling\n", unit)
	case ErrDoesNotExist:
		fmt.Printf("%s does not exist\n", unit)
	case ErrInsufficientPermissions:
		fmt.Printf("permission to enable %s denied\n", unit)
	case ErrBusFailure:
		fmt.Printf("Cannot communicate with the bus\n")
	case nil:
		fmt.Printf("%s enabled successfully\n", unit)
	default:
		fmt.Printf("Error: %v", err)
	}
}

func TestIsActive(t *testing.T) {
	unit := "nginx"
	t.Run(fmt.Sprintf("check active"), func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping in short mode")
		}
		if userString != "root" && userString != "system" {
			t.Skip("skipping superuser test while running as user")
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err := Restart(ctx, unit, Options{UserMode: false})
		if err != nil {
			t.Errorf("Unable to restart %s: %v", unit, err)
		}
		time.Sleep(time.Second)
		isActive, err := IsActive(ctx, unit, Options{UserMode: false})
		if !isActive {
			t.Errorf("IsActive didn't return true for %s", unit)
		}
	})
	t.Run(fmt.Sprintf("check masked"), func(t *testing.T) {
		if userString != "root" && userString != "system" {
			t.Skip("skipping superuser test while running as user")
		}
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		err := Mask(ctx, unit, Options{UserMode: false})
		if err != nil {
			t.Errorf("Unable to mask %s", unit)
		}
		_, err = IsActive(ctx, unit, Options{UserMode: false})
		if err != nil {
			t.Errorf("error is %v, but should have been %v", err, nil)
		}
		Unmask(ctx, unit, Options{UserMode: false})
	})
	t.Run(fmt.Sprintf("check masked"), func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		_, err := IsActive(ctx, "nonexistant", Options{UserMode: false})
		if err != nil {
			t.Errorf("error is %v, but should have been %v", err, ErrDoesNotExist)
		}
	})

}

func TestIsEnabled(t *testing.T) {
	unit := "nginx"
	userMode := false
	if userString != "root" && userString != "system" {
		userMode = true
		unit = "syncthing"
	}
	t.Run(fmt.Sprintf("check enabled"), func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err := Enable(ctx, unit, Options{UserMode: userMode})
		if err != nil {
			t.Errorf("Unable to enable %s: %v", unit, err)
		}
		isEnabled, err := IsEnabled(ctx, unit, Options{UserMode: userMode})
		if !isEnabled {
			t.Errorf("IsEnabled didn't return true for %s", unit)
		}
	})
	t.Run(fmt.Sprintf("check disabled"), func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		err := Disable(ctx, unit, Options{UserMode: userMode})
		if err != nil {
			t.Errorf("Unable to disable %s", unit)
		}
		isEnabled, err := IsEnabled(ctx, unit, Options{UserMode: userMode})
		if err != nil {
			t.Errorf("Error: %v", err)
		}
		if isEnabled {
			t.Errorf("IsEnabled didn't return false for %s", unit)
		}
		Enable(ctx, unit, Options{UserMode: false})
	})
	t.Run(fmt.Sprintf("check masked"), func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		err := Mask(ctx, unit, Options{UserMode: userMode})
		if err != nil {
			t.Errorf("Unable to mask %s", unit)
		}
		isEnabled, err := IsEnabled(ctx, unit, Options{UserMode: userMode})
		if err != ErrMasked {
			t.Errorf("error is %v, but should have been %v", err, ErrMasked)
		}
		if isEnabled {
			t.Errorf("IsEnabled didn't return false for %s", unit)
		}
		Unmask(ctx, unit, Options{UserMode: userMode})
		Enable(ctx, unit, Options{UserMode: userMode})
	})

}

func TestMask(t *testing.T) {
	errCases := []struct {
		unit      string
		err       error
		opts      Options
		runAsUser bool
	}{
		/* Run these tests only as an unpriviledged user */

		//try nonexistant unit in user mode as user
		{"nonexistant", ErrDoesNotExist, Options{UserMode: true}, true},
		// try existing unit in user mode as user
		{"syncthing", nil, Options{UserMode: true}, true},
		// try nonexisting unit in system mode as user
		{"nonexistant", ErrDoesNotExist, Options{UserMode: false}, true},
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
	for _, tc := range errCases {
		t.Run(fmt.Sprintf("%s as %s", tc.unit, userString), func(t *testing.T) {
			if (userString == "root" || userString == "system") && tc.runAsUser {
				t.Skip("skipping user test while running as superuser")
			} else if (userString != "root" && userString != "system") && !tc.runAsUser {
				t.Skip("skipping superuser test while running as user")
			}
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()
			err := Mask(ctx, tc.unit, tc.opts)
			if err != tc.err {
				t.Errorf("error is %v, but should have been %v", err, tc.err)
			}
			Unmask(ctx, tc.unit, tc.opts)
		})
	}
	t.Run(fmt.Sprintf("test double masking existing"), func(t *testing.T) {
		unit := "nginx"
		userMode := false
		if userString != "root" && userString != "system" {
			userMode = true
			unit = "syncthing"
		}
		opts := Options{UserMode: userMode}
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		err := Mask(ctx, unit, opts)
		if err != nil {
			t.Errorf("error on initial masking is %v, but should have been %v", err, nil)
		}
		err = Mask(ctx, unit, opts)
		if err != nil {
			t.Errorf("error on second masking is %v, but should have been %v", err, nil)
		}
		Unmask(ctx, unit, opts)

	})
	t.Run(fmt.Sprintf("test double masking nonexisting"), func(t *testing.T) {
		unit := "nonexistant"
		userMode := false
		if userString != "root" && userString != "system" {
			userMode = true
		}
		opts := Options{UserMode: userMode}
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		err := Mask(ctx, unit, opts)
		if err != ErrDoesNotExist {
			t.Errorf("error on initial masking is %v, but should have been %v", err, ErrDoesNotExist)
		}
		err = Mask(ctx, unit, opts)
		if err != nil {
			t.Errorf("error on second masking is %v, but should have been %v", err, nil)
		}
		Unmask(ctx, unit, opts)
	})

}

func TestRestart(t *testing.T) {
	unit := "nginx"
	userMode := false
	if userString != "root" && userString != "system" {
		userMode = true
		unit = "syncthing"
	}
	opts := Options{UserMode: userMode}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	restarts, err := GetNumRestarts(ctx, unit, opts)
	if err != nil {
		t.Errorf("issue getting number of restarts for %s: %v", unit, err)
	}
	Start(ctx, unit, opts)
	pid, err := GetPID(ctx, unit, opts)
	if err != nil {
		t.Errorf("issue getting MainPID for %s as %s: %v", unit, userString, err)
	}
	syscall.Kill(pid, syscall.SIGKILL)
	for {
		running, err := IsActive(ctx, unit, opts)
		if err != nil {
			t.Errorf("error asserting %s is up: %v", unit, err)
			break
		} else if running {
			break
		}
	}
	secondRestarts, err := GetNumRestarts(ctx, unit, opts)
	if err != nil {
		t.Errorf("issue getting second reading on number of restarts for %s: %v", unit, err)
	}
	if restarts+1 != secondRestarts {
		t.Errorf("Expected restart count to differ by one, but difference was: %d", secondRestarts-restarts)
	}

}

// Runs through all defined Properties in parallel and checks for error cases
func TestShow(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	unit := "nginx"
	opts := Options{
		UserMode: false,
	}
	for _, x := range properties.Properties {
		t.Run(fmt.Sprintf("show property %s", string(x)), func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()
			t.Parallel()
			_, err := Show(ctx, unit, x, opts)
			if err != nil {
				t.Errorf("error is %v, but should have been %v", err, nil)
			}
		})
	}
}

func TestStart(t *testing.T) {
	unit := "nginx"
	userMode := false
	if userString != "root" && userString != "system" {
		userMode = true
		unit = "syncthing"
	}
	opts := Options{UserMode: userMode}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	Stop(ctx, unit, opts)
	for {
		running, err := IsActive(ctx, unit, opts)
		if err != nil {
			t.Errorf("error asserting %s is up: %v", unit, err)
			break
		} else if !running {
			break
		}
	}
	err := Start(ctx, unit, opts)
	if err != nil {
		t.Errorf("error: %v", err)
	}
	for {
		running, err := IsActive(ctx, unit, opts)
		if err != nil {
			t.Errorf("error asserting %s started: %v", unit, err)
			break
		} else if running {
			break
		}
	}

}

func TestStatus(t *testing.T) {

}

func TestStop(t *testing.T) {

}

func TestUnmask(t *testing.T) {

}
