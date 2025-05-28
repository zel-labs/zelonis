package version

import "runtime/debug"

var gitCommit, gitDate string

const ourPath = "github.com/Zel-labs/zelonis" // Path to our module

type VCSInfo struct {
	Commit string // head commit hash
	Date   string // commit time in YYYYMMDD format
	Dirty  bool
}

func VCS() (VCSInfo, bool) {
	if gitCommit != "" {
		// Use information set by the build script if present.
		return VCSInfo{Commit: gitCommit, Date: gitDate}, true
	}
	if buildInfo, ok := debug.ReadBuildInfo(); ok {
		if buildInfo.Main.Path == ourPath {
			return buildInfoVCS(buildInfo)
		}
	}
	return VCSInfo{}, false
}
