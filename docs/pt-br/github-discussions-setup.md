# GitHub Discussions Setup - Shantilly

## 📋 Configuração de Discussions

### Status: Configuração Local Criada

- ✅ `.github/discussions.yml` criado com 4 categorias
- 📋 Prontas para implementação no repositório GitHub

### 📋 Categorias Configuradas

- **Ideas** (💡)
  - Descrição: Compartilhar ideias de features e melhorias
  - Cor: `#f2d604`
  - Tipo: Discussion

- **Q&A** (❓)
  - Descrição: Fazer perguntas e obter ajuda
  - Cor: `#c5def5`
  - Tipo: Question

- **Polls** (📊)
  - Descrição: Pesquisas e feedback da comunidade
  - Cor: `#fbca04`
  - Tipo: Poll

- **General** (💬)
  - Descrição: Discussões gerais e anúncios
  - Cor: `#d4c5f9`
  - Tipo: Discussion

## 🚀 Implementação no GitHub

### Opção 1: Interface Web (Recomendado)

1. Acessar repository `helton-godoy/shantilly`
2. Ir em **Settings** → **Features** → **Discussions**
3. Ativar **Discussions**
4. Configurar categorias manualmente conforme `.github/discussions.yml`

### Opção 2: GitHub CLI (gh)

```bash
# Instalar gh (se não tiver)
# Authenticate
gh auth login

# Habilitar discussions (se necessário)
gh api repos/helton-godoy/shantilly --method PATCH -F discussions_category_required=true

# Criar categorias
gh api repos/helton-godoy/shantilly/discussions/categories \
  --method POST \
  -F name="Ideas 💡" \
  -F about="Share and discuss new features, enhancements, and creative ideas for Shantilly" \
  -F color="f2d604" \
  -F slug="ideas"
```

### Opção 3: GitHub API

```bash
# Usando Personal Access Token
curl -X POST \
  -H "Authorization: token YOUR_TOKEN" \
  -H "Accept: application/vnd.github.v3+json" \
  https://api.github.com/repos/helton-godoy/shantilly/discussions/categories \
  -d '{
    "name": "Ideas 💡",
    "about": "Share and discuss new features, enhancements, and creative ideas for Shantilly", 
    "color": "f2d604",
    "slug": "ideas"
  }'
```

## 📖 Diretrizes de Uso das Discussions

### 💡 Ideas Category

**Propósito**: Brainstorming e planejamento de features
**Quando usar**:

- Propor novas funcionalidades para Shantilly
- Discutir melhorias na arquitetura
- Sugerir integrações e extensões
- Compartilhar visões de longo prazo

**Exemplos de posts**:

- "Nova feature: Sistema de templates personalizáveis"
- "Integração com Supabase para persistência"
- "Roadmap para suportar múltiplas linguagens"

### ❓ Q&A Category

**Propósito**: Suporte e aprendizado
**Quando usar**:

- Dúvidas sobre instalação e uso
- Problemas com configuração
- Explicação de conceitos técnicos
- Como contribuir para o projeto

**Exemplos de posts**:

- "Como criar um novo form component?"
- "Diferença entre runtime-v2 e runtime-migration"
- "Como configurar um ambiente de desenvolvimento?"

### 📊 Polls Category

**Propósito**: Coleta de feedback e decisões
**Quando usar**:

- Priorizar funcionalidades
- Decidir arquitetura
- Validar decisões técnicas
- Medir satisfação da comunidade

**Exemplos de polls**:

- "Qual runtime você prefere?"
- "Priorizar: velocidade vs funcionalidades?"
- "Suportar CLI ou só TUI?"

### 💬 General Category

**Propósito**: Comunicação geral
**Quando usar**:

- Anúncios de releases
- Discussões da comunidade
- Updates de desenvolvimento
- Eventos e meetups

**Exemplos de posts**:

- "v0.1 Released! Check it out"
- "Weekly Development Update"
- "GitHub Repository Reorganized"

## 🎯 Benefícios Esperados

### Para Desenvolvedores

- **Centralização**: Discussões concentradas no repositório
- **Histórico**: Thread discussions vs scattered Discord/Slack
- **Integração**: Links diretos com Issues/PRs
- **Colaboração**: Contributions públicas e transparente

### Para Usuários

- **Acesso**: Não precisa de Discord/Slack para participar
- **Notificações**: Updates via GitHub notifications
- **Pesquisa**: Buscável via GitHub search
- **Documentação**: Discussions se tornam conhecimento público

## 🔄 Próximos Passos

1. **Implementar Discussions** no repositório GitHub
2. **Criar post inicial** explicando categorias e uso
3. **Adicionar no README.md** link para Discussions
4. **Configurar auto-moderation** se necessário
5. **Documentar na Wiki** guidelines detalhadas

---

**⚠️ Note**: Discussions require repository admin privileges to set up initially.
