package main

import (
	"fmt"

	"github.com/wroog-com/demiurge/internal/build"
)

func main() {
	if build.Date != "" {
		fmt.Printf("demi version %s (%s)\n", build.Version, build.Date)
	} else {
		fmt.Printf("demi version %s\n", build.Version)
	}
}
