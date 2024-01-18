//go:build !windows && !js && !darwin
// +build !windows,!js,!darwin

package atomic

import "syscall"

func getSysProcAttrs() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{
		Setsid: true,
	}
}