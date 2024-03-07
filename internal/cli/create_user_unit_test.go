//go:build unit

package cli

import (
	"errors"
	"fmt"
	"testing"
	"yubigo-pass/test"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/stretchr/testify/assert"
)

func TestShouldValidateCorrectInput(t *testing.T) {
	// given
	var correctInput = make([]textinput.Model, 2)
	correctInput[0].SetValue(test.RandomString())
	correctInput[1].SetValue(test.RandomString())

	// when
	err := validateCreateUserModelInputs(correctInput, nil)

	// then
	assert.Nil(t, err)
}

func TestShouldNotValidateIncorrectInputWithEmptyPassword(t *testing.T) {
	// given
	var incorrectInput = make([]textinput.Model, 2)
	incorrectInput[0].SetValue(test.RandomString())
	incorrectInput[1].SetValue("")

	// expected
	expectedError := errors.New("password cannot empty")

	// when
	err := validateCreateUserModelInputs(incorrectInput, nil)

	// then
	assert.Error(t, expectedError, err)
}

func TestShouldNotValidateIncorrectInputWithEmptyUsername(t *testing.T) {
	// given
	var incorrectInput = make([]textinput.Model, 2)
	incorrectInput[0].SetValue("")
	incorrectInput[1].SetValue(test.RandomString())

	// expected
	expectedError := errors.New("username cannot empty")

	// when
	err := validateCreateUserModelInputs(incorrectInput, nil)

	// then
	assert.Error(t, expectedError, err)
}

func TestShouldReturnErrorIfErrorWasPassed(t *testing.T) {
	// given
	var correctInput = make([]textinput.Model, 2)
	correctInput[0].SetValue(test.RandomString())
	correctInput[1].SetValue(test.RandomString())
	passedError := fmt.Errorf("example error")

	// when
	err := validateCreateUserModelInputs(correctInput, passedError)

	// then
	assert.Error(t, passedError, err)
}
