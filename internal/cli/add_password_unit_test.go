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

func TestAddPasswordShouldValidateCorrectInput(t *testing.T) {
	// given
	var correctInput = make([]textinput.Model, 4)
	correctInput[0].SetValue(test.RandomString())
	correctInput[1].SetValue(test.RandomString())
	correctInput[2].SetValue(test.RandomString())
	correctInput[3].SetValue(test.RandomString())

	// when
	err := validateAddPasswordModelInputs(correctInput, nil)

	// then
	assert.Nil(t, err)
}

func TestAddPasswordShouldNotValidateIncorrectInputWithEmptyField(t *testing.T) {
	// given
	var incorrectInput = make([]textinput.Model, 4)
	incorrectInput[0].SetValue(test.RandomString())
	incorrectInput[1].SetValue("")
	incorrectInput[2].SetValue(test.RandomString())
	incorrectInput[3].SetValue(test.RandomString())

	// expected
	expectedError := errors.New("only optional fields can be empty")

	// when
	err := validateAddPasswordModelInputs(incorrectInput, nil)

	// then
	assert.EqualError(t, err, expectedError.Error())
}

func TestAddPasswordShouldReturnErrorIfErrorWasPassed(t *testing.T) {
	// given
	var correctInput = make([]textinput.Model, 4)
	correctInput[0].SetValue(test.RandomString())
	correctInput[1].SetValue(test.RandomString())
	correctInput[2].SetValue(test.RandomString())
	correctInput[3].SetValue(test.RandomString())
	passedError := fmt.Errorf("example error")

	// when
	err := validateAddPasswordModelInputs(correctInput, passedError)

	// then
	assert.EqualError(t, err, passedError.Error())
}
