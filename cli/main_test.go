package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var expectedCommands = map[string]string{
	"init <project> [path]": "development",
	"run [flags]":           "development",
	"update":                "management",
	"version":               "management",
}

func TestAntithesisCommand(t *testing.T) {
	t.Run("Antithesis command", func(t *testing.T) {
		t.Parallel()

		antithesis := AntithesisCommand()
		assert.NotNil(t, antithesis, "antithesis command should not be nil")

		groups := antithesis.Groups()
		assert.Len(t, groups, 2, "must be 2 groups")
		assert.Equal(t, "management", groups[0].ID, "first group must be 'management'")
		assert.Equal(t, "development", groups[1].ID, "second group must be 'development'")

		commands := antithesis.Commands()
		foundCommands := make(map[string]string, 0)
		for _, command := range commands {
			if command.Hidden {
				continue
			}
			foundCommands[command.Use] = command.GroupID
		}
		assert.Equal(t, expectedCommands, foundCommands, "all commands should be found in the antithesis command")
	})
}
