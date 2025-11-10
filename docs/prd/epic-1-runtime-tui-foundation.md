# Epic 1: Runtime TUI Declarativo — Fundação do Runtime (v2.0)

Este épico consolida a visão do shantilly como um Runtime TUI Declarativo orientado a eventos, consumindo um único YAML que define layout, componentes e lógica (`on:` + `run:`), alinhado com [`docs/prd.md`](docs/prd.md:1) e [`docs/architecture.md`](docs/architecture.md:1). Ele substitui a visão anterior centrada exclusivamente em "form" como produto e absorve o v1.0 como implementação do `FormComponent`.

Este documento é a referência de requisitos para os agentes de Arquitetura, Dev, QA e Análise ao implementar a fundação v2.0.

## 1. Estrutura do Épico 1 v2.0

O Épico 1 é dividido em 5 blocos fundamentais:

- 1.1: Core Runtime & Entry (binário único + YAML único).
- 1.2: LayoutManager (layout hierárquico, foco básico, múltiplos componentes).
- 1.3: EventManager + ShantillyEvent + blocos `on:`.
- 1.4: ScriptRunner (run:, args, stdin, update_target, lifecycle).
- 1.5: Modal Stack & Segurança JIT (confirm, prompt_secrets, anti-Trojan YAML).

Cada bloco abaixo inclui contexto, requisitos e rastreabilidade com PRD/Arquitetura.

---

## 1.1 Core Runtime & Entry

Como um SysAdmin/script developer,
Eu quero um binário único `shantilly` capaz de ler um YAML declarativo (via stdin ou `--file`) e inicializar o Runtime TUI Declarativo,
Para que eu possa usar o mesmo runtime genérico em qualquer automação de terminal.

### Escopo

- Empacotamento como binário único, estático, cross-platform.
- Entrada padrão: YAML único, via:
  - `stdin`, ou
  - flag `--file`.
- Gatilho único: runtime orientado a layout + componentes + `on:`.

### Requisitos

1. O binário `shantilly` DEVE:
   - Ler um único documento YAML como fonte de verdade da UI/fluxo.
   - Inicializar o Runtime TUI Declarativo com base nos modelos definidos em [`docs/architecture.md`](docs/architecture.md:141).
2. Não há mais “subcomando `form`” como produto final:
   - Formulários passam a ser um tipo de `ShantillyComponent` utilizável no layout.
3. O comportamento de entrada/saída (NFRs v1.0) permanece:
   - Binário único, rápido, com erros em `stderr` e códigos de saída corretos.

### Rastreabilidade

- PRD v2.0:
  - FR1–FR3, FR7–FR11 (fundação de layout/Componentes/on:/run:).
  - NFR1, NFR2, NFR3, NFR5.
- Arquitetura v2.0:
  - Seções 2, 3, 4 (Config, LayoutNode, Component, Logic, RunAction).
  - Source tree alvo.

---

## 1.2 LayoutManager — Layout Hierárquico e Foco Global

Como um SysAdmin,
Eu quero definir um layout TUI hierárquico com `column`, `row` e `box`, contendo múltiplos componentes,
Para que eu possa construir dashboards complexos, responsivos e utilizáveis pelo teclado.

### Escopo

- Implementar o LayoutManager como:
  - Modelo raiz TEA.
  - Responsável por layout, foco global e integração com a pilha de modais.
- Suporte a:
  - `type: column`, `type: row`, `type: box`.
  - Propriedades `width`, `height`, `flex`.
  - Composição de componentes declarativos dentro de `box`.

### Requisitos Funcionais

1. O parser DEVE suportar `type: column|row|box` com `items` aninhados.
2. O LayoutManager DEVE:
   - Calcular e aplicar dimensões segundo `width`, `height`, `flex`.
   - Reagir a `tea.WindowSizeMsg` recalculando layout de forma fluida e sem flicker (NFR2).
