package atomic

import (
	"os"
)

const (
	Linux   = "linux"
	Darwin  = "darwin"
	Windows = "windows"
)

var (
	DefaultAtomicsPath = os.Getenv("ATOMICS_PATH")
)
