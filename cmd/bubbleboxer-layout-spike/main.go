package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	boxer "github.com/treilik/bubbleboxer"
)

// staticModel é um tea.Model simples que exibe texto estático, mas respeita
// a largura fornecida via tea.WindowSizeMsg, quebrando as linhas sempre que
// necessário para que nenhuma ultrapasse a largura máxima.
type staticModel struct {
	content string
	width   int
}

func (m staticModel) Init() tea.Cmd { return nil }

func (m staticModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
	}
	return m, nil
}

func (m staticModel) View() string {
	if m.width <= 0 {
		return m.content
	}

	var out []string
	for _, line := range strings.Split(m.content, "\n") {
		runes := []rune(line)
		for len(runes) > 0 {
			if len(runes) <= m.width {
				out = append(out, string(runes))
				break
			}
			out = append(out, string(runes[:m.width]))
			runes = runes[m.width:]
		}
		if len(line) == 0 {
			out = append(out, "")
		}
	}
	return strings.Join(out, "\n")
}

// makeVerticalFlexSizeFunc cria uma SizeFunc que distribui a altura
// entre os filhos de forma proporcional aos ratios fornecidos.
func makeVerticalFlexSizeFunc(ratios []int) func(node boxer.Node, height int) []int {
	return func(node boxer.Node, height int) []int {
		if len(node.Children) == 0 {
			return nil
		}
		if len(ratios) != len(node.Children) {
			// fallback: dividir igualmente
			res := make([]int, len(node.Children))
			base := height / len(node.Children)
			for i := range res {
				res[i] = base
			}
			// distribuir resto
			rest := height - base*len(node.Children)
			for i := 0; i < rest && i < len(res); i++ {
				res[i]++
			}
			return res
		}

		total := 0
		for _, r := range ratios {
			if r <= 0 {
				r = 1
			}
			total += r
		}
		res := make([]int, len(ratios))
		remaining := height
		remainingTotal := total
		for i, r := range ratios {
			if r <= 0 {
				r = 1
			}
			if i == len(ratios)-1 {
				res[i] = remaining
			} else {
				v := remaining * r / remainingTotal
				res[i] = v
				remaining -= v
				remainingTotal -= r
			}
		}
		return res
	}
}

// makeHorizontalFlexSizeFunc cria uma SizeFunc que distribui a largura
// entre os filhos de forma proporcional aos ratios fornecidos.
func makeHorizontalFlexSizeFunc(ratios []int) func(node boxer.Node, width int) []int {
	return func(node boxer.Node, width int) []int {
		if len(node.Children) == 0 {
			return nil
		}
		if len(ratios) != len(node.Children) {
			// fallback: dividir igualmente
			res := make([]int, len(node.Children))
			base := width / len(node.Children)
			for i := range res {
				res[i] = base
			}
			rest := width - base*len(node.Children)
			for i := 0; i < rest && i < len(res); i++ {
				res[i]++
			}
			return res
		}

		total := 0
		for _, r := range ratios {
			if r <= 0 {
				r = 1
			}
			total += r
		}
		res := make([]int, len(ratios))
		remaining := width
		remainingTotal := total
		for i, r := range ratios {
			if r <= 0 {
				r = 1
			}
			if i == len(ratios)-1 {
				res[i] = remaining
			} else {
				v := remaining * r / remainingTotal
				res[i] = v
				remaining -= v
				remainingTotal -= r
			}
		}
		return res
	}
}

// buildLayout constrói uma árvore de layout equivalente ao app.layout-flex-demo.yaml
// usando bubbleboxer.Boxer, com proporções verticais 1:4:1 e horizontais 1:2:1.
func buildLayout() boxer.Boxer {
	b := boxer.Boxer{}

	// Cria modelos estáticos equivalentes aos viewports da demo.
	headerModel := staticModel{content: "Layout Flex Demo — header ocupando a largura inteira"}
	leftModel := staticModel{content: "Painel esquerdo (flex: 1)\n\nEsta área deve ocupar 1/4 da largura total."}
	centerModel := staticModel{content: "Painel central (flex: 2)\n\nEsta área deve ocupar 2/4 (metade) da largura total."}
	rightModel := staticModel{content: "Painel direito (flex: 1)\n\nEsta área deve ocupar 1/4 da largura total."}
	footerModel := staticModel{content: "Footer / Barra de status (flex: 1 na coluna)\n\nEsta área demonstra um padrão comum de UI: header + conteúdo + footer."}

	// Leafs
	headerLeaf, _ := b.CreateLeaf("header", headerModel)
	leftLeaf, _ := b.CreateLeaf("left", leftModel)
	centerLeaf, _ := b.CreateLeaf("center", centerModel)
	rightLeaf, _ := b.CreateLeaf("right", rightModel)
	footerLeaf, _ := b.CreateLeaf("footer", footerModel)

	// Linha principal com três colunas 1:2:1.
	mainRow := boxer.Node{
		Children:        []boxer.Node{leftLeaf, centerLeaf, rightLeaf},
		VerticalStacked: false,
		SizeFunc:        makeHorizontalFlexSizeFunc([]int{1, 2, 1}),
	}

	// Coluna raiz: header / main / footer com proporção 1:4:1.
	root := boxer.Node{
		Children:        []boxer.Node{headerLeaf, mainRow, footerLeaf},
		VerticalStacked: true,
		SizeFunc:        makeVerticalFlexSizeFunc([]int{1, 4, 1}),
	}

	b.LayoutTree = root
	return b
}

func main() {
	const (
		width  = 80
		height = 24
	)

	b := buildLayout()
	// Inicializa tamanhos manualmente (sem rodar um programa Bubble Tea completo).
	_ = b.UpdateSize(tea.WindowSizeMsg{Width: width, Height: height})

	fmt.Println(b.View())
}
