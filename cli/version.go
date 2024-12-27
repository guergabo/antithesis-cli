package cli

import (
	"runtime/debug"

	"github.com/spf13/cobra"
)

func versionCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "version",
		Long:    "Print the CLI version",
		Short:   "Print the CLI version",
		GroupID: "management",
		Example: `
# Print the CLI version 
antithesis version		
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Printf("antithesis version %s\n", version())
			return nil
		},
	}
}

func version() string {
	version := "dev"
	if info, ok := debug.ReadBuildInfo(); ok {
		switch info.Main.Version {
		case "":
		case "(devel)":
		default:
			version = info.Main.Version
		}

		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				version += " " + setting.Value
			}
		}
	}
	return version
}
