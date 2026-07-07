package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wroog-com/demiurge/internal/cmdutil"
)

func NewRootCmd(f *cmdutil.Factory) *cobra.Command {
	root := &cobra.Command{
		Use:           "demi",
		Short:         "Terminal-native project awareness for developers",
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	root.CompletionOptions.DisableDefaultCmd = true

	// Wrap flag errors as FlagError so demi.Main prints usage alongside them.
	root.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		return cmdutil.FlagErrorWrap(err)
	})

	root.AddCommand(NewVersionCmd(f))

	return root
}
