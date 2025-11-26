# GitHub Structures Implementation Guide - Shantilly

## 📋 Status: Structures Preparadas para Implementação

### ✅ Material Preparado

- **Labels hierárquicas** configuradas (.github/labels/)
- **Milestones** definidos (.github/milestones.yml)
- **GitHub Projects** estruturado (.github/projects.yml)
- **Templates** personalizados (.github/ISSUE_TEMPLATE/)

### 🏷️ Labels Hierárquicas (15+ labels)

#### Priority Labels

| Label                | Color   | Description                                | Usage                 |
| -------------------- | ------- | ------------------------------------------ | --------------------- |
| `priority::critical` | #d73a4a | Critical priority - must fix immediately   | Must-have fixes       |
| `priority::high`     | #fb8500 | High priority - important for next release | Next release scope    |
| `priority::medium`   | #fbbf24 | Medium priority - nice to have             | Nice to have features |
| `priority::low`      | #10b981 | Low priority - future consideration        | Future planning       |

#### Type Labels  

| Label                 | Color   | Description              | Usage               |
| --------------------- | ------- | ------------------------ | ------------------- |
| `type::feature`       | #1f77b4 | New functionality        | New features        |
| `type::bug`           | #d62728 | Something isn't working  | Bug fixes           |
| `type::refactor`      | #ff7f0e | Refactoring code         | Code improvements   |
| `type::task`          | #2ca02c | Non-code related tasks   | Administrative work |
| `type::security`      | #e377c2 | Security related         | Security concerns   |
| `type::documentation` | #9467bd | Documentation changes    | Docs updates        |
| `type::optimization`  | #bcbd22 | Performance improvements | Performance work    |

#### Area Labels

| Label                  | Color   | Description              | Usage                 |
| ---------------------- | ------- | ------------------------ | --------------------- |
| `area::core`           | #8c564b | Core functionality       | Core system           |
| `area::ui`             | #17becf | User Interface           | UI/UX work            |
| `area::security`       | #9467bd | Security                 | Security features     |
| `area::runtime`        | #bcbd22 | Runtime engine           | Runtime system        |
| `area::integrations`   | #ffbb78 | Third-party integrations | External services     |
| `area::infrastructure` | #98df8a | Infrastructure           | DevOps/Infrastructure |

#### Status Labels

| Label                 | Color   | Description               | Usage                |
| --------------------- | ------- | ------------------------- | -------------------- |
| `status::triage`      | #c7c7c7 | Needs initial assessment  | Newly created issues |
| `status::in-progress` | #6f42c1 | Currently being worked on | Active development   |
| `status::review`      | #fd7e14 | Needs code review         | Ready for review     |
| `status::blocked`     | #dc3545 | Cannot proceed            | Blocked issues       |
| `status::done`        | #28a745 | Completed                 | Finished work        |

#### Epic Labels

| Label      | Color   | Description          | Usage                     |
| ---------- | ------- | -------------------- | ------------------------- |
| `epic::e1` | #6f42c1 | E1 Foundation        | Foundation Sprint         |
| `epic::e2` | #fd7e14 | E2 Advanced Features | Advanced Features         |
| `epic::e3` | #20c997 | E3 Multi-Panel       | Multi-Panel Runtime       |
| `epic::e4` | #ffc107 | E4 ScriptRunner      | ScriptRunner Architecture |
| `epic::e5` | #17a2b8 | E5 Security          | Security & Modal System   |

### 🎯 Milestones Estratégicos (4 milestones)

#### v1.0 Alpha - 19/12/2025

- **Issues**: #001-#008 (Foundation Sprint)
- **Goal**: MVP funcional com foundations básicas
- **Duration**: 3-4 semanas
- **Deliverables**: CLI Foundation, YAML Parsing, TUI Structure, Form Rendering

#### v1.0 Beta - 15/01/2026  

- **Issues**: #009-#010 (Advanced Features)
- **Goal**: Funcionalidades avançadas e UX melhorada
- **Duration**: 2-3 semanas após Alpha
- **Deliverables**: Advanced Form Types, Enhanced Error Handling

#### v2.0 Features - 20/02/2026

- **Issues**: #011-#016 (Runtime Architecture)  
- **Goal**: Runtime architecture completa
- **Duration**: 4-5 semanas
- **Deliverables**: Event Engine, ScriptRunner, Modal Stack, Security Hardening

#### v2.0 Extensions - 15/03/2026

