package cli

import "github.com/spf13/cobra"

func debugCommand() *cobra.Command {
	return &cobra.Command{
		Hidden:  true,
		Use:     "debug",
		Long:    "Start multiverse debugging session",
		Short:   "Start multiverse debugging session",
		GroupID: "development",
		Example: `
# Start multiverse debugging sessions 
antithesis debug		
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Printf("NOT IMPLEMENTED.\n")
			return nil
		},
	}
}