3. Um `box` pode conter um único `component:` declarativo.
4. Foco:
   - O LayoutManager DEVE manter foco global entre componentes.
   - Teclas de navegação de alto nível (ex.: Ctrl+Tab) DEVEM alternar entre componentes/painéis focáveis.
   - Somente o componente focado recebe eventos de entrada relevantes.

### Rastreabilidade

- PRD v2.0:
  - FR1, FR2, FR3; NFR2; metas de UX de dashboards.
- Arquitetura v2.0:
  - Seções 4.2 (LayoutNode), 5.4 (LayoutManager), 7.1 (workflow de foco).

---

## 1.3 EventManager + ShantillyEvent + Blocos `on:`

Como um SysAdmin,
Eu quero que interações em componentes gerem eventos semânticos que disparem ações declaradas no YAML (`on:`),
Para que minha UI seja reativa e ligada à automação.

### Escopo

- Padronizar eventos de UI via `ShantillyEvent`.
- Introduzir EventManager para:
  - Receber eventos de componentes.
  - Resolver regras `on:` declaradas no YAML.
  - Produzir comandos de execução (RunRequest) para o ScriptRunner.

### Requisitos Funcionais

1. Definir `ShantillyEvent` com:
   - `ComponentID`, `Type`, `Payload`.
2. Componentes TUI (list, buttongroup, form etc.) DEVEM emitir ShantillyEvent em vez de acoplamento direto.
3. O YAML DEVE suportar bloco `on:` na raiz:
   - Lista de regras:
     - `event: "component_id:action"` (ex.: `menu:select`, `user_form:submit`).
     - `run: { ... }` obrigatório.
4. O EventManager DEVE:
   - Mapear ShantillyEvent → regra `on:` correspondente.
   - Encaminhar uma `RunRequest` estruturada ao ScriptRunner.
   - Integrar com segurança JIT (ver 1.5) via flags de `confirm` e `prompt_secrets`.

### Rastreabilidade

- PRD v2.0:
  - FR8 (lógica `on:`), FR9–FR11 (integração com run/script/update_target).
- Arquitetura v2.0:
  - Seções 4.6 (Logic, RunAction), 5.4 (EventManager), 7.2 (workflow eventos→scripts).

---

## 1.4 ScriptRunner — Execução, Dados e `update_target`

Como um SysAdmin,
Eu quero mapear eventos a scripts externos com passagem de dados via args/stdin e atualização de viewports,
Para que shantilly orquestre automações em tempo real de forma segura e previsível.

### Escopo

- Implementar ScriptRunner como motor dedicado para `run:`.
- Suportar:
  - `run: { script: "/path/to/script.sh" }`.
  - `args: []string` com templates.
  - `stdin: any` serializado como JSON.
  - `update_target` para vínculo com componentes (ex.: viewport).

### Requisitos Funcionais

1. ScriptRunner DEVE:
   - Executar `run.script` como processo filho.
   - Aplicar templates de `args` e `stdin` com base em dados dos componentes (ex.: `{{ form.username }}`, `{{ component.menu.selected_id }}`).
2. Se `update_target` for definido:
   - O ScriptRunner DEVE enviar saída (`stdout`) para o componente alvo (tipicamente `viewport`).
3. Política de ciclo de vida (FR11):
   - Para um mesmo `update_target`, um novo `run:` DEVE:
     - Enviar SIGTERM ao processo anterior.
     - Garantir limpeza antes de iniciar o novo processo (com timeout e fallback SIGKILL, se definido em arquitetura).
4. Em caso de erro:
   - Emitir eventos de erro (`RuntimeErrorMsg`) para o runtime, nunca `os.Exit` direto.

### Rastreabilidade

- PRD v2.0:
  - FR9, FR10, FR11.
- Arquitetura v2.0:
  - Seções 4.6 (RunAction), 5.4 (ScriptRunner), 7.2 (workflow scripts/viewport).

---

## 1.5 Modal Stack & Segurança JIT (confirm / prompt_secrets / Anti-Trojan YAML)

