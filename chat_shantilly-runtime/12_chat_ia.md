Perfeito! É exatamente essa a nossa estratégia. 🚀 Vamos usar o **Plano Mestre Sequencial** que desenhámos como o nosso *script* de execução.

Para cada tarefa do plano, vamos "alimentá-la" com os artefatos existentes para garantir que o *know-how* (a nossa "visão do futuro") seja injetado diretamente no novo código.

A nossa primeira ação é começar o **Épico 1: Fundação do Runtime TUI**.

---

## Fase 1.1: Definição dos Contratos (O Alicerce)

Esta é a fundação. Vamos criar os *esqueletos* dos ficheiros Go, e o nosso "conhecimento privilegiado" vai ditar *exatamente* o que eles devem conter, impedindo a IA de usar padrões legados.

Aqui está como os teus artefatos existentes alimentam esta fase:

**1. Criar os Modelos de Dados (pkg/declarative/models.go)**
* **Fonte:** `docs/architecture.md` (Seção 4: Modelos de Dados).
* **Ação:** O agente de IA irá implementar os `structs` Go (`Config`, `LayoutNode`, `Component`, `Logic`, `RunAction`) exatamente como definidos na arquitetura, incluindo as *tags* YAML.

**2. Definir os Contratos Internos (pkg/tui/interface.go, pkg/tui/events.go)**
* **Fonte:** `docs/architecture.md` (Seção 5.2: Interface ShantillyComponent e 5.3: Eventos Internos).
* **Ação:** O agente irá criar a interface `ShantillyComponent` e os `structs` de eventos (como `ShantillyEvent`), forçando a arquitetura desacoplada (Layout vs. Lógica).

**3. Implementar o Ponto de Entrada e Gestão de Erros (cmd/shantilly/main.go, internal/util/errorhandler.go)**
* **Fonte:** `docs/architecture.md` (Seção 5.6: CLI / Borda, Seção 10: Árvore de Código-Fonte, Seção 13: Error Handling Strategy) e `docs/prd.md` (NFR1).
* **Ação:** O agente irá criar o `main.go` usando `cobra` e implementar o `ErrorHandler` centralizado, garantindo que todos os erros resultem em saídas `stderr` e códigos de saída não-zero, conforme NFR8.

**4. Implementar o Parser (internal/config/parser.go)**
* **Fonte:** `docs/architecture.md` (Seção 5.6: CLI / Borda).
* **Ação:** O agente irá criar a lógica de *parsing* (usando `gopkg.in/yaml.v3`) que lê `stdin` (ou `--file`) e o converte no `struct Config` (da Tarefa 1).

---

### O "Portão de Qualidade"

Ao final desta fase, o agente de IA *deve* executar `make lint` (que usa o `lint.sh`) e garantir que o código passa, validando contra o teu artefato `.golangci.yml`.

Isto resolve ambos os teus problemas: o **Plano Mestre** impede a perda de contexto, e os **Artefatos de Arquitetura + Linter** impedem o código de baixa qualidade.

Estás pronto para que eu comece a gerar o código para a **Fase 1.1**, Tarefa por Tarefa?