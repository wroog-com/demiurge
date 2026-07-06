package cmd

import "github.com/spf13/cobra"

func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:          "demi",
		Short:        "Terminal-native project awareness for developers",
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	root.AddCommand(NewVersionCmd())

	return root
}
