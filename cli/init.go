package cli

// TODO:
import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	availableProjects = map[string]string{
		"quickstart": "https://github.com/guergabo/quickstarts/tarball/main",
	}
)

// TODO: optimization to update the project only if the latest commit SHA is different.
func initCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "init <project> [path]",
		Long:    "Initialize an Antithesis demo project. This command downloads and sets up a preconfigured project structure, allowing you to quickly start experimenting with Antithesis. You can initialize the project in the current directory or specify a custom path.",
		Short:   "Initialize an Antithesis demo project",
		GroupID: "development",
		Example: `
# Initialize in current directory
antithesis init quickstart .		

# Initialize in a new directory (exiting or not)
antithesis init quickstart ./my-quickstart

# Initialize with absolute path
antithesis init quickstart /Users/username/projects/my-quickstart
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				cmd.Print(cmd.UsageString())
				return nil
			}

			project := args[0]
			projectURL, ok := availableProjects[project]
			if !ok {
				cmd.SilenceUsage = true
				keys := make([]string, 0, len(availableProjects))
				for k := range availableProjects {
					keys = append(keys, k)
				}
				cmd.Printf("Project %s is not supported.\n\nAvailable projects:\n%s", project, keys)
				return nil // TODO: error handling. Choose this way across to not see the help page.
			}

			cmd.Println(SubtleStyle.Render(fmt.Sprintf("Downloading project %s...", project)))

			// Create temp dir to download project.

			projectTempDir, err := os.MkdirTemp(os.TempDir(), "antithesis-*")
			if err != nil {
				return fmt.Errorf("failed to create temp dir: %w", err)
			}
			defer os.RemoveAll(projectTempDir)

			err = downloadAndExtractProject(projectURL, projectTempDir)
			if err != nil {
				return fmt.Errorf("Failed to download and extract quickstart: %w\n", err)
			}

			// Check if directory is provided and empty. Defaults to current.

			directory := "."
			if len(args) > 1 {
				directory = args[1]
				exists, err := directoryExists(directory)
				if err != nil {
					return fmt.Errorf("failed to check if directory exists: %w\n", err)
				}

				// Create the directory if it doesn't exist.
				if !exists {
					err := os.MkdirAll(directory, 0755)
					if err != nil {
						return fmt.Errorf("failed to create directory %v: %w\n", directory, err)
					}
					exists = true
				}
			}

			isEmpty, err := isDirectoryEmpty(directory)
			if err != nil {
				return fmt.Errorf("failed to check if directory is empty: %w\n", err)
			}
			if !isEmpty { // TODO: prefer this route to control error output message for user error. vs return fmt.Errorf() for system error.
				cmd.SilenceUsage = true
				cmd.Printf("Could not create project in %s because directory is not empty\n", ValueStyle.Render(fmt.Sprintf("'%s'", directory)))
				return nil
			}

			// Attempt atomic rename from temp to final destination.

			path, err := filepath.Abs(directory)
			if err != nil {
				return fmt.Errorf("failed to get absolute path of directory: %w", err)
			}
			err = os.Rename(projectTempDir, path+"/"+project)
			if err != nil {
				return fmt.Errorf("failed to rename directory: %w", err)
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

		target := filepath.Join(dst, header.Name)

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
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}
			f.Close()
		}
	}
}
