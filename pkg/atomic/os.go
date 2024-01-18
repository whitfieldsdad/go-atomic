package atomic

import "runtime"

type OS struct {
	Type string `json:"type"`
	Arch string `json:"arch"`
}

func GetOS() OS {
	return OS{
		Type: runtime.GOOS,
		Arch: runtime.GOARCH,
	}
}
