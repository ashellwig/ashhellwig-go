package version

import (
	"flag"
	"runtime"
	"strings"
)

var (
	version      = "v0.1.0"
	metadata     = ""
	gitCommit    = ""
	gitTreeState = ""
)

// BuildInfo is our schema for the build.
type BuildInfo struct {
	Version      string ` json:"version,omitempty"`
	GitCommit    string ` json:"git_commit,omitempty"`
	GitTreeState string ` json:"git_tree_state,omitempty"`
	GoVersion    string ` json:"go_version,omitempty"`
}

// GetVersion returns the version of the application.
func GetVersion() string {
	if metadata == "" {
		return version
	}
	return version + "+" + metadata
}

// GetUserAgent returns a string with the current user agent.
func GetUserAgent() string {
	return "ashellwiggo/" + strings.TrimPrefix(GetVersion(), "v")
}

// Get retrieves the build info for the project.
func Get() BuildInfo {
	v := BuildInfo{
		Version:      GetVersion(),
		GitCommit:    gitCommit,
		GitTreeState: gitTreeState,
		GoVersion:    runtime.Version(),
	}

	if flag.Lookup("test.v") != nil {
		v.GoVersion = ""
	}
	return v
}
