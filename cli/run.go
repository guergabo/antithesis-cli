package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"
	"strings"

	"github.com/spf13/cobra"
)

func runCommand() *cobra.Command {
	var (
		name        string
		notebook    string
		tenant      string
		description string
		username    string
		password    string
		config      string
		images      []string
		duration    int32
		emails      []string
	)

	cmd := &cobra.Command{
		Use:     "run [flags]",
		Long:    "Run an antithesis test. Note: Before running this command, you must first build and push all required images to either a public container registry or to Antithesis' private registry.",
		Short:   "Run an antithesis test",
		GroupID: "development",
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

			if duration < 15 {
				return fmt.Errorf("duration can't be less than 15.")
			}

			for _, email := range emails {
				_, err := mail.ParseAddress(email)
				if err != nil {
					return fmt.Errorf("email not valid: %w", err)
				}
			}

			params := map[string]string{
				"antithesis.test_name":         name,
				"antithesis.config_image":      trim(config),
				"antithesis.images":            trim(strings.Join(images, ";")),
				"antithesis.description":       description,
				"antithesis.report.recipients": trim(strings.Join(emails, ";")),
				"antithesis.duration":          fmt.Sprintf("%d", duration),
			}

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

			prettyPrintRunOutput(cmd, params)
			return nil
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "unique identifier for this test run")
	cmd.Flags().StringVarP(&description, "description", "d", "", "description explaining the purpose of this test run")
	cmd.Flags().StringVarP(&tenant, "tenant", "t", "", "target tenant ID for test execution")
	cmd.Flags().StringVarP(&username, "username", "u", "", "authentication username for accessing test resources")
	cmd.Flags().StringVarP(&password, "password", "p", "", "authentication password for accessing test resources")
	cmd.Flags().StringVarP(&config, "config", "c", "", "url of configuration image containing docker-compose setup")
	cmd.Flags().StringArrayVarP(&images, "image", "i", make([]string, 0), "list of image URLs to process during testing (can specify multiple)")
	cmd.Flags().StringVarP(&notebook, "notebook", "b", "basic_test", "notebook to execute")
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

func prettyPrintRunOutput(cmd *cobra.Command, params map[string]string) {
	name := params["antithesis.test_name"]
	recipients := params["antithesis.report.recipients"]
	duration := params["antithesis.duration"]

	cmd.Printf("\nSuccessfully submitted a request to launch test %s\n\n",
		ValueStyle.Render(name+"!"))
	cmd.Printf("\nYou should get a test report emailed to %s in about %s minutes (plus ~10 min to initialize your environment).\n",
		ValueStyle.Render(recipients),
		ValueStyle.Render(duration))
	cmd.Printf("If you encounter any issues, use the %s command to reach out.\n",
		SubtleStyle.Render("'antithesis contact'")) // TODO: implement.
}

func trim(in string) string {
	return strings.ReplaceAll(in, " ", "")
}

// TODO: environment variable to do deal with updating and secrets. For GitHub Action too.
