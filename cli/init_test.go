package cli

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitCommand(t *testing.T) {
	t.Run("Initialize in existing directory", func(t *testing.T) {
		t.Parallel()

		tempDir := t.TempDir()
		init := initCommand()
		stdout := &bytes.Buffer{}
		init.SetOut(stdout)
		init.SetArgs([]string{"quickstart", tempDir})

		err := init.Execute()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expectedStdout := fmt.Sprintf("Downloading project quickstart...\nProject quickstart was created in %s\n", tempDir)
		assert.Equal(t, expectedStdout, stdout.String())
	})
}
