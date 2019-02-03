// +build !windows

package log

import (
	"golang.org/x/sys/unix"
)

func isWritable(dir string) bool {
	if err := unix.Access(dir, unix.W_OK); err == nil {
		return true
	}
	return false
}
