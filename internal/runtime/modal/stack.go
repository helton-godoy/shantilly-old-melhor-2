// Package modal contém a implementação normativa da Modal Stack única do runtime.
//
// Este pacote é o ÚNICO lugar autorizado para criação, empilhamento e gestão de modais,
// conforme mandatos E1.5-1/2/4 e governança:
//
//   - Implementação exclusiva em internal/runtime/modal/**.
//   - Pilha explícita (push/pop/top) com foco exclusivo no topo.
//   - Enquanto houver modal ativo, interação com plano de fundo deve ser bloqueada.
//   - Integração obrigatória via pipeline declarativo único:
//     ShantillyEvent -> EventManager -> on: -> ModalRequest -> ModalStack (Push/Top/Pop)
//     ModalStack -> ShantillyEvent (modal.*) -> EventManager -> on:
//   - Proibidos modais "soltos" em componentes, LayoutManager ou legado.
//   - Segredos trafegando pela Modal Stack:
//   - Nunca logados,
//   - Nunca persistidos,
//   - Mantidos apenas em memória pelo tempo mínimo necessário.
//
// Referências normativas e QA (não remover):
// - E1.5 — Modal Stack + Security.
// - Gates: docs/qa/gates/1.x.modal-stack.yml, docs/qa/gates/1.x.security-jit-anti-trojan.yml.
// - Arquitetura:
//   - docs/architecture/components.md#6-modal-stack--e15-wave-5-parte-1
//   - docs/architecture/core-workflows.md#workflow-3--seguranca-jit-com-modal-stack-e13-e15
//   - docs/architecture/security.md#4-modal-stack-e-seguranca-jit-e15
//   - docs/architecture/governance-runtime-tui-v2.0.md#4-wave-5--modal-stack--seguranca-jit-e15--parte-1
package modal

import "sync"

// ModalStack gerencia a pilha única de modais do runtime.
//
// Invariantes (verificados por testes e gates):
// - Somente o topo (Top) pode receber foco/eventos.
// - HasActive() true implica bloqueio lógico do plano de fundo pelo chamador.
// - Nenhuma lógica de renderização ou I/O aqui: responsabilidade do LayoutManager e TUI.
// - Nenhum log de conteúdo sensível: este pacote não faz logging.
type ModalStack struct {
	mu     sync.Mutex
	modals []Modal
	// nextID é usado apenas como gerador interno simples quando ModalRequest.ID não é fornecido.
	nextID int
}

// NewStack cria uma nova ModalStack vazia.
func NewStack() *ModalStack {
	return &ModalStack{
		modals: make([]Modal, 0),
		nextID: 1,
	}
}

// Push adiciona um novo modal ao topo da pilha.
//
// Regras:
//   - Se request.ID estiver vazio, é gerado um ID interno estável ("modal-N").
//   - Todos os modais existentes têm Active=false; o novo topo recebe Active=true.
//   - Nenhum efeito colateral fora da pilha é produzido aqui; integração é feita
//     por quem consome o estado (ex.: EventManager/LayoutManager).
func (s *ModalStack) Push(request ModalRequest) Modal {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := request.ID
	if id == "" {
		id = s.nextIDString()
	}

	// Desativa o topo atual, se existir.
	if len(s.modals) > 0 {
		top := s.modals[len(s.modals)-1]
		top.Active = false
		s.modals[len(s.modals)-1] = top
	}

	modal := NewModalFromRequest(request, id)
	modal.Active = true

	s.modals = append(s.modals, modal)

	return modal
}

// Pop remove e retorna o modal do topo da pilha.
//
// Comportamento:
//   - Se a pilha estiver vazia, retorna (Modal{}, false).
//   - Após remover o topo, o novo topo (se existir) é marcado como Active=true.
//   - Não emite eventos diretamente: o chamador decide quando gerar `modal.cancelled`,
//     `modal.confirmed` etc., garantindo isolamento deste pacote.
func (s *ModalStack) Pop() (Modal, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	n := len(s.modals)
	if n == 0 {
		return Modal{}, false
	}

	top := s.modals[n-1]
	s.modals = s.modals[:n-1]

	// Reativa o novo topo, se houver.
	if len(s.modals) > 0 {
		m := s.modals[len(s.modals)-1]
		m.Active = true
		s.modals[len(s.modals)-1] = m
	}

	return top, true
}

// Top retorna um ponteiro para o modal ativo no topo da pilha,
// ou nil quando não há modais.
//
// A estrutura retornada não deve ser mutada diretamente por chamadores externos
// em cenários concorrentes; métodos da stack devem ser usados para mudanças.
func (s *ModalStack) Top() *Modal {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.modals) == 0 {
		return nil
	}
	top := s.modals[len(s.modals)-1]
	return &top
}

// HasActive indica se existe algum modal ativo na pilha.
//
// Enquanto HasActive() for true, motores superiores (LayoutManager, EventManager)
// DEVEM tratar o plano de fundo como bloqueado e encaminhar eventos apenas ao topo.
func (s *ModalStack) HasActive() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.modals) > 0
}

// Size retorna o número de modais atualmente na pilha.
// Uso principal: diagnósticos internos e testes.
func (s *ModalStack) Size() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.modals)
}

// nextIDString gera um ID estável para novos modais.
// Não inclui conteúdo sensível; apenas contador incremental.
func (s *ModalStack) nextIDString() string {
	id := s.nextID
	s.nextID++
	return "modal-" + itoa(id)
}

// itoa é uma implementação mínima local para evitar dependência extra.
// Esta função é trivial e não realiza qualquer formatação rica.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}

	var buf [32]byte
	i := len(buf)

	neg := n < 0
	if neg {
		n = -n
	}

	for n > 0 {
		i--
		buf[i] = byte('0' + (n % 10))
		n /= 10
	}

	if neg {
		i--
		buf[i] = '-'
	}

	return string(buf[i:])
}
