package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

func AntithesisCommand() *cobra.Command {
	cmd := &cobra.Command{
		Version: version(),
		Use:     "antithesis",
		Long: `
█████╗ ███╗   ██╗████████╗██╗████████╗██╗  ██╗███████╗███████╗██╗███████╗
██╔══██╗████╗  ██║╚══██╔══╝██║╚══██╔══╝██║  ██║██╔════╝██╔════╝██║██╔════╝
███████║██╔██╗ ██║   ██║   ██║   ██║   ███████║█████╗  ███████╗██║███████╗
██╔══██║██║╚██╗██║   ██║   ██║   ██║   ██╔══██║██╔══╝  ╚════██║██║╚════██║
██║  ██║██║ ╚████║   ██║   ██║   ██║   ██║  ██║███████╗███████║██║███████║
╚═╝  ╚═╝╚═╝  ╚═══╝   ╚═╝   ╚═╝   ╚═╝   ╚═╝  ╚═╝╚══════╝╚══════╝╚═╝╚══════╝
																					
The entrypoint of the antithesis ecosystem. Build the impossible.`,
		Short: "Antithesis CLI",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Eagerly inform customers of when a new update is available.
			current := version()
			latest, _ := latestVersion() // suppress to not interfere.
			if current == "dev" || current >= latest {
				return nil
			}
			fmt.Printf("A new update is available. To install it, run 'antithesis update'\n")
			// TODO: Environment variables, config.
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
		ID:    "antithesis",
		Title: "Development Commands:",
	})

	cmd.AddCommand(authCommand())   // *
	cmd.AddCommand(configCommand()) // *
	cmd.AddCommand(contactCommand())
	cmd.AddCommand(updateCommand())
	cmd.AddCommand(versionCommand())
	cmd.AddCommand(debugCommand()) // *
	cmd.AddCommand(initCommand())
	cmd.AddCommand(runCommand())

	return cmd
}

func Main() error {
	return AntithesisCommand().ExecuteContext(context.Background())
}
