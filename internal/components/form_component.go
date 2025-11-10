// Package components implementa componentes concretos ShantillyComponent.
//
// Este arquivo define o FormComponent (type: "form") como sandbox LEGADO do v1.x,
// alinhado às regras normativas das Waves 4–7.
//
// Mandatos (derivados de docs/architecture/*.md, AGENTS.md e gates QA):
//
// - Este é o ÚNICO ponto autorizado para reaproveitar `FormConfig`/`internal/tui`.
// - Configurado EXCLUSIVAMENTE via AppConfig v2.0:
//   - Component{ Type: "form", Props: {...} }.
//
// - Qualquer mapeamento Props -> FormConfig/internal/tui é DETALHE INTERNO.
// - Não expõe FormConfig como contrato público.
// - Não cria rotas paralelas de automação:
//   - Toda automação flui via:
//     ShantillyEvent -> EventManager -> on: -> (RunAction | ModalRequest | UpdateState).
//   - Não executa scripts diretamente, não abre modais diretamente,
//     não muta estado global fora do pipeline normativo.
//   - Respeita gates:
//   - 1.x.legacy-formcomponent-encapsulation.yml
//   - 1.x.event-engine.yml
//   - 1.x.modal-stack.yml
//   - 1.x.scriptrunner-and-update-target.yml
//
// IMPORTANTE:
//   - Implementação atual é um "esqueleto" seguro para ativar o ciclo 6,
//     pronto para integrar o backend legado v1.x em iterações futuras.
package components

import (
	tea "github.com/charmbracelet/bubbletea"
)

// ShantillyComponent define o contrato mínimo esperado pelo LayoutManager
// (espelhado de internal/runtime/layout/manager.go).
//
// Mantemos a definição aqui para desacoplar o FormComponent do legado v1.x
// e alinhar com o runtime v2.0.
//
// Qualquer divergência deve seguir docs/architecture/components.md#1-shantillycomponent-contrato-base--e12.
type ShantillyComponent interface {
	ID() string
	Init() tea.Cmd
	Update(msg tea.Msg) (ShantillyComponent, tea.Cmd)
	View() string
	SetBounds(width, height int)
}

// FormComponent é o sandbox LEGADO para formulários (type: "form").
//
//   - É criado a partir de um ID lógico (componentID) e de um mapa de Props
//     proveniente do AppConfig validado.
//   - Internamente poderá mapear Props -> estruturas compatíveis com FormConfig
//     e delegar para o código v1.x (internal/tui) de forma encapsulada.
//   - Externamente, se comporta como um ShantillyComponent comum.
//
// NOTA: Implementação inicial neutra; integrações com o backend v1.x
// serão adicionadas mantendo estas invariantes.
type FormComponent struct {
	componentID string
	props       map[string]interface{}

	width  int
	height int

	// Campos internos futuros:
	// - legacyConfig *config.FormConfig          // LEGACY (detalhe interno)
	// - legacyModel  tea.Model                   // adaptador para internal/tui/model.go
}

// NewFormComponent cria um novo FormComponent sandboxado.
//
// - componentID: ID lógico vindo de AppConfig.Components.
// - props: mapa de propriedades declarativas (ex.: título, campos, etc.).
// - Não aceita FormConfig diretamente (qualquer tentativa deve ser tratada como bug).
func NewFormComponent(componentID string, props map[string]interface{}) *FormComponent {
	if props == nil {
		props = make(map[string]interface{})
	}

	return &FormComponent{
		componentID: componentID,
		props:       props,
	}
}

// ID retorna o identificador lógico do componente.
func (f *FormComponent) ID() string {
	return f.componentID
}

// Init inicializa o estado interno.
//
// Implementação atual é mínima; integração com backend legado será adicionada
// em iteração futura sem quebrar o contrato.
func (f *FormComponent) Init() tea.Cmd {
	// Futuro:
	// - Construir estrutura interna baseada em Props (ou mapear para FormConfig LEGACY).
	// - Inicializar modelo TUI interno adaptado.
	return nil
}

// Update processa mensagens TEA para o FormComponent.
//
// Regras normativas:
// - Não executa scripts diretamente.
// - Não abre modais diretamente.
// - Emite apenas eventos (via tea.Msg/ShantillyEvent) que serão roteados pelo EventManager.
// - Qualquer automação deve ser descrita no AppConfig e resolvida via on:.
func (f *FormComponent) Update(msg tea.Msg) (ShantillyComponent, tea.Cmd) {
	// Placeholder:
	// - No futuro, delegar para backend legado adaptado.
	return f, nil
}

// View retorna a representação textual do formulário.
//
// Implementação inicial devolve placeholder neutro para não interferir
// na experiência até o backend legado ser plugado.
func (f *FormComponent) View() string {
	// Futuro:
	// - Renderizar com base no modelo interno (ex.: huh.Form legado adaptado).
	return ""
}

// SetBounds define os limites de renderização fornecidos pelo LayoutManager.
//
// - Usado para integrar com layout responsivo.
// - Não altera invariantes de segurança ou automação.
func (f *FormComponent) SetBounds(width, height int) {
	f.width = width
	f.height = height
}
