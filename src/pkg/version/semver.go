package version

import (
	"fmt"
	"github.com/coreos/go-semver/semver"
)

var (
	// VersionString gets defined by the build system.
	VersionString = "0.0.0"

	// VersionMajor is the current major version.
	VersionMajor int64

	// VersionMinor is the current minor version.
	VersionMinor int64 = 1

	// VersionPatch is the current patch version.
	VersionPatch int64

	// VersionPre indicates a pre release tag.
	VersionPre = "alpha"

	// VersionDev indicates the current commit.
	VersionDev = "0000000"

	// VersionDate indicates the build date.
	VersionDate = "20190214"

	// Version is the version of the current implementation.
	Version *semver.Version
)

func init() {
	Version = semver.New(VersionString)
	Version.PreRelease = semver.PreRelease(VersionPre)
	Version.Metadata = fmt.Sprintf("git%s.%s", VersionDate, VersionDev)
}
