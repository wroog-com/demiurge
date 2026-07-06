package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wroog-com/demiurge/internal/iostreams"
)

func NewRootCmd(ioStreams *iostreams.IOStreams) *cobra.Command {
	root := &cobra.Command{
		Use:           "demi",
		Short:         "Terminal-native project awareness for developers",
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	root.CompletionOptions.DisableDefaultCmd = true

	root.AddCommand(NewVersionCmd(ioStreams))

	return root
}