Como um responsável por automação segura,
Eu quero que ações sensíveis ou com segredos sejam protegidas por confirmações e prompts seguros,
Para que arquivos YAML não possam executar comandos perigosos sem validação explícita do operador.

### Escopo

- Introduzir pilha de modais (Modal Stack) como parte do LayoutManager.
- Integrar segurança JIT no fluxo `on:`:
  - `confirm: true`
  - `prompt_secrets: [...]`

### Requisitos Funcionais

1. O YAML `on:` DEVE permitir:
   - `confirm: true` para exigir confirmação explícita do usuário antes da execução.
   - `prompt_secrets: ["secret_name", ...]` para coletar segredos via modal TUI.
2. O EventManager, ao encontrar essas flags, DEVE:
   - Emitir mensagens para o LayoutManager abrir modais correspondentes.
3. O LayoutManager DEVE:
   - Manter uma pilha de modais sobre o layout.
   - Bloquear interação com o plano de fundo enquanto um modal estiver ativo.
4. Apenas após confirmação e coleta de segredos:
   - O EventManager DEVE montar a `RunRequest` final e encaminhar ao ScriptRunner.
5. Nenhum segredo DEVE ser:
   - Persistido em disco.
   - Exposto em logs.
6. Anti-Trojan YAML:
   - O design DEVE desencorajar execuções automáticas sem interação do operador:
     - scripts só são executados em resposta a eventos de UI,
     - ações marcadas como sensíveis exigem confirmação/prompt.

### Rastreabilidade

- PRD v2.0:
  - Requisitos de segurança JIT descritos nas metas e NFRs.
- Arquitetura v2.0:
  - Seções 5.4 (EventManager), 5.4/5.5 (Modal stack implícita no LayoutManager), 7.3 (workflow segurança JIT), Security.

---

## 2. Restrições e Decisões Estratégicas

Estas regras são obrigatórias em todos os artefatos e implementações derivados deste épico:

1. v2.0 é exclusivamente sobre o Runtime TUI Declarativo:
   - Layout + Componentes + on:/run: + ScriptRunner + Modal Stack.
2. v1.0:
   - Tratado apenas como base do `FormComponent`.
   - Nunca descrito como produto final separado.
3. Proibição de `os.Exit` no core do runtime:
   - Encerramentos e erros fluem via mensagens/contratos internos.
4. Um processo por `update_target`:
   - Sempre encerrar o anterior antes de iniciar o próximo.
5. YAML como fonte única de verdade:
   - Nenhum outro artefato pode redefinir o modelo conceitual fora do que está em PRD v2.0 + Arquitetura v2.0.

---

## 3. Impacto sobre Épicos Anteriores (v1.x)

Os épicos `epic-1-mvp-core-form-functionality.md` e seguintes são reclassificados como:

- Inputs históricos e material base para:
  - implementação do `FormComponent` (Story 1.4 deste épico),
  - futuras extensões específicas (ex.: runners especializados, UX avançada, i18n).
- Eles NÃO representam mais o roadmap principal.
- O roadmap autoritativo passa a ser:
  - Seção “Epic List (Roadmap v2.0)” de [`docs/prd.md`](docs/prd.md:215),
  - Este arquivo como detalhamento normativo do Épico 1.

---

## 4. Entregáveis Esperados por Agentes

Este épico guia diretamente o trabalho dos demais agentes:

- Architect:
  - Garantir que [`docs/architecture.md`](docs/architecture.md:1) e derivados implementem estes contratos.
- Dev:
  - Implementar LayoutManager, EventManager, ScriptRunner, componentes e modal stack seguindo estas histórias.
- QA:
  - Definir gates e testes (incl. teatest) validados contra 1.1–1.5.
- Analyst / PO:
  - Mapear stories em `docs/stories/` diretamente para estas seções, com rastreio explícito.

Este documento deve ser mantido sincronizado com o PRD v2.0 e a Arquitetura v2.0. Qualquer alteração estrutural no runtime deve primeiro ser refletida aqui.
