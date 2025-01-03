package cli

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"
)

const (
	antithesisDir = "antithesis"
)

var (
	availableProjects = map[string]string{
		"quickstart": "https://github.com/guergabo/quickstarts/tarball/main",
	}
)

// TODO: optimizate to update the project only if the latest commit SHA is different.
func initCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "init <project> [path]",
		Long:    "Initialize an Antithesis demo project. This command downloads and sets up a preconfigured project structure, allowing you to quickly start experimenting with Antithesis. You can initialize the project in the current directory or specify a custom path.",
		Short:   "Initialize an Antithesis demo project",
		GroupID: "development",
		Example: `
# Initialize in current directory
antithesis init quickstart .		

# Initialize in a new directory
antithesis init quickstart ./output

# Initialize with absolute path
antithesis init quickstart /Users/username/projects/output
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			if len(args) == 0 {
				cmd.Print(cmd.UsageString())
				return nil
			}

			project := args[0]
			projectURL, ok := availableProjects[project]
			if !ok {
				keys := strings.Join(maps.Keys(availableProjects), "\n  - ")
				return fmt.Errorf("Project %q is not supported.\n\nAvailable projects:\n  - %s", project, keys)
			}

			cmd.Println(SubtleStyle.Render(fmt.Sprintf("Downloading project %s...", project)))

			// Create temp directory to download project.

			cfg, err := getUserConfigDir()
			if err != nil {
				return fmt.Errorf("failed to get user config directory: %w", err)
			}
			projectTempDir, err := os.MkdirTemp(cfg, "antithesis-*")
			if err != nil {
				return fmt.Errorf("failed to create temp dir: %w", err)
			}
			defer os.RemoveAll(projectTempDir)

			err = downloadAndExtractProject(projectURL, projectTempDir)
			if err != nil {
				return fmt.Errorf("Failed to download and extract quickstart: %w", err)
			}

			// Check if directory is provided and empty. Defaults to current directory.

			directory := "."
			if len(args) > 1 {
				directory = args[1]
				exists, err := directoryExists(directory)
				if err != nil {
					return fmt.Errorf("failed to check if directory exists: %w", err)
				}
				if !exists {
					err := os.MkdirAll(directory, 0755)
					if err != nil {
						return fmt.Errorf("failed to create directory %v: %w", directory, err)
					}
				}
			}
			isEmpty, err := isDirectoryEmpty(directory)
			if err != nil {
				return fmt.Errorf("failed to check if directory is empty: %w", err)
			}
			if !isEmpty {
				return fmt.Errorf("Could not create project in %s because directory is not empty", ValueStyle.Render(fmt.Sprintf("'%s'", directory)))
			}

			// Copy over files.

			path, err := filepath.Abs(directory)
			if err != nil {
				return fmt.Errorf("failed to get absolute path of directory: %w", err)
			}
			targetPath := filepath.Join(path, project)
			err = os.CopyFS(targetPath, os.DirFS(projectTempDir))
			if err != nil {
				return fmt.Errorf("failed to copy directory: %w", err)
			}
			cmd.Println(SuccessStyle.Render(fmt.Sprintf("Project %s was created in %s", project, path)))
			return nil
		},
	}
}

func directoryExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}

func isDirectoryEmpty(path string) (bool, error) {
	dir, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer dir.Close()
	_, err = dir.Readdirnames(1)
	if err == nil {
		return false, nil
	}
	if err == io.EOF {
		return true, nil
	}
	return false, err
}

func downloadAndExtractProject(url, dir string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download project: %s", resp.Status)
	}
	return untar(resp.Body, dir)
}

func untar(r io.Reader, dst string) error {
	gzr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	var topLevelDir string

	for {
		header, err := tr.Next()
		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		case header == nil:
			continue
		}

		// We don't want the top-level directory created by GitHub.
		if topLevelDir == "" {
			parts := strings.Split(header.Name, "/")
			if len(parts) > 1 {
				topLevelDir = parts[0]
			}
		}
		headNameWithoutTop := strings.TrimPrefix(header.Name, topLevelDir+"/")
		target := filepath.Join(dst, headNameWithoutTop)

		switch header.Typeflag {
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return err
			}
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			defer f.Close()
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}
		}
	}
}

func getUserConfigDir() (string, error) {
	var (
		configDir string
		err       error
	)
	switch runtime.GOOS {
	case "windows":
		configDir, err = os.UserConfigDir()
		if err != nil {
			return "", err
		}
	case "darwin":
		configDir, err = os.UserConfigDir()
		if err != nil {
			return "", err
		}
	// Linux.
	default:
		configDir = os.Getenv("XDG_CONFIG_HOME")
		if configDir == "" {
			configDir, err = os.UserConfigDir()
			if err != nil {
				return "", err
			}
		}
	}
	// Lazily create it (if needed).
	antithesisUserConfigDir := filepath.Join(configDir, antithesisDir)
	err = os.MkdirAll(antithesisUserConfigDir, 0755)
	if err != nil {
		return "", err
	}
	return antithesisUserConfigDir, nil
}
