// +build go1.18

package luhn

import (
	"testing"
	"unicode"
)

func FuzzValidate(f *testing.F) {
	// Add seed corpus - various valid and invalid credit card numbers
	testcases := []string{
		"4532015112830366",
		"5425233430109903",
		"371449635398431",
		"6011000990139424",
		"",
		"0",
		"1",
		"59",
		"invalid",
		"1234567890123456",
		"4532 0151 1283 0366",
		"4532-0151-1283-0366",
		"  4532015112830366  ",
		"!@#$%^&*()",
		"aaaaaaaaaaaaaaaa",
		"0000000000000000",
		"9999999999999999",
	}

	for _, tc := range testcases {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, input string) {
		// The function should not panic
		result := Validate(input)

		// Property: result should be deterministic
		result2 := Validate(input)
		if result != result2 {
			t.Errorf("Non-deterministic result for input %q: first=%v, second=%v", input, result, result2)
		}

		// Property: if valid, must have at least 2 digits
		if result {
			digitCount := 0
			for _, r := range input {
				if unicode.IsDigit(r) {
					digitCount++
				}
			}
			if digitCount <= 1 {
				t.Errorf("Validated with only %d digits: %q", digitCount, input)
			}
		}

		// Property: if it contains non-digit, non-space, non-dash characters, it should be invalid
		hasInvalidChar := false
		for _, r := range input {
			if !unicode.IsDigit(r) && r != ' ' && r != '-' {
				hasInvalidChar = true
				break
			}
		}
		if hasInvalidChar && result {
			t.Errorf("Validated despite invalid characters: %q", input)
		}

		// Property: empty or single digit should always be invalid
		cleaned := ""
		for _, r := range input {
			if unicode.IsDigit(r) {
				cleaned += string(r)
			}
		}
		if len(cleaned) <= 1 && result {
			t.Errorf("Validated with %d digits (should need at least 2): %q", len(cleaned), input)
		}
	})
}

func FuzzGenerate(f *testing.F) {
	// Seed corpus
	testcases := []string{
		"453201511283036",
		"542523343010990",
		"37144963539843",
		"",
		"1",
		"12345",
		"4532 0151 1283 036",
		"invalid",
	}

	for _, tc := range testcases {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, input string) {
		// The function should not panic
		result := Generate(input)

		// If input contains non-digits (except spaces and dashes), result should be empty
		hasInvalidChar := false
		for _, r := range input {
			if !unicode.IsDigit(r) && r != ' ' && r != '-' {
				hasInvalidChar = true
				break
			}
		}

		if hasInvalidChar {
			if result != "" {
				t.Errorf("Generate should return empty for invalid input %q, got %q", input, result)
			}
			return
		}

		// If we got a result, it should be valid
		if result != "" {
			if !Validate(result) {
				t.Errorf("Generate(%q) produced invalid number: %q", input, result)
			}

			// The result should start with the cleaned input digits
			cleaned := ""
			for _, r := range input {
				if unicode.IsDigit(r) {
					cleaned += string(r)
				}
			}

			if cleaned != "" && len(result) != len(cleaned)+1 {
				t.Errorf("Generate(%q) should append exactly one digit, got %q", input, result)
			}
		}
	})
}

func FuzzRoundTrip(f *testing.F) {
	// Test that Generate followed by Validate always returns true
	testcases := []string{
		"453201511283036",
		"542523343010990",
		"37144963539843",
		"1",
		"12",
		"123",
		"1234567890",
	}

	for _, tc := range testcases {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, input string) {
		// Clean input to only digits
		cleaned := ""
		for _, r := range input {
			if unicode.IsDigit(r) {
				cleaned += string(r)
			}
		}

		if cleaned == "" {
			return // Skip empty inputs
		}

		// Generate a valid number
		generated := Generate(cleaned)
		if generated == "" {
			return // Generation failed, which is OK for invalid input
		}

		// The generated number must be valid
		if !Validate(generated) {
			t.Errorf("Round trip failed: Generate(%q) = %q, but Validate returned false", cleaned, generated)
		}
	})
}