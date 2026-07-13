package build

import "runtime/debug"

// Stamped at release via -X ldflags; read only where the root command is built.
var (
	Version = "dev"
	Date    = ""
)

func init() {
	if Version == "dev" {
		if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "(devel)" {
			Version = info.Main.Version
		}
	}
}
