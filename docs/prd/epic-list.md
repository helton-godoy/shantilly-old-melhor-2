# Epic List (Roadmap v2.0)

Este arquivo descreve o roadmap de épicos vigente para o Runtime TUI Declarativo v2.0.
Os épicos v1.x focados em "Form Functionality" são considerados material legado de referência e não representam mais o plano principal do produto.

* **Épico 1: Runtime TUI Declarativo — Fundação do Runtime**
  * **Meta:** Construir o motor central: layout (`column`/`row`/`box`), componentes essenciais (`list`, `viewport`, `form`, `buttongroup`), lógica de eventos (`on:`), runner genérico (`run: { script: ... }`), `args`, `stdin`, `update_target` e segurança JIT com pilha de modais.
  * **Detalhamento:** Ver [`docs/prd/epic-1-runtime-tui-foundation.md`](docs/prd/epic-1-runtime-tui-foundation.md:1).

* **Épico 2: Runner Especialista (Ansible)**
  * **Meta:** Introduzir `run: { ansible_playbook: ... }` como runner especializado sobre a fundação genérica, com suporte a `vars`, descoberta de inventário e prompts seguros (ex.: `ask_vault_pass`).

* **Épico 3: Runtime Preditivo**
  * **Meta:** Adicionar componentes de descoberta inteligente (ex.: `playbook_explorer`, `inventory_explorer`) que consomem e apresentam informações no Runtime TUI Declarativo, construídos sobre a fundação do Épico 1.

* **Épico 4: Administração SSH**
  * **Meta:** Permitir servir o Runtime TUI Declarativo via SSH (ex.: integração com `wish`), mantendo o mesmo modelo declarativo e políticas de segurança.

---

## Apêndice: Épicos v1.x (Legado / Referência)

Os épicos abaixo NÃO representam mais o roadmap principal, mas servem como insumo histórico e fonte para a refatoração do `FormComponent` dentro do Runtime TUI Declarativo:

* Epic 1: MVP - Core Form Functionality
* Epic 2: Advanced Form Features & User Experience
* Epic 3: Advanced Layouts & Multi-Panel Forms
* Epic 4: SSH Server Mode & Remote Forms
* Epic 5: Advanced Components & Interactions
* Epic 6: Form Templates & Reusability
* Epic 7: Internationalization & Localization
