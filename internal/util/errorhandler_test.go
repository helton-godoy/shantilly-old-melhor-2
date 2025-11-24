package util

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

func captureStderr(t *testing.T, fn func()) string {
	t.Helper()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Erro ao criar pipe: %v", err)
	}

	oldStderr := os.Stderr
	os.Stderr = w

	fn()

	_ = w.Close()
	os.Stderr = oldStderr

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("Erro ao ler stderr: %v", err)
	}
	_ = r.Close()

	return buf.String()
}

// TestHandle_NilError testa que Handle(nil) sai com código 0
func TestHandle_NilError(t *testing.T) {
	output := captureStderr(t, func() {
		code := Handle(nil)
		if code != ExitSuccess {
			t.Fatalf("Esperado código de saída %d, mas obteve: %d", ExitSuccess, code)
		}
	})

	// Não deve haver saída para stderr
	if output != "" {
		t.Errorf("Esperado nenhuma saída stderr, mas obteve: %s", output)
	}
}

// TestHandle_AbortedError testa que Handle(ErrAborted) sai com código 2
func TestHandle_AbortedError(t *testing.T) {
	output := captureStderr(t, func() {
		code := Handle(ErrAborted)
		if code != ExitCancelled {
			t.Fatalf("Esperado código de saída %d, mas obteve: %d", ExitCancelled, code)
		}
	})

	// Para ErrAborted, não deve haver mensagem de erro impressa
	if strings.Contains(output, "Error:") {
		t.Errorf("Não esperado mensagem de erro para ErrAborted, mas obteve: %s", output)
	}
}

// TestHandle_GenericError testa que Handle(erro genérico) sai com código 1
func TestHandle_GenericError(t *testing.T) {
	output := captureStderr(t, func() {
		code := Handle(errors.New("erro de teste"))
		if code != ExitError {
			t.Fatalf("Esperado código de saída %d, mas obteve: %d", ExitError, code)
		}
	})

	// Deve haver mensagem de erro impressa
	if !strings.Contains(output, "Error: erro de teste") {
		t.Errorf("Esperado 'Error: erro de teste' em stderr, mas obteve: %s", output)
	}
}

// TestHandle_WrappedAbortedError testa que erros encapsulados são tratados corretamente
func TestHandle_WrappedAbortedError(t *testing.T) {
	wrappedErr := fmt.Errorf("operação falhou: %w", ErrAborted)

	output := captureStderr(t, func() {
		code := Handle(wrappedErr)
		if code != ExitCancelled {
			t.Fatalf("Esperado código de saída %d para erro encapsulado, mas obteve: %d", ExitCancelled, code)
		}
	})

	if strings.Contains(output, "Error:") {
		t.Errorf("Não esperado mensagem de erro para ErrAborted encapsulado, mas obteve: %s", output)
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
