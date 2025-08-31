// Package luhn implements the Luhn algorithm for validation of credit card numbers
// and other identification numbers.
package luhn

import (
	"strings"
	"unicode"
)

// Validate checks if a string passes the Luhn algorithm.
// It accepts numbers with spaces and dashes, which are ignored during validation.
// Returns true if the number is valid according to the Luhn algorithm, false otherwise.
func Validate(number string) bool {
	// Clean the input: remove spaces and dashes
	cleaned := strings.ReplaceAll(number, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")

	// Check minimum length (must have more than 1 digit)
	if len(cleaned) <= 1 {
		return false
	}

	// Verify all characters are digits
	for _, r := range cleaned {
		if !unicode.IsDigit(r) {
			return false
		}
	}

	return checksum(cleaned)
}

// checksum performs the actual Luhn algorithm calculation
func checksum(number string) bool {
	sum := 0
	double := false

	// Process digits from right to left
	for i := len(number) - 1; i >= 0; i-- {
		digit := int(number[i] - '0')

		if double {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
		double = !double
	}

	return sum%10 == 0
}

// Generate creates a valid Luhn number by appending a check digit
// to the provided number string.
func Generate(number string) string {
	// Clean the input
	cleaned := strings.ReplaceAll(number, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")

	// Check for empty input
	if cleaned == "" {
		return ""
	}

	// Verify all characters are digits
	for _, r := range cleaned {
		if !unicode.IsDigit(r) {
			return ""
		}
	}

	// Calculate what check digit would make this valid
	// Append a '0' temporarily and calculate
	tempNumber := cleaned + "0"
	sum := 0
	double := true // Start with true because we're starting from second-to-last

	for i := len(tempNumber) - 2; i >= 0; i-- {
		digit := int(tempNumber[i] - '0')

		if double {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
		double = !double
	}

	// Calculate check digit
	checkDigit := (10 - (sum % 10)) % 10

	return cleaned + string(rune('0'+checkDigit))
}