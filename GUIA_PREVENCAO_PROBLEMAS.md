# Guia de Prevenção de Problemas em Componentes TUI


## 🎯 Objetivo

Este guia estabelece padrões e práticas para evitar problemas como os enfrentados na implementação do componente Button, garantindo desenvolvimento eficiente e código robusto.

---


## 📋 Checklist de Desenvolvimento de Componentes


### ✅ Fase 1: Criação do Componente

- [ ] **Estrutura mínima** criada em `internal/components/[nome]/`

- [ ] **Interface ShantillyComponent** implementada

- [ ] **Métodos obrigatórios**: `Init()`, `Update()`, `View()`, `ID()`, `SetDimensions()`

- [ ] **Props struct** definida com valores padrão

- [ ] **Registro no registry** em `internal/runtime/layout/registry.go`


### ✅ Fase 2: Teste Isolado

```yaml
# Criar: test-[nome]-isolado.yaml
version: "1.0"
layout:
    id: root
    type: row
    items:
        - id: center
          type: box
          component: [nome]
components:
    - id: [nome]
      type: [tipo]

```

- [ ] **Componente aparece visualmente**

- [ ] **Sem erros no console**

- [ ] **Renderização básica funcionando**


### ✅ Fase 3: Funcionalidade Básica

- [ ] **Eventos básicos** funcionando (click, focus)

- [ ] **Navegação Tab** entre componentes

- [ ] **Feedback visual** de estado (focado, normal, desabilitado)

- [ ] **Props configuráveis** via YAML


### ✅ Fase 4: Integração Progressiva

```yaml
# Testar: test-[nome]-simples.yaml
layout:
    id: root
    type: row
    items:
        - id: left
          type: box
          component: list
          flex: 1
        - id: center
          type: box
          component: [nome]
          width: 20
        - id: right
          type: box
          component: viewport
          flex: 2

```

- [ ] **Layout com múltiplos componentes**

- [ ] **Espaçamento e dimensões corretas**

- [ ] **Navegação fluida** entre todos componentes


### ✅ Fase 5: Testes de Compatibilidade

- [ ] **Múltiplos tamanhos de terminal** (pequeno, médio, grande)

- [ ] **Diferentes terminais** (GNOME Terminal, VSCode, etc.)

- [ ] **Sem AltScreen issues** (usar `tea.WithOutput(os.Stdout)`)

- [ ] **Performance aceitável** em renderização

---


## 🚨 Padrões de Debug


### 1. Logs Estruturados

```go
// PADRÃO RECOMENDADO
const debugMode = false // false em produção

func (m *Model) View() string {
    if debugMode {
        fmt.Printf("[DEBUG] %s.View() called\n", m.ID())
    }
    
    result := m.renderContent()
    
    if debugMode {
        fmt.Printf("[DEBUG] %s.View() result: %q\n", m.ID(), result)
    }
    
    return result
}

```


### 2. Fallbacks Visuais

```go
// SEMPRE ter fallback visual
func (m *Model) View() string {
    result := m.renderStyled()
    
    // Fallback se estilo falhar
    if result == "" || len(strings.TrimSpace(result)) == 0 {
        return fmt.Sprintf("[FALLBACK: %s]", m.props.Text)
    }
    
    return result
}

```


### 3. Testes Incrementais

```bash
# Sequência obrigatória de testes
./shantilly runtime --file test-[nome]-isolado.yaml
./shantilly runtime --file test-[nome]-simples.yaml  
./shantilly runtime --file app-with-[nome].yaml
./shantilly runtime --file app.yaml

```

---


## 🛡️ Padrões de Runtime


### Configuração Segura do Bubble Tea

```go
// PADRÃO OBRIGATÓRIO
p := tea.NewProgram(
    mainModel,
    tea.WithOutput(os.Stdout),     // ✅ Universalmente compatível
    tea.WithInputTTY(),            // ✅ Input consistente
    // ❌ NUNCA usar tea.WithAltScreen() por padrão
)

```


### Tratamento de Erros

```go
// PADRÃO RECOMENDADO
if _, err := p.Run(); err != nil {
    // Log detalhado do erro
    fmt.Printf("[ERROR] Runtime failed: %v\n", err)
    
    // Tentar fallback
    if fallbackModel := m.createFallback(); fallbackModel != nil {
        fmt.Println("[INFO] Tentando modo fallback...")
        p := tea.NewProgram(fallbackModel, tea.WithOutput(os.Stdout))
        p.Run()
    }
    
    return fmt.Errorf("runtime execution failed: %w", err)
}

```

---


## 🔧 Arquitetura de Componentes


### Estrutura de Diretórios

```
internal/components/[nome]/
├── model.go          # Componente principal
├── props.go          # Estrutura de props
├── styles.go         # Estilos lipgloss
├── events.go         # Handlers de eventos
└── README.md         # Documentação do componente

```


### Interface Obrigatória