- **Issues**: #017-#025 (Extensions & Integrations)
- **Goal**: Integrações e extensões avançadas  
- **Duration**: 3-4 semanas
- **Deliverables**: Ansible Integration, SSH Server, Plugin Architecture

### 📊 GitHub Project: "Shantilly Roadmap"

#### Project Structure

- **Layout**: Table with 7 columns
- **Columns**: Backlog | Prioritized | To Do | In Progress | Review | Blocked | Done
- **Auto-add**: Issues automatically added when created
- **Automation Rules**: 15+ automation rules configured

#### Automation Rules

1. **Auto-move to "In Progress"** when assigned
2. **Auto-move to "Review"** when labeled "status::review"  
3. **Auto-move to "Done"** when closed
4. **Auto-add to "Blocked"** when labeled "status::blocked"
5. **Auto-prioritize** when labeled "priority::critical"
6. **Filter by Milestones** for sprint planning
7. **Group by Epic** for roadmap tracking

#### Views Configuration

- **Sprint View**: Filter by current milestone
- **Epic View**: Group by epic labels
- **Priority View**: Sort by priority labels
- **Team View**: Filter by assignees
- **Blockers View**: Show only blocked issues

## 🚀 Implementação no GitHub

### Opção 1: Interface Web (Recomendado para configuração manual)

#### Labels Implementation

1. Acessar repository `helton-godoy/shantilly`
2. Ir em **Issues** → **Labels** → **New label**
3. Criar labels conforme tabela acima
4. Aplicar cores e descrições correspondentes

#### Milestones Implementation  

1. Ir em **Issues** → **Milestones** → **New milestone**
2. Criar milestones com títulos e deadlines
3. Aplicar descrições e estados

#### Project Implementation

1. **Settings** → **Features** → **Projects** → **Enable projects**
2. **Projects** → **New project** → **Table**
3. Configurar columns conforme estrutura
4. Configurar automations rules
5. Adicionar views personalizadas

### Opção 2: GitHub CLI (gh) - Para Implementação Rápida

#### Instalar e Configurar gh CLI

```bash
# Instalar gh CLI
sudo apt install gh  # Ubuntu/Debian
brew install gh      # macOS
winget install --id GitHub.cli  # Windows

# Autenticar
gh auth login
```

#### Criar Labels em Lote

```bash
#!/bin/bash

# Priority Labels
gh label create "priority::critical" --color "d73a4a" --description "Critical priority - must fix immediately"
gh label create "priority::high" --color "fb8500" --description "High priority - important for next release"  
gh label create "priority::medium" --color "fbbf24" --description "Medium priority - nice to have"
gh label create "priority::low" --color "10b981" --description "Low priority - future consideration"

# Type Labels
gh label create "type::feature" --color "1f77b4" --description "New functionality"
gh label create "type::bug" --color "d62728" --description "Something isn't working"
gh label create "type::refactor" --color "ff7f0e" --description "Refactoring code"
gh label create "type::task" --color "2ca02c" --description "Non-code related tasks"
gh label create "type::security" --color "e377c2" --description "Security related"
gh label create "type::documentation" --color "9467bd" --description "Documentation changes"
gh label create "type::optimization" --color "bcbd22" --description "Performance improvements"

# Area Labels
gh label create "area::core" --color "8c564b" --description "Core functionality"
gh label create "area::ui" --color "17becf" --description "User Interface"
gh label create "area::security" --color "9467bd" --description "Security"
gh label create "area::runtime" --color "bcbd22" --description "Runtime engine"
gh label create "area::integrations" --color "ffbb78" --description "Third-party integrations"
gh label create "area::infrastructure" --color "98df8a" --description "Infrastructure"

# Status Labels
gh label create "status::triage" --color "c7c7c7" --description "Needs initial assessment"
gh label create "status::in-progress" --color "6f42c1" --description "Currently being worked on"
gh label create "status::review" --color "fd7e14" --description "Needs code review"
gh label create "status::blocked" --color "dc3545" --description "Cannot proceed"
gh label create "status::done" --color "28a745" --description "Completed"

# Epic Labels
gh label create "epic::e1" --color "6f42c1" --description "E1 Foundation"
gh label create "epic::e2" --color "fd7e14" --description "E2 Advanced Features"
gh label create "epic::e3" --color "20c997" --description "E3 Multi-Panel"
gh label create "epic::e4" --color "ffc107" --description "E4 ScriptRunner"
gh label create "epic::e5" --color "17a2b8" --description "E5 Security"
```

#### Criar Milestones

