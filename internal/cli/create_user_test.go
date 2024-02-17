//go:build unit

package cli

import (
	"errors"
	"testing"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/stretchr/testify/assert"
)

func TestShouldValidateCorrectInput(t *testing.T) {
	// given
	var correctInput = make([]textinput.Model, 2)
	correctInput[0].SetValue("exampleUsername")
	correctInput[1].SetValue("examplePassword")

	// when
	err := validateInputs(correctInput)

	// then
	assert.Nil(t, err)
}

func TestShouldNotValidateIncorrectInputWithEmptyPassword(t *testing.T) {
	// given
	var incorrectInput = make([]textinput.Model, 2)
	incorrectInput[0].SetValue("exampleUsername")
	incorrectInput[1].SetValue("")

	// expected
	expectedError := errors.New("password cannot empty")

	// when
	err := validateInputs(incorrectInput)

	// then
	assert.Error(t, expectedError, err)
}

func TestShouldNotValidateIncorrectInputWithEmptyUsername(t *testing.T) {
	// given
	var incorrectInput = make([]textinput.Model, 2)
	incorrectInput[0].SetValue("")
	incorrectInput[1].SetValue("examplePassword")

	// expected
	expectedError := errors.New("username cannot empty")

	// when
	err := validateInputs(incorrectInput)

	// then
	assert.Error(t, expectedError, err)
}
