# Coverage Baseline Report - Shantilly v2.0

**Data de geração**: 2025-11-12 (a ser atualizado com dados reais quando Go estiver disponível)  
**Branch**: main  
**Commit**: [a ser preenchido]

---

## Sumário Executivo

Este documento estabelece o baseline de coverage do projeto Shantilly após a implementação da FASE 1 do plano de execução.

**Meta estabelecida**: ≥75% coverage geral  
**Threshold CI**: ≥65% (temporário, será incrementado para 70%)  
**Coverage atual**: [a ser medido com `go test -cover ./...`]

---

## Como Gerar Este Relatório

```bash
# 1. Executar testes com coverage
go test -v -race -coverpkg=./... -covermode=atomic -coverprofile=coverage.txt ./...

# 2. Gerar relatório funcional
go tool cover -func=coverage.txt > docs/qa/coverage-functional.txt

# 3. Gerar relatório HTML (para análise visual)
go tool cover -html=coverage.txt -o docs/qa/coverage.html

# 4. Extrair coverage total
go tool cover -func=coverage.txt | tail -1
```

---

## Coverage por Pacote

### Alta Prioridade (Core Runtime)

| Pacote                    | Coverage | LOC  | Testes | Status                |
|---------------------------|----------|------|--------|-----------------------|
| `internal/runtime/runner` | TBD%     | 596  | 5      | ⚠️ BAIXO - Meta: ≥80% |
| `internal/runtime/event`  | TBD%     | 162  | 5      | 🟡 MÉDIO - Meta: ≥80% |
| `internal/runtime/layout` | TBD%     | ~150 | 5      | 🟡 MÉDIO - Meta: ≥75% |
| `internal/runtime/modal`  | TBD%     | 272  | 12     | ✅ OK - Manter         |
| `pkg/declarative`         | TBD%     | 248  | 9      | ✅ OK - Manter         |

### Média Prioridade (Legacy TUI)

| Pacote                    | Coverage | LOC  | Testes | Status                |
|---------------------------|----------|------|--------|-----------------------|
| `internal/tui`            | TBD%     | 490  | 11     | 🟡 MÉDIO - Meta: ≥70% |
| `internal/config`         | TBD%     | 321  | 14     | ⚠️ BAIXO - Meta: ≥75% |
| `internal/tui/components` | TBD%     | ~200 | 7      | 🟡 MÉDIO - Meta: ≥70% |

### Baixa Prioridade (Utilidades)

| Pacote          | Coverage | LOC  | Testes | Status        |
|-----------------|----------|------|--------|---------------|
| `internal/util` | TBD%     | ~100 | 6      | ✅ OK - Manter |
| `cmd/shantilly` | TBD%     | 85   | 5      | ✅ OK - CLI    |

---

## Áreas Críticas com Coverage Insuficiente

### 1. ScriptRunner (`internal/runtime/runner/runner.go`)
- **Coverage estimado**: ~40-50%
- **Gap identificado**: 
  - Process lifecycle (timeout, cancellation)
  - Concurrent process replacement
  - Signal handling (SIGTERM → SIGKILL)
  - Security policy edge cases
- **Ação**: FASE 2.1 - Adicionar 10+ testes

### 2. Validation (`internal/config/validation.go`)
- **Coverage estimado**: ~50-60%
- **Gap identificado**:
  - Boundary values (min/max)
  - Malformed inputs
  - Edge cases de tipo
- **Ação**: FASE 2.2 - Adicionar 8+ testes

### 3. Form Model (`internal/tui/model.go`)
- **Coverage estimado**: ~55-65%
- **Gap identificado**:
  - State consistency
  - Update() edge cases
  - Validation display
- **Ação**: FASE 2.2 - Expandir testes de integração

---

## Gaps de Coverage por Tipo

### Funcionalidade Não Testada
1. **Concurrency**: Testes de race conditions limitados
2. **Error paths**: Poucos testes de cenários de erro
3. **Edge cases**: Inputs extremos não cobertos
4. **Security**: Apenas 1 teste de SecurityPolicy

### Código Não Testável
1. **View() methods**: Renderização visual difícil de testar unitariamente
2. **Terminal I/O**: Interações diretas com stdout/stderr
3. **Signal handlers**: Comportamento de sistema operacional

---

## Plano de Melhoria

### Fase 2 (Expansão de Testes)
- [ ] **2.1** - ScriptRunner: +10 testes → ≥80% coverage
- [ ] **2.2** - Validation: +8 testes → ≥75% coverage
- [ ] **2.3** - Security: +8 testes → 100% de attack vectors

### Meta de Coverage por Fase
- **Fase 1 (atual)**: Baseline estabelecido, threshold ≥65%
- **Fase 2 (após expansão)**: ≥70% geral
- **Fase 3 (após Waves 2-3)**: ≥72% geral
- **Fase 6 (final)**: ≥75% geral

---

## Threshold Evolution

| Data        | Coverage | Threshold CI | Ação                  |
|-------------|----------|--------------|-----------------------|
| 2025-11-12  | TBD%     | 65%          | Baseline estabelecido |
| Após Fase 2 | ~70%     | 70%          | Incrementar threshold |
| Após Fase 3 | ~72%     | 72%          | Incrementar threshold |
| Após Fase 6 | ≥75%     | 75%          | Meta final atingida   |

---

## Comandos Úteis

```bash
# Verificar coverage de um pacote específico
go test -cover ./internal/runtime/runner

# Gerar coverage apenas para arquivos não testados
go test -cover -coverpkg=./internal/runtime/runner ./internal/runtime/runner/...

# Verificar quais linhas NÃO estão cobertas
go tool cover -html=coverage.txt

# Comparar coverage entre commits
# (requer git diff + parsing customizado)
```

---

## Notas

1. **Threshold inicial (65%)** foi estabelecido conservadoramente para não bloquear CI imediatamente
2. Coverage será incrementado **gradualmente** conforme Fases 2-6 forem completadas
3. **Meta final (75%)** é alcançável e alinhada com projetos Go de qualidade
4. Coverage de **View() methods** não deve ser priorizado (difícil de testar unitariamente)
5. Foco deve ser em **lógica de negócio, validação, security e concurrency**

---

## Referências

- **Plano de Execução**: Documento de plano aprovado
- **Assessment QA**: `docs/qa/CODE-QUALITY-ASSESSMENT-20251112.md`
- **Matriz de Cobertura**: `docs/qa/matrix-epic-1-runtime-tui-coverage.md`
- **Gates QA**: `docs/qa/gates/`

---

**Última atualização**: 2025-11-12  
**Responsável**: QA Team  
**Próxima revisão**: Após Fase 2
