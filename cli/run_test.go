package cli

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockHttpClient struct {
	resp *http.Response
	err  error
}

func NewMockHttpClient(resp *http.Response, err error) HTTPClient {
	return &MockHttpClient{
		resp,
		err,
	}
}

func (c *MockHttpClient) Do(req *http.Request) (*http.Response, error) {
	return c.resp, c.err
}

type TestCase struct {
	name       string
	args       []string
	mockClient HTTPClient
	expected   string
	wantErr    bool
}

func TestRunCommand(t *testing.T) {
	tcs := []TestCase{
		// {
		// 	name: "Valid request",
		// 	args: []string{
		// 		"--name=quickstart",
		// 		"--description=desc",
		// 		"--tenant=tenant",
		// 		"--username=user",
		// 		"--password=pass",
		// 		"--config=config",
		// 		"--image=image1",
		// 		"--image=image2",
		// 		"--duration=30",
		// 		"--email=email1@gmail.com",
		// 	},
		// 	mockClient: NewMockHttpClient(&http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(""))}, nil),
		// 	expected: "\nSuccessfully submitted a request to launch test run 'quickstart'!\n\n\n" +
		// 		"You should get a test report emailed to email1@gmail.com in about 30 minutes " +
		// 		"(plus ~10 min to initialize your environment). So--set your timer!\n\n" +
		// 		"If you encounter any issues, use Antithesis' discord to reach out.\n",
		// 	wantErr: false,
		// },
		{
			name: "Missing required fields",
			args: []string{
				"--name=quickstart",
				"--duration=30",
			},
			wantErr: true,
		},
		{
			name: "Invalid duration",
			args: []string{
				"--name=quickstart",
				"--email=test@example.com",
				"--duration=-1",
			},
			wantErr: true,
		},
		{
			name: "Invalid email",
			args: []string{
				"--name=quickstart",
				"--description=desc",
				"--tenant=tenant",
				"--username=user",
				"--password=pass",
				"--config=config",
				"--image=image1",
				"--image=image2",
				"--duration=30",
				"--email=email1",
			},
			wantErr: true,
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			run := runCommand(tc.mockClient)
			stdout := &bytes.Buffer{}
			run.SetOut(stdout)
			run.SetArgs(tc.args)

			err := run.Execute()
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, stdout.String())
		})
	}
}
