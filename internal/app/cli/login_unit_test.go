//go:build unit

package cli

import (
	"errors"
	"testing"
	"yubigo-pass/test"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/stretchr/testify/assert"
)

func TestLoginShouldValidateCorrectInput(t *testing.T) {
	// given
	var correctInput = make([]textinput.Model, 2)
	correctInput[0].SetValue(test.RandomString())
	correctInput[1].SetValue(test.RandomString())

	// when
	err := validateLoginModelInputs(correctInput)

	// then
	assert.Nil(t, err)
}

func TestLoginShouldNotValidateIncorrectInputWithEmptyPassword(t *testing.T) {
	// given
	var incorrectInput = make([]textinput.Model, 2)
	incorrectInput[0].SetValue(test.RandomString())
	incorrectInput[1].SetValue("")

	// expected
	expectedError := errors.New("username and password cannot be empty")

	// when
	err := validateLoginModelInputs(incorrectInput)

	// then
	assert.EqualError(t, err, expectedError.Error())
}

func TestLoginShouldNotValidateIncorrectInputWithEmptyUsername(t *testing.T) {
	// given
	var incorrectInput = make([]textinput.Model, 2)
	incorrectInput[0].SetValue("")
	incorrectInput[1].SetValue(test.RandomString())

	// expected
	expectedError := errors.New("username and password cannot be empty")

	// when
	err := validateLoginModelInputs(incorrectInput)

	// then
	assert.EqualError(t, err, expectedError.Error())
}

func TestLoginShouldNotValidateIncorrectInputWithEmptyBothFields(t *testing.T) {
	// given
	var incorrectInput = make([]textinput.Model, 2)
	incorrectInput[0].SetValue("")
	incorrectInput[1].SetValue("")

	// expected
	expectedError := errors.New("username and password cannot be empty")

	// when
	err := validateLoginModelInputs(incorrectInput)

	// then
	assert.EqualError(t, err, expectedError.Error())
}
