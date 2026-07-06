package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wroog-com/demiurge/internal/build"
)

func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version",
		Run: func(cmd *cobra.Command, args []string) {
			if build.Date != "" {
				fmt.Printf("demi version %s (%s)\n", build.Version, build.Date)
			} else {
				fmt.Printf("demi version %s\n", build.Version)
			}
		},
	}
}
