package build

import "runtime/debug"

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
