.PHONY: test bench coverage clean help fmt vet fuzz ci all

# Default target
all: test

# Ejecutar tests
test:
	@echo "🧪 Ejecutando tests..."
	@go test -v -race ./...

# Ejecutar benchmarks
bench:
	@echo "📊 Ejecutando benchmarks..."
	@go test -bench=. -benchmem

# Generar cobertura
coverage:
	@echo "📈 Generando reporte de cobertura..."
	@go test -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Reporte generado: coverage.html"
	@go tool cover -func=coverage.out

# Fuzzing
fuzz:
	@echo "🎲 Ejecutando fuzzing por 30s..."
	@go test -fuzz=FuzzValidate -fuzztime=30s

# Limpiar archivos generados
clean:
	@echo "🧹 Limpiando..."
	@rm -f coverage.out coverage.html
	@rm -rf testdata/fuzz

# Verificar formato
fmt:
	@echo "🎨 Verificando formato..."
	@gofmt -l -w .
	@echo "✅ Formato aplicado"

# Análisis estático
vet:
	@echo "🔍 Ejecutando go vet..."
	@go vet ./...
	@echo "✅ Análisis estático completado"

# Ejecutar los ejemplos
examples:
	@echo "📖 Ejecutando ejemplos..."
	@go test -run Example

# CI: todo lo necesario para CI/CD
ci: fmt vet test coverage
	@echo "✅ Pipeline CI completado exitosamente"

# Instalar herramientas de desarrollo
tools:
	@echo "🔧 Instalando herramientas..."
	@go install golang.org/x/perf/cmd/benchstat@latest
	@echo "✅ Herramientas instaladas"

# Comparar benchmarks
bench-compare: bench
	@echo "📊 Para comparar benchmarks:"
	@echo "1. Guarda los resultados actuales: make bench > old.txt"
	@echo "2. Haz tus cambios"
	@echo "3. Ejecuta: make bench > new.txt"
	@echo "4. Compara: benchstat old.txt new.txt"

# Test rápido sin race detection
quick:
	@echo "⚡ Test rápido..."
	@go test ./...

# Ver documentación
doc:
	@echo "📚 Abriendo documentación..."
	@go doc -all

# Ayuda
help:
	@echo "Comandos disponibles:"
	@echo "  make test       - Ejecutar tests con race detection"
	@echo "  make quick      - Ejecutar tests rápidos (sin race detection)"
	@echo "  make bench      - Ejecutar benchmarks"
	@echo "  make coverage   - Generar reporte de cobertura"
	@echo "  make fuzz       - Ejecutar fuzzing por 30s"
	@echo "  make examples   - Ejecutar los ejemplos"
	@echo "  make clean      - Limpiar archivos generados"
	@echo "  make fmt        - Formatear código"
	@echo "  make vet        - Análisis estático"
	@echo "  make ci         - Pipeline CI completo"
	@echo "  make tools      - Instalar herramientas de desarrollo"
	@echo "  make doc        - Ver documentación"
	@echo "  make help       - Mostrar esta ayuda"