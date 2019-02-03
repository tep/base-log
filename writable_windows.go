package log

import "os"

func isWritable(dir string) bool {
	if fi, _ := os.Stat(dir); fi != nil {
		return fi.Mode().Perm()&(1<<(uint(7))) != 0
	}
	return false
}
