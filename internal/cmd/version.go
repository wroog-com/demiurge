package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wroog-com/demiurge/internal/cmdutil"
)

func NewVersionCmd(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version",
		RunE: func(cmd *cobra.Command, args []string) error {
			root := cmd.Root()
			fmt.Fprintf(f.IOStreams.Out, "%s version %s\n", root.DisplayName(), root.Version)
			return nil
		},
	}
}

// versionString strips the v prefix some install channels report so every
// channel prints the same form.
func versionString(version, date string) string {
	version = strings.TrimPrefix(version, "v")
	if date == "" {
		return version
	}
	return fmt.Sprintf("%s (%s)", version, date)
}
