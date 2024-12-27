package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

// TODO: mention must build and push images first.

func runCommand() *cobra.Command {
	var (
		name        string
		notebook    string // for free trial will always be basic_test
		tenant      string
		description string
		username    string
		password    string
		config      string
		images      []string
		duration    int32 // TODO: (min 15 - in minutes).
		emails      []string
	)

	// TODO: environment variable to do deal with updating...for github action too...

	cmd := &cobra.Command{
		Use:     "run [flags]",
		Long:    "Run an antithesis test.",
		Short:   "Run an antithesis test with the specified configuration/",
		GroupID: "antithesis",
		Example: `
# Run a test.
antithesis run \
  --name='quickstart' \
  --description='Running a quick antithesis test.' \
  --tenant='tenant' \
  --username='username' \
  --password='password' \
  --config='us-central1-docker.pkg.dev/molten-verve-216720/ant-pdogfood-repository/config:v1' \
  --image='us-central1-docker.pkg.dev/molten-verve-216720/ant-pdogfood-repository/order:v1' \
  --image='us-central1-docker.pkg.dev/molten-verve-216720/ant-pdogfood-repository/payment:v1' \
  --image='us-central1-docker.pkg.dev/molten-verve-216720/ant-pdogfood-repository/test-template:v1' \
  --image='docker.io/postgres:16' \
  --image='docker.io/nats:latest' \
  --image='docker.io/stripemock/stripe-mock:latest' \
  --duration=15 \
  --email='gguergabo@gmail.com'`,
		RunE: func(cmd *cobra.Command, args []string) error {
			url := fmt.Sprintf("https://%s.antithesis.com/api/v1/launch_experiment/%s", tenant, notebook)

			// TODO: trim white space.
			params := map[string]string{
				// "antithesis.integrations.type": "cli", (cause error with validting parameter, but just stops)
				"antithesis.test_name":         name,
				"antithesis.config_image":      config,
				"antithesis.images":            strings.Join(images, ";"),
				"antithesis.description":       description,
				"antithesis.report.recipients": strings.Join(emails, ";"),
				"antithesis.duration":          fmt.Sprintf("%d", duration),
			}

			// fmt.Printf("%v\n", params)
			// return nil

			body := &struct {
				Params map[string]string `json:"params"`
			}{params}

			jsonBody, err := json.Marshal(body)
			if err != nil {
				return fmt.Errorf("failed to marshal request body: %v\n", err)
			}

			req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
			if err != nil {
				return fmt.Errorf("failed to create request: %v\n", err)
			}

			req.SetBasicAuth(username, password)
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				return fmt.Errorf("failed to send request: %v\n", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != 200 {
				return fmt.Errorf("received non-200 response: %d", resp.StatusCode)
			}

			prettyPrintTestParameters(params)

			return nil
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "unique identifier for this test run")
	cmd.Flags().StringVarP(&description, "description", "d", "", "description explaining the purpose of this test run")
	cmd.Flags().StringVarP(&notebook, "notebook", "b", "basic_test", "notebook to execute")
	cmd.Flags().StringVarP(&tenant, "tenant", "t", "", "target tenant ID for test execution")
	cmd.Flags().StringVarP(&username, "username", "u", "", "authentication username for accessing test resources")
	cmd.Flags().StringVarP(&password, "password", "p", "", "authentication password for accessing test resources")
	cmd.Flags().StringVarP(&config, "config", "c", "", "url of configuration image containing docker-compose setup")
	cmd.Flags().StringArrayVarP(&images, "image", "i", make([]string, 0), "list of image URLs to process during testing (can specify multiple)")
	cmd.Flags().Int32VarP(&duration, "duration", "m", 15, "maximum test runtime in minutes (minimum is 15)")
	cmd.Flags().StringArrayVarP(&emails, "email", "e", make([]string, 0), "email addresses to notify with test results (can specify multiple)")

	requiredFlags := []string{
		"name",
		"tenant",
		"username",
		"password",
		"config",
		"image",
		"email",
	}

	for _, flag := range requiredFlags {
		_ = cmd.MarkFlagRequired(flag)
	}

	return cmd
}

// TODO: improve.
func prettyPrintTestParameters(params map[string]string) {
	name := params["antithesis.test_name"]
	recipients := params["antithesis.report.recipients"]
	duration := params["antithesis.duration"]

	fmt.Printf("\nSuccessfully submitted a request to launch test %s\n\n",
		lipgloss.NewStyle().Foreground(magentaColor).Render(name+"!"))
	fmt.Printf("\nYou should get a test report emailed to %s in about %s minutes (plus ~10 min to initialize the test).\n",
		lipgloss.NewStyle().Foreground(magentaColor).Render(recipients),
		lipgloss.NewStyle().Bold(true).Render(duration))
	fmt.Printf("If you encounter any issues, use the %s command to reach out.\n",
		lipgloss.NewStyle().Bold(true).Render("antithesis contact"))
}
