package version

import (
	"fmt"
	"runtime"
)

const (
	versionFormat = "Build version: %s; Commit: %s; Build date: %s; Go version: %s; OS/Arch: %s/%s"
)

type BuildVersion struct {
	Version   string
	Commit    string
	BuildDate string
}

func (b BuildVersion) String() string {
	return fmt.Sprintf(versionFormat, b.Version, b.Commit, b.BuildDate, runtime.Version(), runtime.GOOS, runtime.GOARCH)
}
