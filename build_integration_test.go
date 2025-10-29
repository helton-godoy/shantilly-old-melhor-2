package main

import (
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"testing"
)

// TestBuildBasic verifica se o target 'build' do Makefile funciona
func TestBuildBasic(t *testing.T) {
	// Limpar binário anterior
	cleanCmd := exec.Command("make", "clean")
	if err := cleanCmd.Run(); err != nil {
		t.Fatalf("Falha ao executar make clean: %v", err)
	}

	// Verificar que o binário foi removido
	if _, err := os.Stat("shantilly"); !os.IsNotExist(err) {
		t.Error("Binário shantilly ainda existe após make clean")
	}

	// Executar build
	buildCmd := exec.Command("make", "build")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Falha ao executar make build: %v", err)
	}

	// Verificar que o binário foi criado
	if _, err := os.Stat("shantilly"); os.IsNotExist(err) {
		t.Fatal("Binário shantilly não foi criado após make build")
	}

	// Verificar que é executável
	if runtime.GOOS != "windows" {
		if stat, err := os.Stat("shantilly"); err != nil || stat.Mode()&0o111 == 0 {
			t.Error("Binário shantilly não é executável")
		}
	}
}

// TestBuildCrossCompilation verifica se os targets de cross-compilação funcionam
func TestBuildCrossCompilation(t *testing.T) {
	// Limpar binários anteriores
	cleanCmd := exec.Command("make", "clean")
	if err := cleanCmd.Run(); err != nil {
		t.Fatalf("Falha ao executar make clean: %v", err)
	}

	// Executar build-all
	buildAllCmd := exec.Command("make", "build-all")
	if err := buildAllCmd.Run(); err != nil {
		t.Fatalf("Falha ao executar make build-all: %v", err)
	}

	// Verificar que todos os binários foram criados
	expectedBinaries := []string{
		"bin/linux/amd64/shantilly",
		"bin/linux/arm64/shantilly",
		"bin/darwin/amd64/shantilly",
		"bin/darwin/arm64/shantilly",
		"bin/windows/amd64/shantilly.exe",
	}

	for _, binary := range expectedBinaries {
		if _, err := os.Stat(binary); os.IsNotExist(err) {
			t.Errorf("Binário esperado não encontrado: %s", binary)
		}
	}
}

// TestBinaryExecutability verifica se os binários gerados são executáveis na plataforma atual
func TestBinaryExecutability(t *testing.T) {
	// Este teste só funciona na plataforma linux/amd64
	if runtime.GOOS != "linux" || runtime.GOARCH != "amd64" {
		t.Skip("Teste de executabilidade só funciona em linux/amd64")
	}

	binaryPath := "bin/linux/amd64/shantilly"

	// Verificar que o binário existe
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		t.Fatalf("Binário %s não encontrado", binaryPath)
	}

	// Testar comando --help
	cmd := exec.Command(binaryPath, "--help")
	if err := cmd.Run(); err != nil {
		t.Errorf("Falha ao executar %s --help: %v", binaryPath, err)
	}

	// Testar execução básica (deve mostrar usage e sair com sucesso)
	cmd = exec.Command(binaryPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("Comando sem argumentos deveria sair com sucesso, mas falhou: %v", err)
	}
	if len(output) == 0 {
		t.Error("Comando sem argumentos deveria produzir saída de uso")
	}
}

// TestBuildDirectoryStructure verifica se a estrutura de diretórios está correta
func TestBuildDirectoryStructure(t *testing.T) {
	// Verificar que o diretório bin existe
	if _, err := os.Stat("bin"); os.IsNotExist(err) {
		t.Fatal("Diretório bin não existe")
	}

	// Verificar estrutura de subdiretórios
	expectedDirs := []string{
		"bin/linux",
		"bin/linux/amd64",
		"bin/linux/arm64",
		"bin/darwin",
		"bin/darwin/amd64",
		"bin/darwin/arm64",
		"bin/windows",
		"bin/windows/amd64",
	}

	for _, dir := range expectedDirs {
		if info, err := os.Stat(dir); os.IsNotExist(err) {
			t.Errorf("Diretório esperado não encontrado: %s", dir)
		} else if !info.IsDir() {
			t.Errorf("%s deveria ser um diretório", dir)
		}
	}
}

// TestStaticLinking verifica se os binários são estaticamente linkados
func TestStaticLinking(t *testing.T) {
	// Este teste só funciona em sistemas Linux com ldd disponível
	if runtime.GOOS != "linux" {
		t.Skip("Teste de linking estático só funciona em Linux")
	}

	binaryPath := "bin/linux/amd64/shantilly"

	// Verificar que o binário existe
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		t.Fatalf("Binário %s não encontrado", binaryPath)
	}

	// Executar ldd para verificar linking
	cmd := exec.Command("ldd", binaryPath)
	output, err := cmd.CombinedOutput()
	// ldd retorna exit code 1 para binários estáticos, o que é esperado
	if err != nil && err.Error() != "exit status 1" {
		t.Fatalf("Falha ao executar ldd: %v", err)
	}

	outputStr := string(output)
	matched1, _ := regexp.MatchString("not a dynamic executable", outputStr)
	matched2, _ := regexp.MatchString("statically linked", outputStr)
	matched3, _ := regexp.MatchString("não é um executável dinâmico", outputStr)
	if !matched1 && !matched2 && !matched3 {
		t.Errorf("Binário não parece ser estaticamente linkado. Saída ldd: %s", outputStr)
	}
}

// TestBinarySizes verifica se os binários têm tamanhos razoáveis
func TestBinarySizes(t *testing.T) {
	binaries := map[string]int64{
		"bin/linux/amd64/shantilly":       8 * 1024 * 1024, // 8MB
		"bin/linux/arm64/shantilly":       8 * 1024 * 1024, // 8MB
		"bin/darwin/amd64/shantilly":      8 * 1024 * 1024, // 8MB
		"bin/darwin/arm64/shantilly":      8 * 1024 * 1024, // 8MB
		"bin/windows/amd64/shantilly.exe": 8 * 1024 * 1024, // 8MB
	}

	for binary, maxSize := range binaries {
		info, err := os.Stat(binary)
		if os.IsNotExist(err) {
			t.Errorf("Binário não encontrado: %s", binary)
			continue
		}
		if err != nil {
			t.Errorf("Erro ao verificar tamanho de %s: %v", binary, err)
			continue
		}

		if info.Size() > maxSize {
			t.Errorf("Binário %s é muito grande: %d bytes (máximo esperado: %d bytes)",
				binary, info.Size(), maxSize)
		}

		if info.Size() < 1024*1024 { // 1MB
			t.Errorf("Binário %s é muito pequeno: %d bytes (mínimo esperado: 1MB)",
				binary, info.Size())
		}
	}
}
