package cli

import "github.com/spf13/cobra"

func updateCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "update",
		Long:    "Update the CLI to the latest version",
		Short:   "Update the CLI to the latest version",
		GroupID: "management",
		Example: `
# Update the CLI to the latest version 
antithesis update		
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Printf("NOT IMPLEMENTED.\n")
			return nil
		},
	}
}
