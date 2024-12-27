package cli

// TODO:
import "github.com/spf13/cobra"

func initCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "init <template> [path]",
		Long:    "Initialize a new Antithesis project",
		Short:   "Initialize a new Antithesis project",
		GroupID: "antithesis",
		Example: `
# Initialize a new Antithesis project 
antithesis init quickstart .		
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Printf("NOT IMPLEMENTED.\n")
			return nil
		},
	}
}
