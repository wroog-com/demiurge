package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wroog-com/demiurge/internal/build"
	"github.com/wroog-com/demiurge/internal/iostreams"
)

func NewVersionCmd(ioStreams *iostreams.IOStreams) *cobra.Command {
	return &cobra.Command{
		Use:    "version",
		Short:  "Print the version",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if build.Date != "" {
				fmt.Fprintf(ioStreams.Out, "demi version %s (%s)\n", build.Version, build.Date)
			} else {
				fmt.Fprintf(ioStreams.Out, "demi version %s\n", build.Version)
			}
			return nil
		},
	}
}
