package main

import (
	"fmt"

	"github.com/76creates/stickers/flexbox"
	"github.com/charmbracelet/lipgloss"
)

// renderStickersFlexDemo monta um layout simples usando a biblioteca
// 76creates/stickers/flexbox para inspecionar visualmente a distribuição
// de espaço (1:4:1 nas linhas e 1:2:1 nas colunas da linha principal).
// Este código é um spike isolado e não interfere no runtime declarativo.
func renderStickersFlexDemo(width, height int) string {
	box := flexbox.New(width, height)

	// Header (linha 1, razão 1).
	headerRow := box.NewRow()
	headerCell := flexbox.NewCell(1, 1)
	headerCell.SetContent("Layout Flex Demo — header ocupando a largura inteira")
	headerCell.SetStyle(
		headerCell.GetStyle().
			Padding(1).
			Border(lipgloss.NormalBorder()),
	)
	headerRow.AddCells(headerCell)

	// Main (linha 2, razão 4) com três colunas 1:2:1.
	mainRow := box.NewRow()
	leftCell := flexbox.NewCell(1, 4)
	leftCell.SetContent(`Painel esquerdo (flex: 1)

Esta área deve ocupar 1/4 da largura total.`)
	leftCell.SetStyle(
		leftCell.GetStyle().
			Padding(1),
	)

	centerCell := flexbox.NewCell(2, 4)
	centerCell.SetContent(`Painel central (flex: 2)

Esta área deve ocupar 2/4 (metade) da largura total.`)
	centerCell.SetStyle(
		centerCell.GetStyle().
			Padding(1).
			Border(lipgloss.NormalBorder()),
	)

	rightCell := flexbox.NewCell(1, 4)
	rightCell.SetContent(`Painel direito (flex: 1)

Esta área deve ocupar 1/4 da largura total.`)
	rightCell.SetStyle(
		rightCell.GetStyle().
			Padding(1),
	)

	mainRow.AddCells(leftCell, centerCell, rightCell)

	// Footer (linha 3, razão 1).
	footerRow := box.NewRow()
	footerCell := flexbox.NewCell(1, 1)
	footerCell.SetContent(`Footer / Barra de status (flex: 1 na coluna)

Esta área demonstra um padrão comum de UI: header + conteúdo + footer.`)
	footerCell.SetStyle(
		footerCell.GetStyle().
			Padding(1).
			Border(lipgloss.NormalBorder()),
	)
	footerRow.AddCells(footerCell)

	box.AddRows([]*flexbox.Row{headerRow, mainRow, footerRow})

	return box.Render()
}

func main() {
	const (
		width  = 80
		height = 24
	)

	out := renderStickersFlexDemo(width, height)
	fmt.Println(out)
}
