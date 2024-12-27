package cli

import "github.com/spf13/cobra"

func debugCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "debug",
		Long:    "Start multiverse debugging session",
		Short:   "Start multiverse debugging session",
		GroupID: "antithesis",
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
