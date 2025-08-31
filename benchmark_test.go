package luhn

import (
	"fmt"
	"testing"
)

// Basic benchmarks
func BenchmarkValidate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Validate("4532015112830366")
	}
}

func BenchmarkValidateWithSpaces(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Validate("4532 0151 1283 0366")
	}
}

func BenchmarkValidateWithDashes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Validate("4532-0151-1283-0366")
	}
}

func BenchmarkValidateInvalid(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Validate("1234567890123456")
	}
}

// Benchmark different card number lengths
func BenchmarkValidateSizes(b *testing.B) {
	sizes := []struct {
		name   string
		number string
	}{
		{"14digits", "30569309025904"},        // Diners Club
		{"15digits", "371449635398431"},       // American Express
		{"16digits", "4532015112830366"},      // Visa
		{"16digitsSpaces", "4532 0151 1283 0366"},
		{"16digitsDashes", "4532-0151-1283-0366"},
		{"19digits", "6011000990139424543"},   // Some Discover cards
	}

	for _, size := range sizes {
		b.Run(size.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				Validate(size.number)
			}
		})
	}
}

// Benchmark the Generate function
func BenchmarkGenerate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Generate("453201511283036")
	}
}

func BenchmarkGenerateWithSpaces(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Generate("4532 0151 1283 036")
	}
}

// Benchmark with memory allocation tracking
func BenchmarkValidateAllocs(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		Validate("4532015112830366")
	}
}

func BenchmarkValidateWithSpacesAllocs(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		Validate("4532 0151 1283 0366")
	}
}

// Parallel benchmark to test concurrent performance
func BenchmarkValidateParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Validate("4532015112830366")
		}
	})
}

// Table-driven benchmark with setup
func BenchmarkValidateTable(b *testing.B) {
	testCases := []string{
		"4532015112830366",
		"5425233430109903",
		"371449635398431",
		"6011000990139424",
	}

	b.ResetTimer() // Reset timer after setup

	for i := 0; i < b.N; i++ {
		for _, tc := range testCases {
			Validate(tc)
		}
	}
}

// Benchmark to compare performance of valid vs invalid numbers
func BenchmarkValidateComparison(b *testing.B) {
	benchmarks := []struct {
		name   string
		number string
		valid  bool
	}{
		{"ValidVisa", "4532015112830366", true},
		{"InvalidVisa", "4532015112830367", false},
		{"ValidWithSpaces", "4532 0151 1283 0366", true},
		{"InvalidWithSpaces", "4532 0151 1283 0367", false},
		{"ShortValid", "59", true},
		{"ShortInvalid", "58", false},
		{"InvalidLetters", "4532a15112830366", false},
	}

	for _, bm := range benchmarks {
		b.Run(fmt.Sprintf("%s_%v", bm.name, bm.valid), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				result := Validate(bm.number)
				if result != bm.valid {
					b.Fatalf("unexpected result: got %v, want %v", result, bm.valid)
				}
			}
		})
	}
}