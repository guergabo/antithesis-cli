package cli

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	hashi_version "github.com/hashicorp/go-version"
	"github.com/spf13/cobra"
)

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
			cmd.SilenceUsage = true

			current := version()
			latest, err := latestVersion()
			if err != nil {
				return fmt.Errorf("failed to get latest version: %w", err)
			}

			cmd.Println(SubtleStyle.Render(fmt.Sprintf("Current version: %s, latest version: %s", current, latest)))
			if current == "dev" {
				cmd.Printf("You're compiling from source.\n")
				return nil
			}

			c, err := hashi_version.NewVersion(current)
			if err != nil {
				return fmt.Errorf("failed to get current version: %w", err)
			}
			l, err := hashi_version.NewVersion(latest)
			if err != nil {
				return fmt.Errorf("failed to get latest version: %w", err)
			}
			if c.GreaterThanOrEqual(l) {
				cmd.Printf("version %s is already latest\n", ValueStyle.Render(current))
				return nil
			}

			var confirm string
			cmd.Println(WarningStyle.Render(fmt.Sprintf("Do you want to perform the update to version %s?", latest)))
			cmd.Printf("Only %s will be accepted to approve.\n", ValueStyle.Render("'yes'"))
			cmd.Printf("Enter a value: ")
			_, _ = fmt.Scanln(&confirm)

			if strings.ToLower(confirm) != "yes" {
				return nil
			}
			err = updateCLI()
			if err != nil {
				return err
			}
			cmd.Println(SuccessStyle.Render(fmt.Sprintf("Antithesis has been sucessfully upgraded to %s", latest)))
			return nil
		},
	}
}

func isHomebrew() bool {
	binary, err := os.Executable()
	if err != nil {
		return false
	}
	brewExe, err := exec.LookPath("brew")
	if err != nil {
		return false
	}
	brewPrefixBS, err := exec.Command(brewExe, "--prefix").Output()
	if err != nil {
		return false
	}
	brewBinPrefix := filepath.Join(strings.TrimSpace(string(brewPrefixBS)), "bin") + string(filepath.Separator)
	return strings.HasPrefix(binary, brewBinPrefix)
}

func updateCLI() error {
	updateCommand := "brew update && brew upgrade antithesis"
	if !isHomebrew() {
		// TODO: support non-homebrew installation (curl, nix).
		return fmt.Errorf("currently only support automatic updates with homebrew installations")
	}
	cmd := exec.Command("sh", "-c", updateCommand)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run update command: %w", err)
	}
	return nil
}

func latestVersion() (string, error) {
	resp, err := http.Get("https://api.github.com/repos/guergabo/homebrew-antithesis/releases/latest")
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
	return strings.TrimPrefix(release.TagName, "v"), nil
}