```go
type ShantillyComponent interface {
    Init() tea.Cmd
    Update(msg tea.Msg) (tuiapi.ShantillyComponent, tea.Cmd)
    View() string
    ID() string
    SetDimensions(w, h int)
}

// Adicional para componentes focáveis
type Focusable interface {
    Focus(focused bool)
}

```


### Registro Padrão

```go
// Em internal/runtime/layout/registry.go
case "[tipo]":
    props := [Nome]Props{
        Text: getString(comp.Props, "text", "Default"),
        // Outras props com defaults
    }
    return [nome].New(comp.ID, r.Theme, props)

```

---


## 📊 Testes Automatizados


### 1. Teste de Renderização

```go
func TestComponentView(t *testing.T) {
    comp := NewComponent("test", theme, Props{Text: "Hello"})
    
    view := comp.View()
    assert.NotEmpty(t, view, "View() não pode retornar vazio")
    assert.Contains(t, view, "Hello", "View() deve conter o texto")
}

```


### 2. Teste de Eventos

```go
func TestComponentEvents(t *testing.T) {
    comp := NewComponent("test", theme, Props{})
    
    // Teste foco
    comp.Focus(true)
    assert.True(t, comp.IsFocused(), "Componente deve aceitar foco")
    
    // Teste click
    model, cmd := comp.Update(tea.KeyMsg{Type: tea.KeyEnter})
    assert.NotNil(t, cmd, "Click deve gerar comando")
}

```


### 3. Teste de Layout

```go
func TestComponentLayout(t *testing.T) {
    comp := NewComponent("test", theme, Props{Text: "Hello"})
    
    // Teste dimensões
    comp.SetDimensions(20, 3)
    view := comp.View()
    
    lines := strings.Split(view, "\n")
    assert.LessOrEqual(t, len(lines), 3, "Altura deve respeitar SetDimensions")
}

```

---


## 🚨 Problemas Comuns e Soluções


### Problema 1: Componente Invisível
**Sintomas:** Funciona mas não aparece visualmente  
**Causa:** `tea.WithAltScreen()` incompatibilidade  
**Solução:** Usar `tea.WithOutput(os.Stdout)`


### Problema 2: Layout Quebrado
**Sintomas:** Sobreposição ou espaçamento incorreto  
**Causa:** Cálculo incorreto de dimensões  
**Solução:** Implementar `SetDimensions()` corretamente


### Problema 3: Eventos Não Funcionam
**Sintomas:** Cliques não têm efeito  
**Causa:** Componente não registrado ou evento não mapeado  
**Solução:** Verificar registry e eventos YAML


### Problema 4: Performance Ruim
**Sintomas:** Renderização lenta  
**Causa:** View() sendo chamado excessivamente  
**Solução:** Cache de renderização quando possível

---


## 📋 Processo de Code Review


### Checklist para PRs

- [ ] **Componente isolado** testado e funcionando

- [ ] **Integração com layout** validada

- [ ] **Múltiplos terminais** testados

- [ ] **Sem debug logs** em produção

- [ ] **Documentação** atualizada

- [ ] **Testes automatizados** incluídos

- [ ] **Fallbacks visuais** implementados

- [ ] **Performance** aceitável


### Perguntas Obrigatórias
1. "O componente aparece visualmente em todos os cenários?"
2. "A navegação Tab funciona corretamente?"
3. "Os eventos são disparados e tratados?"
4. "O layout é responsivo?"
5. "Existem fallbacks para casos de erro?"

---


## 🎓 Formação da Equipe


### Tópicos Obrigatórios
1. **Bubble Tea fundamentals**
2. **Lip Gloss styling patterns**
3. **Shantilly architecture**
4. **Debugging techniques**
5. **Testing strategies**


### Exercícios Práticos
1. Criar componente simples do zero
2. Debugar componente quebrado
3. Implementar layout complexo
4. Otimizar performance
5. Escrever testes automatizados

---


## 🔄 Processo Contínuo


### Revisão Semanal

- [ ] **Novos componentes** seguem padrões

- [ ] **Problemas recorrentes** identificados

- [ ] **Documentação** atualizada

- [ ] **Testes** executados


### Melhoria Contínua

- [ ] **Métricas de bugs** por componente

- [ ] **Tempo de debug** reduzido

- [ ] **Qualidade de código** melhorada

- [ ] **Feedback da equipe** coletado

---


## 📞 Contato e Suporte


### Emergências
- **Slack #shantilly-dev** para problemas críticos
- **GitHub Issues** para bugs documentados
- **Code Review** para padrões


### Recursos
- **Documentação interna:** `/docs/components/`
- **Exemplos:** `/examples/components/`
- **Templates:** `/templates/component/`

---

**Status:** ✅ **ATIVO** - Guia em uso contínuo pela equipe  
**Última atualização:** 2025-11-28  
**Responsável:** Equipe de Desenvolvimento Shantilly
