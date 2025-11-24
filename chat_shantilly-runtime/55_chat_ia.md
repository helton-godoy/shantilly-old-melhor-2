# Contexto – Next Wave 2 (Streaming viewport)

Data: 2025-11-14
Branch: `feat/runtime-migration`

## Estado atual

- Runtime declarativo v2.0 rodando via:
  - `./shantilly runtime --file app.example.yaml`
- Menu principal do exemplo contém:
  - `1) Verificar encapsulamento do legado (script real)`
  - `2) Não fazer nada (placeholder)`
  - `3) Rodar script de Hello (demo rápida do viewport)`

### Mudança aplicada nesta sessão

Arquivo modificado:
- `internal/components/viewport/model.go`

Alteração principal:
- No `Update`, ao receber `tuiapi.ScriptStdoutMsg`:
  - **Antes:** `m.rawContent = msg.Line` (sobrescrevia o conteúdo anterior).
  - **Agora:** `m.rawContent += msg.Line` (acumula o texto de cada mensagem).
- O restante do fluxo permanece igual:
  - `m.rawContentType = "text"`
  - Usa `wrapPlainText` se já houver largura conhecida.
  - `m.viewport.SetContent(m.content)` + `m.viewport.GotoBottom()`.

Efeito prático (com o estado atual do runner):
- O `ScriptRunner` + `tea_adapter` ainda enviam **uma** `ScriptStdoutMsg` por execução de script, com stdout/stderr já agregados.
- O viewport agora passa a **acumular blocos** de saída por execução (histórico da sessão), em vez de mostrar apenas o último resultado.

Build atual:
- `go build ./...` está **limpo**.
- `./shantilly runtime --file app.example.yaml` sobe o menu e executa scripts normalmente.

## Comportamentos considerados para o viewport

1. **Sobrescrever sempre** (mostrar só o último resultado)
   - Pro: painel limpo, ideal para validações rápidas (ex: `check-legacy-encapsulation.sh`).
   - Contra: não guarda histórico da sessão.

2. **Acumular por execução** (estado atual após esta wave)
   - Pro: histórico de tudo o que foi rodado na sessão, útil para troubleshooting e comparação antes/depois.
   - Contra: pode poluir quando scripts geram muito output.

3. **Streaming linha a linha em tempo real** (futuro próximo da wave 2)
   - Pro: UX ideal para scripts longos (Ansible, migrações, etc.), com feedback contínuo.
   - Contra: implementação mais complexa (pipes, goroutines, EventSink) e precisa decidir se combina com histórico ou limpa a cada run.

## Próximos passos sugeridos (Next Wave 2)

1. **Refatorar ScriptRunner para streaming real**
   - Em `internal/runtime/runner/runner.go`:
     - Trocar o modelo bufferizado (`stdoutBuf`/`stderrBuf`) por `StdoutPipe`/`StderrPipe` + `bufio.Scanner`.
     - Emitir eventos de output parciais (por linha ou chunk) via `EventSink.EmitUpdate`.
   - Em `internal/runtime/runner/tea_adapter.go`:
     - Remover o `bufferingSink` ou usá-lo apenas para casos específicos.
     - Passar a transformar cada `UpdateTargetUpdate` em uma `tui.ScriptStdoutMsg` separada (não só uma no final).

2. **Definir política de limpeza para o viewport**
   - Opções:
     - Limpar o viewport antes de cada nova execução (painel de execução atual).
     - Manter acumulado, mas talvez com um limite ou separadores claros entre execuções.
   - Decisão pode ser por componente (ex.: um viewport "log" acumula; outro viewport "resultado" sobrescreve).

3. **Criar scripts de demo para experimentar os modos** (exemplos sugeridos)
   - `scripts/status_once.sh` – curto, ideal para ver diferença entre sobrescrever x acumular.
   - `scripts/multi_step_demo.sh` – bloco com vários passos, bom para visualizar histórico.
   - `scripts/slow_counter.sh` – script lento para testar streaming linha a linha quando o runner for refatorado.

## Como retomar o trabalho no outro computador

1. Clonar o repositório e puxar o branch:

```bash
git clone https://github.com/helton-godoy/shantilly.git
cd shantilly
git checkout feat/runtime-migration
```

2. Testar o estado atual:

```bash
go build ./...
./shantilly runtime --file app.example.yaml
```

3. Confirmar comportamento do viewport:
- Usar o menu do `app.example.yaml`:
  - Rodar `3) Rodar script de Hello` algumas vezes.
  - Rodar `1) Verificar encapsulamento do legado`.
- Esperado: viewport acumula blocos de saída das execuções, em vez de mostrar só o último.

4. Próxima etapa de implementação:
- Começar pela refatoração do `ScriptRunner` + `tea_adapter` para streaming real.
- Depois, decidir configuração fina do viewport (limpar/acomodar histórico) com base nesses casos de uso.

---

## Decisão de nomenclatura para binding `input` → script (env vars)

Contexto:
- O runtime v2.0 passa a ter componentes de entrada (`type: input`) que podem alimentar scripts.
- Queremos um padrão de nomenclatura **claro, coerente e identitário**, inspirado em ferramentas como o Ansible.

Referência Ansible (inspiração conceitual, não cópia literal):
- Variáveis de playbook/facts: `lower_snake_case` (ex.: `deploy_version`, `http_port`).
- Variáveis internas de ambiente: prefixo `ANSIBLE_` em `UPPER_SNAKE_CASE` (ex.: `ANSIBLE_CONFIG`).
- Inputs interativos (`vars_prompt`) deixam o autor do playbook escolher o nome da variável, usado depois como `{{ deploy_version }}`.

### Decisão atual (Wave: binding automático input → env)

Para o binding inicial de inputs do Shantilly para scripts via variáveis de ambiente:

- **Prefixo reservado do runtime:** `SHANTILLY_`.
- Para componentes `type: input`, o runtime exporta automaticamente:

  - `SHANTILLY_INPUT_<ID_EM_UPPER_SNAKE>`

  Onde:
  - `ID_EM_UPPER_SNAKE` é derivado do `components[].id` (ex.: `input_box` → `INPUT_BOX`).

Exemplo:

```yaml
components:
  - id: param_input
    type: input
    props:
      label: "Parâmetro"
      placeholder: "Digite e pressione Enter…"
```

Ao submeter (`param_input:submit`), o runtime definirá para o script:

```bash
SHANTILLY_INPUT_PARAM_INPUT="<valor_digitado>"
```

Scripts podem consumir isso via `"$SHANTILLY_INPUT_PARAM_INPUT"`.

### Plano futuro: nome explícito via `env_name` (inspirado em `vars_prompt` do Ansible)

Para aproximar ainda mais da ergonomia do Ansible, planejamos introduzir um campo opcional em `props`:

```yaml
components:
  - id: param_input
    type: input
    props:
      label: "Parâmetro"
      env_name: "DEPLOY_VERSION"
```

Regras planejadas:
- Quando `env_name` estiver presente e válido:
  - O valor do input será exportado como `DEPLOY_VERSION=<valor>` para o script.
  - O binding automático `SHANTILLY_INPUT_PARAM_INPUT` poderá ser mantido por compatibilidade ou tornado opcional (a decidir na Wave correspondente).
- Quando `env_name` não estiver definido:
  - Vale o padrão automático `SHANTILLY_INPUT_<ID_EM_UPPER_SNAKE>`.

Este bloco serve como referência de governança para futuras Waves que implementarem efetivamente o binding `input` → env/stdin no `EventManager` + `ScriptRunner`.
