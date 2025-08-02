package main

import (
	"fmt"
	"runtime/debug"
)

var version = "dev"

func Version() string {
	commit := "unknown"
	date := "unknown"
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			switch setting.Key {
			case "vcs.revision":
				commit = setting.Value[:8] // Short commit hash
			case "vcs.time":
				date = setting.Value
			}
		}
	}

	return fmt.Sprintf("%s-%s @ %s", version, commit, date)
}
