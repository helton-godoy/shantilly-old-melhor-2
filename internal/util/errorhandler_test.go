package util

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// TestHandle_NilError testa que Handle(nil) sai com código 0
func TestHandle_NilError(t *testing.T) {
	if os.Getenv("TEST_HANDLE_NIL") == "1" {
		Handle(nil)
		return
	}

	// Execute o teste em um subprocesso
	cmd := exec.Command(os.Args[0], "-test.run=TestHandle_NilError")
	cmd.Env = append(os.Environ(), "TEST_HANDLE_NIL=1")

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	// Handle(nil) deve sair com código 0, então não esperamos erro
	if err != nil {
		t.Fatalf("Esperado código de saída 0, mas obteve erro: %v", err)
	}

	// Não deve haver saída para stderr
	if stderr.Len() > 0 {
		t.Errorf("Esperado nenhuma saída stderr, mas obteve: %s", stderr.String())
	}
}

// TestHandle_AbortedError testa que Handle(ErrAborted) sai com código 2
func TestHandle_AbortedError(t *testing.T) {
	if os.Getenv("TEST_HANDLE_ABORTED") == "1" {
		Handle(ErrAborted)
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestHandle_AbortedError")
	cmd.Env = append(os.Environ(), "TEST_HANDLE_ABORTED=1")

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()

	// Esperamos código de saída 2 (ExitCancelled)
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		if exitErr.ExitCode() != ExitCancelled {
			t.Errorf("Esperado código de saída %d, mas obteve %d", ExitCancelled, exitErr.ExitCode())
		}
	} else {
		t.Fatalf("Esperado exit error com código %d, mas obteve: %v", ExitCancelled, err)
	}

	// Para ErrAborted, não deve haver mensagem de erro impressa
	stderrStr := stderr.String()
	if strings.Contains(stderrStr, "Error:") {
		t.Errorf("Não esperado mensagem de erro para ErrAborted, mas obteve: %s", stderrStr)
	}
}

// TestHandle_GenericError testa que Handle(erro genérico) sai com código 1
func TestHandle_GenericError(t *testing.T) {
	if os.Getenv("TEST_HANDLE_GENERIC") == "1" {
		err := errors.New("erro de teste")
		Handle(err)
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestHandle_GenericError")
	cmd.Env = append(os.Environ(), "TEST_HANDLE_GENERIC=1")

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()

	// Esperamos código de saída 1 (ExitError)
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		if exitErr.ExitCode() != ExitError {
			t.Errorf("Esperado código de saída %d, mas obteve %d", ExitError, exitErr.ExitCode())
		}
	} else {
		t.Fatalf("Esperado exit error com código %d, mas obteve: %v", ExitError, err)
	}

	// Deve haver mensagem de erro impressa
	stderrStr := stderr.String()
	if !strings.Contains(stderrStr, "Error: erro de teste") {
		t.Errorf("Esperado 'Error: erro de teste' em stderr, mas obteve: %s", stderrStr)
	}
}

// TestHandle_WrappedAbortedError testa que erros encapsulados são tratados corretamente
func TestHandle_WrappedAbortedError(t *testing.T) {
	if os.Getenv("TEST_HANDLE_WRAPPED") == "1" {
		wrappedErr := fmt.Errorf("operação falhou: %w", ErrAborted)
		Handle(wrappedErr)
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestHandle_WrappedAbortedError")
	cmd.Env = append(os.Environ(), "TEST_HANDLE_WRAPPED=1")

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()

	// Mesmo encapsulado, deve reconhecer ErrAborted e sair com código 2
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		if exitErr.ExitCode() != ExitCancelled {
			t.Errorf("Esperado código de saída %d para erro encapsulado, mas obteve %d", ExitCancelled, exitErr.ExitCode())
		}
	} else {
		t.Fatalf("Esperado exit error com código %d, mas obteve: %v", ExitCancelled, err)
	}
}

// TestErrAborted verifica que a variável ErrAborted está definida corretamente
func TestErrAborted(t *testing.T) {
	if ErrAborted == nil {
		t.Fatal("ErrAborted não deve ser nil")
	}

	expectedMsg := "operation aborted by user"
	if ErrAborted.Error() != expectedMsg {
		t.Errorf("Esperado mensagem '%s', mas obteve '%s'", expectedMsg, ErrAborted.Error())
	}
}

// TestExitCodes verifica que as constantes de código de saída estão definidas corretamente
func TestExitCodes(t *testing.T) {
	tests := []struct {
		name     string
		code     int
		expected int
	}{
		{"ExitSuccess", ExitSuccess, 0},
		{"ExitError", ExitError, 1},
		{"ExitCancelled", ExitCancelled, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.code != tt.expected {
				t.Errorf("%s: esperado %d, mas obteve %d", tt.name, tt.expected, tt.code)
			}
		})
	}
}
