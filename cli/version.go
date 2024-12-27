package cli

import (
	"encoding/json"
	"fmt"
	"net/http"
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
			// By default, GoReleaser will set the following 3 ldflags: https://goreleaser.com/cookbooks/using-main.version/
			version = info.Main.Version
		}

		// TODO: not neeed
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				version += " " + setting.Value
			}
		}
	}
	return strings.TrimPrefix(version, "v")
}

func latestVersion() (string, error) {
	resp, err := http.Get("https://api.github.com/repos/guergabo/antithesis-cli/releases/latest")
	if err != nil {
		return "", fmt.Errorf("failed to fetch latest release: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch latest release: HTTP %d", resp.StatusCode)
	}

	release := struct {
		TagName string `json:"tag_name"`
	}{}

	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", fmt.Errorf("failed to decode release response: %w", err)
	}

	// Normalize by removing 'v' prefix if present
	version := strings.TrimPrefix(release.TagName, "v")
	return version, nil
}
