//go:build go1.18
// +build go1.18

package luhn

import (
	"testing"
	"unicode"
)

func FuzzValidate(f *testing.F) {
	// Añadir corpus semilla - varios números de tarjeta de crédito válidos e inválidos
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
		// La función no debe entrar en pánico
		result := Validate(input)

		// Propiedad: el resultado debe ser determinístico
		result2 := Validate(input)
		if result != result2 {
			t.Errorf("Resultado no determinístico para la entrada %q: primero=%v, segundo=%v", input, result, result2)
		}

		// Propiedad: si es válido, debe tener al menos 2 dígitos
		if result {
			digitCount := 0
			for _, r := range input {
				if unicode.IsDigit(r) {
					digitCount++
				}
			}
			if digitCount <= 1 {
				t.Errorf("Validado con solo %d dígitos: %q", digitCount, input)
			}
		}

		// Propiedad: si contiene caracteres que no son dígitos, espacios o guiones, debe ser inválido
		hasInvalidChar := false
		for _, r := range input {
			if !unicode.IsDigit(r) && r != ' ' && r != '-' {
				hasInvalidChar = true
				break
			}
		}
		if hasInvalidChar && result {
			t.Errorf("Validado a pesar de caracteres inválidos: %q", input)
		}

		// Propiedad: vacío o un solo dígito siempre debe ser inválido
		cleaned := ""
		for _, r := range input {
			if unicode.IsDigit(r) {
				cleaned += string(r)
			}
		}
		if len(cleaned) <= 1 && result {
			t.Errorf("Validado con %d dígitos (debería necesitar al menos 2): %q", len(cleaned), input)
		}
	})
}

func FuzzGenerate(f *testing.F) {
	// Corpus semilla
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
		// La función no debe entrar en pánico
		result := Generate(input)

		// Si la entrada contiene no-dígitos (excepto espacios y guiones), el resultado debe estar vacío
		hasInvalidChar := false
		for _, r := range input {
			if !unicode.IsDigit(r) && r != ' ' && r != '-' {
				hasInvalidChar = true
				break
			}
		}

		if hasInvalidChar {
			if result != "" {
				t.Errorf("Generate debería retornar vacío para entrada inválida %q, obtuvo %q", input, result)
			}
			return
		}

		// Si obtuvimos un resultado, debe ser válido
		if result != "" {
			if !Validate(result) {
				t.Errorf("Generate(%q) produjo un número inválido: %q", input, result)
			}

			// El resultado debe empezar con los dígitos limpios de la entrada
			cleaned := ""
			for _, r := range input {
				if unicode.IsDigit(r) {
					cleaned += string(r)
				}
			}

			if cleaned != "" && len(result) != len(cleaned)+1 {
				t.Errorf("Generate(%q) debería añadir exactamente un dígito, obtuvo %q", input, result)
			}
		}
	})
}

func FuzzRoundTrip(f *testing.F) {
	// Prueba que Generate seguido de Validate siempre retorna true
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
		// Limpiar entrada para obtener solo dígitos
		cleaned := ""
		for _, r := range input {
			if unicode.IsDigit(r) {
				cleaned += string(r)
			}
		}

		if cleaned == "" {
			return // Saltar entradas vacías
		}

		// Generar un número válido
		generated := Generate(cleaned)
		if generated == "" {
			return // La generación falló, lo cual está bien para entrada inválida
		}

		// El número generado debe ser válido
		if !Validate(generated) {
			t.Errorf("Prueba de ida y vuelta falló: Generate(%q) = %q, pero Validate retornó false", cleaned, generated)
		}
	})
}