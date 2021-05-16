package systemctl

import (
	"context"
	"fmt"
	"os"
	"os/user"
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

func TestEnable(t *testing.T) {
	testCases := []struct {
		unit      string
		err       error
		opts      Options
		runAsUser bool
	}{
		// Run these tests only as a user

		//try nonexistant unit in user mode as user
		{"nonexistant", ErrDoesNotExist, Options{usermode: true}, true},
		// try existing unit in user mode as user
		{"syncthing", nil, Options{usermode: true}, true},
		// try nonexisting unit in system mode as user
		{"nonexistant", ErrInsufficientPermissions, Options{usermode: false}, true},
		// try existing unit in system mode as user
		{"nginx", ErrInsufficientPermissions, Options{usermode: false}, true},

		// Run these tests only as a superuser

		// try nonexistant unit in system mode as system
		{"nonexistant", ErrDoesNotExist, Options{usermode: false}, false},
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
			err := Enable(ctx, tc.unit, tc.opts)
			if err != tc.err {
				t.Errorf("error is %v, but should have been %v", err, tc.err)
			}
		})
	}

}

func TestDisable(t *testing.T) {
	testCases := []struct {
		unit      string
		err       error
		opts      Options
		runAsUser bool
	}{
		// Run these tests only as a user

		//try nonexistant unit in user mode as user
		{"nonexistant", ErrDoesNotExist, Options{usermode: true}, true},
		// try existing unit in user mode as user
		{"syncthing", nil, Options{usermode: true}, true},
		// try nonexisting unit in system mode as user
		{"nonexistant", ErrInsufficientPermissions, Options{usermode: false}, true},
		// try existing unit in system mode as user
		{"nginx", ErrInsufficientPermissions, Options{usermode: false}, true},

		// Run these tests only as a superuser

		// try nonexistant unit in system mode as system
		{"nonexistant", ErrDoesNotExist, Options{usermode: false}, false},
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
			err := Disable(ctx, tc.unit, tc.opts)
			if err != tc.err {
				t.Errorf("error is %v, but should have been %v", err, tc.err)
			}
		})
	}

}

// Runs through all defined properties and checks for error cases
func TestShow(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	unit := "nginx"
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	opts := Options{
		usermode: false,
	}
	for _, x := range properties.Properties {
		t.Run(fmt.Sprintf("show property %s", string(x)), func(t *testing.T) {
			_, err := Show(ctx, unit, x, opts)
			if err != nil {
				t.Errorf("error is %v, but should have been %v", err, nil)
			}
		})
	}
}
