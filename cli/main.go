package cli

import (
	"context"

	"github.com/spf13/cobra"
)

// TODO: error handling.
// TODO: contact command.
// TODO: colors align.
// TODO: make init less nested.
// TODO: better latestVersion() to map with brew release.
// TODO: port all 3 repos to antithesishq and update 'guergabo' and url stuff.
func AntithesisCommand() *cobra.Command {
	cmd := &cobra.Command{
		Version: version(),
		Use:     "antithesis",
		Long: ValueStyle.Render(`
█████╗ ███╗   ██╗████████╗██╗████████╗██╗  ██╗███████╗███████╗██╗███████╗
██╔══██╗████╗  ██║╚══██╔══╝██║╚══██╔══╝██║  ██║██╔════╝██╔════╝██║██╔════╝
███████║██╔██╗ ██║   ██║   ██║   ██║   ███████║█████╗  ███████╗██║███████╗
██╔══██║██║╚██╗██║   ██║   ██║   ██║   ██╔══██║██╔══╝  ╚════██║██║╚════██║
██║  ██║██║ ╚████║   ██║   ██║   ██║   ██║  ██║███████╗███████║██║███████║
╚═╝  ╚═╝╚═╝  ╚═══╝   ╚═╝   ╚═╝   ╚═╝   ╚═╝  ╚═╝╚══════╝╚══════╝╚═╝╚══════╝
																					
The entrypoint of the antithesis ecosystem. Build the impossible.`),
		Short: "Antithesis CLI",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Environment variables, config.
			// Eagerly inform customers of when a new update is available.
			if cmd.Name() == "update" {
				return nil
			}
			current := version()
			latest, _ := latestVersion() // Suppress error to not block usage.
			if current == "dev" || current >= latest {
				return nil
			}
			cmd.Printf("A new update is available. To install it, run %s\n", SubtleStyle.Render("'antithesis update'"))
			return nil
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
		ID:    "development",
		Title: "Development Commands:",
	})

	cmd.AddCommand(authCommand())
	cmd.AddCommand(configCommand())
	cmd.AddCommand(contactCommand())
	cmd.AddCommand(updateCommand())
	cmd.AddCommand(versionCommand())
	cmd.AddCommand(debugCommand())
	cmd.AddCommand(initCommand())
	cmd.AddCommand(runCommand())

	return cmd
}

func Main() error {
	return AntithesisCommand().ExecuteContext(context.Background())
}
