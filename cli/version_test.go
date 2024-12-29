package cli

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	versionOutput = "antithesis version dev\n"
)

func TestVersionCommand(t *testing.T) {
	t.Run("Print CLI version", func(t *testing.T) {
		t.Parallel()

		version := versionCommand()
		stdout := &bytes.Buffer{}
		version.SetOut(stdout)

		err := version.Execute()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		assert.Equal(t, versionOutput, stdout.String())
	})
}
