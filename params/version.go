package params

import "fmt"

const (
	VersionMajor = 1          // Major version component of the current release
	VersionMinor = 14         // Minor version component of the current release
	VersionPatch = 9          // Patch version component of the current release
	VersionMeta  = "unstable" // Version metadata to append to the version string
)

var VersionWithMeta = func() string {
	v := Version
	if VersionMeta != "" {
		v += "-" + VersionMeta
	}
	return v
}()

// Version holds the textual version string.
var Version = func() string {
	return fmt.Sprintf("%d.%d.%d", VersionMajor, VersionMinor, VersionPatch)
}()

func VersionWithCommit(gitCommit, gitDate string) string {
	vsn := VersionWithMeta
	if len(gitCommit) >= 8 {
		vsn += "-" + gitCommit[:8]
	}
	if (VersionMeta != "stable") && (gitDate != "") {
		vsn += "-" + gitDate
	}
	return vsn
}
