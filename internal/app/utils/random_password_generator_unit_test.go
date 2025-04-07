//go:build unit

package utils

import (
	"fmt"
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGeneratePassword_Defaults(t *testing.T) {
	testCases := []struct {
		name           string
		inputLength    int
		expectedLength int
	}{
		{"Length 0", 0, DefaultLength},
		{"Length -5", -5, DefaultLength},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			length := tc.inputLength
			includeLower := true // Need at least one set

			// when
			password, err := GeneratePassword(length, includeLower, false, false, false)

			// then
			require.NoError(t, err)
			assert.Len(t, password, tc.expectedLength)
		})
	}
}

func TestGeneratePassword_NoCharsets(t *testing.T) {
	// given
	length := 10
	includeLower := false
	includeUpper := false
	includeDigits := false
	includeSymbols := false

	// when
	_, err := GeneratePassword(length, includeLower, includeUpper, includeDigits, includeSymbols)

	// then
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no character sets selected")
}

func TestGeneratePassword_CharsetInclusionAndLength(t *testing.T) {
	testCases := []struct {
		name           string
		length         int
		includeLower   bool
		includeUpper   bool
		includeDigits  bool
		includeSymbols bool
		expectedLength int
		checkLower     bool
		checkUpper     bool
		checkDigits    bool
		checkSymbols   bool
		allowedChars   string
	}{
		{
			name:           "All Charsets Default Length",
			length:         DefaultLength,
			includeLower:   true,
			includeUpper:   true,
			includeDigits:  true,
			includeSymbols: true,
			expectedLength: DefaultLength,
			checkLower:     true,
			checkUpper:     true,
			checkDigits:    true,
			checkSymbols:   true,
			allowedChars:   lowercaseChars + uppercaseChars + digitChars + symbolChars,
		},
		{
			name:           "Specific Length All Charsets",
			length:         30,
			includeLower:   true,
			includeUpper:   true,
			includeDigits:  true,
			includeSymbols: true,
			expectedLength: 30,
			checkLower:     true,
			checkUpper:     true,
			checkDigits:    true,
			checkSymbols:   true,
			allowedChars:   lowercaseChars + uppercaseChars + digitChars + symbolChars,
		},
		{
			name:           "Only Lowercase",
			length:         15,
			includeLower:   true,
			includeUpper:   false,
			includeDigits:  false,
			includeSymbols: false,
			expectedLength: 15,
			checkLower:     true,
			checkUpper:     false,
			checkDigits:    false,
			checkSymbols:   false,
			allowedChars:   lowercaseChars,
		},
		{
			name:           "Only Uppercase",
			length:         15,
			includeLower:   false,
			includeUpper:   true,
			includeDigits:  false,
			includeSymbols: false,
			expectedLength: 15,
			checkLower:     false,
			checkUpper:     true,
			checkDigits:    false,
			checkSymbols:   false,
			allowedChars:   uppercaseChars,
		},
		{
			name:           "Only Digits",
			length:         15,
			includeLower:   false,
			includeUpper:   false,
			includeDigits:  true,
			includeSymbols: false,
			expectedLength: 15,
			checkLower:     false,
			checkUpper:     false,
			checkDigits:    true,
			checkSymbols:   false,
			allowedChars:   digitChars,
		},
		{
			name:           "Only Symbols",
			length:         15,
			includeLower:   false,
			includeUpper:   false,
			includeDigits:  false,
			includeSymbols: true,
			expectedLength: 15,
			checkLower:     false,
			checkUpper:     false,
			checkDigits:    false,
			checkSymbols:   true,
			allowedChars:   symbolChars,
		},
		{
			name:           "Lower and Digits",
			length:         25,
			includeLower:   true,
			includeUpper:   false,
			includeDigits:  true,
			includeSymbols: false,
			expectedLength: 25,
			checkLower:     true,
			checkUpper:     false,
			checkDigits:    true,
			checkSymbols:   false,
			allowedChars:   lowercaseChars + digitChars,
		},
		{
			name:           "Upper and Symbols",
			length:         25,
			includeLower:   false,
			includeUpper:   true,
			includeDigits:  false,
			includeSymbols: true,
			expectedLength: 25,
			checkLower:     false,
			checkUpper:     true,
			checkDigits:    false,
			checkSymbols:   true,
			allowedChars:   uppercaseChars + symbolChars,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			length := tc.length
			includeLower := tc.includeLower
			includeUpper := tc.includeUpper
			includeDigits := tc.includeDigits
			includeSymbols := tc.includeSymbols

			// when
			password, err := GeneratePassword(length, includeLower, includeUpper, includeDigits, includeSymbols)

			// then
			require.NoError(t, err)
			assert.Len(t, password, tc.expectedLength)

			if tc.checkLower {
				assert.True(t, containsAny(password, lowercaseChars))
			}
			if tc.checkUpper {
				assert.True(t, containsAny(password, uppercaseChars))
			}
			if tc.checkDigits {
				assert.True(t, containsAny(password, digitChars))
			}
			if tc.checkSymbols {
				assert.True(t, containsAny(password, symbolChars))
			}

			for _, char := range password {
				assert.True(t, strings.ContainsRune(tc.allowedChars, char), fmt.Sprintf("Password contains disallowed char '%c'", char))
			}

			if len(password) > 1 {
				firstChar := password[0]
				allSame := true
				for i := 1; i < len(password); i++ {
					if password[i] != firstChar {
						allSame = false
						break
					}
				}
				assert.False(t, allSame)
			}
		})
	}
}

