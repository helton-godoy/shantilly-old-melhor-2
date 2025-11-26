# Shantilly - Runtime TUI Declarativo

[![Build and Test](https://github.com/helton-godoy/shantilly/actions/workflows/build.yml/badge.svg)](https://github.com/helton-godoy/shantilly/actions/workflows/build.yml)
[![Lint](https://github.com/helton-godoy/shantilly/actions/workflows/lint.yml/badge.svg)](https://github.com/helton-godoy/shantilly/actions/workflows/lint.yml)
[![codecov](https://codecov.io/gh/helton-godoy/shantilly/branch/main/graph/badge.svg)](https://codecov.io/gh/helton-godoy/shantilly)
[![Go Report Card](https://goreportcard.com/badge/github.com/helton-godoy/shantilly)](https://goreportcard.com/report/github.com/helton-godoy/shantilly)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.21-00ADD8.svg?logo=go)](https://golang.org/doc/devel/release)
[![Platform](https://img.shields.io/badge/platform-linux%20macos%20windows-lightgrey.svg?logo=linux&logoColor=white)](https://github.com/helton-godoy/shantilly/releases)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Contributors](https://img.shields.io/github/contributors/helton-godoy/shantilly.svg)](https://github.com/helton-godoy/shantilly/graphs/contributors)
[![GitHub stars](https://img.shields.io/github/stars/shantilly.svg?style=social&label=Star)](https://github.com/helton-godoy/shantilly)
[![GitHub forks](https://img.shields.io/github/forks/shantilly.svg?style=social&label=Fork)](https://github.com/helton-godoy/shantilly)

## 🌐 Site Oficial

**Visite [helton-godoy.github.io/shantilly](https://helton-godoy.github.io/shantilly) para documentação completa, demonstrações interativas e guias de início rápido.**

---

> **Runtime TUI Declarativo orientado a eventos** que consome YAML único para definir layout, componentes e lógica de automação.

Shantilly é uma ferramenta CLI escrita em Go que funciona como um intérprete para definições declarativas de TUI (YAML). Permite que desenvolvedores de scripts criem interfaces TUI modernas, ricas e portáveis diretamente de seus scripts shell.

## 🎯 Características Principais

- ✨ **Abordagem Declarativa**: Foque no YAML para definir a UI, separando *descrição* da *implementação*
- 🏗️ **Ecossistema Moderno**: Usa bibliotecas como `bubbletea`, `lipgloss`, `huh` para componentes ricos
- 📦 **Binário Único Portável**: Compilado estaticamente para Linux, macOS e Windows
- 🔗 **Integração Shell**: Projetado para ler configuração via `stdin` e retornar dados estruturados (JSON)
- 🎨 **Componentes Avançados**: Formulários com validação, seletores, layouts organizados, feedback visual
- 🖱️ **Suporte a Mouse**: Interação avançada com componentes interativos
- 🛡️ **Seguro por Padrão**: Execução sandboxed e validação de entrada

## 🚀 Instalação Rápida

### Via Homebrew (macOS/Linux)

```bash
brew tap helton-godoy/shantilly
brew install shantilly
```

### Via Go

```bash
go install github.com/helton-godoy/shantilly@latest
```

### Download Binário

Baixe o binário mais recente em [Releases](https://github.com/helton-godoy/shantilly/releases)

### Verificação

```bash
shantilly --version
```

## ⚡ Quick Start

### Exemplo Básico - Formulário Simples

```bash
# Crie um arquivo YAML
cat > formulario.yaml << 'EOF'
title: "Dados do Usuário"
fields:
  - key: "nome"
    label: "Nome"
    type: "input"
    required: true
    placeholder: "Seu nome completo"

  - key: "email"
    label: "Email"
    type: "input"
    required: true
    pattern: "email"
    placeholder: "seu@email.com"

  - key: "idade"
    label: "Idade"
    type: "number"
    required: true
    min: 0
    max: 120
EOF

# Execute o formulário
cat formulario.yaml | shantilly form

# Saída JSON
{
  "nome": "João Silva",
  "email": "joao@exemplo.com",
  "idade": 30
}
```

### Exemplo Avançado - Formulário com Validação

```yaml
title: "Configuração Avançada"
theme:
  accent: "cyan"
  primary: "blue"

fields:
  - key: "projeto"
    label: "Nome do Projeto"
    type: "input"
    required: true
    minLength: 3
    placeholder: "meu-projeto"

  - key: "ambiente"
    label: "Ambiente"
    type: "select"
    required: true
    options:
      - value: "dev"
        label: "Desenvolvimento"
      - value: "staging"
        label: "Staging"
      - value: "prod"
        label: "Produção"

  - key: "recursos"
    label: "Recursos Necessários"
    type: "multiselect"
    required: true
    options:
      - value: "database"
        label: "Banco de Dados"
      - value: "cache"
        label: "Cache Redis"
      - value: "cdn"
        label: "CDN"

  - key: "confirmar"
    label: "Confirma as configurações?"
    type: "confirm"
    required: true
```

### Integração com Scripts

```bash
#!/bin/bash

# Coleta dados do usuário via Shantilly
DADOS=$(cat config.yaml | shantilly form)

# Extrai valores via jq
PROJETO=$(echo "$DADOS" | jq -r '.projeto')
AMBIENTE=$(echo "$DADOS" | jq -r '.ambiente')

# Executa comando baseado nos dados
if [ "$AMBIENTE" = "prod" ]; then
    echo "Deploy para produção do projeto: $PROJETO"
    deploy_projeto.sh "$PROJETO"
else
    echo "Deploy para $AMBIENTE do projeto: $PROJETO"
    deploy_dev.sh "$PROJETO" "$AMBIENTE"
fi
```

## 📖 Documentação Completa

> **Nota:** A documentação está disponível em [Português (PT-BR)](../pt-br/) e [Inglês (EN)](../en/).

- **[User Guide](../pt-br/user-guide.md)** - Guia completo para usuários
- **[Developer Guide](../pt-br/contributing.md)** - Como contribuir para o projeto
- **[API Reference](../pt-br/api-reference.md)** - Documentação detalhada da API
- **[Examples](../../examples/)** - Exemplos práticos de uso
- **[Architecture](../pt-br/architecture.md)** - Arquitetura do projeto

## 🛠️ Componentes Suportados

### Campos de Formulário

- `input` - Campo de texto simples
- `textarea` - Campo de texto multilinha  
- `select` - Seleção de uma opção
- `multiselect` - Seleção múltipla de opções
- `confirm` - Confirmação sim/não
- `note` - Texto informativo (somente leitura)
- `number` - Campo numérico com validação
- `date` - Campo de data com formatos flexíveis
- `file` - Seleção de arquivo com filtro de tipos

### Regras de Validação

- `required` - Campo obrigatório
- `min`/`max` - Limites numéricos
- `minLength`/`maxLength` - Limites de comprimento
- `pattern` - Padrões especiais (ex: email)
- `fileTypes` - Tipos de arquivo permitidos

### Recursos Avançados

- Layouts hierárquicos (columns, rows, boxes)
- Sistema de temas personalizável
- Suporte completo a mouse
- Modais e janelas flutuantes
- Animações suaves
- Execução de scripts integrada

## 🎯 Para Desenvolvedores

### Setup do Ambiente de Desenvolvimento

```bash
# Clone o repositório
git clone https://github.com/helton-godoy/shantilly.git
cd shantilly

# Instale dependências
go mod tidy

# Execute os linters
golangci-lint run ./...

# Execute os testes
go test ./...

# Execute o programa
go run .
```

### Guia para Agentes de IA

Olá! Você foi encarregado de desenvolver este projeto. Siga estas diretrizes para garantir que seu código seja de alta qualidade:

#### 1. Padrões de Qualidade do Código

A qualidade é garantida pelo `golangci-lint` com regras definidas em `.golangci.yml`:

- **Tratamento de Erros:** Todos os erros devem ser tratados
- **Performance:** Pré-alocação de slices é obrigatória
- **Complexidade:** Funções devem ser curtas e focadas
- **Estilo:** Código deve seguir convenções idiomáticas do Go

#### 2. Fluxo de Desenvolvimento

1. **Desenvolva:** Modifique arquivos Go, atenção aos comentários `// TODO: IA`
2. **Verificação CI:** GitHub Actions executam:
   - **Lint:** Verifica qualidade do código
   - **Build:** Garante compilação e testes
3. **Correção:** Ambos os checks devem estar verdes

#### 3. Processo de PR

- Use branch `feat/issue-XXX-descricao`
- Commits seguem [Conventional Commits](https://conventionalcommits.org/)
- Adicione testes para novas funcionalidades
- Documentação atualizada se necessário

## 🚀 Roadmap

- [x] **v1.0 MVP** - Formulários básicos e integração shell
- [ ] **v1.1** - Layouts avançados e componentes adicionais
- [ ] **v1.2** - Sistema de temas e personalização
- [ ] **v1.3** - Mouse support e interatividade avançada
- [ ] **v2.0** - Modo SSH e interface administrativa remota
- [ ] **v2.1** - Internacionalização e suporte multilíngue
- [ ] **v3.0** - Integrações avançadas e plugins

Veja [ROADMAP.md](../ROADMAP.md) para detalhes completos.

## 🤝 Contribuindo

Contribuições são bem-vindas! Por favor:

1. Leia nosso [Código de Conduta](../pt-br/code-of-conduct.md)
2. Siga o [Guia de Contribuição](../pt-br/contributing.md)
3. Crie issues para bugs e features
4. Use templates apropriados para issues
5. Mantenha qualidade do código com linters

## 📞 Suporte e Comunidade

- **Issues** - [Reportar bugs e solicitar funcionalidades](https://github.com/helton-godoy/shantilly/issues)
- **Discussions** - [Discussões e Q&A](https://github.com/helton-godoy/shantilly/discussions)
- **Security** - [Reportar vulnerabilidades](../../SECURITY.md)
- **Wiki** - [Documentação expandida](https://github.com/helton-godoy/shantilly/wiki)

## 📄 Licença

Este projeto está licenciado sob a Licença MIT - veja o arquivo [LICENSE](../../LICENSE) para detalhes.

## 🙏 Agradecimentos

- [Charmbracelet](https://charm.sh/) pelo ecossistema TUI incrível
- [Bubbletea](https://github.com/charmbracelet/bubbletea) - Framework TUI
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Sistema de styling
- [Huh](https://github.com/charmbracelet/huh) - Componentes de formulário
- Comunidade Go pelo suporte incrível

---

### ⭐ Apoie o projeto

Se este projeto te ajudou, considere dar uma estrela.

Feito com ❤️ pela comunidade Shantilly

Deploy trigger sáb 22 nov 2025 18:48:09 -04
