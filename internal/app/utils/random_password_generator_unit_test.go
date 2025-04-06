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

// --- Test GetStrengthStyle ---

func TestGetStrengthStyle(t *testing.T) {
	testCases := []struct {
		name          string
		score         int
		expectedColor lipgloss.TerminalColor
	}{
		{name: "Score 0 Very Weak", score: 0, expectedColor: lipgloss.Color("9")},
		{name: "Score 1 Weak", score: 1, expectedColor: lipgloss.Color("208")},
		{name: "Score 2 Fair", score: 2, expectedColor: lipgloss.Color("11")},
		{name: "Score 3 Good", score: 3, expectedColor: lipgloss.Color("10")},
		{name: "Score 4 Strong", score: 4, expectedColor: lipgloss.Color("82")},
		{name: "Score -1 Below Range", score: -1, expectedColor: lipgloss.NoColor{}},
		{name: "Score 5 Above Range", score: 5, expectedColor: lipgloss.NoColor{}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			score := tc.score

			// then
			var actualColor lipgloss.TerminalColor = lipgloss.NoColor{}
			switch score {
			case 0:
				actualColor = lipgloss.Color("9")
			case 1:
				actualColor = lipgloss.Color("208")
			case 2:
				actualColor = lipgloss.Color("11")
			case 3:
				actualColor = lipgloss.Color("10")
			case 4:
				actualColor = lipgloss.Color("82")
			}
			assert.Equal(t, tc.expectedColor, actualColor)
		})
	}
}

// --- Test GetStrengthText ---

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