func containsAny(s, chars string) bool {
	return strings.ContainsAny(s, chars)
}

func TestGetStrengthText(t *testing.T) {
	testCases := []struct {
		name         string
		score        int
		expectedText string
	}{
		{name: "Score 0", score: 0, expectedText: "Very Weak"},
		{name: "Score 1", score: 1, expectedText: "Weak"},
		{name: "Score 2", score: 2, expectedText: "Fair"},
		{name: "Score 3", score: 3, expectedText: "Good"},
		{name: "Score 4", score: 4, expectedText: "Strong"},
		{name: "Score -1 Below Range", score: -1, expectedText: ""},
		{name: "Score 5 Above Range", score: 5, expectedText: ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			score := tc.score

			// when
			text := GetStrengthText(score)

			// then
			assert.Equal(t, tc.expectedText, text)
		})
	}
}

func TestGetStrengthStyle(t *testing.T) {
	testCases := []struct {
		name          string
		score         int
		expectedStyle lipgloss.Style
	}{
		{name: "Score 0", score: 0, expectedStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("9"))},
		{name: "Score 1", score: 1, expectedStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("208"))},
		{name: "Score 2", score: 2, expectedStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("11"))},
		{name: "Score 3", score: 3, expectedStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("10"))},
		{name: "Score 4", score: 4, expectedStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("82"))},
		{name: "Score -1 Below Range", score: -1, expectedStyle: lipgloss.NewStyle()},
		{name: "Score 5 Above Range", score: 5, expectedStyle: lipgloss.NewStyle()},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			score := tc.score

			// when
			style := GetStrengthStyle(score)

			// then
			assert.Equal(t, tc.expectedStyle, style)
		})
	}
}

func TestShouldAddMissingCharacters(t *testing.T) {
	testCases := []struct {
		name         string
		password     string
		characterSet string
	}{
		{name: "Lowercase", password: "aaaaa", characterSet: lowercaseChars},
		{name: "Uppercase", password: "AAAAA", characterSet: uppercaseChars},
		{name: "Digits", password: "11111", characterSet: digitChars},
		{name: "Symbols", password: "!!!!!", characterSet: symbolChars},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			changedPassword := addMissingCharacters([]byte(tc.password), len(tc.password), tc.characterSet)

			assert.NotEqual(t, tc.password, changedPassword, "Password slice should have been modified")
		})
	}
}
