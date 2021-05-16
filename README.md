# systemctl

This library aims at providing idiomatic `systemctl` bindings for go developers, in order to make it easier to write system tooling using golang.
This tool tries to take guesswork out of arbitrarily shelling out to `systemctl` by providing a structured, throuroughly-tested wrapper for the `systemctl` functions most-likely to be used in a system program.

If your system isn't running (or targeting another system running) `systemctl`, this library will be of little use to you.
In fact, if `systemctl` isn't found in the `PATH`, this library will panic.

## What is systemctl

`systemctl` is a command-line program which grants the user control over the systemd system and service manager.

`systemctl` may be used to introspect and control the state of the "systemd" system and service manager. Please refer to `systemd(1)` for an introduction into the basic concepts and functionality this tool manages.

## Supported systemctl functions

- [ ] `systemctl daemon-reload`
- [x] `systemctl disable`
- [x] `systemctl enable`
- [ ] `systemctl is-active`
- [ ] `systemctl is-enabled`
- [ ] `systemctl is-failed`
- [ ] `systemctl mask`
- [ ] `systemctl restart`
- [ ] `systemctl show`
- [ ] `systemctl start`
- [ ] `systemctl status`
- [ ] `systemctl stop`
- [ ] `systemctl unmask`

## Helper functionality

- [x] Get start time of a service (`ExecMainStartTimestamp`) as a `Time` type
- [x] Get current memory in bytes (`MemoryCurrent`) an an int
- [x] Get the PID of the main process (`MainPID`) as an int


## Useful errors

All functions return a predefinied error type, and it is highly recommended these errors are handled properly.

## Context support

All calls into this library support go's `context` functionality.
Therefore, blocking calls can time out according to the callee's needs, and the returned error should be checked to see if a timeout occurred (`ErrExecTimeout`).

## TODO

- [ ] Add additional bindings for systemctl options I (the author) don't use frequently (or ever) for others to use.
- [ ] Set up `go test` testing

## Simple example

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/taigrr/systemctl"
)

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    // Equivalent to `systemctl enable nginx` with a 10 second timeout
    opts := Options{
        usermode: false,
    }
    err := Enable(ctx, unit, opts)
    if err != nil {
        log.Fatalf("unable to enable unit %s: %v", "nginx", err)
    }
}
```

## License

This project is licensed under the 0BSD License, written by [Rob Landley](https://github.com/landley).
As such, you may use this library without restriction or attribution, but please don't pass it off as your own.
Attribution, though not required, is appreciated.

By contributing, you agree all code submitted also falls under the License.

## External resources

- [Official systemctl documentation](https://www.man7.org/linux/man-pages/man1/systemctl.1.html)
- [Inspiration for this repo](https://github.com/Ullaakut/nmap/)
