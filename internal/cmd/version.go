package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wroog-com/demiurge/internal/build"
	"github.com/wroog-com/demiurge/internal/cmdutil"
)

func NewVersionCmd(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:    "version",
		Short:  "Print the version",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprint(f.IOStreams.Out, Format(build.Version, build.Date))
			return nil
		},
	}
}

func Format(version, buildDate string) string {
	if buildDate != "" {
		return fmt.Sprintf("demi version %s (%s)\n", version, buildDate)
	}
	return fmt.Sprintf("demi version %s\n", version)
}
