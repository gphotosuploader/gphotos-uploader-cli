package cli_test

import (
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/cli"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestNewCommand(t *testing.T) {
	t.Run("Should return error when using --silent and --debug at the same time", func(t *testing.T) {
		cmd := cli.NewCommand()
		cmd.SetOut(io.Discard)
		cmd.SetErr(io.Discard)
		cmd.SetArgs([]string{"version", "--silent", "--debug"})

		assert.Error(t, cmd.Execute())
	})

	t.Run("Should return success when using --silent", func(t *testing.T) {
		// TODO: Assert that nothing is written to the output when using --silent.
		cmd := cli.NewCommand()
		cmd.SetOut(io.Discard)
		cmd.SetErr(io.Discard)
		cmd.SetArgs([]string{"version", "--silent"})

		assert.NoError(t, cmd.Execute())
	})

	t.Run("Should return success when using --debug", func(t *testing.T) {
		// TODO: Assert that log message is written when using --debug.
		cmd := cli.NewCommand()
		cmd.SetOut(io.Discard)
		cmd.SetErr(io.Discard)
		cmd.SetArgs([]string{"version", "--debug"})

		assert.NoError(t, cmd.Execute())
	})
}