```bash
# v1.0 Alpha
gh milestone create "v1.0 Alpha" \
  --title "v1.0 Alpha" \
  --description "MVP funcional com foundations básicas" \
  --due-on "2025-12-19"

# v1.0 Beta  
gh milestone create "v1.0 Beta" \
  --title "v1.0 Beta" \
  --description "Funcionalidades avançadas e UX melhorada" \
  --due-on "2026-01-15"

# v2.0 Features
gh milestone create "v2.0 Features" \
  --title "v2.0 Features" \
  --description "Runtime architecture completa" \
  --due-on "2026-02-20"

# v2.0 Extensions
gh milestone create "v2.0 Extensions" \
  --title "v2.0 Extensions" \
  --description "Integrações e extensões avançadas" \
  --due-on "2026-03-15"
```

#### Criar GitHub Project

```bash
# Criar projeto via gh CLI
gh project create "Shantilly Roadmap" --owner "helton-godoy/shantilly"

# Adicionar colunas (via web interface)
# - Backlog
# - Prioritized  
# - To Do
# - In Progress
# - Review
# - Blocked
# - Done
```

### Opção 3: GitHub API Script (Para importação completa)

#### Script de Configuração Automatizada

```bash
#!/bin/bash

# GitHub Structures Configuration Script
REPO="helton-godoy/shantilly" 
TOKEN="YOUR_GITHUB_TOKEN"

# Function to create label
create_label() {
    local name="$1"
    local color="$2" 
    local description="$3"
    
    curl -X POST \
      -H "Authorization: token $TOKEN" \
      -H "Accept: application/vnd.github.v3+json" \
      https://api.github.com/repos/$REPO/labels \
      -d "{
        \"name\": \"$name\",
        \"color\": \"$color\",
        \"description\": \"$description\"
      }"
}

# Function to create milestone
create_milestone() {
    local title="$1"
    local description="$2"
    local due_on="$3"
    
    curl -X POST \
      -H "Authorization: token $TOKEN" \
      -H "Accept: application/vnd.github.v3+json" \
      https://api.github.com/repos/$REPO/milestones \
      -d "{
        \"title\": \"$title\",
        \"description\": \"$description\",
        \"due_on\": \"$due_on\"
      }"
}

# Executar configurações
echo "Creating labels..."
create_label "priority::critical" "d73a4a" "Critical priority"
create_label "priority::high" "fb8500" "High priority"
create_label "type::feature" "1f77b4" "New functionality"
create_label "area::core" "8c564b" "Core functionality"

echo "Creating milestones..."
create_milestone "v1.0 Alpha" "MVP funcional com foundations básicas" "2025-12-19T23:59:59Z"
create_milestone "v1.0 Beta" "Funcionalidades avançadas e UX melhorada" "2026-01-15T23:59:59Z"

echo "Configuration complete!"
```

## 🎯 Próximos Passos

### Imediato (Sequência Recomendada)

1. **Labels**: Criar todas as labels hierárquicas primeiro
2. **Milestones**: Configurar milestones com deadlines
3. **Project**: Criar "Shantilly Roadmap" project
4. **Automations**: Configurar automation rules
5. **Views**: Criar views personalizadas

### Após Estruturas Configuradas

1. **Validar estrutura** - testar labels, milestones, project
2. **Importar issues** - usar issues estruturadas preparadas
3. **Sincronizar** - mover issues para project correto
4. **Configurar automations** - ativar auto-labeling e auto-movement

### Validação Final

- **Labels aplicados** em todas as 25 issues
- **Milestones assignados** por sprint
- **Issues no project** "Shantilly Roadmap"
- **Automation rules** funcionando
- **Views configuradas** para sprint planning

## 🎯 Benefícios Esperados

### Para o Projeto

- **Organização**: Issues categorizadas e priorizadas
- **Visibilidade**: Roadmap claro e progress tracking
- **Automation**: Redução de trabalho manual
- **Collaboration**: Equipe pode se auto-organizar

### Para Contribuidores

- **Clareza**: Template e guidelines claras
- **Self-service**: Auto-assign e tracking
- **Focus**: Milestones guiam desenvolvimento
- **Quality**: Automation mantém padrões

### Para Gestão

- **Metrics**: Progress tracking e velocity
- **Planning**: Sprint planning estruturado
- **Transparency**: Status visível para stakeholders
- **Efficiency**: Menos overhead administrativo

---

**⚠️ Important**: Implementar labels primeiro, depois milestones, depois project. Sempre testar automations antes de usar em produção.
