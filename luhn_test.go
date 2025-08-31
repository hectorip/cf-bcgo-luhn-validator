package luhn

import (
	"testing"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name     string
		number   string
		expected bool
	}{
		// Valid credit card numbers
		{"valid Visa", "4532015112830366", true},
		{"valid Visa with spaces", "4532 0151 1283 0366", true},
		{"valid Visa with dashes", "4532-0151-1283-0366", true},
		{"valid MasterCard", "5425233430109903", true},
		{"valid MasterCard with spaces", "5425 2334 3010 9903", true},
		{"valid American Express", "371449635398431", true},
		{"valid Discover", "6011000990139424", true},
		{"valid Diners Club", "30569309025904", true},
		{"valid JCB", "3530111333300000", true},

		// Invalid numbers
		{"invalid checksum", "4532015112830367", false},
		{"invalid MasterCard", "5425233430109904", false},
		{"single digit", "5", false},
		{"single zero", "0", false},
		{"empty string", "", false},
		{"only spaces", "   ", false},
		{"only dashes", "---", false},
		{"letters in number", "4532a15112830366", false},
		{"special characters", "4532@151#1283$366", false},
		{"mixed invalid chars", "4532-0151-ABCD-0366", false},
		{"too short", "1", false},

		// Edge cases
		{"two digits valid", "59", true},
		{"two digits invalid", "58", false},
		{"all zeros", "0000000000000000", true}, // Actually valid per Luhn algorithm
		{"number with trailing spaces", "4532015112830366  ", true},
		{"number with leading spaces", "  4532015112830366", true},
		{"mixed spaces and dashes", "4532 0151-1283 0366", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Validate(tt.number)
			if result != tt.expected {
				t.Errorf("Validate(%q) = %v, want %v", tt.number, result, tt.expected)
			}
		})
	}
}

func TestChecksum(t *testing.T) {
	// Test the internal checksum function with clean numbers
	testCases := []struct {
		name     string
		number   string
		expected bool
	}{
		{"valid Visa", "4532015112830366", true},
		{"valid MasterCard", "5425233430109903", true},
		{"valid Amex", "371449635398431", true},
		{"invalid number", "1234567890123456", false},
		{"another invalid", "4532015112830367", false},
		{"two digit valid", "59", true},
		{"two digit invalid", "58", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := checksum(tc.number)
			if result != tc.expected {
				t.Errorf("checksum(%q) = %v, want %v", tc.number, result, tc.expected)
			}
		})
	}
}

func TestGenerate(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedSuffix string // The last digit should be this
		shouldBeValid  bool
	}{
		{"generate for partial Visa", "453201511283036", "6", true},
		{"generate for partial MasterCard", "542523343010990", "3", true},
		{"generate with spaces", "4532 0151 1283 036", "6", true},
		{"generate with dashes", "4532-0151-1283-036", "6", true},
		{"invalid input with letters", "4532a151", "", false},
		{"empty input", "", "", false},
		{"single digit", "5", "9", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Generate(tt.input)

			if !tt.shouldBeValid {
				if result != "" {
					t.Errorf("Generate(%q) should return empty string for invalid input, got %q", tt.input, result)
				}
				return
			}

			// Check if the generated number ends with expected suffix
			if len(result) > 0 && result[len(result)-1:] != tt.expectedSuffix {
				t.Errorf("Generate(%q) = %q, expected to end with %s", tt.input, result, tt.expectedSuffix)
			}

			// Verify the generated number is valid
			if !Validate(result) {
				t.Errorf("Generate(%q) = %q, which is not a valid Luhn number", tt.input, result)
			}
		})
	}
}

// Test concurrent access to ensure thread safety
func TestValidateConcurrent(t *testing.T) {
	numbers := []string{
		"4532015112830366",
		"5425233430109903",
		"371449635398431",
		"1234567890123456",
	}

	// Run validation concurrently
	done := make(chan bool)
	for i := 0; i < 100; i++ {
		go func(n int) {
			number := numbers[n%len(numbers)]
			_ = Validate(number)
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 100; i++ {
		<-done
	}
}