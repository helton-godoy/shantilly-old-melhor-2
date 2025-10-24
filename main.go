package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

// --- MODELO ---
// O modelo (model) é o estado da sua aplicação.
// TODO: IA, adicione os campos necessários para o estado da sua aplicação aqui.
// Ex:- loading bool
//   - err     error
//   - items   []string
//   - choice  int
type model struct {
	cursor  int
	choices []string
}

// initialModel é o estado inicial da aplicação quando ela é iniciada.
func initialModel() model {
	return model{
		choices: []string{"Comprar cenouras", "Comprar aipo", "Comprar couve"},
	}
}

// --- INIT ---
// Init é o primeiro comando a ser executado quando o programa inicia.
// É aqui que você pode carregar dados ou iniciar tickers.
func (m model) Init() tea.Cmd {
	// Retorna nil significa que nenhum comando será executado inicialmente.
	return nil
}

// --- UPDATE ---
// Update é chamado toda vez que uma "mensagem" (msg) é recebida.
// Mensagens são eventos, como pressionamento de teclas, respostas de I/O, etc.
// TODO: IA, a lógica principal da sua aplicação vai dentro deste switch.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// `tea.KeyMsg` é uma mensagem enviada quando uma tecla é pressionada.
	case tea.KeyMsg:
		// O `switch` aninhado trata as teclas específicas.
		switch msg.String() {
		// Teclas para sair do programa.
		case "ctrl+c", "q":
			return m, tea.Quit

		// Teclas "up" e "k" movem o cursor para cima.
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// Teclas "down" e "j" movem o cursor para baixo.
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		// Tecla "enter" seleciona o item.
		case "enter", " ":
			// TODO: IA, adicione a lógica para quando um item é selecionado.
			// Por enquanto, apenas saímos.
			return m, tea.Quit
		}
	}

	// Retorna o modelo atualizado e nenhum novo comando.
	return m, nil
}

// --- VIEW ---
// View é responsável por renderizar a interface do usuário.
// É chamado toda vez que o modelo é atualizado.
// TODO: IA, modifique esta função para renderizar a UI da sua aplicação.
func (m model) View() string {
	s := "O que devemos comprar no mercado?\n\n"

	for i, choice := range m.choices {
		cursor := " " // O cursor padrão é um espaço.
		if m.cursor == i {
			cursor = ">" // O cursor aponta para o item selecionado.
		}
		s += fmt.Sprintf("%s [ ] %s\n", cursor, choice)
	}

	s += "\nUse as setas para navegar. Pressione Q para sair.\n"
	return s
}

// --- MAIN ---
// A função principal que inicia o programa Bubble Tea.
func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Ops, houve um erro: %v", err)
		os.Exit(1)
	}
}
