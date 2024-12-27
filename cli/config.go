package cli

import "github.com/spf13/cobra"

func configCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "config",
		Long:    "Manage your CLI configuration",
		Short:   "Manage your CLI configuration",
		GroupID: "management",
		Example: `
# Manage your CLI configuration 
antithesis config [path | set]		
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Printf("NOT IMPLEMENTED.\n")
			return nil
		},
	}
}
