# GitHub Roadmap Issues - Shantilly Project

## Issues de Roadmap Evolutivo Baseadas no PRD

### Foundation Sprint (v1.0 Alpha - 3-4 semanas)

**Issue #001: CLI Foundation & Project Structure**
- Type: feature | Priority: high | Area: core | Story Points: 5 | Epic: E1 - Runtime TUI Foundation
- Description: Implementar estrutura CLI básica com comandos foundation, configuração inicial e estrutura de diretórios para o projeto Shantilly.
- Acceptance Criteria: CLI aceita comandos básicos (help, version, init) | Estrutura de diretórios criada automaticamente | Configuração inicial gerada

**Issue #002: YAML Configuration Parsing**
- Type: feature | Priority: high | Area: core | Story Points: 8 | Epic: E1 - Runtime TUI Foundation
- Description: Implementar parser YAML robusto para configurações de formulário declarativo com validação e tipagem.
- Acceptance Criteria: Parser YAML funcional para AppConfig | Validação de schemas | Suporte a tipos complexos

**Issue #003: Basic TUI Structure with Bubble Tea**
- Type: feature | Priority: high | Area: ui | Story Points: 8 | Epic: E1 - Runtime TUI Foundation
- Description: Criar estrutura básica da TUI usando Bubble Tea com modelo principal e renderização inicial.
- Acceptance Criteria: Aplicação TUI básica funcional | Renderização de componentes simples | Loop principal operacional

**Issue #004: Form Rendering with huh**
- Type: feature | Priority: high | Area: ui | Story Points: 10 | Epic: E1 - Runtime TUI Foundation
- Description: Implementar renderização de formulários usando a biblioteca huh com suporte a campos básicos.
- Acceptance Criteria: Renderização de formulários básicos | Suporte a campos de texto, checkbox, select | Interação com usuário funcional

**Issue #005: Keyboard Navigation & Interaction**
- Type: feature | Priority: high | Area: ui | Story Points: 6 | Epic: E1 - Runtime TUI Foundation
- Description: Implementar navegação por teclado e interações básicas na interface TUI.
- Acceptance Criteria: Navegação por Tab/Shift+Tab | Enter/Space para ativação | Escape para cancelamento

**Issue #006: Submission & JSON Output**
- Type: feature | Priority: high | Area: core | Story Points: 5 | Epic: E1 - Runtime TUI Foundation
- Description: Implementar funcionalidade de submissão de formulários e output JSON estruturado.
- Acceptance Criteria: Coleta de dados do formulário | Geração de JSON válido | Output para stdout/file

**Issue #007: Basic Styling & Alignment**
- Type: feature | Priority: medium | Area: ui | Story Points: 4 | Epic: E1 - Runtime TUI Foundation
- Description: Implementar sistema básico de styling e alinhamento para melhor apresentação visual.
- Acceptance Criteria: Tema básico implementado | Alinhamento correto de elementos | Feedback visual adequado

**Issue #008: Static Build & Distribution**
- Type: task | Priority: medium | Area: infrastructure | Story Points: 3 | Epic: E1 - Runtime TUI Foundation
- Description: Configurar build estático e sistema de distribuição para múltiplas plataformas.
- Acceptance Criteria: Build para Linux/Windows/macOS | Binários gerados automaticamente | Sistema de release configurado

### Advanced Features Sprint (v1.0 Beta)

**Issue #009: Advanced Form Types & Validation**
- Type: feature | Priority: high | Area: core | Story Points: 8 | Epic: E2 - Advanced Form Features
- Description: Implementar tipos avançados de campos e sistema de validação robusto.
- Acceptance Criteria: Campos: date, file, number, email | Validação em tempo real | Mensagens de erro customizadas

**Issue #010: Enhanced Error Handling & UX**
- Type: feature | Priority: medium | Area: ui | Story Points: 6 | Epic: E2 - Advanced Form Features
- Description: Implementar tratamento de erros abrangente e melhorias de experiência do usuário.
- Acceptance Criteria: Error handling centralizado | Mensagens de erro claras | Recuperação graceful de erros

