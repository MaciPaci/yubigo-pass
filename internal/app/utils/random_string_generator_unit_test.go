//go:build unit

package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomStringWithLength(t *testing.T) {
	// given
	testCases := []struct {
		name   string
		length int
	}{
		{
			"should generate random string with length 1",
			1,
		},
		{
			"should generate random string with length 16",
			16,
		},
		{
			"should generate random string with length 32",
			32,
		},
		{
			"should generate random string with length 64",
			64,
		},
		{
			"should generate random string with length 128",
			128,
		},
		{
			"should generate random string with length 1000000",
			1000000,
		},
	}

	for _, testCase := range testCases {
		t.Run(
			testCase.name, func(t *testing.T) {
				// when
				generatedString, err := RandomStringWithLength(testCase.length)

				// then
				assert.Nil(t, err)
				assert.Equal(t, len(generatedString), testCase.length)
			})

	}
}
