package resolver

import (
	"runtime"
)

var downloadURLMap = map[string]map[string]string{
	"darwin": {
		"amd64": "https://github.com/actions/runner/releases/download/v2.311.0/actions-runner-osx-x64-2.311.0.tar.gz",
		"arm64": "https://github.com/actions/runner/releases/download/v2.311.0/actions-runner-osx-arm64-2.311.0.tar.gz",
	},
	"linux": {
		"amd64": "https://github.com/actions/runner/releases/download/v2.311.0/actions-runner-linux-x64-2.311.0.tar.gz",
		"arm":   "https://github.com/actions/runner/releases/download/v2.311.0/actions-runner-linux-arm-2.311.0.tar.gz",
		"arm64": "https://github.com/actions/runner/releases/download/v2.311.0/actions-runner-linux-arm64-2.311.0.tar.gz",
	},
	"windows": {
		"amd64": "https://github.com/actions/runner/releases/download/v2.311.0/actions-runner-win-x64-2.311.0.zip -OutFile actions-runner-win-x64-2.311.0.zip",
		"arm64": "https://github.com/actions/runner/releases/download/v2.311.0/actions-runner-win-arm64-2.311.0.zip -OutFile actions-runner-win-arm64-2.311.0.zip",
	},
}

func RunnerDownloadURL() string {
	return downloadURLMap[runtime.GOOS][runtime.GOARCH]
}
