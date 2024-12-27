package cli

import "github.com/spf13/cobra"

func contactCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "contact",
		Long:    "Reach out to the developers of Antithesis for help or feedback",
		Short:   "Get help or give feedback",
		GroupID: "management",
		Example: `
# Contact Antithesis developers 
antithesis contact [bookmeeting | feedback]		
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Printf("NOT IMPLEMENTED.\n")
			return nil
		},
	}
}
