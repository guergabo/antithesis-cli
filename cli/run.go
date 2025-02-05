package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func runCommand(c HTTPClient) *cobra.Command {
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
# Run a test with 2 microservices and 3 infrastructure dependencies.
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
			cmd.SilenceUsage = true

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
				"antithesis.config_image":      trimWhitespace(config),
				"antithesis.images":            trimWhitespace(strings.Join(images, ";")),
				"antithesis.description":       description,
				"antithesis.report.recipients": trimWhitespace(strings.Join(emails, ";")),
				"antithesis.duration":          fmt.Sprintf("%d", duration),
			}

			body := &struct {
				Params map[string]string `json:"params"`
			}{params}

			jsonBody, err := json.Marshal(body)
			if err != nil {
				return fmt.Errorf("failed to marshal request body: %v", err)
			}

			req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
			if err != nil {
				return fmt.Errorf("failed to create request: %v", err)
			}

			req.SetBasicAuth(username, password)
			req.Header.Set("Content-Type", "application/json")

			resp, err := c.Do(req)
			if err != nil {
				return fmt.Errorf("failed to send request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != 200 {
				switch resp.StatusCode {
				case 403:
					return fmt.Errorf("access forbidden (HTTP 403): please verify your tenant, username, and password are correct")
				default:
					return fmt.Errorf("unexpected non-200 status code: %d", resp.StatusCode)
				}
			}

			prettyPrintRunOutput(cmd, params, duration)
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
	cmd.Flags().Int32VarP(&duration, "duration", "m", 15, "maximum test runtime in minutes (minimum is 15, the longer the deeper)")
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

func prettyPrintRunOutput(cmd *cobra.Command, params map[string]string, duration_mins int32) {
	name := params["antithesis.test_name"]
	recipients := params["antithesis.report.recipients"]

	duration := (time.Duration(duration_mins) * time.Minute) + (10 * time.Minute)
	finishesApprox := time.Now().Add(duration)

	cmd.Printf("\n%s\n\n",
		SuccessStyle.Render(fmt.Sprintf("Successfully submitted a request to launch test run '%s'!", name)))
	cmd.Printf("You should receive a test report emailed to %s around %s.\n",
		ValueStyle.Render(recipients),
		ValueStyle.Render(finishesApprox.Format("Jan 2 3:04PM")))
	cmd.Printf("(Thats roughly %s from now including setup)\n\n", ValueStyle.Render(duration.String()))
	cmd.Printf("If you encounter any issues, use %s to reach out.\n",
		ValueStyle.Render("Antithesis' discord"))
}

func trimWhitespace(in string) string {
	return strings.ReplaceAll(in, " ", "")
}
