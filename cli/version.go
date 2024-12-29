package cli

import (
	"runtime/debug"
	"strings"

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
			cmd.Printf("antithesis version %s\n", ValueStyle.Render(version()))
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
			// By default, GoReleaser will set the following 3 ldflags: https://goreleaser.com/cookbooks/using-main.version/
			version = info.Main.Version
		}
	}
	return strings.TrimPrefix(version, "v")
}