### Runtime Architecture Sprint (v2.0 Features)

**Issue #011: Basic Multi-Panel Navigation**
- Type: feature | Priority: high | Area: ui | Story Points: 10 | Epic: E3 - Multi-Panel Runtime
- Description: Implementar sistema de navegação entre múltiplos painéis/páginas.
- Acceptance Criteria: Sistema de tabs/painéis | Navegação fluida | Estado preservado entre painéis

**Issue #012: Event Engine & on: Routing**
- Type: feature | Priority: high | Area: core | Story Points: 12 | Epic: E3 - Event-Driven Runtime
- Description: Implementar engine de eventos e sistema de roteamento declarativo com sintaxe on:.
- Acceptance Criteria: EventManager funcional | Sintaxe on: implementada | Roteamento baseado em eventos

**Issue #013: ScriptRunner Execution Declarativa (E14)**
- Type: feature | Priority: high | Area: runtime | Story Points: 15 | Epic: E4 - ScriptRunner Architecture
- Description: Implementar ScriptRunner para execução declarativa de scripts com isolamento e segurança.
- Acceptance Criteria: Execução de scripts seguros | Isolamento por processo | Sistema de update_target

**Issue #014: Modal Stack & JIT Security (E15)**
- Type: feature | Priority: high | Area: security | Story Points: 10 | Epic: E5 - Security & Modal System
- Description: Implementar pilha de modais com segurança JIT e validação anti-trojan.
- Acceptance Criteria: ModalStack funcional | Validação JIT de segurança | Prevenção de ataques

**Issue #015: Encapsulamento FormComponent Legacy**
- Type: refactor | Priority: medium | Area: legacy | Story Points: 8 | Epic: E6 - Legacy Encapsulation
- Description: Encapsular código legado v1.x do FormComponent em sandbox isolado.
- Acceptance Criteria: Código legado encapsulado | Interface limpa exposta | Migração gradual

**Issue #016: Hardening Security JIT Anti-Trojan**
- Type: security | Priority: critical | Area: security | Story Points: 12 | Epic: E7 - Security Hardening
- Description: Implementar hardening de segurança com validação JIT e proteção anti-trojan.
- Acceptance Criteria: Whitelist explícita de ações | Validação forte de AppConfig | Políticas de segurança rigorosas

### Extensions & Integrations

**Issue #017: Ansible Runner Integration**
- Type: feature | Priority: high | Area: integrations | Story Points: 12 | Epic: E8 - Ansible Integration
- Description: Implementar runner especializado para Ansible com suporte a playbooks e vars.
- Acceptance Criteria: Integração com Ansible | Suporte a vault secrets | Inventário discovery

**Issue #018: SSH Server Mode Implementation**
- Type: feature | Priority: medium | Area: server | Story Points: 15 | Epic: E9 - SSH Server Mode
- Description: Implementar modo servidor SSH para execução remota do runtime TUI.
- Acceptance Criteria: Servidor SSH funcional | Autenticação segura | TUI remoto operacional

**Issue #019: Predictive Runtime Components**
- Type: feature | Priority: medium | Area: ai | Story Points: 10 | Epic: E10 - Predictive Runtime
- Description: Implementar componentes de descoberta inteligente (playbook_explorer, inventory_explorer).
- Acceptance Criteria: Discovery automático | Sugestões inteligentes | Cache de metadados

**Issue #020: Advanced Layout System**
- Type: feature | Priority: medium | Area: ui | Story Points: 8 | Epic: E11 - Advanced Layouts
- Description: Implementar sistema avançado de layout com column/row/box.
- Acceptance Criteria: Layouts complexos | Responsividade | Temas customizáveis

### Additional Enhancement Issues

