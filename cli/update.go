package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// TODO: update pre-command.

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
				fmt.Errorf("failed to get latest version: %w", err)
			}

			fmt.Printf("Current version: %s, latest version: %s\n", current, latest)
			if current == "dev" {
				fmt.Printf("You're compiling from source.\n")
				return nil
			}
			if current >= latest {
				fmt.Printf("version %s is already latest\n", version)
			}

			var confirm string
			fmt.Printf("Do you want to perform the update to version %s?\n", latest)
			fmt.Printf("Only 'yes' will be accepted to approve.\n")
			fmt.Printf("Enter a value: ")
			fmt.Scanln(&confirm)

			if strings.ToLower(confirm) != "yes" {
				return nil
			}

			err = updateCLI()
			if err != nil {
				return err
			}

			fmt.Printf("Antithesis has been sucessfully upgraded to %s\n", latest)
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

	// Homewbrew's installation prefix.
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
