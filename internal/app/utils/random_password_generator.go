package utils

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// --- Password Generation Constants ---
const (
	lowercaseChars = "abcdefghijklmnopqrstuvwxyz"
	uppercaseChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digitChars     = "0123456789"
	symbolChars    = "!@#$%^&*()-_=+[]{}|;:,.<>/?"
	DefaultLength  = 20 // Default password length
)

// GeneratePassword creates a random password with specified criteria.
func GeneratePassword(length int, includeLower, includeUpper, includeDigits, includeSymbols bool) (string, error) {
	if length <= 0 {
		length = DefaultLength
	}

	var charSet string
	if includeLower {
		charSet += lowercaseChars
	}
	if includeUpper {
		charSet += uppercaseChars
	}
	if includeDigits {
		charSet += digitChars
	}
	if includeSymbols {
		charSet += symbolChars
	}

	if charSet == "" {
		return "", errors.New("no character sets selected for password generation")
	}

	password := make([]byte, length)
	max := big.NewInt(int64(len(charSet)))

	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", fmt.Errorf("failed to generate random number: %w", err)
		}
		password[i] = charSet[num.Int64()]
	}

	missingChars := map[bool]string{
		includeLower:   lowercaseChars,
		includeUpper:   uppercaseChars,
		includeDigits:  digitChars,
		includeSymbols: symbolChars,
	}

	for include, chars := range missingChars {
		if include && !strings.ContainsAny(string(password), chars) {
			password = addMissingCharacters(password, length, chars)
		}
	}

	// Fisher-Yates shuffle
	for i := range password {
		jBig, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			return "", fmt.Errorf("failed to generate random index for shuffle: %w", err)
		}
		j := jBig.Int64()
		password[i], password[j] = password[j], password[i]
	}

	return string(password), nil
}

func addMissingCharacters(password []byte, length int, characterSet string) []byte {
	idx, _ := rand.Int(rand.Reader, big.NewInt(int64(length)))
	charIdx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(characterSet))))
	password[idx.Int64()] = lowercaseChars[charIdx.Int64()]
	return password
}

// GetStrengthStyle returns the style for the password strength score.
func GetStrengthStyle(score int) lipgloss.Style {
	switch score {
	case 0:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	case 1:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("208"))
	case 2:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	case 3:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	case 4:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("82"))
	default:
		return lipgloss.NewStyle()
	}
}

// GetStrengthText returns the text representation of the password strength score.
func GetStrengthText(score int) string {
	switch score {
	case 0:
		return "Very Weak"
	case 1:
		return "Weak"
	case 2:
		return "Fair"
	case 3:
		return "Good"
	case 4:
		return "Strong"
	default:
		return ""
	}
}
