package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wroog-com/demiurge/internal/build"
	"github.com/wroog-com/demiurge/internal/cmdutil"
)

func NewRootCmd(f *cmdutil.Factory) *cobra.Command {
	root := &cobra.Command{
		Use:           "demi",
		Short:         "Terminal-native project lifecycle companion",
		Long:          `demi is a terminal-native project lifecycle companion. It surfaces project state and context so you always know where things stand, handles filesystem operations so you never have to remember where things live, and lets you spin up new project versions with no friction.`,
		Version:       versionString(build.Version, build.Date),
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	// Keep framework output (--version, help, usage) on the streams commands write to.
	root.SetOut(f.IOStreams.Out)
	root.SetErr(f.IOStreams.Err)

	root.CompletionOptions.DisableDefaultCmd = true

	// Declared explicitly so no -v shorthand is claimed; the flag is still handled by name.
	root.Flags().Bool("version", false, "Print the version")

	// Wrap flag errors as FlagError so demi.Main prints usage alongside them.
	root.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		return cmdutil.FlagErrorWrap(err)
	})

	root.AddCommand(NewVersionCmd(f))

	return root
}
