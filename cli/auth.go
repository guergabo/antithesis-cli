package cli

import "github.com/spf13/cobra"

func authCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "auth",
		Long:    "Authenticate with Antithesis",
		Short:   "Authenticate with Antithesis",
		GroupID: "management",
		Example: `
# Authenticate with Antithesis 
antithesis auth [login | logout | whoami]		
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Printf("NOT IMPLEMENTED.\n")
			return nil
		},
	}
}
