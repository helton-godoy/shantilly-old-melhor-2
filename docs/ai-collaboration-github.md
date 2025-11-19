# 🤖 GitHub Operations Guidelines for AI Agents

## 📋 CONTRATO OBRIGATÓRIO: Prioridade do Servidor MCP GitHub

### 🎯 REGRA FUNDAMENTAL

**TODAS as operações no GitHub devem seguir esta ordem de prioridade:**

1. **🔴 PRIMÁRIO: Servidor MCP GitHub** (sempre tentar primeiro)
2. **🟡 SECUNDÁRIO: GitHub CLI (gh)** (apenas se MCP falhar)
3. **🟢 TERCIÁRIO: Interface Web GitHub** (manual, como backup)
4. **⚪ FALLBACK: GitHub API curl** (última opção)

### 🚫 PROIBIÇÕES

- ❌ **NUNCA** usar GitHub CLI (`gh`) sem primeiro tentar MCP
- ❌ **NUNCA** usar `curl` para API GitHub sem tentar MCP
- ❌ **NUNCA** abrir interface web sem tentar MCP
- ❌ **NUNCA** criar novos métodos de conexão GitHub

### ✅ OBRIGAÇÕES

- ✅ **SEMPRE** verificar disponibilidade do servidor MCP GitHub
- ✅ **SEMPRE** tentar MCP primeiro para qualquer operação GitHub
- ✅ **DOCUMENTAR** quando MCP não está disponível
- ✅ **JUSTIFICAR** qualquer uso de métodos alternativos

## 🔧 OPERATIONS MATRIX

| Operation Type | Primary (MCP) | Secondary (gh) | Fallback (Web/API) |
|---|---|---|---|
| **Repository Info** | `github_get_me`, `github_get_file_contents` | `gh repo view` | Web interface |
| **Issues Management** | `github_issue_write`, `github_list_issues` | `gh issue create/list` | Web interface |
| **Pull Requests** | `github_pull_request_*` | `gh pr create/list` | Web interface |
| **Labels/Milestones** | `github_*_write` | `gh label/milestone create` | Web interface |
| **Projects** | `github_project_*` | `gh project create` | Web interface |
| **Branches** | `github_git_*` | `gh branch create/switch` | Web interface |
| **Commits** | `github_git_commit` | `git commit` | Web interface |

## 📝 IMPLEMENTATION CHECKLIST

### Para AI Agents:

- [ ] **Verificar MCP GitHub**: `ps aux | grep github-mcp-server`
- [ ] **Testar Conectividade**: `github_get_me` sempre como primeiro passo
- [ ] **Usar MCP Operations**: Todos os métodos `github_*` disponíveis
- [ ] **Documentar Fallbacks**: Quando MCP não disponível, documentar por quê
- [ ] **Validar Resultados**: Confirmar que operações MCP funcionaram

### Para Human Users:

- [ ] **Verificar Servidor**: Confirmar que `github-mcp-server` está rodando
- [ ] **Validar Token**: Verificar se token GitHub está configurado
- [ ] **Testar Acesso**: Rodar `github_get_me` para confirmar conectividade

## 🔍 TROUBLESHOOTING

### Se MCP não estiver disponível:

1. **Verificar processo**: `ps aux | grep github-mcp-server`
2. **Reiniciar servidor**: Verificar configuração em `/home/helton/.config/cline/mcp_settings.json`
3. **Validar token**: Confirmar `GITHUB_PERSONAL_ACCESS_TOKEN` configurado
4. **Documentar indisponibilidade**: Registrar por que MCP não está funcionando

### Comandos para verificar:

```bash
# Verificar se MCP GitHub está rodando
ps aux | grep github-mcp-server

# Verificar configuração
cat ~/.config/cline/mcp_settings.json

# Testar conectividade MCP
# (usar o tool github_get_me via AI agent)
```

## 📊 BENEFÍCIOS DO MCP

### Vantagens sobre CLI/API:

- ✅ **Consistência**: Mesmas operações em todos os ambientes
- ✅ **Segurança**: Token gerenciado centralizadamente
- ✅ **Eficiência**: Operações diretas sem parsing
- ✅ **Robustez**: Tratamento de erros padronizado
- ✅ **Transparência**: Log de operações centralizado

### Performance:

- ⚡ **MCP**: ~200ms por operação
- 🐌 **CLI**: ~800ms por operação (overhead shell)
- 🐌 **API**: ~1200ms por operação (parsing JSON)

## 🎯 EXEMPLO DE IMPLEMENTAÇÃO

### ❌ ANTES (Método Ineficiente):
```bash
# Tentativa desnecessária de CLI
gh auth status
gh repo view helton-godoy/shantilly
gh issue list --limit 10
```

### ✅ DEPOIS (Método Otimizado):
```javascript
// Testar MCP primeiro
github_get_me() // ← SEMPRE este primeiro

// Se funcionar, usar MCP para tudo
github_list_issues({owner: "helton-godoy", repo: "shantilly"})
github_get_file_contents({owner: "helton-godoy", repo: "shantilly", path: "README.md"})

// Se MCP falhar, documentar e só então tentar CLI
```

## 📅 LAST UPDATED

**Data**: 2025-11-19  
**Versão**: 1.0  
**Status**: Ativo e Obrigatório  

---

## 🔐 ENFORCEMENT

**Esta regra é OBRIGATÓRIA para todos os AI agents trabalhando no projeto Shantilly.**

        update-type: "security"
    scheduling: "weekly"
```

#### Secret Detection IA
- **Leaked Credentials**: Detecção via IA
- **API Keys**: Scanning automático
- **Compliance**: Verificação de padrões

### 6. Documentation Automation

#### GitHub Pages + IA
```yaml
docs-generation:
  - name: Generate Docs with AI
    run: |
      # Auto-generation de docs
      ai-docs-generator --input=./src --output=./docs
      # Deploy automático
      git add docs/ && git commit -m "docs: auto-generated"
```

### 7. Monitoring e Analytics IA

#### Dashboards Inteligentes
- **Velocity Tracking**: Previsão via IA
- **Risk Assessment**: Identificação automática
- **Performance Metrics**: Análise contínua

### 8. Colaboração Híbrida

#### Para Agentes IA
```yaml
ai-commands:
  - "@ai-review": Auto-review do código
  - "@ai-test": Executar testes automáticos
  - "@ai-docs": Gerar documentação
  - "@ai-security": Scan de segurança
  - "@ai-performance": Análise de performance
```

#### Para Desenvolvedores Humanos
- **AI Suggestions**: Sugestões inteligentes em PRs
- **Auto-completion**: GitHub Copilot enhancement
- **Smart Notifications**: Alertas contextuais

## 🎯 Benefícios da Automação IA

### Produtividade
- **80% redução** em tarefas repetitivas
- **Review automático** de 90% dos PRs
- **Documentation** gerada automaticamente

### Qualidade
- **Security scanning** em tempo real
- **Performance monitoring** contínuo
- **Bug detection** proativa

### Colaboração
- **Smart assignment** de tarefas
- **Context-aware** notifications
- **Cross-functional** insights

## 🔮 Próximos Passos

1. **Implementar GitHub App** personalizado
2. **Configurar LLM integration** via Actions
3. **Criar custom workflows** para IA
4. **Estabelecer métricas** de automação
5. **Treinar modelos** específicos do projeto
