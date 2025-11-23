# GitHub Issues Implementation Guide - Shantilly

## 📋 Status: Issues Preparadas para Importação

### ✅ Material Preparado
- **25 issues estruturadas** organizadas por épicos (github_roadmap_issues.md)
- **Labels hierárquicas** configuradas (.github/labels/)
- **Milestones** definidos (.github/milestones.yml)
- **Templates** personalizados (.github/ISSUE_TEMPLATE/)

### 📋 Estrutura das Issues

#### Foundation Sprint (v1.0 Alpha - 3-4 semanas)
- **Issues #001-#008**: CLI Foundation, YAML Parsing, TUI Structure, Form Rendering, Navigation, Submission, Styling, Build

#### Advanced Features Sprint (v1.0 Beta)
- **Issues #009-#010**: Advanced Form Types, Enhanced Error Handling

#### Runtime Architecture Sprint (v2.0 Features)
- **Issues #011-#016**: Multi-Panel Navigation, Event Engine, ScriptRunner, Modal Stack, Legacy Encapsulation, Security Hardening

#### Extensions & Integrations
- **Issues #017-#025**: Ansible Integration, SSH Server, Predictive Components, Layout System, Plugin Architecture, Multi-tenant, Security Policies, Performance, Documentation Portal

## 🚀 Implementação no GitHub

### Opção 1: Interface Web (Recomendado para teste inicial)

#### Passo 1: Criar Issues Manualmente
1. Acessar repository `helton-godoy/shantilly`
2. Ir em **Issues** → **New issue**
3. Usar templates apropriados para cada tipo
4. Aplicar labels e milestons correspondentes

#### Passo 2: Aplicar Labels em Lote
1. Configurar labels hierárquicas conforme `.github/labels/`
2. Labels de prioridade: `priority::critical`, `priority::high`, etc.
3. Labels de tipo: `type::feature`, `type::bug`, etc.
4. Labels de área: `area::core`, `area::ui`, etc.

### Opção 2: GitHub CLI (gh) - Para Importação em Lote

#### Instalação e Setup
```bash
# Instalar gh (se não tiver)
# Ubuntu/Debian
sudo apt install gh

# macOS
brew install gh

# Windows (via winget)
winget install --id GitHub.cli

# Autenticar
gh auth login
```

#### Criar Labels
```bash
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

# Area Labels
gh label create "area::core" --color "8c564b" --description "Core functionality"
gh label create "area::ui" --color "17becf" --description "User Interface"
gh label create "area::security" --color "9467bd" --description "Security"
gh label create "area::runtime" --color "bcbd22" --description "Runtime engine"
```

#### Criar Milestones
```bash
# v1.0 Alpha
gh milestone create "v1.0 Alpha" --title "v1.0 Alpha" --description "MVP funcional com foundations básicas" --due-on "2025-12-19"

# v1.0 Beta
gh milestone create "v1.0 Beta" --title "v1.0 Beta" --description "Funcionalidades avançadas e UX melhorada" --due-on "2026-01-15"

# v2.0 Features
gh milestone create "v2.0 Features" --title "v2.0 Features" --description "Runtime architecture completa" --due-on "2026-02-20"

# v2.0 Extensions
gh milestone create "v2.0 Extensions" --title "v2.0 Extensions" --description "Integrações e extensões avançadas" --due-on "2026-03-15"
```

#### Criar Issues via CLI
```bash
# Exemplo: CLI Foundation Issue
gh issue create \
  --title "CLI Foundation & Project Structure" \
  --body "$(cat <<'EOF'
## Description
Implementar estrutura CLI básica com comandos foundation, configuração inicial e estrutura de diretórios para o projeto Shantilly.

## Acceptance Criteria
- CLI aceita comandos básicos (help, version, init)
- Estrutura de diretórios criada automaticamente  
- Configuração inicial gerada

## Labels
- type::feature
- priority::high
- area::core
- epic::e1

## Story Points
5
EOF
)" \
  --label "type::feature,priority::high,area::core,epic::e1" \
  --milestone "v1.0 Alpha"
```

### Opção 3: GitHub API Script (Recomendado para importação completa)

#### Script de Importação em Lote
```bash
#!/bin/bash

# GitHub Issues Import Script
REPO="helton-godoy/shantilly"
TOKEN="YOUR_GITHUB_TOKEN"  # Personal Access Token

# Function to create issue
create_issue() {
    local title="$1"
    local body="$2"
    local labels="$3"
    local milestone="$4"
    
    curl -X POST \
      -H "Authorization: token $TOKEN" \
      -H "Accept: application/vnd.github.v3+json" \
      https://api.github.com/repos/$REPO/issues \
      -d "{
        \"title\": \"$title\",
        \"body\": $(echo "$body" | jq -Rs .),
        \"labels\": [$labels],
        \"milestone\": $(echo "$milestone" | jq -Rs .)
      }"
}

# Criar Foundation Sprint Issues (001-008)
for i in {001..008}; do
    # Issues específicas com dados detalhados
    create_issue \
      "CLI Foundation & Project Structure" \
      "## Description\nImplementar estrutura CLI básica..." \
      "[\"type::feature\",\"priority::high\",\"area::core\",\"epic::e1\"]" \
      "1"  # milestone ID
done
```

### Opção 4: Import via CSV/JSON

#### Preparar Dados para Importação
```json
[
  {
    "title": "CLI Foundation & Project Structure",
    "body": "## Description\nImplementar estrutura CLI básica...",
    "labels": ["type::feature", "priority::high", "area::core"],
    "milestone": "v1.0 Alpha",
    "assignees": []
  }
]
```

## 📊 GitHub Project: "Shantilly Roadmap"

### Configuração Recomendada
- **Columns**: Backlog | Prioritized | To Do | In Progress | Review | Blocked | Done
- **Automation Rules**:
  - Auto-move to "In Progress" when assigned
  - Auto-move to "Review" when labeled "status::review"
  - Auto-move to "Done" when closed
  - Auto-add to "Blocked" when labeled "status::blocked"

### Como Configurar
1. **Settings** → **Features** → **Projects** → **Enable projects**
2. **Create a project** → **Table** layout
3. **Add automation** → **Move issues when...**
4. **Configure views** → **Filter por Labels/Milestones**

## 🎯 Próximos Passos

### Imediato
1. **Escolher método de importação** (recomendo Option 2 - gh CLI)
2. **Configurar labels e milestones** primeiro
3. **Importar issues em batches** (Foundation → Advanced → Extensions)
4. **Configurar GitHub Project** com automações

### Após Import
1. **Validar estrutura** - verificar labels, milestones, templates
2. **Importar no GitHub Project** - adicionar todas as issues
3. **Configurar automações** - auto-labeling, auto-movement
4. **Documentar workflow** - guidelines para contributors

### Benefícios Esperados
- **Visibilidade**: Roadmap claro e organizado
- **Triage**: Issues automaticamente categorizadas
- **Progress Tracking**: Milestones com deadlines
- **Automation**: Redução manual de trabalho
- **Contributor Experience**: Templates e guidelines claras

---

**⚠️ Important**: Sempre fazer backup antes de importação em massa. Usar gh CLI para melhor controle e rollback se necessário.
