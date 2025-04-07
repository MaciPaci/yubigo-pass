//go:build unit

package cli

import (
	"testing"
	"yubigo-pass/test"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestInput() textinput.Model {
	ti := textinput.New()
	return ti
}

func TestAddPasswordShouldValidateCorrectInput(t *testing.T) {
	inputs := make([]textinput.Model, 4)
	for i := range inputs {
		inputs[i] = newTestInput()
	}
	inputs[0].SetValue(test.RandomString())
	inputs[1].SetValue(test.RandomString())
	inputs[2].SetValue(test.RandomString())
	inputs[3].SetValue(test.RandomString())

	err := validateAddPasswordModelInputs(inputs)

	assert.NoError(t, err, "Validation should pass for correct input")
}

func TestAddPasswordShouldValidateCorrectInputWithEmptyURL(t *testing.T) {
	inputs := make([]textinput.Model, 4)
	for i := range inputs {
		inputs[i] = newTestInput()
	}
	inputs[0].SetValue(test.RandomString())
	inputs[1].SetValue(test.RandomString())
	inputs[2].SetValue(test.RandomString())
	inputs[3].SetValue("")

	err := validateAddPasswordModelInputs(inputs)

	assert.NoError(t, err, "Validation should pass with empty URL")
}

func TestAddPasswordShouldNotValidateEmptyRequiredField(t *testing.T) {
	expectedErrorMsg := "title, username, and password fields cannot be empty"

	testCases := []struct {
		name        string
		title       string
		username    string
		password    string
		expectedErr string
	}{
		{
			name:        "Empty Title",
			title:       "",
			username:    test.RandomString(),
			password:    test.RandomString(),
			expectedErr: expectedErrorMsg,
		},
		{
			name:        "Empty Username",
			title:       test.RandomString(),
			username:    "",
			password:    test.RandomString(),
			expectedErr: expectedErrorMsg,
		},
		{
			name:        "Empty Password",
			title:       test.RandomString(),
			username:    test.RandomString(),
			password:    "",
			expectedErr: expectedErrorMsg,
		},
		{
			name:        "Empty Title and Username",
			title:       "",
			username:    "",
			password:    test.RandomString(),
			expectedErr: expectedErrorMsg,
		},
		{
			name:        "Empty All Required",
			title:       "",
			username:    "",
			password:    "",
			expectedErr: expectedErrorMsg,
		},
		{
			name:        "Whitespace Title",
			title:       "   ",
			username:    test.RandomString(),
			password:    test.RandomString(),
			expectedErr: expectedErrorMsg,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			inputs := make([]textinput.Model, 4)
			for i := range inputs {
				inputs[i] = newTestInput()
			}
			inputs[0].SetValue(tc.title)
			inputs[1].SetValue(tc.username)
			inputs[2].SetValue(tc.password)
			inputs[3].SetValue(test.RandomString())

			err := validateAddPasswordModelInputs(inputs)

			require.Error(t, err, "Expected an error for empty required field")
			assert.EqualError(t, err, tc.expectedErr, "Error message mismatch")
		})
	}
}
