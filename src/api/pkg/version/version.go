// Package version provides build version information.
package version

import (
	"fmt"
	"runtime"
)

// The following fields are populated at build date using -ldflags -X.
// Note that DATE is omitted for reproducible builds
var (
	buildVersion  = "unknown"
	buildUser     = "unknown"
	buildHost     = "unknown"
	buildStatus   = "unknown"
	buildDate     = "unknown"
	buildCompiler = "unknown"
	buildPlatform = "unknown"
)

// BuildInfo describes version information about the binary build.
type BuildInfo struct {
	Version     string `json:"version"`
	User        string `json:"user"`
	Host        string `json:"host"`
	GoVersion   string `json:"goVersion"`
	BuildStatus string `json:"status"`
	BuildDate   string `json:"buildDate"`
	Compiler    string `json:"compiler"`
	Platform    string `json:"platform"`
}

var (
	// Info exports the build version information.
	Info BuildInfo
)

// String produces a single-line version info
//
// This looks like:
//
// ```
// user@host-<version>-<build status>-<compiler>-<platform>
// ```
func (b BuildInfo) String() string {
	return fmt.Sprintf("%v@%v-%v-%v-%v-%v-%v",
		b.User,
		b.Host,
		b.Version,
		b.BuildStatus,
		b.BuildDate,
		b.Compiler,
		b.Platform)
}

// LongForm returns a multi-line version information
//
// This looks like:
//
// ```
// Version: <version>
// User: user@host
// GoVersion: go1.10.2
// BuildStatus: <build status>
// BuildDate: <build date>
// Compiler: <compiler>
// Platform: <platform>
// ```
func (b BuildInfo) LongForm() string {
	return fmt.Sprintf(`
	  Version: %v
		User: %v@%v
		GoVersion: %v
		BuildStatus: %v
		BuildDate: %v
		Compiler: %v
		Platform: %v
		`,
		b.Version,
		b.User,
		b.Host,
		b.GoVersion,
		b.BuildStatus,
		b.BuildDate,
		b.Compiler,
		b.Platform)
}

func init() {
	Info = BuildInfo{
		Version:     Version.String(),
		User:        buildUser,
		Host:        buildHost,
		GoVersion:   runtime.Version(),
		BuildStatus: buildStatus,
		BuildDate:   buildDate,
		Compiler:    runtime.Compiler,
		Platform:    fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}