**Issue #021: Plugin Architecture**
- Type: feature | Priority: low | Area: extensibility | Story Points: 12 | Epic: E12 - Plugin System
- Description: Implementar arquitetura de plugins para extensibilidade futura.
- Acceptance Criteria: Sistema de plugins | API de extensão | Marketplace básico

**Issue #022: Multi-tenant Support**
- Type: feature | Priority: low | Area: enterprise | Story Points: 10 | Epic: E13 - Multi-tenancy
- Description: Implementar suporte multi-tenant para uso empresarial.
- Acceptance Criteria: Isolamento de dados | Configurações por tenant | Billing integration

**Issue #023: Enhanced Security Policies**
- Type: security | Priority: medium | Area: security | Story Points: 8 | Epic: E14 - Security Policies
- Description: Implementar políticas de segurança avançadas e compliance.
- Acceptance Criteria: RBAC básico | Audit logging | Compliance reports

**Issue #024: Performance Optimization**
- Type: optimization | Priority: medium | Area: performance | Story Points: 6 | Epic: E15 - Performance
- Description: Otimizações de performance para grandes formulários e datasets.
- Acceptance Criteria: Lazy loading | Virtual scrolling | Memory optimization

**Issue #025: Documentation Portal**
- Type: documentation | Priority: medium | Area: docs | Story Points: 5 | Epic: E16 - Documentation
- Description: Criar portal de documentação completo com exemplos interativos.
- Acceptance Criteria: Documentação interativa | Exemplos funcionais | Search integrado

---

## Labels Hierárquicas Sugeridas

### Priority Labels
- priority::critical - 🔴 Crítico
- priority::high - 🟠 Alto
- priority::medium - 🟡 Médio
- priority::low - 🟢 Baixo

### Type Labels
- type::feature - ✨ Nova funcionalidade
- type::bug - 🐛 Bug
- type::refactor - 🔄 Refatoração
- type::task - 📋 Tarefa
- type::security - 🔒 Segurança
- type::documentation - 📚 Documentação
- type::optimization - ⚡ Otimização

### Area Labels
- area::core - 🏗️ Core
- area::ui - 🖥️ Interface
- area::security - 🔐 Segurança
- area::runtime - ⚙️ Runtime
- area::integrations - 🔗 Integrações
- area::server - 🖥️ Servidor
- area::ai - 🤖 IA/ML
- area::infrastructure - 🏗️ Infraestrutura

### Status Labels
- status::triage - 🔍 Triagem
- status::in-progress - 🚧 Em progresso
- status::review - 👀 Revisão
- status::blocked - 🚫 Bloqueado
- status::done - ✅ Concluído

### Epic Labels
- epic::e1 - E1 Foundation
- epic::e2 - E2 Advanced Features
- epic::e3 - E3 Multi-Panel
- epic::e4 - E4 ScriptRunner
- epic::e5 - E5 Security
- epic::e6 - E6 Legacy
- epic::e7 - E7 Hardening
- epic::e8+ - E8+ Extensions

---

## Milestones Sugeridos

### v1.0 Alpha (3-4 semanas)
- Issues: #001-#008
- Goal: MVP funcional com foundations básicas

### v1.0 Beta (2-3 semanas após Alpha)
- Issues: #009-#010
- Goal: Funcionalidades avançadas e UX melhorada

### v2.0 Features (4-5 semanas)
- Issues: #011-#016
- Goal: Runtime architecture completa

### v2.0 Extensions (3-4 semanas)
- Issues: #017-#025
- Goal: Integrações e extensões avançadas

---

## GitHub Project: "Shantilly Roadmap"

Columns sugeridos:
- 📋 Backlog
- 🚧 In Progress  
- 👀 Review
- ✅ Done
- 🚫 Blocked

Automation rules:
- Auto-move to "In Progress" when assigned
- Auto-move to "Review" when labeled "status::review"
- Auto-move to "Done" when closed
- Auto-add to "Blocked" when labeled "status::blocked"
