# Contributing to Shantilly

Primeiro, muito obrigado pelo interesse em contribuir para o Shantilly! Suas contribuições são valiosas para tornar o Shantilly ainda melhor.

## 🚀 Como Contribuir

### Reportar Bugs

Se você encontrou um bug, por favor [crie uma issue](https://github.com/helton-godoy/shantilly/issues/new/choose) usando o template "🐛 Bug Report". Inclua:

- Versão do Shantilly que você está usando
- Sistema operacional
- Passos para reproduzir o bug
- Configuração YAML (se aplicável)
- Logs de erro relevantes

### Sugerir Funcionalidades

Para sugerir uma nova funcionalidade, [crie uma issue](https://github.com/helton-godoy/shantilly/issues/new/choose) usando o template "✨ Feature Request". Inclua:

- Problema que a funcionalidade resolveria
- Solução proposta
- User stories
- Critérios de aceitação

### Contribuições de Código

#### Pré-requisitos

- Go 1.21+ instalado
- Conhecimento básico de Go e TUI development
- Ambiente de desenvolvimento configurado

#### Setup do Ambiente

1. **Fork e Clone:**
   ```bash
   git clone https://github.com/SEU_USUARIO/shantilly.git
   cd shantilly
   ```

2. **Adicione o upstream:**
   ```bash
   git remote add upstream https://github.com/helton-godoy/shantilly.git
   ```

3. **Instale dependências:**
   ```bash
   go mod tidy
   ```

#### Padrões de Desenvolvimento

1. **Branches:**
   - Use `feat/issue-XXX-descricao` para novas funcionalidades
   - Use `fix/issue-XXX-descricao` para correções
   - Use `docs/issue-XXX-descricao` para documentação

2. **Commits:**
   - Siga [Conventional Commits](https://www.conventionalcommits.org/):
     - `feat: add new form validation feature`
     - `fix: resolve yaml parsing error`
     - `docs: update installation guide`

3. **Código:**
   - Execute `golangci-lint run ./...` antes de commit
   - Mantenha funções curtas e focadas
   - Adicione testes para novas funcionalidades

#### Processo de Pull Request

1. **Crie uma branch:**
   ```bash
   git checkout -b feat/issue-XXX-minha-feature
   ```

2. **Desenvolva e teste:**
   ```bash
   # Execute linters
   golangci-lint run ./...
   
   # Execute testes
   go test ./...
   
   # Teste manualmente com examples
   go run . form examples/basic_form.yaml
   ```

3. **Commit suas mudanças:**
   ```bash
   git add .
   git commit -m "feat: implement feature description"
   ```

4. **Push e PR:**
   ```bash
   git push origin feat/issue-XXX-minha-feature
   ```

5. **Crie Pull Request** com:
   - Título descritivo seguindo conventional commits
   - Descrição detalhada das mudanças
   - Screenshots se aplicável
   - Referência à issue relacionada

#### Critérios de Aceitação do PR

- ✅ Código passa em todos os linters (`golangci-lint`)
- ✅ Todos os testes passam
- ✅ Cobertura de teste não diminui
- ✅ Documentação atualizada se necessário
- ✅ Mínimo 1 approval de reviewer
- ✅ Status checks do GitHub Actions passam

## 📋 Padrões de Código

### Qualidade do Código

O projeto usa `golangci-lint` para garantir qualidade. O arquivo `.golangci.yml` contém todas as regras. Principais pontos:

- **Tratamento de Erros:** Todos os erros devem ser tratados
- **Performance:** Pré-alocação de slices é obrigatória
- **Complexidade:** Funções devem ser curtas e focadas
- **Estilo:** Código deve seguir convenções idiomáticas do Go

### Estrutura do Projeto

```
shantilly/
├── cmd/                    # Comandos CLI
│   └── shantilly/         # Entry point
├── internal/              # Código interno
│   ├── config/           # Configuração e parsing
│   ├── runtime/          # Runtime TUI
│   ├── tui/              # Componentes TUI
│   └── util/             # Utilitários
├── pkg/                   # Bibliotecas públicas
├── docs/                  # Documentação
├── examples/              # Exemplos YAML
└── scripts/               # Scripts auxiliares
```

### Testes

- **Unit Tests:** Testes de funções individuais
- **Integration Tests:** Testes de fluxo end-to-end
- **Example Tests:** Validação de exemplos YAML

Execute todos os testes:
```bash
go test ./...
```

### Documentação

- **README.md:** Visão geral e quickstart
- **docs/:** Documentação detalhada
- **examples/:** Exemplos práticos
- **Comments:** Código bem comentado

## 🎯 Áreas que Precisam de Contribuição

### Prioridade Alta

- [ ] **Componentes TUI Avançados:** Date picker, file browser
- [ ] **Layout System:** Multi-panel layouts, tabs
- [ ] **Performance:** Otimização para grandes formulários
- [ ] **Cross-platform:** Melhorias no Windows/WSL

### Prioridade Média

- [ ] **Temas:** Sistema de theming avançado
- [ ] **Animations:** Transições suaves
- [ ] **SSH Mode:** Suporte para uso remoto
- [ ] **i18n:** Internacionalização

### Iniciante-Friendly

- [ ] **Documentação:** Melhorar guias e exemplos
- [ ] **Tests:** Adicionar testes para cobertura
- [ ] **Examples:** Criar mais exemplos YAML
- [ ] **CI/CD:** Otimizar workflows

## 📞 Comunicação

- **Issues:** Para bugs e features
- **Discussions:** Para perguntas e brainstorming
- **PR Reviews:** Feedback em pull requests

## 📜 Código de Conduta

Este projeto segue nosso [Código de Conduta](CODE_OF_CONDUCT.md). Ao participar, você concorda em seguir estas diretrizes.

## 🎉 Reconhecimento

Contribuidores serão reconhecidos em:
- README.md do projeto
- Releases notes
- Documentação

Obrigado por tornar o Shantilly melhor! 🚀
