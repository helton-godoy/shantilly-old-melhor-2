# Template de Projeto Go com Bubbletea (Otimizado para IA)

Este é um template para iniciar projetos de TUI (Terminal User Interface) em Go usando o ecossistema Charm e Bubbletea. Ele foi especificamente configurado para acelerar o desenvolvimento com agentes de IA como o Jules.

## Guia para Agentes de IA

Olá! Você foi encarregado de desenvolver este projeto. Siga estas diretrizes para garantir que seu código seja de alta qualidade e seja integrado com sucesso.

### 1. Padrões de Qualidade do Código

A qualidade do código é garantida pelo `golangci-lint`. As regras estão definidas no arquivo `.golangci.yml`. Este arquivo é a **fonte da verdade** para os padrões de codificação deste projeto. Estude-o antes de começar. As regras cobrem:

- **Tratamento de Erros:** Todos os erros devem ser tratados.
- **Performance:** Práticas como pré-alocação de slices são obrigatórias.
- **Complexidade:** Funções devem ser curtas e focadas.
- **Estilo:** O código deve seguir as convenções idiomáticas do Go, reforçadas por `gofumpt` e `stylecheck`.

### 2. Fluxo de Trabalho de Desenvolvimento

1. **Gere o código:** Modifique os arquivos Go, prestando atenção aos comentários-guia `// TODO: IA`.
2. **Verificação de Qualidade (CI):** Ao enviar seu código (via Pull Request), dois processos automáticos (GitHub Actions) serão executados:
   - **Lint (`.github/workflows/lint.yml`):** Verifica a qualidade do código contra as regras do `.golangci.yml`.
   - **Build (`.github/workflows/build.yml`):** Garante que o código compila e os testes passam.
3. **Correção de Erros:** **Seu trabalho só é considerado concluído quando ambos os checks ("Lint" e "Build") passarem (ficarem verdes).** Se um deles falhar, analise o log de erro fornecido pelo GitHub Actions, corrija o problema e envie o código novamente.

### 3. Ponto de Partida

Comece analisando o arquivo `main.go`. Ele contém um esqueleto funcional de uma aplicação Bubbletea com comentários indicando onde sua lógica deve ser implementada.

---

## Comandos para Desenvolvedores Humanos

- **Para executar o programa:**
  
  ```bash
  go run .
  ```

- **Para executar os linters localmente (altamente recomendado antes de um commit):**
  
  ```bash
  golangci-lint run ./...
  ```

- **Para executar os testes:**
  
  ```bash
  go test ./...
  ```

- **Para adicionar novas dependências:**
  
  ```bash
  go get github.com/novas-dependencias
  ```

- **Para atualizar dependências:**
  
  ```bash
  go get -u ./...
  ```

- **Para formatar o código:**
  
  ```bash
  gofumpt -w .
  ```

- **Para verificar a formatação do código:**
  
  ```bash
  gofumpt -l .
  ```
