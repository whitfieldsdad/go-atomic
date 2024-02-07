//go:build darwin || linux
// +build darwin linux

package atomic

import "golang.org/x/sys/unix"

func getOSVersion() string {
	var uts unix.Utsname
	if err := unix.Uname(&uts); err != nil {
		return ""
	}
	return unix.ByteSliceToString(uts.Release[:])
}
