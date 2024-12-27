package cli

import (
	"context"

	"github.com/spf13/cobra"
)

func AntithesisCommand() *cobra.Command {
	cmd := &cobra.Command{
		Version: version(),
		Use:     "antithesis",
		Long:    "The entrypoint of the antithesis ecosystem.",
		Short:   "Antithesis CLI",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil // Environment variables + Check for latest release.
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddGroup(&cobra.Group{
		ID:    "management",
		Title: "Management Commands:",
	})
	cmd.AddGroup(&cobra.Group{
		ID:    "antithesis",
		Title: "Development Commands:",
	})

	// cmd.AddCommand(authCommand())
	// cmd.AddCommand(configCommand())
	cmd.AddCommand(contactCommand())
	cmd.AddCommand(updateCommand())
	cmd.AddCommand(versionCommand())
	// cmd.AddCommand(debugCommand())
	cmd.AddCommand(initCommand())
	cmd.AddCommand(runCommand())

	return cmd
}

func Main() error {
	return AntithesisCommand().ExecuteContext(context.Background())
}
