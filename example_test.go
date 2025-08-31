package luhn_test

import (
	"fmt"
	"github.com/example/luhn-validator"
)

func ExampleValidate() {
	// Validar una tarjeta Visa
	valid := luhn.Validate("4532015112830366")
	fmt.Println(valid)
	// Output: true
}

func ExampleValidate_withSpaces() {
	// Los espacios son ignorados durante la validación
	valid := luhn.Validate("4532 0151 1283 0366")
	fmt.Println(valid)
	// Output: true
}

func ExampleValidate_withDashes() {
	// Los guiones también son ignorados
	valid := luhn.Validate("4532-0151-1283-0366")
	fmt.Println(valid)
	// Output: true
}

func ExampleValidate_invalid() {
	// Un número con checksum inválido
	valid := luhn.Validate("1234567890123456")
	fmt.Println(valid)
	// Output: false
}

func ExampleValidate_tooShort() {
	// Los números de un solo dígito no son válidos
	valid := luhn.Validate("5")
	fmt.Println(valid)
	// Output: false
}

func ExampleGenerate() {
	// Generar un número válido agregando el dígito de control
	number := luhn.Generate("453201511283036")
	fmt.Println(number)
	fmt.Println(luhn.Validate(number))
	// Output:
	// 4532015112830366
	// true
}

func ExampleGenerate_withFormatting() {
	// El formato es ignorado al generar
	number := luhn.Generate("4532 0151 1283 036")
	fmt.Println(number)
	// Output: 4532015112830366
}