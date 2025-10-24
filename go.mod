// TODO: Altere o caminho do módulo para o repositório do seu projeto.
// Ex: module github.com/seu-usuario/seu-projeto
module github.com/your-username/your-project-name

// Go 1.24 é o requisito mínimo para ser compatível com as bibliotecas
// mais recentes do ecossistema Charm (bubbletea, bubbles, etc.).
go 1.24.2

// Adicione aqui as dependências que você precisa, por exemplo:
require (
	github.com/charmbracelet/bubbles v0.18.0
	github.com/charmbracelet/bubbletea v0.26.4
	github.com/charmbracelet/lipgloss v0.11.0
)

// O restante das dependências (// indirect) será gerenciado pelo 'go mod tidy'.