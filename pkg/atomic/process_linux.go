package atomic

import "os"

func isElevated() bool {
	return os.Getuid() == 0
}
