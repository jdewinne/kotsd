package cli

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ExecuteVersionCommand(t *testing.T) {

	actual := new(bytes.Buffer)
	cmd := RootCmd()
	cmd.SetOut(actual)
	cmd.SetErr(actual)
	cmd.SetArgs([]string{"--version"})
	cmd.Execute()

	expected := "kotsd version 0.0.1\n"

	assert.Equal(t, expected, actual.String(), "actual is not expected")
}
