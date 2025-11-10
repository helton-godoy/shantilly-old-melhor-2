# **SHANTILLY RUNTIME 02 - BMad**

[https://gemini.google.com/share/ad7254acd532](https://gemini.google.com/share/ad7254acd532)

As respostas abaixo foram geradas por um Gem, de acordo com as instruções personalizadas pelo criador.
 [Saiba maisAbre em uma nova janela](https://support.google.com/gemini?p=shared_gem) 8 de novembro de 2025 às 09:31 Publicada em 8 de novembro de 2025 às 18:43

Abaixo está minha última interação com o John.

Me ajude com a próxima etapa, por favor. Direcione para o agente correto.

```
[cite\_start]Como John (PM) 📋, concluí a tarefa `create-doc`[cite: 
79]. [cite\_start]Consolidei todas as nossas decisões, refinamentos e o 
novo *roadmap* [cite: 1686, 1683] num único Documento de Requisitos do 
Produto (PRD) v2.0.



[cite\_start]Este documento formaliza a nossa nova visão do **"Runtime 
TUI Declarativo"** [cite: 1686, 1683] [cite\_start]e serve como a "fonte
 da verdade" para o Arquiteto (Winston) [cite: 14] iniciar a próxima 
fase de design.



Obrigado pela excelente colaboração na definição desta visão\!



-----



# Documento de Requisitos do Produto (PRD) v2.0: Shantilly



## 1\. Metas e Contexto de Fundo



### Metas (Goals)



* [cite\_start]Oferecer uma alternativa declarativa (YAML) e moderna 
ao `dialog`/`whiptail` para *todos* os scripts shell (Bash, PowerShell, 
etc.)[cite: 363, 1686, 1683].

* [cite\_start]Servir como um **"Runtime TUI Declarativo"** que gere 
layouts (`column`/`row`), componentes (`list`, `form`, `viewport`) e 
lógica de automação (`on:`)[cite: 1686, 1683].

* [cite\_start]Fornecer uma **Fundação Genérica** (`run: { script: 
... }`) para scripts arbitrários, garantindo integração via passagem de 
dados por template (`args:` e `stdin:`)[cite: 1686, 1683].

* [cite\_start]Fornecer ***runners* Especialistas** otimizados 
(Pós-MVP) para ferramentas de DevOps como Ansible (`ansible_playbook:`) 
e, futuramente, Terraform, incluindo descoberta preditiva de inventário e
 playbooks[cite: 1686, 1683].

* [cite\_start]Manter a portabilidade (binário estático único) e a integração de pipeline (stdin/stdout)[cite: 363].



### Contexto de Fundo (Background Context)



[cite\_start]O desenvolvimento do MVP v1.0 (Épicos 1 e 2) [cite: 1531, 
1640] [cite\_start]validou a pilha de tecnologia (`huh`, `bubbletea`) 
[cite: 1149, 1150, 1151, 1152, 1153, 1154, 1155, 1156, 1157, 1158, 1159,
 1160, 1609, 1610, 1611, 1612, 1613, 1614, 1615, 1616, 1617, 1618, 1619,
 1620, 1621, 1622, 1623, 1624, 1625, 1626, 1627, 1628, 1629, 1630, 1631,
 1632, 1633, 1634, 1635, 1636, 1637, 1638, 1639, 1641, 1642, 1643, 1644,
 1645, 1646, 1647, 1648, 1649, 1650, 1651, 1652, 1653, 1654, 1655, 1656,
 1657, 1658, 1659, 1660, 1661, 1662, 1663, 1664, 1665, 1666, 1667, 1668,
 1669, 1670, 1671, 1672, 1673, 1674, 1675, 1676, 1677, 1678, 1679, 1680,
 1681, 1682, 1683, 1684, 1685, 1686] e resolveu o caso de uso de 
"formulário simples".



[cite\_start]Esta v2.0 pivota dessa ferramenta de "formulário único" 
para um "Runtime TUI" completo, inspirado no "Appsmith" (para UI 
declarativa) e "Ansible" (para lógica de eventos `on:`)[cite: 1686, 
1683]. [cite\_start]Esta arquitetura permite a criação de dashboards TUI
 complexos e multi-componente (layouts, menus, viewports) que 
*orquestram* automações de backend, em vez de serem apenas chamados por 
elas[cite: 1686, 1683].



### Change Log



| Data | Versão | Descrição | Autor |

| :--- | :--- | :--- | :--- |

| 08/11/2025 | 2.0.0 | Rascunho inicial do PRD v2.0, redefinindo o projeto como um "Runtime TUI Declarativo". | John (PM) |



-----



## 2\. Requisitos (PRD v2.0)



### Funcionais (FRs) - O Runtime TUI



[cite\_start]Estes requisitos definem o nosso novo MVP: a "Fundação Genérica"[cite: 1686, 1683].



* [cite\_start]**FR1 (Layout):** O `shantilly` DEVE analisar e 
renderizar uma estrutura de layout hierárquica definida em YAML, usando 
os tipos `type: column`, `type: row`, e `type: box`[cite: 1686, 1683].

* [cite\_start]**FR2 (Estilo/Flex):** O layout DEVE suportar 
propriedades de dimensionamento como `height: <int>`, `width: 
'N%'`, e `flex: <int>` para controlar o espaço[cite: 1686, 1683].

* [cite\_start]**FR3 (Componentes Embutidos):** O `shantilly` DEVE 
suportar a definição de componentes de UI diretamente dentro de um `box`
 usando a chave `component:` (o foco do nosso MVP)[cite: 1686, 1683].

* [cite\_start]**FR4 (Componente: `viewport`):** DEVE suportar 
`component: { type: viewport }`, capaz de exibir `source: { type: 
static, content: "..." }` (incluindo markdown) e `source: { type: 
command, exec: "..." }` (para streaming de stdout)[cite: 1686, 1683].

* [cite\_start]**FR5 (Componente: `list`):** DEVE suportar 
`component: { type: list }`, com um `id:` de grupo e `items:` (cada um 
com `id:` e `text`), e DEVE emitir um evento `list_id:select`[cite: 
1686, 1683].

* [cite\_start]**FR6 (Componente: `buttongroup`):** DEVE suportar 
`component: { type: buttongroup }`, com um `id:` de grupo, `items:` (com
 `id`, `label`, `role`), e DEVE emitir um evento 
`buttongroup_id:press`[cite: 1686, 1683].

* [cite\_start]**FR7 (Componente: `form`):** DEVE suportar 
`component: { type: form }`, que contém `fields:` (usando a sintaxe 
`huh` já validada no v1.0 [cite: 1149, 1150, 1151, 1152, 1153, 1154, 
1155, 1156, 1157, 1158, 1159, 1160, 1609, 1610, 1611, 1612, 1613, 1614, 
1615, 1616, 1617, 1618, 1619, 1620, 1621, 1622, 1623, 1624, 1625, 1626, 
1627, 1628, 1629, 1630, 1631, 1632, 1633, 1634, 1635, 1636, 1637, 1638, 
1639]) e `actions:`. [cite\_start]DEVE emitir um evento `form_id:submit`
 contendo o *payload* de dados do formulário[cite: 1686, 1683].

* [cite\_start]**FR8 (Lógica: `on:`):** O `shantilly` DEVE analisar 
um bloco `on:` na raiz do YAML para definir a lógica de automação[cite: 
1686, 1683].

* [cite\_start]**FR9 (Ação: `script`):** O bloco `on:` DEVE suportar o
 *runner* de fundação: `run: { script: "/path/to/script.sh" }`[cite: 
1686, 1683].

* [cite\_start]**FR10 (Fluxo de Dados):** O *runner* `script:` DEVE 
suportar duas chaves para passagem de dados[cite: 1686, 1683]:

1. **`args: []string`**: Uma lista de *strings* que serão passadas
 como argumentos de linha de comando para o script, com suporte para 
*templates* (ex: `{{ form.field_name }}`).

2. **`stdin: any`**: Um objeto (ex: `{{ form }}`) que o 
`shantilly` irá serializar como JSON e passar para o `stdin` do 
*script*.

* [cite\_start]**FR11 (Ciclo de Vida do Target):** O bloco `run:` 
DEVE suportar uma chave `update_target: "id_do_viewport"`[cite: 1686, 
1683]. [cite\_start]Se um novo evento `run:` for disparado para o 
*mesmo* `update_target`, o `shantilly` DEVE primeiro **terminar (enviar 
`SIGTERM`)** o processo anterior antes de iniciar o novo[cite: 1686, 
1683].



### Não Funcionais (NFRs) - O Runtime TUI



* [cite\_start]**NFR1 (Fundação v1.0):** Todos os NFRs do PRD v1.0 
permanecem válidos: binário estático único [cite: 1531][cite\_start], 
cross-platform (Linux, macOS, Windows) [cite: 1531][cite\_start], 
escrito em Go [cite: 1531][cite\_start], arranque rápido (\<500ms) 
[cite: 1531][cite\_start], e gestão de erros com `stderr` e códigos de 
saída não-zero[cite: 1531].

* [cite\_start]**NFR2 (Layout Fluido):** O motor de layout 
(`column`/`row`/`box`) DEVE responder a mensagens de redimensionamento 
do terminal (`tea.WindowSizeMsg`) e re-calcular o layout de forma fluida
 e "flicker-free" (sem piscar)[cite: 1686, 1683].

* [cite\_start]**NFR3 (Precedência de Conteúdo):** O `shantilly` DEVE
 seguir a "Lógica de Precedência Unificada" para conteúdo de componentes
 (1º: CLI `--set`, 2º: YAML `model:`/`component:`, 3º: Vazio)[cite: 
1686, 1683].

* [cite\_start]**NFR4 (Descoberta Preditiva - Ansible):** Para a Fase
 2/3, o `playbook_explorer` DEVE filtrar "ruído" (pastas `roles/`, 
`tasks/`) [cite: 1686, 1683] [cite\_start]e o `inventory_explorer` DEVE 
usar `ansible-inventory` como "Oráculo"[cite: 1686, 1683].

* [cite\_start]**NFR5 (Ficheiro-Sombra):** O ficheiro de catálogo 
(`.shantilly.yml`) DEVE ser opcional e usado apenas para *refinar* a 
descoberta automática, não sendo obrigatório[cite: 1686, 1683].



-----



## 3\. Metas de Design da Interface do Usuário



### Visão Geral da UX (User Experience)



A UX deve ser a de um **"Runtime TUI Declarativo"**. [cite\_start]A 
interface não é mais um formulário linear único[cite: 1531], mas sim um 
*dashboard* composto, definido inteiramente pelo YAML. [cite\_start]A 
experiência deve ser semelhante ao Appsmith[cite: 1686, 1683]: limpa, 
responsiva (ao terminal) e orientada a componentes.



### Paradigmas Chave de Interação



* **Orientada a Eventos (Nova):** A interação principal não é linear.
 [cite\_start]O utilizador seleciona itens em listas (`list:select`) 
[cite: 1686, 1683] [cite\_start]ou pressiona botões 
(`buttongroup:press`) [cite: 1686, 1683][cite\_start], que disparam 
ações no bloco `on:`[cite: 1686, 1683].

* [cite\_start]**Foco no Teclado (Mantido):** A navegação DEVE 
continuar a ser primariamente baseada no teclado (Tab, Setas, 
Enter)[cite: 1531].

* [cite\_start]**Feedback Imediato (Mantido):** O componente focado 
DEVE ser claramente destacado[cite: 1531]. [cite\_start]O 
`update_target` (FR11) DEVE exibir o *output* de comandos em tempo 
real[cite: 1686, 1683].

* [cite\_start]**Gestão de Foco Global (Nova):** A UI DEVE ter um 
mecanismo claro para indicar qual painel/componente (ex: `sidebar` vs 
`content`) está "em foco", e DEVE fornecer navegação intuitiva *entre* 
painéis (ex: Ctrl+Tab)[cite: 1686, 1683].

* **Ligação de Dados (Nova):** A UI DEVE ser reativa. 
[cite\_start]Componentes (ex: um `viewport` estático) DEVEM ser capazes 
de exibir dados de outros componentes (ex: `Olá, {{ form.username 
}}`)[cite: 1686, 1683].



### Ecrãs e Vistas Principais



Não há ecrãs "pré-definidos". Os ecrãs são *definidos dinamicamente* 
pelo utilizador através do **`layout` YAML** (FR1). [cite\_start]A UI é 
uma composição de `type: column`, `type: row`, e `type: box`[cite: 1686,
 1683].



### Alinhamento e Layout (Nova Visão)



[cite\_start]O layout linear do MVP v1.0 [cite: 1531] está obsoleto. O novo requisito é:



* O `shantilly` DEVE renderizar com precisão o layout `column`/`row` definido pelo utilizador.

* [cite\_start]O `shantilly` DEVE respeitar as propriedades de 
dimensionamento (`height`, `width`, `flex`) para distribuir o 
espaço[cite: 1686, 1683].

* [cite\_start]O `shantilly` DEVE responder a mensagens de 
redimensionamento do terminal (`tea.WindowSizeMsg`) e re-calcular o 
layout fluido e "flicker-free" (NFR2, Meta de UI Refinada)[cite: 1686, 
1683].



### Acessibilidade, Branding, Plataformas Alvo



[cite\_start]Estes requisitos permanecem os mesmos do PRD v1.0 (WCAG 
AA, estética Charmbracelet, terminais modernos em 
Linux/macOS/Windows)[cite: 1531].



-----



## 4\. Assunções Técnicas



### Estrutura do Repositório: Monorepo



* [cite\_start]**Assunção:** Continuamos a usar um Monorepo Go 
simples (`cmd/`, `internal/`)[cite: 1531, 1600]. [cite\_start]Esta 
estrutura modular provou ser eficaz nos Épicos 1 e 2[cite: 1149, 1150, 
1151, 1152, 1153, 1154, 1155, 1156, 1157, 1158, 1159, 1160, 1609, 1610, 
1611, 1612, 1613, 1614, 1615, 1616, 1617, 1618, 1619, 1620, 1621, 1622, 
1623, 1624, 1625, 1626, 1627, 1628, 1629, 1630, 1631, 1632, 1633, 1634, 
1635, 1636, 1637, 1638, 1639, 1641, 1642, 1643, 1644, 1645, 1646, 1647, 
1648, 1649, 1650, 1651, 1652, 1653, 1654, 1655, 1656, 1657, 1658, 1659, 
1660, 1661, 1662, 1663, 1664, 1665, 1666, 1667, 1668, 1669, 1670, 1671, 
1672, 1673, 1674, 1675, 1676, 1677, 1678, 1679, 1680, 1681, 1682, 1683, 
1684, 1685, 1686].



### Arquitetura de Serviço: Runtime TUI Declarativo



* [cite\_start]**Assunção (Nova):** A arquitetura evoluiu de um 
"Monolithic CLI Application" [cite: 1531] [cite\_start]para um 
**"Runtime TUI Declarativo"**[cite: 1686, 1683]. [cite\_start]O 
`shantilly` é um motor que consome YAML, renderiza uma UI composta e 
gere um ciclo de vida de eventos (`on:`) para orquestrar automações de 
*scripts*[cite: 1686, 1683].



### Requisitos de Teste: Teste de Integração TUI



* [cite\_start]**Assunção (Atualizada):** O foco "apenas em Testes de
 Unidade" do v1.0 [cite: 1531] é insuficiente. [cite\_start]Devido à 
complexidade do layout (FR1) e da gestão de foco (Metas de UI), os 
**Testes de Integração TUI** (usando `teatest`) são agora um requisito 
central para validar o layout fluido e a navegação entre painéis[cite: 
1600].



### Estrutura YAML Esperada (A Nova Fonte da Verdade)



* [cite\_start]**Assunção (Nova):** A estrutura YAML do v1.0 (lista 
simples de `fields:`) [cite: 1531] está obsoleta. [cite\_start]A nova 
assunção de arquitetura é o YAML "Appsmith-style" [cite: 1686, 1683] que
 definimos, composto por **Layout**, **Componentes** e **Lógica**:



<!-- end list -->



```yaml

# 1. LAYOUT (Define o "onde")

[cite_start]# [cite: 1686, 1683]

type: column

items:

- type: row

flex: 1

items:

- type: box

id: "sidebar"

width: "30%"

# 2. COMPONENTE (Define o "o quê")

[cite_start]# [cite: 1686, 1683]

component:

type: list

id: "menu"

items:

- { id: "users", text: "Gerir Utilizadores" }

- type: box

id: "content"

width: "70%"

component:

type: form

id: "user_form"

fields:

- { name: "username", label: "Nome", type: "input" }

actions:

type: buttongroup

items:

- { id: "submit", label: "Criar", role: "primary" }



# 3. LÓGICA (Define o "como")

[cite_start]# [cite: 1686, 1683]

on:

- event: "user_form:submit"

run:

script: "/opt/scripts/create_user.sh"

# Nomenclatura refinada (Opção C)

args:

- "--mode=production"

stdin: "{{ form }}" # Passa o payload JSON para o stdin

# Ciclo de vida refinado (FR11)

update_target: "log_viewport"
```

### Assunções Técnicas Adicionais

* [cite\_start]**Pilha de Tecnologias (Mantida e Validada):** Go 
  (1.24.2+), `spf13/cobra`, `charmbracelet/bubbletea`, 
  `charmbracelet/lipgloss`, `charmbracelet/huh`, `gopkg.in/yaml.v3`[cite: 
  1600, 1686, 1683].

* [cite\_start]**Bibliotecas Relevantes (Nova):** A arquitetura 
  dependerá de `charmbracelet/bubbles` (para `list`, `viewport`), 
  `charmbracelet/glamour` (para markdown), `rmhubbert/bubbletea-overlay` 
  (para modais Fase 2), e inspiração de `76creates/stickers` (para 
  layout)[cite: 1686, 1683].

* [cite\_start]**Build e Distribuição (Mantido):** `GoReleaser` para binários estáticos cross-platform[cite: 1600].

-----

## 5\. Lista de Épicos (Roadmap v2.0)

[cite\_start]Este *roadmap* substitui a lista de épicos do PRD v1.0[cite: 1531].

* **Épico 1: Fundação do Runtime TUI (Genérico)**

* [cite\_start]**Meta:** Construir o motor central: o layout 
  (`column`/`row`/`box`) [cite: 1686, 1683][cite\_start], os componentes 
  essenciais (`list`, `viewport`, `form`, `buttongroup`) [cite: 1686, 
  1683][cite\_start], e a lógica de eventos (`on:`, `run: { script: ... 
  }`)[cite: 1686, 1683]. [cite\_start](Este épico absorve todo o trabalho 
  já concluído no v1.0 [cite: 1149, 1150, 1151, 1152, 1153, 1154, 1155, 
  1156, 1157, 1158, 1159, 1160, 1609, 1610, 1611, 1612, 1613, 1614, 1615, 
  1616, 1617, 1618, 1619, 1620, 1621, 1622, 1623, 1624, 1625, 1626, 1627, 
  1628, 1629, 1630, 1631, 1632, 1633, 1634, 1635, 1636, 1637, 1638, 1639, 
  1641, 1642, 1643, 1644, 1645, 1646, 1647, 1648, 1649, 1650, 1651, 1652, 
  1653, 1654, 1655, 1656, 1657, 1658, 1659, 1660, 1661, 1662, 1663, 1664, 
  1665, 1666, 1667, 1668, 1669, 1670, 1671, 1672, 1673, 1674, 1675, 1676, 
  1677, 1678, 1679, 1680, 1681, 1682, 1683, 1684, 1685, 1686]).

* **Épico 2: O Runner Especialista (Ansible Fase 2)**

* [cite\_start]**Meta:** Implementar o *runner* de conveniência 
  `run: { ansible_playbook: ... }`, focando na gestão de `vars:` e no 
  popup modal `ask_vault_pass: true`[cite: 1686, 1683].

* **Épico 3: O Runtime Preditivo (Ansible Fase 3)**

* [cite\_start]**Meta:** Implementar os componentes 
  `playbook_explorer` (Magia 1: descobrir playbooks) e 
  `inventory_explorer` (Magia 2: descobrir inventário)[cite: 1686, 1683].

* **Épico 4: Administração SSH (Visão de Longo Prazo)**

* [cite\_start]**Meta:** Integrar o `charmbracelet/wish` [cite: 
  1686, 1683] [cite\_start]para servir o Runtime TUI sobre SSH[cite: 363].

-----

## 6\. Detalhes do Épico 1: Fundação do Runtime TUI (Genérico)

[cite\_start]**Meta do Épico:** Construir o motor central do 
`shantilly`: o motor de layout (`column`/`row`/`box`) [cite: 1686, 
1683][cite\_start], os componentes essenciais de dashboard (`list`, 
`viewport`, `form`, `buttongroup`) [cite: 1686, 1683][cite\_start], e a 
lógica de eventos (`on:`, `run: { script: ... }`)[cite: 1686, 1683]. 
[cite\_start]Este épico irá refatorar o trabalho concluído do v1.0 para 
que ele funcione como o componente `type: form` [cite: 1686, 1683] 
dentro deste novo runtime.

### Estória 1.1: O Motor de Layout (Renderização)

[cite\_start]**Como um** SysAdmin, **Eu quero** definir um layout TUI 
usando `column`, `row`, e `box` no meu YAML[cite: 1686, 1683], **Para 
que** eu possa criar dashboards complexos e organizados.

#### Critérios de Aceitação

1. O parser DEVE suportar as chaves `type: column`, `type: row`, e `type: box` (FR1).

2. O motor de renderização (`lipgloss`) DEVE respeitar as propriedades
   `height: <int>`, `width: 'N%'`, e `flex: <int>` (FR2).

3. O layout DEVE recalcular-se fluidamente (sem piscar) ao receber uma
   mensagem de redimensionamento (`tea.WindowSizeMsg`) (NFR2).

4. Um `box` DEVE renderizar o seu `component: { type: static, content: "..." }` (para testes de layout).

### Estória 1.2: O Motor de Lógica (Eventos e Ações)

**Como um** SysAdmin, **Eu quero** que a minha UI TUI possa "ouvir" 
eventos e executar ações (`scripts`) em resposta, **Para que** o meu 
dashboard seja interativo e possa orquestrar automações.

#### Critérios de Aceitação

1. O `shantilly` DEVE analisar um bloco `on:` na raiz do YAML (FR8).

2. O motor DEVE suportar o *runner* de fundação: `run: { script: "/path/to/script.sh" }` (FR9).

3. O `run:` DEVE suportar `update_target: "id_do_viewport"`, direcionando o `stdout` do script para o *viewport* alvo (FR11).

4. [cite\_start]O motor DEVE garantir que, se um novo script for 
   direcionado para um `update_target` já ocupado, o script anterior seja 
   terminado (SIGTERM) antes de o novo começar (Refinamento FR11)[cite: 
   1686, 1683].

### Estória 1.3: Componentes Essenciais de Display (List, Viewport, Button)

**Como um** SysAdmin, **Eu quero** usar componentes de `list` (para 
menus), `viewport` (para saída de log) e `buttongroup` (para ações), 
**Para que** eu possa construir um dashboard funcional.

#### Critérios de Aceitação

1. [cite\_start]DEVE implementar `component: { type: viewport }` 
   (FR4), incluindo `source: { type: command, exec: "..." }` (ex: `tail 
   -f`) e `content_type: markdown`[cite: 1686, 1683].

2. [cite\_start]DEVE implementar `component: { type: list }` (FR5), 
   que emite um evento `list_id:select` quando um item é selecionado[cite: 
   1686, 1683].

3. [cite\_start]DEVE implementar `component: { type: buttongroup }` 
   (FR6), que emite um evento `buttongroup_id:press` (com o `item.id`) 
   quando um botão é pressionado[cite: 1686, 1683].

4. [cite\_start]A navegação por teclado DEVE permitir "saltar" entre 
   estes novos painéis/componentes (Refinamento de Meta de UI)[cite: 1686, 
   1683].

### Estória 1.4: Integração do Componente `form` (Absorção do v1.0)

**Como um** SysAdmin, **Eu quero** usar o `type: form` (que já 
construímos no v1.0) como um componente *dentro* do meu novo layout, 
**Para que** eu possa coletar dados de forma organizada.

#### Critérios de Aceitação

1. [cite\_start]Refatorar o código dos Épicos 1 e 2 (v1.0) [cite: 
   1531, 1640] para que funcione como um `component: { type: form }` (FR7).

2. O `form` DEVE renderizar e funcionar corretamente quando colocado dentro de um `box` do layout.

3. [cite\_start]Quando a ação `id: "submit"` do formulário [cite: 
   1686, 1683] for pressionada, o componente `form` DEVE emitir um evento 
   `form_id:submit`.

4. O *payload* do evento `form_id:submit` DEVE conter o JSON de dados do formulário (o output do v1.0).

### Estória 1.5: O Fluxo de Dados (Args & Stdin)

**Como um** SysAdmin, **Eu quero** passar os dados coletados no meu 
`form` (ou a seleção de uma `list`) para os meus `scripts` de forma 
robusta, **Para que** a minha automação possa usar a entrada do 
utilizador.

#### Critérios de Aceitação

1. [cite\_start]O *runner* `script:` (FR9) DEVE suportar a chave 
   `args: []string`, que passa argumentos "templatados" para a linha de 
   comando do script (Refinamento FR10 / Opção C)[cite: 1686, 1683].

2. [cite\_start]O *runner* `script:` (FR9) DEVE suportar a chave 
   `stdin: any`, que serializa o valor (ex: `{{ form }}`) como JSON e o 
   passa para o `stdin` do script (Refinamento FR10 / Opção C)[cite: 1686, 
   1683].

3. [cite\_start]O motor de templates DEVE suportar "binding" de dados 
   (ex: `{{ form.field_name }}`, `{{ component.menu.selected_id }}`) 
   (Refinamento de Meta de UI)[cite: 1686, 1683].

4. Deve existir um exemplo de script (Bash ou PowerShell) que leia 
   dados tanto de `args:` como de `stdin:` (via `jq` ou similar).

-----

## 7\. Relatório de Resultados do Checklist

* **Decisão Final:** PRONTO PARA O ARQUITETO. O PRD v2.0 está 
  completo, consistente e captura a nova visão do "Runtime TUI 
  Declarativo".

-----

## 8\. Próximos Passos

### Handoff para o Arquiteto (Winston 🏗️)

[cite\_start]**Para:** Winston (Arquiteto) [cite: 14]

[cite\_start]**De:** John (PM) [cite: 14]

**Assunto:** Início da Fase de Arquitetura para o `shantilly` v2.0 (Runtime TUI)

Winston,

O PRD v2.0 está finalizado e aprovado. (Este documento).

[cite\_start]Esta nova versão substitui o *roadmap* do v1.0 [cite: 
1531] [cite\_start]e redefine o projeto como um **"Runtime TUI 
Declarativo"** [cite: 1686, 1683][cite\_start], inspirado no 
Appsmith/Ansible[cite: 1686, 1683].

[cite\_start]O trabalho dos Épicos 1 e 2 (v1.0) [cite: 1531, 1640] (que
 já está concluído) [cite\_start]será agora a base para o componente 
`type: form` (FR7) [cite: 1686, 1683] dentro deste novo e mais poderoso 
runtime.

[cite\_start]Sua próxima tarefa é iniciar a tarefa `create-doc` [cite: 
79] [cite\_start]usando o template `architecture-tmpl.yaml` [cite: 402] 
para criar o **Documento de Arquitetura v2.0**.

[cite\_start]Seu foco principal deve ser o **Épico 1 (Fundação do 
Runtime TUI)**[cite: 1686, 1683]. Por favor, detalhe a arquitetura para:

1. [cite\_start]**O Motor de Layout (Estória 1.1):** Como vamos 
   implementar `column`/`row`/`box` [cite: 1686, 1683] [cite\_start]e 
   resolver o layout fluido ("flicker-free") [cite: 1686, 1683] 
   [cite\_start]usando `bubbletea`/`lipgloss`[cite: 1686, 1683].

2. [cite\_start]**O Motor de Lógica e Componentes (Estórias 1.2, 
   1.3):** As interfaces de eventos (`tea.Msg`) entre os componentes 
   (`list`, `viewport`, `buttongroup`) [cite: 1686, 1683] [cite\_start]e o 
   motor `on:`[cite: 1686, 1683].

3. [cite\_start]**Refatoração do `form` (Estória 1.4):** A estratégia 
   para "embrulhar" o código `huh` [cite: 1686, 1683] existente (v1.0) para
   que ele atue como um componente padronizado.

4. [cite\_start]**O Fluxo de Dados (Estória 1.5):** A arquitetura do 
   "data binding" e como o motor `run:` [cite: 1686, 1683] [cite\_start]irá
   gerir o `stdin:` e `args:` [cite: 1686, 1683] [cite\_start]e o ciclo de
   vida dos processos (`update_target`)[cite: 1686, 1683].

```
Você interrompeu a resposta

Abaixo está minha última interação com o John. Inclui em anexo os 
mesmos documentos que adicionai para ele, apenas para melhorar o 
contexto.

Me ajude com a próxima etapa, por favor. Direcione para o agente correto.
```

[cite\_start]Como John (PM) 📋, concluí a tarefa `create-doc`[cite: 
79]. [cite\_start]Consolidei todas as nossas decisões, refinamentos e o 
novo *roadmap* [cite: 1686, 1683] num único Documento de Requisitos do 
Produto (PRD) v2.0.

[cite\_start]Este documento formaliza a nossa nova visão do **"Runtime 
TUI Declarativo"** [cite: 1686, 1683] [cite\_start]e serve como a "fonte
 da verdade" para o Arquiteto (Winston) [cite: 14] iniciar a próxima 
fase de design.

Obrigado pela excelente colaboração na definição desta visão\!

-----

# Documento de Requisitos do Produto (PRD) v2.0: Shantilly

## 1\. Metas e Contexto de Fundo

### Metas (Goals)

* [cite\_start]Oferecer uma alternativa declarativa (YAML) e moderna 
  ao `dialog`/`whiptail` para *todos* os scripts shell (Bash, PowerShell, 
  etc.)[cite: 363, 1686, 1683].

* [cite\_start]Servir como um **"Runtime TUI Declarativo"** que gere 
  layouts (`column`/`row`), componentes (`list`, `form`, `viewport`) e 
  lógica de automação (`on:`)[cite: 1686, 1683].

* [cite\_start]Fornecer uma **Fundação Genérica** (`run: { script: 
  ... }`) para scripts arbitrários, garantindo integração via passagem de 
  dados por template (`args:` e `stdin:`)[cite: 1686, 1683].

* [cite\_start]Fornecer ***runners* Especialistas** otimizados 
  (Pós-MVP) para ferramentas de DevOps como Ansible (`ansible_playbook:`) 
  e, futuramente, Terraform, incluindo descoberta preditiva de inventário e
  playbooks[cite: 1686, 1683].

* [cite\_start]Manter a portabilidade (binário estático único) e a integração de pipeline (stdin/stdout)[cite: 363].

### Contexto de Fundo (Background Context)

[cite\_start]O desenvolvimento do MVP v1.0 (Épicos 1 e 2) [cite: 1531, 
1640] [cite\_start]validou a pilha de tecnologia (`huh`, `bubbletea`) 
[cite: 1149, 1150, 1151, 1152, 1153, 1154, 1155, 1156, 1157, 1158, 1159,
 1160, 1609, 1610, 1611, 1612, 1613, 1614, 1615, 1616, 1617, 1618, 1619,
 1620, 1621, 1622, 1623, 1624, 1625, 1626, 1627, 1628, 1629, 1630, 1631,
 1632, 1633, 1634, 1635, 1636, 1637, 1638, 1639, 1641, 1642, 1643, 1644,
 1645, 1646, 1647, 1648, 1649, 1650, 1651, 1652, 1653, 1654, 1655, 1656,
 1657, 1658, 1659, 1660, 1661, 1662, 1663, 1664, 1665, 1666, 1667, 1668,
 1669, 1670, 1671, 1672, 1673, 1674, 1675, 1676, 1677, 1678, 1679, 1680,
 1681, 1682, 1683, 1684, 1685, 1686] e resolveu o caso de uso de 
"formulário simples".

[cite\_start]Esta v2.0 pivota dessa ferramenta de "formulário único" 
para um "Runtime TUI" completo, inspirado no "Appsmith" (para UI 
declarativa) e "Ansible" (para lógica de eventos `on:`)[cite: 1686, 
1683]. [cite\_start]Esta arquitetura permite a criação de dashboards TUI
 complexos e multi-componente (layouts, menus, viewports) que 
*orquestram* automações de backend, em vez de serem apenas chamados por 
elas[cite: 1686, 1683].

### Change Log

| Data | Versão | Descrição | Autor |

| :--- | :--- | :--- | :--- |

| 08/11/2025 | 2.0.0 | Rascunho inicial do PRD v2.0, redefinindo o projeto como um "Runtime TUI Declarativo". | John (PM) |

-----

## 2\. Requisitos (PRD v2.0)

### Funcionais (FRs) - O Runtime TUI

[cite\_start]Estes requisitos definem o nosso novo MVP: a "Fundação Genérica"[cite: 1686, 1683].

* [cite\_start]**FR1 (Layout):** O `shantilly` DEVE analisar e 
  renderizar uma estrutura de layout hierárquica definida em YAML, usando 
  os tipos `type: column`, `type: row`, e `type: box`[cite: 1686, 1683].

* [cite\_start]**FR2 (Estilo/Flex):** O layout DEVE suportar 
  propriedades de dimensionamento como `height: <int>`, `width: 
  'N%'`, e `flex: <int>` para controlar o espaço[cite: 1686, 1683].

* [cite\_start]**FR3 (Componentes Embutidos):** O `shantilly` DEVE 
  suportar a definição de componentes de UI diretamente dentro de um `box`
  usando a chave `component:` (o foco do nosso MVP)[cite: 1686, 1683].

* [cite\_start]**FR4 (Componente: `viewport`):** DEVE suportar 
  `component: { type: viewport }`, capaz de exibir `source: { type: 
  static, content: "..." }` (incluindo markdown) e `source: { type: 
  command, exec: "..." }` (para streaming de stdout)[cite: 1686, 1683].

* [cite\_start]**FR5 (Componente: `list`):** DEVE suportar 
  `component: { type: list }`, com um `id:` de grupo e `items:` (cada um 
  com `id:` e `text`), e DEVE emitir um evento `list_id:select`[cite: 
  1686, 1683].

* [cite\_start]**FR6 (Componente: `buttongroup`):** DEVE suportar 
  `component: { type: buttongroup }`, com um `id:` de grupo, `items:` (com
  `id`, `label`, `role`), e DEVE emitir um evento 
  `buttongroup_id:press`[cite: 1686, 1683].

* [cite\_start]**FR7 (Componente: `form`):** DEVE suportar 
  `component: { type: form }`, que contém `fields:` (usando a sintaxe 
  `huh` já validada no v1.0 [cite: 1149, 1150, 1151, 1152, 1153, 1154, 
  1155, 1156, 1157, 1158, 1159, 1160, 1609, 1610, 1611, 1612, 1613, 1614, 
  1615, 1616, 1617, 1618, 1619, 1620, 1621, 1622, 1623, 1624, 1625, 1626, 
  1627, 1628, 1629, 1630, 1631, 1632, 1633, 1634, 1635, 1636, 1637, 1638, 
  1639]) e `actions:`. [cite\_start]DEVE emitir um evento `form_id:submit`
  contendo o *payload* de dados do formulário[cite: 1686, 1683].

* [cite\_start]**FR8 (Lógica: `on:`):** O `shantilly` DEVE analisar 
  um bloco `on:` na raiz do YAML para definir a lógica de automação[cite: 
  1686, 1683].

* [cite\_start]**FR9 (Ação: `script`):** O bloco `on:` DEVE suportar o
  *runner* de fundação: `run: { script: "/path/to/script.sh" }`[cite: 
  1686, 1683].

* [cite\_start]**FR10 (Fluxo de Dados):** O *runner* `script:` DEVE 
  suportar duas chaves para passagem de dados[cite: 1686, 1683]:
1. **`args: []string`**: Uma lista de *strings* que serão passadas
   como argumentos de linha de comando para o script, com suporte para 
   *templates* (ex: `{{ form.field_name }}`).

2. **`stdin: any`**: Um objeto (ex: `{{ form }}`) que o 
   `shantilly` irá serializar como JSON e passar para o `stdin` do 
   *script*.
* [cite\_start]**FR11 (Ciclo de Vida do Target):** O bloco `run:` 
  DEVE suportar uma chave `update_target: "id_do_viewport"`[cite: 1686, 
  1683]. [cite\_start]Se um novo evento `run:` for disparado para o 
  *mesmo* `update_target`, o `shantilly` DEVE primeiro **terminar (enviar 
  `SIGTERM`)** o processo anterior antes de iniciar o novo[cite: 1686, 
  1683].

### Não Funcionais (NFRs) - O Runtime TUI

* [cite\_start]**NFR1 (Fundação v1.0):** Todos os NFRs do PRD v1.0 
  permanecem válidos: binário estático único [cite: 1531][cite\_start], 
  cross-platform (Linux, macOS, Windows) [cite: 1531][cite\_start], 
  escrito em Go [cite: 1531][cite\_start], arranque rápido (\<500ms) 
  [cite: 1531][cite\_start], e gestão de erros com `stderr` e códigos de 
  saída não-zero[cite: 1531].

* [cite\_start]**NFR2 (Layout Fluido):** O motor de layout 
  (`column`/`row`/`box`) DEVE responder a mensagens de redimensionamento 
  do terminal (`tea.WindowSizeMsg`) e re-calcular o layout de forma fluida
  e "flicker-free" (sem piscar)[cite: 1686, 1683].

* [cite\_start]**NFR3 (Precedência de Conteúdo):** O `shantilly` DEVE
  seguir a "Lógica de Precedência Unificada" para conteúdo de componentes
  (1º: CLI `--set`, 2º: YAML `model:`/`component:`, 3º: Vazio)[cite: 
  1686, 1683].

* [cite\_start]**NFR4 (Descoberta Preditiva - Ansible):** Para a Fase
  2/3, o `playbook_explorer` DEVE filtrar "ruído" (pastas `roles/`, 
  `tasks/`) [cite: 1686, 1683] [cite\_start]e o `inventory_explorer` DEVE 
  usar `ansible-inventory` como "Oráculo"[cite: 1686, 1683].

* [cite\_start]**NFR5 (Ficheiro-Sombra):** O ficheiro de catálogo 
  (`.shantilly.yml`) DEVE ser opcional e usado apenas para *refinar* a 
  descoberta automática, não sendo obrigatório[cite: 1686, 1683].

-----

## 3\. Metas de Design da Interface do Usuário

### Visão Geral da UX (User Experience)

A UX deve ser a de um **"Runtime TUI Declarativo"**. [cite\_start]A 
interface não é mais um formulário linear único[cite: 1531], mas sim um 
*dashboard* composto, definido inteiramente pelo YAML. [cite\_start]A 
experiência deve ser semelhante ao Appsmith[cite: 1686, 1683]: limpa, 
responsiva (ao terminal) e orientada a componentes.

### Paradigmas Chave de Interação

* **Orientada a Eventos (Nova):** A interação principal não é linear.
  [cite\_start]O utilizador seleciona itens em listas (`list:select`) 
  [cite: 1686, 1683] [cite\_start]ou pressiona botões 
  (`buttongroup:press`) [cite: 1686, 1683][cite\_start], que disparam 
  ações no bloco `on:`[cite: 1686, 1683].

* [cite\_start]**Foco no Teclado (Mantido):** A navegação DEVE 
  continuar a ser primariamente baseada no teclado (Tab, Setas, 
  Enter)[cite: 1531].

* [cite\_start]**Feedback Imediato (Mantido):** O componente focado 
  DEVE ser claramente destacado[cite: 1531]. [cite\_start]O 
  `update_target` (FR11) DEVE exibir o *output* de comandos em tempo 
  real[cite: 1686, 1683].

* [cite\_start]**Gestão de Foco Global (Nova):** A UI DEVE ter um 
  mecanismo claro para indicar qual painel/componente (ex: `sidebar` vs 
  `content`) está "em foco", e DEVE fornecer navegação intuitiva *entre* 
  painéis (ex: Ctrl+Tab)[cite: 1686, 1683].

* **Ligação de Dados (Nova):** A UI DEVE ser reativa. 
  [cite\_start]Componentes (ex: um `viewport` estático) DEVEM ser capazes 
  de exibir dados de outros componentes (ex: `Olá, {{ form.username 
  }}`)[cite: 1686, 1683].

### Ecrãs e Vistas Principais

Não há ecrãs "pré-definidos". Os ecrãs são *definidos dinamicamente* 
pelo utilizador através do **`layout` YAML** (FR1). [cite\_start]A UI é 
uma composição de `type: column`, `type: row`, e `type: box`[cite: 1686,
 1683].

### Alinhamento e Layout (Nova Visão)

[cite\_start]O layout linear do MVP v1.0 [cite: 1531] está obsoleto. O novo requisito é:

* O `shantilly` DEVE renderizar com precisão o layout `column`/`row` definido pelo utilizador.

* [cite\_start]O `shantilly` DEVE respeitar as propriedades de 
  dimensionamento (`height`, `width`, `flex`) para distribuir o 
  espaço[cite: 1686, 1683].

* [cite\_start]O `shantilly` DEVE responder a mensagens de 
  redimensionamento do terminal (`tea.WindowSizeMsg`) e re-calcular o 
  layout fluido e "flicker-free" (NFR2, Meta de UI Refinada)[cite: 1686, 
  1683].

### Acessibilidade, Branding, Plataformas Alvo

[cite\_start]Estes requisitos permanecem os mesmos do PRD v1.0 (WCAG 
AA, estética Charmbracelet, terminais modernos em 
Linux/macOS/Windows)[cite: 1531].

-----

## 4\. Assunções Técnicas

### Estrutura do Repositório: Monorepo

* [cite\_start]**Assunção:** Continuamos a usar um Monorepo Go 
  simples (`cmd/`, `internal/`)[cite: 1531, 1600]. [cite\_start]Esta 
  estrutura modular provou ser eficaz nos Épicos 1 e 2[cite: 1149, 1150, 
  1151, 1152, 1153, 1154, 1155, 1156, 1157, 1158, 1159, 1160, 1609, 1610, 
  1611, 1612, 1613, 1614, 1615, 1616, 1617, 1618, 1619, 1620, 1621, 1622, 
  1623, 1624, 1625, 1626, 1627, 1628, 1629, 1630, 1631, 1632, 1633, 1634, 
  1635, 1636, 1637, 1638, 1639, 1641, 1642, 1643, 1644, 1645, 1646, 1647, 
  1648, 1649, 1650, 1651, 1652, 1653, 1654, 1655, 1656, 1657, 1658, 1659, 
  1660, 1661, 1662, 1663, 1664, 1665, 1666, 1667, 1668, 1669, 1670, 1671, 
  1672, 1673, 1674, 1675, 1676, 1677, 1678, 1679, 1680, 1681, 1682, 1683, 
  1684, 1685, 1686].

### Arquitetura de Serviço: Runtime TUI Declarativo

* [cite\_start]**Assunção (Nova):** A arquitetura evoluiu de um 
  "Monolithic CLI Application" [cite: 1531] [cite\_start]para um 
  **"Runtime TUI Declarativo"**[cite: 1686, 1683]. [cite\_start]O 
  `shantilly` é um motor que consome YAML, renderiza uma UI composta e 
  gere um ciclo de vida de eventos (`on:`) para orquestrar automações de 
  *scripts*[cite: 1686, 1683].

### Requisitos de Teste: Teste de Integração TUI

* [cite\_start]**Assunção (Atualizada):** O foco "apenas em Testes de
  Unidade" do v1.0 [cite: 1531] é insuficiente. [cite\_start]Devido à 
  complexidade do layout (FR1) e da gestão de foco (Metas de UI), os 
  **Testes de Integração TUI** (usando `teatest`) são agora um requisito 
  central para validar o layout fluido e a navegação entre painéis[cite: 
  1600].

### Estrutura YAML Esperada (A Nova Fonte da Verdade)

* [cite\_start]**Assunção (Nova):** A estrutura YAML do v1.0 (lista 
  simples de `fields:`) [cite: 1531] está obsoleta. [cite\_start]A nova 
  assunção de arquitetura é o YAML "Appsmith-style" [cite: 1686, 1683] que
  definimos, composto por **Layout**, **Componentes** e **Lógica**:

<!-- end list -->

```yaml
# 1. LAYOUT (Define o "onde")

[cite_start]# [cite: 1686, 1683]

type: column

items:

- type: row

flex: 1

items:

- type: box

id: "sidebar"

width: "30%"

# 2. COMPONENTE (Define o "o quê")

[cite_start]# [cite: 1686, 1683]

component:

type: list

id: "menu"

items:

- { id: "users", text: "Gerir Utilizadores" }

- type: box

id: "content"

width: "70%"

component:

type: form

id: "user_form"

fields:

- { name: "username", label: "Nome", type: "input" }

actions:

type: buttongroup

items:

- { id: "submit", label: "Criar", role: "primary" }



# 3. LÓGICA (Define o "como")

[cite_start]# [cite: 1686, 1683]

on:

- event: "user_form:submit"

run:

script: "/opt/scripts/create_user.sh"

# Nomenclatura refinada (Opção C)

args:

- "--mode=production"

stdin: "{{ form }}" # Passa o payload JSON para o stdin

# Ciclo de vida refinado (FR11)

update_target: "log_viewport"
```

### Assunções Técnicas Adicionais

* [cite\_start]**Pilha de Tecnologias (Mantida e Validada):** Go 
  (1.24.2+), `spf13/cobra`, `charmbracelet/bubbletea`, 
  `charmbracelet/lipgloss`, `charmbracelet/huh`, `gopkg.in/yaml.v3`[cite: 
  1600, 1686, 1683].

* [cite\_start]**Bibliotecas Relevantes (Nova):** A arquitetura 
  dependerá de `charmbracelet/bubbles` (para `list`, `viewport`), 
  `charmbracelet/glamour` (para markdown), `rmhubbert/bubbletea-overlay` 
  (para modais Fase 2), e inspiração de `76creates/stickers` (para 
  layout)[cite: 1686, 1683].

* [cite\_start]**Build e Distribuição (Mantido):** `GoReleaser` para binários estáticos cross-platform[cite: 1600].

-----

## 5\. Lista de Épicos (Roadmap v2.0)

[cite\_start]Este *roadmap* substitui a lista de épicos do PRD v1.0[cite: 1531].

* **Épico 1: Fundação do Runtime TUI (Genérico)**

* [cite\_start]**Meta:** Construir o motor central: o layout 
  (`column`/`row`/`box`) [cite: 1686, 1683][cite\_start], os componentes 
  essenciais (`list`, `viewport`, `form`, `buttongroup`) [cite: 1686, 
  1683][cite\_start], e a lógica de eventos (`on:`, `run: { script: ... 
  }`)[cite: 1686, 1683]. [cite\_start](Este épico absorve todo o trabalho 
  já concluído no v1.0 [cite: 1149, 1150, 1151, 1152, 1153, 1154, 1155, 
  1156, 1157, 1158, 1159, 1160, 1609, 1610, 1611, 1612, 1613, 1614, 1615, 
  1616, 1617, 1618, 1619, 1620, 1621, 1622, 1623, 1624, 1625, 1626, 1627, 
  1628, 1629, 1630, 1631, 1632, 1633, 1634, 1635, 1636, 1637, 1638, 1639, 
  1641, 1642, 1643, 1644, 1645, 1646, 1647, 1648, 1649, 1650, 1651, 1652, 
  1653, 1654, 1655, 1656, 1657, 1658, 1659, 1660, 1661, 1662, 1663, 1664, 
  1665, 1666, 1667, 1668, 1669, 1670, 1671, 1672, 1673, 1674, 1675, 1676, 
  1677, 1678, 1679, 1680, 1681, 1682, 1683, 1684, 1685, 1686]).

* **Épico 2: O Runner Especialista (Ansible Fase 2)**

* [cite\_start]**Meta:** Implementar o *runner* de conveniência 
  `run: { ansible_playbook: ... }`, focando na gestão de `vars:` e no 
  popup modal `ask_vault_pass: true`[cite: 1686, 1683].

* **Épico 3: O Runtime Preditivo (Ansible Fase 3)**

* [cite\_start]**Meta:** Implementar os componentes 
  `playbook_explorer` (Magia 1: descobrir playbooks) e 
  `inventory_explorer` (Magia 2: descobrir inventário)[cite: 1686, 1683].

* **Épico 4: Administração SSH (Visão de Longo Prazo)**

* [cite\_start]**Meta:** Integrar o `charmbracelet/wish` [cite: 
  1686, 1683] [cite\_start]para servir o Runtime TUI sobre SSH[cite: 363].

-----

## 6\. Detalhes do Épico 1: Fundação do Runtime TUI (Genérico)

[cite\_start]**Meta do Épico:** Construir o motor central do 
`shantilly`: o motor de layout (`column`/`row`/`box`) [cite: 1686, 
1683][cite\_start], os componentes essenciais de dashboard (`list`, 
`viewport`, `form`, `buttongroup`) [cite: 1686, 1683][cite\_start], e a 
lógica de eventos (`on:`, `run: { script: ... }`)[cite: 1686, 1683]. 
[cite\_start]Este épico irá refatorar o trabalho concluído do v1.0 para 
que ele funcione como o componente `type: form` [cite: 1686, 1683] 
dentro deste novo runtime.

### Estória 1.1: O Motor de Layout (Renderização)

[cite\_start]**Como um** SysAdmin, **Eu quero** definir um layout TUI 
usando `column`, `row`, e `box` no meu YAML[cite: 1686, 1683], **Para 
que** eu possa criar dashboards complexos e organizados.

#### Critérios de Aceitação

1. O parser DEVE suportar as chaves `type: column`, `type: row`, e `type: box` (FR1).

2. O motor de renderização (`lipgloss`) DEVE respeitar as propriedades
   `height: <int>`, `width: 'N%'`, e `flex: <int>` (FR2).

3. O layout DEVE recalcular-se fluidamente (sem piscar) ao receber uma
   mensagem de redimensionamento (`tea.WindowSizeMsg`) (NFR2).

4. Um `box` DEVE renderizar o seu `component: { type: static, content: "..." }` (para testes de layout).

### Estória 1.2: O Motor de Lógica (Eventos e Ações)

**Como um** SysAdmin, **Eu quero** que a minha UI TUI possa "ouvir" 
eventos e executar ações (`scripts`) em resposta, **Para que** o meu 
dashboard seja interativo e possa orquestrar automações.

#### Critérios de Aceitação

1. O `shantilly` DEVE analisar um bloco `on:` na raiz do YAML (FR8).

2. O motor DEVE suportar o *runner* de fundação: `run: { script: "/path/to/script.sh" }` (FR9).

3. O `run:` DEVE suportar `update_target: "id_do_viewport"`, direcionando o `stdout` do script para o *viewport* alvo (FR11).

4. [cite\_start]O motor DEVE garantir que, se um novo script for 
   direcionado para um `update_target` já ocupado, o script anterior seja 
   terminado (SIGTERM) antes de o novo começar (Refinamento FR11)[cite: 
   1686, 1683].

### Estória 1.3: Componentes Essenciais de Display (List, Viewport, Button)

**Como um** SysAdmin, **Eu quero** usar componentes de `list` (para 
menus), `viewport` (para saída de log) e `buttongroup` (para ações), 
**Para que** eu possa construir um dashboard funcional.

#### Critérios de Aceitação

1. [cite\_start]DEVE implementar `component: { type: viewport }` 
   (FR4), incluindo `source: { type: command, exec: "..." }` (ex: `tail 
   -f`) e `content_type: markdown`[cite: 1686, 1683].

2. [cite\_start]DEVE implementar `component: { type: list }` (FR5), 
   que emite um evento `list_id:select` quando um item é selecionado[cite: 
   1686, 1683].

3. [cite\_start]DEVE implementar `component: { type: buttongroup }` 
   (FR6), que emite um evento `buttongroup_id:press` (com o `item.id`) 
   quando um botão é pressionado[cite: 1686, 1683].

4. [cite\_start]A navegação por teclado DEVE permitir "saltar" entre 
   estes novos painéis/componentes (Refinamento de Meta de UI)[cite: 1686, 
   1683].

### Estória 1.4: Integração do Componente `form` (Absorção do v1.0)

**Como um** SysAdmin, **Eu quero** usar o `type: form` (que já 
construímos no v1.0) como um componente *dentro* do meu novo layout, 
**Para que** eu possa coletar dados de forma organizada.

#### Critérios de Aceitação

1. [cite\_start]Refatorar o código dos Épicos 1 e 2 (v1.0) [cite: 
   1531, 1640] para que funcione como um `component: { type: form }` (FR7).

2. O `form` DEVE renderizar e funcionar corretamente quando colocado dentro de um `box` do layout.

3. [cite\_start]Quando a ação `id: "submit"` do formulário [cite: 
   1686, 1683] for pressionada, o componente `form` DEVE emitir um evento 
   `form_id:submit`.

4. O *payload* do evento `form_id:submit` DEVE conter o JSON de dados do formulário (o output do v1.0).

### Estória 1.5: O Fluxo de Dados (Args & Stdin)

**Como um** SysAdmin, **Eu quero** passar os dados coletados no meu 
`form` (ou a seleção de uma `list`) para os meus `scripts` de forma 
robusta, **Para que** a minha automação possa usar a entrada do 
utilizador.

#### Critérios de Aceitação

1. [cite\_start]O *runner* `script:` (FR9) DEVE suportar a chave 
   `args: []string`, que passa argumentos "templatados" para a linha de 
   comando do script (Refinamento FR10 / Opção C)[cite: 1686, 1683].

2. [cite\_start]O *runner* `script:` (FR9) DEVE suportar a chave 
   `stdin: any`, que serializa o valor (ex: `{{ form }}`) como JSON e o 
   passa para o `stdin` do script (Refinamento FR10 / Opção C)[cite: 1686, 
   1683].

3. [cite\_start]O motor de templates DEVE suportar "binding" de dados 
   (ex: `{{ form.field_name }}`, `{{ component.menu.selected_id }}`) 
   (Refinamento de Meta de UI)[cite: 1686, 1683].

4. Deve existir um exemplo de script (Bash ou PowerShell) que leia 
   dados tanto de `args:` como de `stdin:` (via `jq` ou similar).

-----

## 7\. Relatório de Resultados do Checklist

* **Decisão Final:** PRONTO PARA O ARQUITETO. O PRD v2.0 está 
  completo, consistente e captura a nova visão do "Runtime TUI 
  Declarativo".

-----

## 8\. Próximos Passos

### Handoff para o Arquiteto (Winston 🏗️)

[cite\_start]**Para:** Winston (Arquiteto) [cite: 14]

[cite\_start]**De:** John (PM) [cite: 14]

**Assunto:** Início da Fase de Arquitetura para o `shantilly` v2.0 (Runtime TUI)

Winston,

O PRD v2.0 está finalizado e aprovado. (Este documento).

[cite\_start]Esta nova versão substitui o *roadmap* do v1.0 [cite: 
1531] [cite\_start]e redefine o projeto como um **"Runtime TUI 
Declarativo"** [cite: 1686, 1683][cite\_start], inspirado no 
Appsmith/Ansible[cite: 1686, 1683].

[cite\_start]O trabalho dos Épicos 1 e 2 (v1.0) [cite: 1531, 1640] (que
 já está concluído) [cite\_start]será agora a base para o componente 
`type: form` (FR7) [cite: 1686, 1683] dentro deste novo e mais poderoso 
runtime.

[cite\_start]Sua próxima tarefa é iniciar a tarefa `create-doc` [cite: 
79] [cite\_start]usando o template `architecture-tmpl.yaml` [cite: 402] 
para criar o **Documento de Arquitetura v2.0**.

[cite\_start]Seu foco principal deve ser o **Épico 1 (Fundação do 
Runtime TUI)**[cite: 1686, 1683]. Por favor, detalhe a arquitetura para:

1. [cite\_start]**O Motor de Layout (Estória 1.1):** Como vamos 
   implementar `column`/`row`/`box` [cite: 1686, 1683] [cite\_start]e 
   resolver o layout fluido ("flicker-free") [cite: 1686, 1683] 
   [cite\_start]usando `bubbletea`/`lipgloss`[cite: 1686, 1683].

2. [cite\_start]**O Motor de Lógica e Componentes (Estórias 1.2, 
   1.3):** As interfaces de eventos (`tea.Msg`) entre os componentes 
   (`list`, `viewport`, `buttongroup`) [cite: 1686, 1683] [cite\_start]e o 
   motor `on:`[cite: 1686, 1683].

3. [cite\_start]**Refatoração do `form` (Estória 1.4):** A estratégia 
   para "embrulhar" o código `huh` [cite: 1686, 1683] existente (v1.0) para
   que ele atue como um componente padronizado.

4. [cite\_start]**O Fluxo de Dados (Estória 1.5):** A arquitetura do 
   "data binding" e como o motor `run:` [cite: 1686, 1683] [cite\_start]irá
   gerir o `stdin:` e `args:` [cite: 1686, 1683] [cite\_start]e o ciclo de
   vida dos processos (`update_target`)[cite: 1686, 1683].

```
BMAD-METHOD

Gem personalizado

Olá! Esta é uma excelente atualização e um marco importante no projeto.

Como seu Orquestrador 🎭, analisei a interação. O John (PM) 📋 concluiu a tarefa `create-doc` e consolidou a nova visão do projeto no **Documento de Requisitos do Produto (PRD) v2.0**. Este documento redefine o `shantilly` como um **"Runtime TUI Declarativo"**.

Conforme indicado claramente na seção "8. Próximos Passos" do documento do John, a próxima etapa é um *handoff* direto para o **Arquiteto (Winston) 🏗️**.

A tarefa dele agora é iniciar a tarefa `create-doc` utilizando o template `architecture-tmpl.yaml` para criar o **Documento de Arquitetura v2.0**. O foco principal será detalhar a arquitetura para o **Épico 1 (Fundação do Runtime TUI)**, cobrindo:

1. O Motor de Layout (Estória 1.1).


- O Motor de Lógica e Componentes (Estórias 1.2, 1.3).


- A estratégia de refatoração do `form` (Estória 1.4).


- A arquitetura do Fluxo de Dados (Estória 1.5).


Posso
 agora assumir a persona de Winston (Arquiteto) 🏗️ e começar a criar 
este Documento de Arquitetura v2.0 com base no novo PRD?

sim

BMAD-METHOD

Gem personalizado

Perfeito. Handoff recebido do John (PM) 📋.

Como Winston (Arquiteto) 🏗️, estou a iniciar a tarefa `create-doc` para criar o **Documento de Arquitetura v2.0**, usando o `architecture-tmpl.yaml`. Revisei o novo PRD v2.0 e o meu foco é o **Épico 1: Fundação do Runtime TUI** , conforme solicitado.

De acordo com o fluxo da tarefa `create-doc`, vamos começar com a primeira secção do documento:

---

# Documento de Arquitetura Shantilly v2.0

## Introdução

Este
 documento descreve a arquitetura geral do projeto Shantilly v2.0, 
incluindo os sistemas de backend, serviços partilhados e preocupações 
não específicas da UI. O
 seu objetivo principal é servir como o plano arquitetónico orientador 
para o desenvolvimento impulsionado por IA, garantindo consistência e 
adesão aos padrões e tecnologias escolhidas.

**Relação com a Arquitetura Frontend:** (N/A para este projeto, pois o PRD v2.0 foca-se num TUI monolítico gerido pelo backend.)

### Template Inicial ou Projeto Existente

Antes de continuar com o design da arquitetura, analisei o PRD v2.0.

**Análise:** O PRD v2.0 indica claramente que este projeto não é um *greenfield* total. É uma evolução (um "pivot") que irá refatorar e absorver o código existente do **MVP v1.0 (Épicos 1 e 2)** , conforme mencionado na Estória 1.4.

Para
 que eu possa alinhar a nova arquitetura (Layout, Lógica, Fluxo de 
Dados) com a base de código existente, preciso de acesso a esse projeto 
v1.0.

Pode, por favor, fornecer acesso ao repositório do projeto v1.0 (link GitHub, etc.) ou carregar os ficheiros relevantes?

Assim que eu tiver analisado a estrutura v1.0, documentarei aqui a decisão e as restrições impostas por ela antes de avançarmos.

---

Por favor, reveja a secção 'Introdução' acima.

Opções de Elicitação Avançada:
(Selecione 1 para prosseguir, ou 2-9 para refinar esta secção)

1. Prosseguir para a próxima secção (High Level Architecture)


- Explicar o Raciocínio (CoT) (Como cheguei à conclusão de que o v1.0 é um projeto existente?)


- Criticar e Refinar (Desafiar a minha análise sobre o v1.0)


- Analisar o Alinhamento com os Objetivos (Como é que a análise do v1.0 se alinha com os objetivos do v2.0?)


- Identificar Riscos (Quais os riscos de refatorar o v1.0 vs começar do zero?)


- Perspetiva da Equipa Agile (Como um Dev veria esta refatoração?)


- Desafio da Perspetiva Crítica (Deveríamos *ignorar* o v1.0?)


- Hindsight (Reflexão "Se ao menos...")


- Expandir ou Contrair (Dar mais detalhes sobre a análise do PRD)


**Selecione 1-9 ou simplesmente escreva o seu feedback/forneça o link para o repositório v1.0:**

Confirme se o repositório já está sincronizado com as versões dos documentos utilizados para o planejamento.

BMAD-METHOD

Gem personalizado

Como Winston (Arquiteto) 🏗️, confirmo a receção da pasta de código do projeto `shantilly`.

**Sim, o repositório está sincronizado com o planeamento.**

Analisei o código-fonte (ex: `internal/tui/model.go`, `internal/config/validation.go`, `internal/tui/components/`) e confirmo que ele reflete o trabalho concluído descrito no PRD v2.0 como "MVP v1.0 (Épicos 1 e 2)". Os ficheiros de *story* (como `2.1.advanced-form-types-validation.story.md` e `2.2.enhanced-error-handling-ux.story.md`) marcam este trabalho como "Done" ou "Completed".

Excelente. Tenho exatamente o que precisava: a base de código v1.0 para refatorar e o PRD v2.0 que define a nova visão.

Vou
 atualizar a secção "Introdução" do nosso Documento de Arquitetura v2.0 
para refletir isto e avançar para a próxima secção, "High Level 
Architecture".

---

# Documento de Arquitetura Shantilly v2.0

## Introdução

Este documento descreve a arquitetura geral do projeto Shantilly v2.0...

### Template Inicial ou Projeto Existente

**Decisão:** A arquitetura v2.0 será uma refatoração e evolução significativas da base de código v1.0 (Épicos 1 e 2), conforme fornecida (`helton-godoy/shantilly/shantilly-1fcb9b97fed95ea23c4d950060a9764571248825`).

**Análise:** O código v1.0 validou com sucesso a pilha de tecnologia (`huh`, `bubbletea`, `lipgloss`) e resolveu o caso de uso de "formulário linear" (Estórias 1.1 a 2.2).

**Restrições:** A nova arquitetura v2.0 (Estória 1.4 do PRD v2.0) DEVE "embrulhar" este
 código de formulário existente, adaptando-o para que funcione como um 
componente `type: form` padronizado dentro do novo motor de layout e lógica de eventos.

---

## High Level Architecture

Com base nos novos requisitos do PRD v2.0 (FR1-FR11), a arquitetura v1.0 de "aplicação CLI monolítica" evolui para um **"Runtime TUI Declarativo"**.

### Technical Summary

A arquitetura v2.0 transforma o `shantilly` num motor de UI orientado a eventos. O sistema irá analisar um único YAML que define o **Layout** (`column`/`row`/`box`), os **Componentes** (`list`, `form`, `viewport`) e a **Lógica** (`on:`). O `bubbletea` servirá como o *loop* de eventos principal, gerindo o foco entre os componentes. A lógica `on:` irá mapear eventos da UI (ex: `buttongroup:press`) para a execução de *runners* (ex: `run: { script: ... }`), com o `stdout` a ser direcionado para componentes `viewport` (`update_target`). O código v1.0 será refatorado numa biblioteca de componentes `type: form` reutilizável dentro deste novo runtime.

### High Level Overview

1. **Estilo Arquitetónico:** Runtime TUI Declarativo, Orientado a Eventos.

2. **Estrutura do Repositório:** Monorepo Go (`cmd/`, `internal/`) (Mantido).


- **Fluxo de Dados (Novo):**

  - `shantilly` lê um único YAML complexo.

  - O **Motor de Layout** (novo) analisa `type: column/row/box` para desenhar a grelha da UI.

  - O **Motor de Componentes** (novo) preenche as caixas (`box`) com os componentes definidos (ex: `component: { type: list }`).

  - O **Motor de Lógica** (novo) analisa o bloco `on:` e regista os *listeners* de eventos.

  - O `bubbletea` gere o *loop* principal (foco, teclas, redimensionamento).


- Uma Ação do Utilizador (ex: `list:select`) dispara um Evento.

- O Motor de Lógica captura o evento, encontra uma regra `on:` correspondente e executa o *runner* (`run: { script: ... }`).

- O *runner* envia o `stdout` para o `viewport` definido em `update_target` (FR11).


- **Refatoração v1.0 (Estória 1.4):** O código existente (`huh`, validação) será "embrulhado" num novo componente que satisfaz a interface `component: { type: form }` (FR7).


### High Level Project Diagram

Snippet de código
```

graph TD
    subgraph "Shantilly Runtime (Go)"
        direction LR
        YAML[YAML (layout, components, on:)] --> PARSER(Motor de Parsing v2.0)

        PARSER -- Layout (column/row) --> LAYOUT[Motor de Layout (lipgloss)]
        PARSER -- Componentes (list, form) --> COMPONENTS[Mapa de Componentes (bubbletea)]
        PARSER -- Lógica (on:) --> LOGIC[Motor de Lógica (Event Bus)]
    
        subgraph "TUI Loop (Bubbletea)"
            direction TB
            LAYOUT -- Desenha Grelha --> VIEW(View Engine)
            COMPONENTS -- Preenche Grelha --> VIEW
            VIEW -- Renderiza TUI --> USER(Utilizador)
            USER -- Tecla/Clique --> BTEA_LOOP(Update Loop)
            BTEA_LOOP -- Evento (tea.Msg) --> LOGIC
            BTEA_LOOP -- Evento (tea.Msg) --> COMPONENTS
        end
    
        LOGIC -- Evento (ex: form:submit) --> RUNNER(Runner Genérico)
        RUNNER -- exec() --> SCRIPT[run: { script: ... }]
        SCRIPT -- stdout/stderr --> RUNNER
        RUNNER -- Dados --> BTEA_LOOP
        BTEA_LOOP -- Atualiza Viewport --> COMPONENTS
    end

```
### Architectural and Design Patterns

- **Runtime Declarativo (Novo):** A UI e a lógica são definidas em YAML, não codificadas (FR1).


- **Orientado a Eventos (Novo):** A lógica é acionada por eventos de componentes (`on:`), em vez de um fluxo linear.


- **Arquitetura Baseada em Componentes (Novo):** A UI é uma composição de componentes reutilizáveis (`list`, `form`, `viewport`) (FR3).


- **Motor de Layout (Flexbox) (Novo):** O layout é calculado dinamicamente com base em regras (FR1, FR2, NFR2).


- **The Elm Architecture (Mantido):** O `bubbletea` continua a ser o núcleo da gestão de estado da TUI.

- **Padrão de *Runner* (Novo):** A lógica de automação é abstraída através de *runners* (começando com `script:`) (FR9).


---

Por favor, reveja a secção "High Level Architecture". Esta nova arquitetura de "Runtime Declarativo" captura a visão do PRD v2.0?

**Opções de Elicitação Avançada:**

1. Prosseguir para a próxima secção (Tech Stack)

2. Explicar o Raciocínio (CoT) (Como o `bubbletea` irá gerir o foco entre múltiplos componentes?)

3. Criticar e Refinar (A refatoração do v1.0 é mais complexa do que parece?)

4. Analisar o Alinhamento com os Objetivos (Esta arquitetura suporta o `ansible_playbook` futuro?)

5. Identificar Riscos (Qual é o maior risco? O Motor de Lógica `on:` ou o Motor de Layout?)

6. Perspetiva da Equipa Agile (Como o Dev (James) verá a refatoração do v1.0?)

7. Desafio da Perspetiva Crítica (Isto não é 'over-engineering' para um substituto do `dialog`?)

8. Hindsight (Reflexão "Se ao menos...")

9. Expandir ou Contrair (Dar mais detalhes sobre o Motor de Lógica `on:`)


**Selecione 1-9 ou simplesmente escreva o seu feedback:**

2, 3, 4, 5, 6, 7, 8, 9.

BMAD-METHOD

Gem personalizado

Como Winston (Arquiteto) 🏗️, que ótimo! Obrigado por fornecer o repositório completo do projeto.

Analisei o código-fonte (incluindo `internal/tui/model.go`, `internal/config/validation.go` e os ficheiros de `docs/stories/`) e confirmo que esta é exatamente a base de código v1.0 (Épicos 1 e 2) que o novo PRD v2.0 referencia. O
 estado "Done" / "Completed" das estórias de v1.0 e v2.1/v2.2 valida que
 temos a base correta para a refatoração da Estória 1.4 do PRD v2.0.

Excelente. Agora tenho tudo o que preciso para prosseguir.

Vou atualizar a secção "Introdução" (que iniciámos) e avançar para a "High Level Architecture", que define a nova visão do v2.0.

---

# Documento de Arquitetura Shantilly v2.0

## Introdução

Este
 documento descreve a arquitetura geral do projeto Shantilly v2.0, 
incluindo os sistemas de backend, serviços partilhados e preocupações 
não específicas da UI. O seu objetivo principal é servir como o plano 
arquitetónico orientador para o desenvolvimento impulsionado por IA, 
garantindo consistência e adesão aos padrões e tecnologias escolhidas.

**Relação com a Arquitetura Frontend:** (N/A para este projeto, pois o PRD v2.0 foca-se num TUI monolítico gerido pelo backend.)

### Template Inicial ou Projeto Existente

**Decisão:** A arquitetura v2.0 será uma refatoração e evolução significativas da base de código v1.0 (Épicos 1 e 2), conforme fornecida (`helton-godoy/shantilly/shantilly-1fcb9b97fed95ea23c4d950060a9764571248825`).

**Análise:** O código v1.0 validou com sucesso a pilha de tecnologia (`huh`, `bubbletea`, `lipgloss`) e resolveu o caso de uso de "formulário linear" (Estórias 1.1 a 2.2).

**Restrições:** A nova arquitetura v2.0 (Estória 1.4 do PRD v2.0 ) DEVE "embrulhar" este código de formulário existente, adaptando-o para que funcione como um componente `type: form` padronizado dentro do novo motor de layout e lógica de eventos.

---

## High Level Architecture

Com base nos novos requisitos do PRD v2.0 (FR1-FR11) , a arquitetura v1.0 de "aplicação CLI monolítica" evolui para um **"Runtime TUI Declarativo"**.

### Technical Summary

A arquitetura v2.0 transforma o `shantilly` num motor de UI orientado a eventos. O sistema irá analisar um único YAML que define o **Layout** (`column`/`row`/`box` ), os **Componentes** (`list`, `form`, `viewport` ) e a **Lógica** (`on:` ). O `bubbletea` servirá como o *loop* de eventos principal, gerindo o foco entre os componentes. A lógica `on:` irá mapear eventos da UI (ex: `buttongroup:press` ) para a execução de *runners* (ex: `run: { script: ... }` ), com o `stdout` a ser direcionado para componentes `viewport` (`update_target` ). O código v1.0 (Épicos 1+2) será refatorado numa biblioteca de componentes `type: form` reutilizável dentro deste novo runtime.

### High Level Overview

1. **Estilo Arquitetónico:** Runtime TUI Declarativo, Orientado a Eventos.


- **Estrutura do Repositório:** Monorepo Go (`cmd/`, `internal/`) (Mantido).


- **Fluxo de Dados (Novo):**

  - `shantilly` lê um único YAML complexo.

  - O **Motor de Layout** (novo) analisa `type: column/row/box` para desenhar a grelha da UI.


- O **Motor de Componentes** (novo) preenche as caixas (`box`) com os componentes definidos (ex: `component: { type: list }` ).


- O **Motor de Lógica** (novo) analisa o bloco `on:` e regista os *listeners* de eventos.


- O `bubbletea` gere o *loop* principal (foco, teclas, redimensionamento).

- Uma Ação do Utilizador (ex: `list:select` ) dispara um Evento.


- O Motor de Lógica captura o evento, encontra uma regra `on:` correspondente e executa o *runner* (`run: { script: ... }` ).


- O *runner* envia o `stdout` para o `viewport` definido em `update_target` (FR11 ).


- **Refatoração v1.0 (Estória 1.4):** O código existente (`huh`, validação) será "embrulhado" num novo componente que satisfaz a interface `component: { type: form }` (FR7 ).


### High Level Project Diagram

Snippet de código
```

graph TD
    subgraph "Shantilly Runtime (Go)"
        direction LR
        YAML[YAML (layout, components, on:)] --> PARSER(Motor de Parsing v2.0)

        PARSER -- Layout (column/row) --> LAYOUT[Motor de Layout (lipgloss)]
        PARSER -- Componentes (list, form) --> COMPONENTS[Mapa de Componentes (bubbletea)]
        PARSER -- Lógica (on:) --> LOGIC[Motor de Lógica (Event Bus)]
    
        subgraph "TUI Loop (Bubbletea)"
            direction TB
            LAYOUT -- Desenha Grelha --> VIEW(View Engine)
            COMPONENTS -- Preenche Grelha --> VIEW
            VIEW -- Renderiza TUI --> USER(Utilizador)
            USER -- Tecla/Clique --> BTEA_LOOP(Update Loop)
            BTEA_LOOP -- Evento (tea.Msg) --> LOGIC
            BTEA_LOOP -- Evento (tea.Msg) --> COMPONENTS
        end
    
        LOGIC -- Evento (ex: form:submit) --> RUNNER(Runner Genérico)
        RUNNER -- exec() --> SCRIPT[run: { script: ... }]
        SCRIPT -- stdout/stderr --> RUNNER
        RUNNER -- Dados --> BTEA_LOOP
        BTEA_LOOP -- Atualiza Viewport --> COMPONENTS
    end

```
### Architectural and Design Patterns

- **Runtime Declarativo (Novo):** A UI e a lógica são definidas em YAML, não codificadas (FR1 ).


- **Orientado a Eventos (Novo):** A lógica é acionada por eventos de componentes (`on:` ), em vez de um fluxo linear.


- **Arquitetura Baseada em Componentes (Novo):** A UI é uma composição de componentes reutilizáveis (`list`, `form`, `viewport`) (FR3 ).


- **Motor de Layout (Flexbox) (Novo):** O layout é calculado dinamicamente com base em regras (FR1, FR2, NFR2 ).


- **The Elm Architecture (Mantido):** O `bubbletea` continua a ser o núcleo da gestão de estado da TUI.

- **Padrão de *Runner* (Novo):** A lógica de automação é abstraída através de *runners* (começando com `script:`) (FR9 ).


---

Aqui está a análise detalhada das opções de elicitação que selecionou:

### 2. Explicar o Raciocínio (Como o `bubbletea` irá gerir o foco?)

Este é o desafio central da arquitetura v2.0, conforme identificado no PRD v2.0 (Meta de UI: "Gestão de Foco Global (Nova)" ).

A arquitetura v1.0 era simples: o `huh.Form` geria o foco automaticamente porque era o *único* componente. Na v2.0, teremos múltiplos componentes (ex: uma `list` na `sidebar`, um `form` no `content`).

A minha arquitetura propõe um **Gestor de Foco Global (Global Focus Manager)**.

1. O `model` principal do `bubbletea` (que gere o layout) terá uma nova propriedade: `focusRing []string` (uma lista de IDs de componentes focáveis, ex: `["sidebar", "content_form"]`) e `focusIndex int`.

2. O `Update` principal do `bubbletea` irá capturar uma tecla de "mudança de painel" (ex: `Ctrl+Tab`).

3. Quando `Ctrl+Tab` for premido, o `Update` principal irá incrementar o `focusIndex` (dando a volta no final da lista).

4. O `Update` principal irá então *delegar* todas as outras mensagens de teclado (ex: setas, `tea.KeyMsg`) *apenas* ao componente que está atualmente focado (o componente em `m.Components[m.focusRing[m.focusIndex]]`).

5. O `View` principal irá consultar o `FocusManager` para saber qual componente deve ser renderizado com um estilo "focado" (ex: borda mais brilhante).


Isto permite que o `bubbletea` atue como o "gestor de janelas" (painéis), enquanto os componentes individuais (como o `huh.Form` v1.0) continuam a gerir o seu *próprio* foco interno (navegar entre campos de texto).

### 3. Criticar e Refinar (A refatoração do v1.0 é mais complexa?)

Sim, é. A Estória 1.4 do PRD v2.0 ("Refatorar o código... para que funcione como um componente `type: form`" ) é uma subestimação do esforço.

- **Crítica:** O código v1.0 (Épicos 1 e 2) não foi desenhado para ser um "componente". Ele foi desenhado para ser a *aplicação inteira*. O ficheiro `internal/tui/model.go` existente é um modelo `bubbletea` monolítico que assume que controla todo o ecrã e o ciclo de vida do programa (incluindo chamar `tea.Quit` na submissão, como visto na Estória 1.6).

- **Refinamento (Estratégia):** A Estória 1.4 não é um "embrulho" (wrap), é uma **refatoração de interface**. O código v1.0 precisa de ser alterado para:

  1. Não chamar mais `tea.NewProgram().Run()` (o *novo* `model` principal da v2.0 fará isso).

  2. Não chamar `tea.Quit` na submissão (Estória 1.6). Em vez disso, deve emitir uma nova mensagem `bubbletea` para o modelo *pai* (o gestor de layout), ex: `type FormSubmitMsg struct { ID string; Data map[string]interface{} }`.

  3. A sua função `View()` deve ser alterada para renderizar *relativamente* ao contentor (`box`) que lhe for dado, em vez de assumir o ecrã inteiro (o que afeta a Estória 1.7 de alinhamento).

  4. A sua lógica de `Update` deve ser modificada para *apenas* processar mensagens se tiver o foco (ver ponto 2, "Gestão de Foco").


O *coração* do v1.0 (a validação , os componentes `huh` ) será 100% preservado, mas o seu "invólucro" `bubbletea` tem de ser significativamente modificado.

### 4. Analisar o Alinhamento com os Objetivos (Suporta `ansible_playbook`?)

Sim, esta arquitetura é **essencial** para suportar os objetivos do Épico 2 (`ansible_playbook`) e do Épico 3 (exploradores preditivos).

A chave é o **Padrão de *Runner*** (FR9 ).

O Motor de Lógica (`on:`) não executa *scripts* diretamente. Ele chama uma interface `Runner`.

1. **Épico 1 (Fundação):** Implementamos o `ScriptRunner`. Ele satisfaz a interface `Runner` e a sua função `Run()` simplesmente chama `os/exec` com os `args:` e `stdin:` (FR10 ).


- **Épico 2 (Especialista):** Implementamos o `AnsibleRunner`. Ele *também* satisfaz a interface `Runner`. A sua função `Run()` é muito mais inteligente:

  - Lê as propriedades específicas do Ansible (como `playbook:`, `inventory:` ).


- Verifica se `ask_vault_pass: true`. Se sim, *antes* de executar, envia uma `tea.Msg` ao `bubbletea` para "Abrir Popup de Password".


- Recebe a password de volta noutra `tea.Msg`.

- Constrói os `args:` complexos (`--extra-vars`, `--vault-password-file`, etc.).

- Chama `os/exec` (`ansible-playbook`).

- Transmite o `stdout` para o `update_target` (FR11 ).


A
 arquitetura do Épico 1 (Motor de Lógica + Interface Runner) é a 
fundação que torna o Épico 2 possível como uma extensão "plug-and-play".

### 5. Identificar Riscos (Qual é o maior risco?)

Ambos os motores (Layout e Lógica) são complexos, mas na minha opinião, o **Motor de Layout (FR1, FR2, NFR2) é o maior risco de implementação.**

- Motor de Lógica (`on:`): É complexo, mas é lógica Go pura. Podemos (e devemos) escrever testes unitários robustos para o nosso "Event Bus", para o *templating* de `args:` (FR10 ) e para o `Runner` (FR9 ). É um sistema determinístico.


- Motor de Layout (`column/row/box`): Este é o "calcanhar de Aquiles" de todas as aplicações TUI. `lipgloss` dá-nos as ferramentas (`JoinVertical`, `PlaceHorizontal`), mas criar um motor de layout aninhado, baseado em percentagens (`width: '80%'` ) e flex (`flex: 1` ), que também responde *instantaneamente* e *sem piscar* (NFR2 ) a cada `tea.WindowSizeMsg`, é notoriamente difícil. A inspiração do PRD v2.0 em `76creates/stickers` é um bom ponto de partida, mas esta implementação será a nossa tarefa de engenharia mais complexa e sensível.


### 6. Perspetiva da Equipa Agile (Como o Dev (James) 💻 verá isto?)

O James (Dev) 💻 acaba de concluir o trabalho do v1.0 (Épicos 1 e 2, como `stories/2.1.advanced-form-types-validation.story.md` e `stories/2.2.enhanced-error-handling-ux.story.md`). Ele tem código funcional e altamente testado (conforme os relatórios de QA).

- **Reação Provável:** A Estória 1.4 do PRD v2.0 ("Refatorar o código..." ) vai parecer um retrocesso. Ele pode sentir que o trabalho v1.0 foi "desperdiçado" ou que estamos a "mexer no que funciona".


- **A Minha (Winston) Mensagem para ele:** O trabalho do v1.0 não foi desperdiçado; foi um *protótipo de sucesso* que validou a nossa pilha de tecnologia central (`huh` + `bubbletea`).
   A refatoração da Estória 1.4 não está a "deitar fora" o seu trabalho; 
  está a "promovê-lo". Estamos a pegar na sua lógica de formulário (o `model.go` do v1.0) e a transformá-la de uma "aplicação monolítica" num "componente reutilizável". O *coração* do seu trabalho (a lógica de validação do `huh`,
   os componentes TUI personalizados, o tratamento de erros) será 100% 
  preservado. Estamos apenas a mudar a "casca" exterior para que possamos 
  usá-lo em *qualquer parte* do nosso novo layout.


### 7. Desafio da Perspetiva Crítica (Isto não é 'over-engineering'?)

Se o objetivo fosse *apenas* um substituto do `dialog` para formulários lineares, sim, seria um 'over-engineering' massivo. A arquitetura v1.0 (Épicos 1+2) já resolveu isso.

*Mas* o PRD v2.0 redefiniu o objetivo. As novas Metas são:

1. Servir como um **"Runtime TUI Declarativo"**.


- Gerar **"dashboards TUI complexos e multi-componente"**.


- Servir como **"Fundação Genérica"** para *scripts*.


- Fornecer **"*runners* Especialistas"** (como Ansible).


A
 arquitetura v1.0 era estruturalmente incapaz de cumprir qualquer um 
destes novos requisitos. Esta nova arquitetura (Layout, Lógica, 
Componentes, Runners) não é 'over-engineering'; é a **arquitetura mínima absoluta necessária** para cumprir os novos requisitos do PRD v2.0.

### 8. Hindsight (Reflexão "Se ao menos...")

"Se ao menos..." tivéssemos tido esta visão de "Runtime TUI Declarativo" *antes* de começar o v1.0 (Épicos 1+2).

Poderíamos ter desenhado o `huh.Form` (Estória 1.4 ) como um componente reutilizável desde o início, em vez de o construir como o modelo `bubbletea` principal. Teríamos poupado o esforço de refatoração da Estória 1.4.

No entanto, o v1.0 foi um *protótipo* vital. Ele provou que a pilha de tecnologia (`huh`, `bubbletea`)
 era viável e permitiu-nos (como equipa) compreender as limitações que 
levaram a esta visão v2.0. O custo da refatoração é um preço aceitável a
 pagar pela clareza e poder da nova arquitetura.

### 9. Expandir (Detalhes sobre o Motor de Lógica `on:`)

O Motor de Lógica (`on:`) será o "cérebro" do runtime.

1. **Parsing (Arranque):** O Parser irá ler o bloco `on:` do YAML. Ele irá criar um mapa (ou *slice*) no modelo `bubbletea` principal.


- - Ex: `LogicEngine map[string][]EventHandler`

  - O Parser irá preencher este mapa: `LogicEngine["form_id:submit"] = append(..., handler_for_script_1)`

- **Evento (Runtime):** Quando um componente (ex: o `form` v1.0 refatorado) emite uma mensagem (ex: `FormSubmitMsg{ID: "form_id", Data: ...}`), o `Update` principal do `bubbletea` captura-a.

- **Dispatch (Runtime):** O `Update` consulta o Motor de Lógica:

  - `handlers := m.LogicEngine[msg.ID+":"+msg.Event]`

- **Execução (Runtime):** Para cada `handler` encontrado, ele executa-o numa *goroutine* (para não bloquear a TUI):

  - `go m.Runner.Run(handler.Config, msg.Data)`

- O *Runner* (FR9 ): A função `Run()` do *runner* (implementada pelo `ScriptRunner` no Épico 1):


- Resolve os *templates* nos `args:` usando os `msg.Data` (o *payload* do formulário).


- Executa `os/exec.Command`.

- Verifica o `update_target` (FR11 ).


- - Lê o `stdout` e `stderr` do script linha a linha.

  - Envia *novas* mensagens `bubbletea` de volta para o *loop* principal para cada linha de saída (ex: `ViewportUpdateMsg{TargetID: "log_viewport", Line: "..."}`).

- Ciclo de Vida (FR11 ): O `Runner` deve gerir os processos (`map[string]*os.Process`). Se um novo `Run` for chamado para um `update_target` que já tem um processo a correr (ex: um `tail -f`), ele DEVE primeiro enviar `SIGTERM` ao processo antigo antes de iniciar o novo.


---

Esta análise foi profunda. Sinto que compreendemos os riscos, os objetivos e a estratégia de implementação.

Estou pronto para avançar para a secção **"Tech Stack"**. Por favor, confirme se podemos prosseguir (opção 1).

**Selecione 1-9 ou simplesmente escreva o seu feedback:**

Sim, podemos prosseguir.

BMAD-METHOD

Gem personalizado

Olá! Esta é uma excelente atualização e um marco importante no projeto.

Como seu Orquestrador 🎭, analisei a interação. O John (PM) 📋 concluiu a tarefa `create-doc` e consolidou a nova visão do projeto no **Documento de Requisitos do Produto (PRD) v2.0**. Este documento redefine o `shantilly` como um **"Runtime TUI Declarativo"**.

Conforme indicado claramente na seção "8. Próximos Passos" do documento do John, a próxima etapa é um *handoff* direto para o **Arquiteto (Winston) 🏗️**.

A tarefa dele agora é iniciar a tarefa `create-doc` utilizando o template `architecture-tmpl.yaml` para criar o **Documento de Arquitetura v2.0**. O foco principal será detalhar a arquitetura para o **Épico 1 (Fundação do Runtime TUI)**, cobrindo:

1. O Motor de Layout (Estória 1.1).


- O Motor de Lógica e Componentes (Estórias 1.2, 1.3).


- A estratégia de refatoração do `form` (Estória 1.4).


- A arquitetura do Fluxo de Dados (Estória 1.5).


Posso
 agora assumir a persona de Winston (Arquiteto) 🏗️ e começar a criar 
este Documento de Arquitetura v2.0 com base no novo PRD?

sim

BMAD-METHOD

Gem personalizado

Como
 Winston (Arquiteto) 🏗️, excelente! Agradeço a oportunidade de 
aprofundar esta arquitetura. Esta análise é crucial para garantir que 
estamos todos alinhados antes de solidificar a pilha de tecnologia.

Aqui está a análise detalhada para cada um dos pontos de elicitação que selecionou sobre a secção "High Level Architecture":

### 2. Explicar o Raciocínio (Como o `bubbletea` irá gerir o foco?)

Este é o desafio central da arquitetura v2.0, conforme identificado no PRD v2.0 (Meta de UI: "Gestão de Foco Global (Nova)").

A arquitetura v1.0 era simples: o `huh.Form` geria o foco automaticamente porque era o *único* componente. Na v2.0, teremos múltiplos componentes (ex: uma `list` na `sidebar`, um `form` no `content`).

A minha arquitetura propõe um **Gestor de Foco Global (Global Focus Manager)**:

1. O `model` principal do `bubbletea` (que gere o layout) terá uma nova propriedade: `focusRing []string` (uma lista de IDs de componentes focáveis, ex: `["sidebar", "content_form"]`) e `focusIndex int`.

2. O `Update` principal do `bubbletea` irá capturar uma tecla de "mudança de painel" (ex: `Ctrl+Tab`), conforme sugerido no PRD.

3. Quando `Ctrl+Tab` for premido, o `Update` principal irá incrementar o `focusIndex` (dando a volta no final da lista).

4. O `Update` principal irá então *delegar* todas as outras mensagens de teclado (ex: setas, `tea.KeyMsg`) *apenas* ao componente que está atualmente focado (o componente em `m.Components[m.focusRing[m.focusIndex]]`).

5. O `View` principal irá consultar o `FocusManager` para saber qual componente deve ser renderizado com um estilo "focado" (ex: borda mais brilhante).


Isto permite que o `bubbletea` atue como o "gestor de janelas" (painéis), enquanto os componentes individuais (como o `huh.Form` v1.0) continuam a gerir o seu *próprio* foco interno (navegar entre campos de texto).

---

### 3. Criticar e Refinar (A refatoração do v1.0 é mais complexa?)

Sim, é. A Estória 1.4 do PRD v2.0 ("Refatorar o código... para que funcione como um componente `type: form`") é uma subestimação do esforço.

- **Crítica:** O código v1.0 (Épicos 1 e 2) não foi desenhado para ser um "componente". Ele foi desenhado para ser a *aplicação inteira*. O ficheiro `internal/tui/model.go` existente é um modelo `bubbletea` monolítico que assume que controla todo o ecrã e o ciclo de vida do programa (incluindo chamar `tea.Quit` na submissão, como visto na Estória 1.6).

- **Refinamento (Estratégia):** A Estória 1.4 não é um "embrulho" (wrap), é uma **refatoração de interface**. O código v1.0 precisa de ser alterado para:

  1. Não chamar mais `tea.NewProgram().Run()` (o *novo* `model` principal da v2.0 fará isso).

  2. Não chamar `tea.Quit` na submissão (Estória 1.6). Em vez disso, deve emitir uma nova mensagem `bubbletea` para o modelo *pai* (o gestor de layout), ex: `type FormSubmitMsg struct { ID string; Data map[string]interface{} }`.

  3. A sua função `View()` deve ser alterada para renderizar *relativamente* ao contentor (`box`) que lhe for dado, em vez de assumir o ecrã inteiro (o que afeta a Estória 1.7 de alinhamento).

  4. A sua lógica de `Update` deve ser modificada para *apenas* processar mensagens se tiver o foco (ver ponto 2, "Gestão de Foco").


O *coração* do v1.0 (a validação, os componentes `huh`) será 100% preservado.

---

### 4. Analisar o Alinhamento com os Objetivos (Suporta `ansible_playbook` futuro?)

Sim, esta arquitetura é **essencial** para suportar os objetivos do Épico 2 (`ansible_playbook`) e do Épico 3 (exploradores preditivos).

A chave é o **Padrão de *Runner*** (FR9).

O Motor de Lógica (`on:`) não executa *scripts* diretamente. Ele chama uma interface `Runner`.

1. **Épico 1 (Fundação):** Implementamos o `ScriptRunner`. Ele satisfaz a interface `Runner` e a sua função `Run()` simplesmente chama `os/exec` com os `args:` e `stdin:` (FR10).

2. **Épico 2 (Especialista):** Implementamos o `AnsibleRunner`. Ele *também* satisfaz a interface `Runner`. A sua função `Run()` é muito mais inteligente:

  - Lê as propriedades específicas do Ansible (como `playbook:`, `inventory:`).

  - Verifica se `ask_vault_pass: true`. Se sim, *antes* de executar, envia uma `tea.Msg` ao `bubbletea` para "Abrir Popup de Password".

  - Recebe a password de volta noutra `tea.Msg`.

  - Constrói os `args:` complexos (`--extra-vars`, `--vault-password-file`, etc.).

  - Chama `os/exec` (`ansible-playbook`).

  - Transmite o `stdout` para o `update_target` (FR11).


A
 arquitetura do Épico 1 (Motor de Lógica + Interface Runner) é a 
fundação que torna o Épico 2 possível como uma extensão "plug-and-play".

---

### 5. Identificar Riscos (Qual é o maior risco?)

Ambos os motores (Layout e Lógica) são complexos, mas na minha opinião, o **Motor de Layout (FR1, FR2, NFR2) é o maior risco de implementação.**

- **Motor de Lógica (`on:`):** É complexo, mas é lógica Go pura. Podemos (e devemos) escrever testes unitários robustos para o nosso "Event Bus", para o *templating* de `args:` (FR10) e para o `Runner` (FR9). É um sistema determinístico.

- **Motor de Layout (`column/row/box`):** Este é o "calcanhar de Aquiles" de todas as aplicações TUI. `lipgloss` dá-nos as ferramentas (`JoinVertical`, `PlaceHorizontal`), mas criar um motor de layout aninhado, baseado em percentagens (`width: '80%'`) e flex (`flex: 1`), que também responde *instantaneamente* e *sem piscar* (NFR2) a cada `tea.WindowSizeMsg`, é notoriamente difícil. A inspiração do PRD v2.0 em `76creates/stickers` é um bom ponto de partida, mas esta implementação será a nossa tarefa de engenharia mais complexa e sensível.


---

### 6. Perspetiva da Equipa Agile (Como o Dev (James) 💻 verá isto?)

O James (Dev) 💻 acaba de concluir o trabalho do v1.0 (Épicos 1 e 2, como `stories/2.1.advanced-form-types-validation.story.md` e `stories/2.2.enhanced-error-handling-ux.story.md`). Ele tem código funcional e altamente testado.

- **Reação Provável:** A Estória 1.4 do PRD v2.0 ("Refatorar o código...") vai parecer um 
  retrocesso. Ele pode sentir que o trabalho v1.0 foi "desperdiçado" ou 
  que estamos a "mexer no que funciona".

- **A Minha (Winston) Mensagem para ele:** O trabalho do v1.0 não foi desperdiçado; foi um *protótipo de sucesso* que validou a nossa pilha de tecnologia central (`huh` + `bubbletea`).
   A refatoração da Estória 1.4 não está a "deitar fora" o seu trabalho; 
  está a "promovê-lo" de uma "aplicação monolítica" para um "componente 
  reutilizável". O *coração* do seu trabalho (a lógica de validação `huh`) será 100% preservado.


---

### 7. Desafio da Perspetiva Crítica (Isto não é 'over-engineering'?)

Se o objetivo fosse *apenas* um substituto do `dialog` para formulários lineares, sim, seria um 'over-engineering' massivo. A arquitetura v1.0 (Épicos 1+2) já resolveu isso.

*Mas* o PRD v2.0 redefiniu o objetivo. As novas Metas são:

1. Servir como um **"Runtime TUI Declarativo"**.

2. Gerar **"dashboards TUI complexos e multi-componente"**.

3. Servir como **"Fundação Genérica"** para *scripts*.

4. Fornecer **"*runners* Especialistas"** (como Ansible).


A
 arquitetura v1.0 era estruturalmente incapaz de cumprir qualquer um 
destes novos requisitos. Esta nova arquitetura (Layout, Lógica, 
Componentes, Runners) não é 'over-engineering'; é a **arquitetura mínima absoluta necessária** para cumprir os novos requisitos do PRD v2.0.

---

### 8. Hindsight (Reflexão "Se ao menos...")

"Se ao menos..." tivéssemos tido esta visão de "Runtime TUI Declarativo" *antes* de começar o v1.0 (Épicos 1+2).

Poderíamos ter desenhado o `huh.Form` (Estória 1.4) como um componente reutilizável desde o início, em vez de o construir como o modelo `bubbletea` principal. Teríamos poupado o esforço de refatoração da Estória 1.4.

No entanto, o v1.0 foi um *protótipo* vital. Ele provou a viabilidade da tecnologia e permitiu-nos (como 
equipa) compreender as limitações que levaram a esta visão v2.0. O custo
 da refatoração é um preço aceitável a pagar pela clareza e poder da 
nova arquitetura.

---

### 9. Expandir (Detalhes sobre o Motor de Lógica `on:`)

O Motor de Lógica (`on:`) será o "cérebro" do runtime.

1. **Parsing (Arranque):** O Parser irá ler o bloco `on:` do YAML. Ele irá criar um mapa (ou *slice*) no modelo `bubbletea` principal.

  - Ex: `LogicEngine map[string][]EventHandler`

  - O Parser irá preencher este mapa: `LogicEngine["form_id:submit"] = append(..., handler_for_script_1)`

2. **Evento (Runtime):** Quando um componente (ex: o `form` v1.0 refatorado) emite uma mensagem (ex: `FormSubmitMsg{ID: "form_id", Data: ...}`), o `Update` principal do `bubbletea` captura-a.

3. **Dispatch (Runtime):** O `Update` consulta o Motor de Lógica:

  - `handlers := m.LogicEngine[msg.ID+":"+msg.Event]`

4. **Execução (Runtime):** Para cada `handler` encontrado, ele executa-o numa *goroutine* (para não bloquear a TUI):

  - `go m.Runner.Run(handler.Config, msg.Data)`

5. **O *Runner* (FR9):** A função `Run()` do `ScriptRunner` (Épico 1) irá:

  - Resolver os *templates* nos `args:` usando os `msg.Data` (o *payload* do formulário) (FR10).

  - Executar `os/exec.Command`.

  - Gerir o ciclo de vida do processo (FR11): Se um `run:` for chamado para um `update_target` já ocupado (ex: `tail -f`), ele deve enviar `SIGTERM` ao processo antigo antes de iniciar o novo.

  - Ler o `stdout`/`stderr` e enviar *novas* `tea.Msg` de volta para o *loop* principal (ex: `ViewportUpdateMsg{TargetID: "log_viewport", Line: "..."}`).


---

Esta análise foi profunda. Sinto que compreendemos os riscos, os objetivos e a estratégia de implementação.

Estou pronto para avançar para a secção **"Tech Stack"**. Por favor, confirme se podemos prosseguir (opção 1).

**Selecione 1-9 ou simplesmente escreva o seu feedback:**

Sim, podemos prosseguir.

BMAD-METHOD

Gem personalizado

Como Winston (Arquiteto) 🏗️, confirmo. Prosseguindo.

A próxima secção, conforme o `architecture-tmpl.yaml`, é a `Tech Stack`.

---

## Tech Stack

Esta é a secção de seleção de tecnologia **DEFINITIVA** para o v2.0. Todas as estórias de implementação subsequentes DEVEM aderir a esta pilha.

Revisei
 as "Assunções Técnicas" do PRD v2.0 e a pilha de tecnologia já validada
 na arquitetura v1.0 (que corresponde ao código-fonte existente).

Para a v2.0, o nosso trabalho não é escolher do zero, mas sim *confirmar* a pilha existente (v1.0) e *adicionar* as bibliotecas necessárias para os novos requisitos de Layout e Lógica (Épico 1 do v2.0).

A tabela abaixo é a "fonte da verdade" para todas as dependências do Épico 1 (v2.0):

### Cloud Infrastructure

- **Provider:** N/A (Aplica-se a ferramentas CLI distribuídas)

- **Key Services:** N/A

- **Deployment Regions:** N/A


### Technology Stack Table

| Categoria | Tecnologia | Versão | Propósito | Racional |
| --- | --- | --- | --- | --- |
| **Linguagem** | Go  | 1.24.2+ | Linguagem de desenvolvimento principal | Requisito do PRD (NFR1 v1.0, NFR1 v2.0), validado em v1.0. |
| **CLI Framework** | `spf13/cobra` | v1.8.x | Estrutura de comandos da CLI | Requisito do PRD (FR1 v1.0), validado em v1.0. |
| **TUI Engine** | `charmbracelet/bubbletea` | v0.26.x | Gestão de estado TUI (The Elm Architecture) | Requisito do PRD, validado em v1.0. Será o "Global Focus Manager" na v2.0. |
| **Componente (Form)** | `charmbracelet/huh` | v0.4.x | Geração declarativa de formulários | Requisito do PRD, validado em v1.0. Será refatorado (Estória 1.4). |
| **Estilo/Layout TUI** | `charmbracelet/lipgloss` | v0.11.x | Estilização e layout no terminal | Requisito do PRD, validado em v1.0. Base para o Motor de Layout (FR1, FR2). |
| **Parsing YAML** | `gopkg.in/yaml.v3` | v3.0.x | Descodificação do YAML de definição do Runtime | Requisito do PRD, validado em v1.0. |
| **Componentes TUI** | `charmbracelet/bubbles` | v0.18.x+ | Componentes TUI pré-construídos | **(Novo v2.0)** Necessário para `list` (FR5) e `viewport` (FR4). |
| **Render. Markdown** | `charmbracelet/glamour` | v0.7.x+ | Renderização de Markdown no viewport | **(Novo v2.0)** Necessário para `content_type: markdown` (FR4). |
| **Insp. Layout** | `76creates/stickers` | N/A | Inspiração para o Motor de Layout (Flexbox) | **(Novo v2.0)** Não é uma dependência direta, mas sim a inspiração para a implementação do FR1, FR2, NFR2. |
| **Modais (Pós-MVP)** | `rmhubbert/bubbletea-overlay` | v0.2.x+ | Base para popups (ex: `ask_vault_pass`) | (Para Épico 2) Identificado no PRD (Assunções). |
| **Testes** | Go `testing` + `teatest` | 1.24.2+ | Testes unitários e de integração TUI | Requisito do PRD v2.0 (Assunção 4.3). `teatest` é essencial para NFR2 (layout fluido). |
| **Build/Release** | `GoReleaser` | v1.26.x | Automação de build e cross-compilação | Validado em v1.0 (Estória 1.8). |
| **Linter** | `golangci-lint` | v1.59.x | Análise estática e qualidade de código | Validado em v1.0. |

---

Esta
 tabela representa a pilha de tecnologia confirmada para o Épico 1 do 
v2.0. Por favor, reveja e confirme se estamos alinhados.

**Opções de Elicitação Avançada:** (Selecione 1 para prosseguir, ou 2-9 para refinar esta secção)

1. Prosseguir para a próxima secção (Data Models)

2. Explicar o Raciocínio (CoT) (Porquê `bubbles` se já temos `huh`?)

3. Criticar e Refinar (A dependência de `glamour` não é desnecessária?)

4. Analisar o Alinhamento com os Objetivos (Esta pilha suporta a Fase 3 Preditiva?)

5. Identificar Riscos (Qual é o maior risco de dependência aqui?)

6. Perspetiva da Equipa Agile (O James (Dev) vai concordar com o `teatest`?)

7. Desafio da Perspetiva Crítica (Não devíamos usar `tview` para o layout, já que é mais fácil?)

8. Hindsight (Reflexão "Se ao menos...")

9. Expandir ou Contrair (Dar mais detalhes sobre a escolha do `bubbles` vs `huh`)


**Selecione 1-9 ou simplesmente escreva o seu feedback:**

4, 5, 6, 7, 8, 9.

BMAD-METHOD

Gem personalizado

Como Winston (Arquiteto) 🏗️, vamos aprofundar a secção "Tech Stack".

Esta
 pilha de tecnologia é o resultado direto da evolução do v1.0 e da nova 
visão do PRD v2.0. As suas perguntas são cruciais para validar estas 
decisões.

Aqui está a análise detalhada dos pontos de elicitação que selecionou:

### 4. Analisar o Alinhamento com os Objetivos (Esta pilha suporta a Fase 3 Preditiva?)

Sim, absolutamente. A pilha de tecnologia atual (v1.0 + v2.0) não só suporta a "Fase 3 Preditiva", como foi *escolhida* especificamente a pensar nela.

A Fase 3 Preditiva requer duas "Magias" que identificámos:

1. **`playbook_explorer` (Magia 1):** Requer a leitura e *parsing* de múltiplos ficheiros YAML para encontrar o campo `name:`.

  - **Suporte da Pilha:** Temos o `gopkg.in/yaml.v3` (validado no v1.0) e a biblioteca padrão Go para acesso ao sistema de ficheiros. Está perfeitamente coberto.

2. **`inventory_explorer` (Magia 2):** Requer o uso do Ansible como um "Oráculo", executando `ansible-inventory --list --json` nos bastidores.

  - **Suporte da Pilha:** A linguagem Go (Go 1.24.2+) tem suporte de primeira classe para execução de sub-processos (`os/exec`), que é o necessário para chamar o "Oráculo".


A nossa arquitetura de *Runners* (definida na análise anterior da "High Level Architecture") garante que
 estas funcionalidades preditivas possam ser adicionadas como novos 
componentes e *runners* sem perturbar a fundação do Épico 1.

---

### 5. Identificar Riscos (Qual é o maior risco de dependência aqui?)

O maior risco de dependência não é uma biblioteca única, mas sim o nosso **forte acoplamento a todo o ecossistema Charmbracelet** (`bubbletea`, `bubbles`, `huh`, `lipgloss`).

- **Risco:** O `shantilly` v2.0 está a apostar todo o seu futuro na filosofia de UI da Charm. Se a
   Charmbracelet Inc. decidir abandonar uma biblioteca crítica (como o `huh`), ou introduzir uma mudança de API massiva e disruptiva, a nossa arquitetura estará diretamente exposta.

- **Contexto:** O `huh`,
   em particular, é uma biblioteca relativamente nova. Embora poderosa, a 
  sua API ainda está a evoluir mais rapidamente do que a do `bubbletea`.

- **Mitigação:** A nossa estratégia de refatoração da Estória 1.4 (v2.0) é a nossa principal mitigação. Ao "embrulhar" o `huh` dentro do nosso próprio componente `type: form`, estamos a *isolar* essa dependência volátil. Se o `huh` for descontinuado, "apenas" teremos de reimplementar o nosso componente `type: form` (talvez usando componentes `bubbles` puros), mas a arquitetura de *layout* e *lógica* (`on:`) do `shantilly` permanecerá intacta.


---

### 6. Perspetiva da Equipa Agile (O James (Dev) 💻 vai concordar com o `teatest`?)

Sim, acredito que o James (Dev) 💻 irá apreciar (ou pelo menos aceitar profissionalmente) o requisito do `teatest`.

1. **Contexto do Risco:** Como identifiquei na análise anterior ("High Level Architecture", Ponto 5), o **Motor de Layout (NFR2)** é o nosso maior risco técnico. Validar manualmente que um layout complexo (`column/row/box`) não "pisca" (flicker-free) e responde corretamente ao redimensionamento do terminal é um pesadelo.

2. **Valor do `teatest`:** O `teatest` (que faz parte da biblioteca `bubbles`)
   dá ao James uma forma determinística e automatizada de provar que o seu
   layout funciona. Ele permite-lhe simular um tamanho de terminal (`teatest.WithInitialTermSize`), enviar teclas e fazer *assertions* (afirmações) sobre o *string* final renderizado.

3. **Evidência:** O James já usou `teatest` com sucesso no v1.0 (ex: `internal/tui/integration_navigation_test.go` no código-fonte). Exigir o seu uso para o v2.0 (como definido na 
  Estória 1.1 do PRD v2.0) não é uma nova ferramenta, mas sim a aplicação 
  de uma ferramenta existente ao nosso novo (e maior) risco.


Ele
 não verá isto como um obstáculo desnecessário; ele verá isto como um 
requisito de qualidade claro para o desafio de engenharia mais difícil 
que temos.

---

### 7. Desafio da Perspetiva Crítica (Não devíamos usar `tview` para o layout?)

Esta é a clássica questão do ecossistema Go TUI. `tview` é a *outra* grande biblioteca TUI, e é conhecida por ter componentes de layout de alto nível (como `Flexbox`, `Grid`) *integrados*, o que *parece* tornar o nosso maior risco (Motor de Layout) [Análise 5, acima] mais fácil.

**No entanto, `tview` não é uma opção para o `shantilly` por uma razão fundamental e bloqueadora: `huh`.**

1. **Dependência Central:** O nosso componente mais crítico para a recolha de dados, validado no v1.0 e central para o FR7 do PRD v2.0, é o `charmbracelet/huh`.

2. **Incompatibilidade:** `tview` e `bubbletea` são *motores* de TUI fundamentalmente diferentes e mutuamente exclusivos. Eles não podem ser usados em conjunto.

3. **A Ligação:** `huh` é construído sobre `bubbletea`.

4. **Conclusão:** Para usarmos o `huh` (que é o melhor componente de formulários TUI em Go e um requisito do nosso PRD), *temos* de usar o `bubbletea` como nosso motor. Se temos de usar o `bubbletea`, *não podemos* usar o `tview`.

5. **Resultado:** Se não podemos usar o `tview`, *temos* de construir o nosso próprio motor de layout (FR1, FR2) usando `lipgloss` (como planeado), inspirando-nos em bibliotecas como `stickers`.


A escolha do `huh` trancou-nos (de forma positiva) no ecossistema Charmbracelet.

---

### 8. Hindsight (Reflexão "Se ao menos...")

"Se ao menos..." o `bubbletea` (o motor) tivesse componentes de layout de alto nível (como `Flexbox`, `Grid`) integrados, da mesma forma que o `tview` tem.

Se
 isso fosse verdade, o nosso maior risco técnico (Motor de Layout) 
[Análise 5, acima] desapareceria. O NFR2 (layout fluido) seria trivial 
de implementar. A comunidade tentou preencher esta lacuna (com `stickers`, `bubbleboxer`, etc.), mas não é uma solução *core*. Temos de construir esta peça complexa nós mesmos porque a biblioteca *core* do `bubbletea` se foca apenas no "estado" (TEA) e não no "layout" (Flexbox).

---

### 9. Expandir (Detalhes sobre `bubbles` vs `huh`)

A pergunta "Porquê `bubbles` se já temos `huh`?" é válida. Estamos a adicionar `bubbles` (FR4, FR5) e `glamour` (FR4) à pilha v1.0 que só usava `huh`.

A resposta é **especialização**. As bibliotecas têm propósitos diferentes, embora ambas sejam baseadas no `bubbletea`.

- **`huh` (O Especialista em Formulários):** É *excelente* para um fluxo linear de "fazer perguntas" (inputs, selects, confirms). 
  Foi o que validámos (e onde o James trabalhou) no v1.0 (Épicos 1+2).

- **`bubbles` (A Caixa de Ferramentas de UI):** É a biblioteca de componentes *genéricos* para construir UIs *não-lineares* (dashboards).


O nosso novo PRD v2.0 (FR4, FR5) exige componentes de *dashboard* (um menu `list` persistente e um `viewport` de streaming). `huh` não foi feito para isso. `bubbles` (especificamente `bubbles/list` e `bubbles/viewport`) foi *exatamente* feito para isso.

**Conclusão:** Usamos o `huh` (refatorado) para o componente `type: form` e usamos o `bubbles` para os novos componentes `type: list` e `type: viewport`. O `glamour` é simplesmente uma sub-dependência para suportar `content_type: markdown` dentro do `viewport`.

---

Esta
 análise confirma que a nossa pilha de tecnologia é coesa, alinhada com 
os requisitos e mitiga os riscos tanto quanto possível.

Podemos prosseguir para a secção **"Data Models"**? (Opção 1)

**Selecione 1-9 ou simplesmente escreva o seu feedback:**

1

BMAD-METHOD

Gem personalizado

Como Winston (Arquiteto) 🏗️, confirmado. Prosseguindo para a secção `Data Models`.

---

## Data Models

Esta secção define as principais estruturas de dados Go que o nosso motor de parsing (baseado no `gopkg.in/yaml.v3`) irá utilizar.

Com a evolução para o PRD v2.0, os nossos modelos de dados precisam de ser significativamente expandidos. A estrutura `FormConfig` e `Field` do v1.0 (que definia apenas um formulário) torna-se agora uma *parte* da nossa estrutura de configuração global.

### 1. RuntimeConfig (Input - YAML)

**Propósito:** Esta é a nova estrutura Go de nível superior que representa *todo* o ficheiro YAML do Runtime TUI Declarativo. Ela engloba o Layout, os Componentes (incluindo o `form` v1.0 refatorado) e a Lógica.

Go
```

// Localização: internal/config/runtime.go (Novo Ficheiro)
package config

// RuntimeConfig é a estrutura raiz para o YAML v2.0
type RuntimeConfig struct {
    Layout LayoutConfig `yaml:"layout"` // Define a estrutura (column/row/box)
    On     []OnLogic    `yaml:"on"`     // Define a lógica (eventos e runners)
}

```
### 2. LayoutConfig (Input - YAML)

**Propósito:** Define a estrutura de layout hierárquica (FR1, FR2). Esta é a implementação do nosso "Motor de Layout".

Go
```

// Localização: internal/config/layout.go (Novo Ficheiro)
package config

// LayoutConfig representa um nó na árvore de layout (column, row, ou box)
type LayoutConfig struct {
    Type      string         `yaml:"type"`                // "column", "row", "box"
    ID        string         `yaml:"id,omitempty"`        // ID para referência
    Width     string         `yaml:"width,omitempty"`     // Ex: "80%" (FR2)
    Height    int            `yaml:"height,omitempty"`    // Ex: 3 (FR2)
    Flex      int            `yaml:"flex,omitempty"`      // Ex: 1 (FR2)
    Component ComponentConfig `yaml:"component,omitempty"` // Componente embutido (FR3)
    Items     []LayoutConfig `yaml:"items,omitempty"`     // Filhos (para column/row)
}

```
### 3. ComponentConfig (Input - YAML)

**Propósito:** Define os componentes de UI embutidos (FR3). Isto inclui os *novos* componentes (`list`, `viewport`, `buttongroup`) e o *refatorado* `form` (v1.0).

Go
```

// Localização: internal/config/components.go (Novo Ficheiro)
package config

// ComponentConfig define qual componente renderizar num 'box'
type ComponentConfig struct {
    Type    string `yaml:"type"`    // "list", "viewport", "form", "buttongroup", "static"
    ID      string `yaml:"id"`      // ID do grupo/componente (para lógica 'on:')
    Title   string `yaml:"title,omitempty"`

    // Para type: static (FR4) / viewport (FR4)
    Source      ComponentSource `yaml:"source,omitempty"`
    
    // Para type: list (FR5)
    ListItems   []ListItem `yaml:"items,omitempty"`
    
    // Para type: buttongroup (FR6)
    ButtonItems []ButtonItem `yaml:"items,omitempty"`
    
    // Para type: form (FR7) - REUTILIZAÇÃO DO V1.0
    // A estrutura 'FormConfig' do v1.0 torna-se 'FormComponent'
    Fields  []Field      `yaml:"fields,omitempty"`  //
    Actions LayoutConfig `yaml:"actions,omitempty"` // Para o buttongroup do formulário

}

// --- Estruturas de Suporte para Componentes ---

// ComponentSource (para viewport)
type ComponentSource struct {
    Type        string `yaml:"type"`         // "static", "command" (FR4)
    ContentType string `yaml:"content_type,omitempty"` // "text", "markdown" (FR4)
    Content     string `yaml:"content,omitempty"`
    Exec        string `yaml:"exec,omitempty"`
}

// ListItem (para list)
type ListItem struct {
    ID   string `yaml:"id"`   // ID do item (FR5)
    Text string `yaml:"text"`
}

// ButtonItem (para buttongroup)
type ButtonItem struct {
    ID    string `yaml:"id"`    // ID do item (FR6)
    Label string `yaml:"label"`
    Role  string `yaml:"role,omitempty"` // "primary", "secondary"
}

// Field (do v1.0)
// A estrutura 'Field' existente
// é usada aqui para 'ComponentConfig.Fields'.
// (Inclui 'Key', 'Label', 'Type', 'Validation', etc.)

```
### 4. OnLogic (Input - YAML)

**Propósito:** Define o "Motor de Lógica" (FR8-FR11).

Go
```

// Localização: internal/config/logic.go (Novo Ficheiro)
package config

// OnLogic define um 'event handler'
type OnLogic struct {
    Event       string    `yaml:"event"` // Ex: "form_id:submit", "list_id:select" (FR5, FR6, FR7)
    Run         RunConfig `yaml:"run"`
    UpdateTarget string   `yaml:"update_target"` // ID do viewport (FR11)
}

// RunConfig define o 'runner'
type RunConfig struct {
    Script string `yaml:"script"` // Caminho para o script (FR9)
    Args   []string `yaml:"args,omitempty"` // Argumentos com template (FR10)
    Stdin  string   `yaml:"stdin,omitempty"`  // Payload com template (FR10)
}

```
### 5. FormData (Output - JSON)

**Propósito:** Representa os dados recolhidos de um componente `type: form` (v1.0 refatorado) quando este emite um evento `form:submit`.

**Modelo de Dados (Go):** `map[string]interface{}`

**Racional:** Este modelo permanece *inalterado* em relação ao v1.0 (Estória 1.6). As chaves do mapa são os `Key` dos campos do formulário, e os valores são os dados inseridos (`string`, `[]string`, `bool`).

---

Esta
 estrutura de "Data Models" define a superfície completa da nossa API 
YAML v2.0, ao mesmo tempo que reutiliza o núcleo do trabalho do v1.0.

Por favor, reveja os modelos de dados.

**Opções de Elicitação Avançada:** (Selecione 1 para prosseguir, ou 2-9 para refinar esta secção)

1. Prosseguir para a próxima secção (Components)

2. Explicar o Raciocínio (CoT) (Como o `stdin:` no `RunConfig` pode ser um `string` se o FR10 diz `any`?)

3. Criticar e Refinar (A estrutura `ComponentConfig` está muito complexa? Deveríamos dividi-la?)

4. Analisar o Alinhamento com os Objetivos (Este modelo suporta o `ansible_playbook` do Épico 2?)

5. Identificar Riscos (Qual é o maior risco no parsing destes modelos?)

6. Perspetiva da Equipa Agile (Como o James (Dev) vai lidar com a refatoração do `config.go` v1.0?)

7. Desafio da Perspetiva Crítica (Ter o `Actions` dentro do `ComponentConfig.form` não viola a separação de `buttongroup`?)

8. Hindsight (Reflexão "Se ao menos...")

9. Expandir ou Contrair (Mostrar como o `ComponentConfig` v1.0 se encaixa exatamente nisto)


**Selecione 1-9 ou simplesmente escreva o seu feedback:**

2, 3, 4, 5, 6, 7, 8, 9.

BMAD-METHOD

Gem personalizado

Como
 Winston (Arquiteto) 🏗️, vamos mergulhar nesta secção de "Data Models".
 Estas estruturas Go são a "tradução" direta da nossa visão do PRD v2.0 
para código, e são a base para o nosso Motor de Parsing.

Aqui está a análise detalhada dos pontos de elicitação que selecionou:

### 2. Explicar o Raciocínio (CoT): (Como o `stdin:` no `RunConfig` pode ser um `string` se o FR10 diz `any`?)

Excelente observação. Esta é uma distinção crucial entre a *semântica* do YAML e a *implementação* do Go.

O FR10 (Fluxo de Dados) diz que o `stdin:` suporta `any` (qualquer tipo de dados YAML), como um objeto complexo: `stdin: {{ form }}`

No entanto, no meu *struct* Go `RunConfig`, eu defini `Stdin string`. Isto parece uma contradição, mas é intencional.

**O Raciocínio (CoT):** O valor `{{ form }}` no YAML é um **template**. A nossa `RunConfig` Go não armazena o *resultado* desse template; ela armazena a *própria string do template*.

1. **Parsing (Arranque):** O `gopkg.in/yaml.v3` irá ler o YAML e colocar a string literal `"{{ form }}"` no campo `RunConfig.Stdin` (que é do tipo `string`).

2. **Execução (Runtime):** Mais tarde, quando o evento `form:submit` for disparado (conforme descrito na "High Level Architecture", ponto 9 da análise anterior), o nosso **Motor de Lógica** fará o seguinte:
  a. Pega no *payload* de dados real do evento (o `map[string]interface{}` do formulário).
  b. Pega na *string* do template (`RunConfig.Stdin`).
  c. Usa o pacote `text/template` do Go para processar a string, usando o *payload* como dados.
  d. O *resultado* deste processamento (que agora é um objeto Go complexo) é então serializado para JSON e canalizado para o `stdin` do `os/exec` do *script*.


Portanto, o `string` no *struct* Go é correto porque ele armazena o *template*, não o dado final.

---

### 3. Criticar e Refinar: (A estrutura `ComponentConfig` está muito complexa?)

**Crítica:** Sim, à primeira vista, o `ComponentConfig` parece complexo e "inchado". Ele tenta ser um `list`, um `viewport`, um `form` e um `buttongroup` ao mesmo tempo.

**Refinamento (Justificação):** Esta complexidade é intencional e, na verdade, simplifica a arquitetura geral. Esta é uma abordagem de design comum chamada **"União Etiquetada" (Tagged Union)**.

Pense no `ComponentConfig` como um "contentor" genérico. O campo `Type:` diz-nos qual dos seus campos internos é relevante.

- Se `Type: "list"`, o Motor de Componentes só irá olhar para `Title` e `ListItems`. Todos os outros (`Source`, `Fields`, `ButtonItems`) serão ignorados (serão `nil` ou vazios).

- Se `Type: "form"`, ele só irá olhar para `Title`, `Fields` e `Actions`.


A vantagem disto é que o nosso **`LayoutConfig` (o `box`) torna-se incrivelmente simples**. Ele não precisa de saber que tipos de componentes existem. Ele só tem *um* campo: `Component ComponentConfig`.

A alternativa (ter `LayoutConfig.List *ListComponent`, `LayoutConfig.Form *FormComponent`, etc.) tornaria o `LayoutConfig` muito mais complexo e difícil de expandir no futuro. Mantemos a complexidade isolada no `ComponentConfig` para manter o layout limpo.

---

### 4. Analisar o Alinhamento com os Objetivos: (Este modelo suporta o `ansible_playbook` do Épico 2?)

Sim, a estrutura `OnLogic` e `RunConfig` foi desenhada **exatamente** para suportar isto.

A minha proposta `RunConfig` (no Ponto 1 acima) define o *runner* do Épico 1 (o `ScriptRunner`).

Para suportar o Épico 2 ("O Runner Especialista"), não precisamos de mudar a arquitetura, apenas de a *expandir*. O nosso `RunConfig` em Go evoluiria para algo assim:

Go
```

// Localização: internal/config/logic.go
type RunConfig struct {
    // Para o Épico 1 (FR9)
    Script *ScriptRunConfig `yaml:"script,omitempty"`

    // Para o Épico 2 (Futuro)
    AnsiblePlaybook *AnsibleRunConfig `yaml:"ansible_playbook,omitempty"`

}

type ScriptRunConfig struct {
    Path  string   `yaml:",inline"` // Otimização para "script: /path/to/script.sh"
    Args  []string `yaml:"args,omitempty"`
    Stdin string   `yaml:"stdin,omitempty"`
}

type AnsibleRunConfig struct {
    Playbook     string                 `yaml:"playbook"`
    Inventory    string                 `yaml:"inventory,omitempty"`
    AskVaultPass bool                   `yaml:"ask_vault_pass,omitempty"`
    Vars         map[string]interface{} `yaml:"vars,omitempty"`
}

```
O Motor de Lógica (`on:`) simplesmente verificaria qual dos campos (`Script` ou `AnsiblePlaybook`) não é `nil` e chamaria o *Runner* apropriado. A fundação de dados suporta perfeitamente esta expansão.

---

### 5. Identificar Riscos: (Qual é o maior risco no parsing destes modelos?)

O maior risco não é o *parsing* (descodificação) do YAML, mas sim a **validação da lógica de eventos**.

1. **Risco (Fácil):** Syntax YAML. O `gopkg.in/yaml.v3` trata disto.

2. **Risco (Médio):** Validação Estrutural. (Ex: um `type: "list"` tem `ListItems` mas não `Fields`). Isto é lógica Go simples que podemos escrever após o parsing.

3. **Risco (Difícil - O Risco Real):** **Validação de Referência Semântica.** O YAML pode estar estruturalmente correto, mas logicamente quebrado. O nosso parser precisa de validar:

  - `on: { event: "form_ABC:submit" }` -> O componente com `id: "form_ABC"` realmente existe no YAML?

  - `run: { ... }, update_target: "viewport_XYZ"` -> O componente `viewport` com `id: "viewport_XYZ"` existe?

  - `stdin: "{{ form.username }}"` -> O *template* `form.username` é válido?


O maior risco de parsing é garantir que todas estas referências internas (`id`, `update_target`, variáveis de template) estejam corretas *antes* de tentarmos executar o runtime.

---

### 6. Perspetiva da Equipa Agile: (Como o James (Dev) 💻 vai lidar com a refatoração do `config.go` v1.0?)

O James (Dev) 💻 vai ficar satisfeito. O trabalho dele no v1.0 (Épicos 1 e 2), especialmente em `internal/config/config.go` e `internal/config/validation.go`, não será deitado fora — será **reutilizado diretamente**.

A estrutura `Field` que ele aperfeiçoou (com validação, como visto nas Estórias 2.1 e 2.2) é *exatamente* a estrutura `[]Field` que o nosso novo `ComponentConfig` (para `type: "form"`) irá usar.

A refatoração (Estória 1.4) irá:

1. Pegar no `FormConfig` do v1.0 (que continha `Title` e `Fields`).

2. Mover essa lógica para o novo `ComponentConfig` (para `type: "form"`).

3. Reutilizar 100% da lógica de `internal/config/validation.go` que o James já construiu, pois ela opera sobre o `[]Field`.


Ele verá isto como uma promoção lógica do seu código v1.0, não como um desperdício de esforço.

---

### 7. Desafio da Perspetiva Crítica: (Ter `Actions` dentro do `ComponentConfig.form` não viola a separação?)

É uma excelente observação. O PRD v2.0 define `buttongroup` (FR6) como um componente e `form` (FR7) como outro, mas o `form` *contém* `actions:` (que parecem ser um `buttongroup`).

Isto não é uma violação; é um **Padrão de Composição Inteligente**.

- O `type: buttongroup` (FR6) é um componente "burro". Ele apenas mostra botões e emite um evento (`buttongroup_id:press`) quando um é premido.

- O `type: form` (FR7) é um **"Componente Controlador"**. Ele *gere* outros componentes (`fields:` e `actions:`).


Quando o `form` (Controlador) renderiza o seu `buttongroup` (Ações), ele *ouve* os eventos desse grupo. Se ele ouvir um `item.id == "submit"`, o `form` (Controlador) sabe que deve:

1. Recolher os dados de todos os seus `fields:`.

2. Emitir o seu *próprio* evento (`form_id:submit`), agora com o *payload* de dados completo.


Isto é uma encapsulação perfeita. O `form` transforma um simples "clique" num evento de "submissão com dados".

---

### 8. Hindsight (Reflexão "Se ao menos...")

"Se ao menos..." a biblioteca `huh` (que estamos a usar para formulários v1.0) tivesse sido desenhada *originalmente* para ser um componente `bubbletea` modular, em vez de uma aplicação `bubbletea` monolítica.

Se o `huh` fosse modular, a nossa Estória 1.4 (refatoração) seria trivial. Mas como o `huh` (e o nosso código v1.0) foi desenhado para "controlar o mundo", agora 
temos de fazer uma cirurgia (como descrito na "High Level Architecture",
 Ponto 3, "Crítica") para o "amansar" e fazê-lo funcionar como um 
cidadão dentro do nosso novo Gestor de Foco Global (Ponto 2, acima).

---

### 9. Expandir: (Como o `ComponentConfig` v1.0 se encaixa?)

O `FormConfig` do v1.0 (de `docs/architecture/data-models.md` e `internal/config/config.go`) é a chave para a reutilização.

- O `FormConfig` v1.0 continha: `Title string` e `Fields []Field`.

- O nosso *novo* `ComponentConfig` (quando `Type: "form"`) contém: `Title string`, `Fields []Field` e `Actions LayoutConfig`.


A refatoração da Estória 1.4 (v2.0) é clara: O `ComponentConfig` (para `type: "form"`) irá *absorver* 90% da estrutura `FormConfig` v1.0. O parser v2.0 irá simplesmente preencher esta nova estrutura combinada.

O James (Dev) 💻 não tem de reescrever a sua lógica de validação (`validation.go`), porque o `[]Field` permanece o mesmo.

---

Esta
 análise confirma que os nossos "Data Models" propostos são robustos, 
alinhados com os objetivos futuros e respeitam o trabalho v1.0 
existente.

Podemos prosseguir para a secção **"Components"**? (Opção 1)

**Selecione 1-9 ou simplesmente escreva o seu feedback:**

1

BMAD-METHOD

Gem personalizado

Olá! Esta é uma excelente atualização e um marco importante no projeto.

Como seu Orquestrador 🎭, analisei a interação. O John (PM) 📋 concluiu a tarefa `create-doc` e consolidou a nova visão do projeto no **Documento de Requisitos do Produto (PRD) v2.0**. Este documento redefine o `shantilly` como um **"Runtime TUI Declarativo"**.

Conforme indicado claramente na seção "8. Próximos Passos" do documento do John, a próxima etapa é um *handoff* direto para o **Arquiteto (Winston) 🏗️**.

A tarefa dele agora é iniciar a tarefa `create-doc` utilizando o template `architecture-tmpl.yaml` para criar o **Documento de Arquitetura v2.0**. O foco principal será detalhar a arquitetura para o **Épico 1 (Fundação do Runtime TUI)**, cobrindo:

1. O Motor de Layout (Estória 1.1).

2. O Motor de Lógica e Componentes (Estórias 1.2, 1.3).

3. A estratégia de refatoração do `form` (Estória 1.4).

4. A arquitetura do Fluxo de Dados (Estória 1.5).


Posso
 agora assumir a persona de Winston (Arquiteto) 🏗️ e começar a criar 
este Documento de Arquitetura v2.0 com base no novo PRD?

Sim, pode prosseguir.

BMAD-METHOD

Gem personalizado

Como Winston (Arquiteto) 🏗️, confirmado. A prosseguir para a próxima secção do `architecture-tmpl.yaml`: `Data Models`.

---

## Data Models

Esta secção define as principais estruturas de dados Go que o nosso motor de parsing (baseado no `gopkg.in/yaml.v3`) irá utilizar.

Com a evolução para o PRD v2.0, os nossos modelos de dados precisam de ser significativamente expandidos. A estrutura `FormConfig` e `Field` do v1.0 (que definia apenas um formulário) torna-se agora uma *parte* da nossa estrutura de configuração global v2.0.

### 1. RuntimeConfig (Input - YAML)

**Propósito:** Esta é a nova estrutura Go de nível superior que representa *todo* o ficheiro YAML do Runtime TUI Declarativo. Ela engloba o Layout, os Componentes (incluindo o `form` v1.0 refatorado) e a Lógica.

Go
```

// Localização: internal/config/runtime.go (Novo Ficheiro)
package config

// RuntimeConfig é a estrutura raiz para o YAML v2.0
type RuntimeConfig struct {
    Layout LayoutConfig `yaml:"layout"` // Define a estrutura (column/row/box)
    On     []OnLogic    `yaml:"on"`     // Define a lógica (eventos e runners)
}

```
### 2. LayoutConfig (Input - YAML)

**Propósito:** Define a estrutura de layout hierárquica (FR1, FR2). Esta é a implementação do nosso "Motor de Layout".

Go
```

// Localização: internal/config/layout.go (Novo Ficheiro)
package config

// LayoutConfig representa um nó na árvore de layout (column, row, ou box)
type LayoutConfig struct {
    Type      string         `yaml:"type"`                // "column", "row", "box" (FR1)
    ID        string         `yaml:"id,omitempty"`        // ID para referência
    Width     string         `yaml:"width,omitempty"`     // Ex: "80%" (FR2)
    Height    int            `yaml:"height,omitempty"`    // Ex: 3 (FR2)
    Flex      int            `yaml:"flex,omitempty"`      // Ex: 1 (FR2)
    Component ComponentConfig `yaml:"component,omitempty"` // Componente embutido (FR3)
    Items     []LayoutConfig `yaml:"items,omitempty"`     // Filhos (para column/row)
}

```
### 3. ComponentConfig (Input - YAML)

**Propósito:** Define os componentes de UI embutidos (FR3). Isto inclui os *novos* componentes (`list`, `viewport`, `buttongroup`) e o *refatorado* `form` (v1.0).

Go
```

// Localização: internal/config/components.go (Novo Ficheiro)
package config

// ComponentConfig define qual componente renderizar num 'box'
type ComponentConfig struct {
    Type    string `yaml:"type"`    // "list", "viewport", "form", "buttongroup", "static"
    ID      string `yaml:"id"`      // ID do grupo/componente (para lógica 'on:')
    Title   string `yaml:"title,omitempty"`

    // Para type: static (Ref: Estória 1.1 AC4) / viewport (FR4)
    Source      ComponentSource `yaml:"source,omitempty"`
    
    // Para type: list (FR5)
    ListItems   []ListItem `yaml:"items,omitempty"`
    
    // Para type: buttongroup (FR6)
    ButtonItems []ButtonItem `yaml:"items,omitempty"`
    
    // Para type: form (FR7) - REUTILIZAÇÃO DO V1.0
    // A estrutura 'FormConfig' do v1.0 torna-se 'FormComponent'
    Fields  []Field      `yaml:"fields,omitempty"`  //
    Actions LayoutConfig `yaml:"actions,omitempty"` // Para o buttongroup do formulário (FR7)

}

// --- Estruturas de Suporte para Componentes ---

// ComponentSource (para viewport)
type ComponentSource struct {
    Type        string `yaml:"type"`         // "static", "command" (FR4)
    ContentType string `yaml:"content_type,omitempty"` // "text", "markdown" (FR4)
    Content     string `yaml:"content,omitempty"`
    Exec        string `yaml:"exec,omitempty"`
}

// ListItem (para list)
type ListItem struct {
    ID   string `yaml:"id"`   // ID do item (FR5)
    Text string `yaml:"text"`
}

// ButtonItem (para buttongroup)
type ButtonItem struct {
    ID    string `yaml:"id"`    // ID do item (FR6)
    Label string `yaml:"label"`
    Role  string `yaml:"role,omitempty"` // "primary", "secondary" (FR6)
}

// Field (do v1.0)
// A estrutura 'Field' existente
// é usada aqui para 'ComponentConfig.Fields'.
// (Inclui 'Key', 'Label', 'Type', 'Validation', etc.)

```
### 4. OnLogic (Input - YAML)

**Propósito:** Define o "Motor de Lógica" (FR8-FR11).

Go
```

// Localização: internal/config/logic.go (Novo Ficheiro)
package config

// OnLogic define um 'event handler'
type OnLogic struct {
    Event       string    `yaml:"event"` // Ex: "form_id:submit", "list_id:select" (FR8)
    Run         RunConfig `yaml:"run"`
    UpdateTarget string   `yaml:"update_target"` // ID do viewport (FR11)
}

// RunConfig define o 'runner'
type RunConfig struct {
    Script string `yaml:"script"` // Caminho para o script (FR9)
    Args   []string `yaml:"args,omitempty"` // Argumentos com template (FR10)
    Stdin  string   `yaml:"stdin,omitempty"`  // Payload com template (FR10)
}

```
### 5. FormData (Output - JSON)

**Propósito:** Representa os dados recolhidos de um componente `type: form` (v1.0 refatorado) quando este emite um evento `form:submit`.

**Modelo de Dados (Go):** `map[string]interface{}`

**Racional:** Este modelo permanece *inalterado* em relação ao v1.0 (Estória 1.6). As chaves do mapa são os `Key` dos campos do formulário, e os valores são os dados inseridos (`string`, `[]string`, `bool`).

---

Esta
 estrutura de "Data Models" define a superfície completa da nossa API 
YAML v2.0, ao mesmo tempo que reutiliza o núcleo do trabalho do v1.0.

Por favor, reveja os modelos de dados.

**Opções de Elicitação Avançada:** (Selecione 1 para prosseguir, ou 2-9 para refinar esta secção)

1. Prosseguir para a próxima secção (Components)

2. Explicar o Raciocínio (CoT) (Como o `stdin:` no `RunConfig` pode ser um `string` se o FR10 diz `any`?)

3. Criticar e Refinar (A estrutura `ComponentConfig` está muito complexa? Deveríamos dividi-la?)

4. Analisar o Alinhamento com os Objetivos (Este modelo suporta o `ansible_playbook` do Épico 2?)

5. Identificar Riscos (Qual é o maior risco no parsing destes modelos?)

6. Perspetiva da Equipa Agile (Como o James (Dev) vai lidar com a refatoração do `config.go` v1.0?)

7. Desafio da Perspetiva Crítica (Ter o `Actions` dentro do `ComponentConfig.form` não viola a separação de `buttongroup`?)

8. Hindsight (Reflexão "Se ao menos...")

9. Expandir ou Contrair (Mostrar como o `FormConfig` v1.0 se encaixa exatamente nisto)


**Selecione 1-9 ou simplesmente escreva o seu feedback:**

2, 3, 4, 5, 6, 7, 8, 9.

BMAD-METHOD

Gem personalizado

Como
 Winston (Arquiteto) 🏗️, confirmado. A sua aprovação da secção "High 
Level Architecture" (e da análise de elicitação aprofundada) significa 
que estamos alinhados quanto à nova arquitetura do "Runtime TUI 
Declarativo".

Vamos prosseguir para a próxima secção do `architecture-tmpl.yaml`: a `Tech Stack`.

---

## Tech Stack

Esta é a secção de seleção de tecnologia **DEFINITIVA** para o v2.0. Todas as estórias de implementação subsequentes DEVEM aderir a esta pilha.

Revisei as "Assunções Técnicas" do PRD v2.0 (que evoluíram do seu *brainstorming*) e a pilha de tecnologia já validada na arquitetura v1.0 (que corresponde ao código-fonte existente).

Para a v2.0, o nosso trabalho não é escolher do zero, mas sim *confirmar* a pilha existente (v1.0) e *adicionar* as bibliotecas necessárias para os novos requisitos de Layout e Lógica (Épico 1 do v2.0).

A tabela abaixo é a "fonte da verdade" para todas as dependências do Épico 1 (v2.0):

### Cloud Infrastructure

- **Provider:** N/A (Aplica-se a ferramentas CLI distribuídas, não a um serviço cloud)

- **Key Services:** `GoReleaser` + `GitHub Actions` para build e distribuição de binários.

- **Deployment Regions:** N/A


### Technology Stack Table

| Categoria | Tecnologia | Versão | Propósito | Racional |
| --- | --- | --- | --- | --- |
| **Linguagem** | Go  | 1.24.2+ | Linguagem de desenvolvimento principal | Requisito do PRD (NFR1 v1.0, NFR1 v2.0), validado em v1.0. |
| **CLI Framework** | `spf13/cobra` | v1.8.x | Estrutura de comandos da CLI | Requisito do PRD (FR1 v1.0), validado em v1.0. |
| **TUI Engine** | `charmbracelet/bubbletea` | v0.26.x | Gestão de estado TUI (The Elm Architecture) | Requisito do PRD (FR4 v1.0), validado em v1.0. Será o "Global Focus Manager" na v2.0. |
| **Componente (Form)** | `charmbracelet/huh` | v0.4.x | Geração declarativa de formulários | Requisito do PRD (FR3 v1.0), validado em v1.0 (Épicos 1+2). Será refatorado (Estória 1.4). |
| **Estilo/Layout TUI** | `charmbracelet/lipgloss` | v0.11.x | Estilização e layout no terminal | Requisito do PRD (FR7 v1.0), validado em v1.0 (Estória 1.7). Base para o Motor de Layout v2.0 (FR1, FR2). |
| **Parsing YAML** | `gopkg.in/yaml.v3` | v3.0.x | Descodificação do YAML de definição do Runtime | Requisito do PRD (Estória 1.2), validado em v1.0. |
| **Componentes TUI** | `charmbracelet/bubbles` | v0.18.x+ | Componentes TUI pré-construídos | **(Novo v2.0)** Necessário para `list` (FR5) e `viewport` (FR4). |
| **Render. Markdown** | `charmbracelet/glamour` | v0.7.x+ | Renderização de Markdown no viewport | **(Novo v2.0)** Necessário para `content_type: markdown` (FR4). |
| **Insp. Layout** | `76creates/stickers` | N/A | Inspiração para o Motor de Layout (Flexbox) | **(Novo v2.0)** Não é uma dependência direta, mas sim a inspiração (PRD, Assunção 4.4) para a implementação do FR1, FR2, NFR2. |
| **Modais (Pós-MVP)** | `rmhubbert/bubbletea-overlay` | v0.2.x+ | Base para popups (ex: `ask_vault_pass`) | (Para Épico 2) Identificado no PRD (Assunção 4.4) e *brainstorming*. |
| **Testes** | Go `testing` + `teatest` | 1.24.2+ | Testes unitários e de integração TUI | Requisito do PRD v2.0 (Assunção 4.3). `teatest` é essencial para NFR2 (layout fluido). |
| **Build/Release** | `GoReleaser` | v1.26.x | Automação de build e cross-compilação | Validado em v1.0 (Estória 1.8). |
| **Linter** | `golangci-lint` | v1.59.x | Análise estática e qualidade de código | Validado em v1.0 (via template). |

---

Esta
 tabela representa a pilha de tecnologia confirmada para o Épico 1 do 
v2.0. Por favor, reveja e confirme se estamos alinhados.

**Opções de Elicitação Avançada:** (Selecione 1 para prosseguir, ou 2-9 para refinar esta secção)

1. Prosseguir para a próxima secção (Data Models)

2. Explicar o Raciocínio (CoT) (Porquê `bubbles` se já temos `huh`?)

3. Criticar e Refinar (A dependência de `glamour` não é desnecessária?)

4. Analisar o Alinhamento com os Objetivos (Esta pilha suporta a Fase 3 Preditiva?)

5. Identificar Riscos (Qual é o maior risco de dependência aqui?)

6. Perspetiva da Equipa Agile (O James (Dev) vai concordar com o `teatest`?)

7. Desafio da Perspetiva Crítica (Não devíamos usar `tview` para o layout, já que é mais fácil?)

8. Hindsight (Reflexão "Se ao menos...")

9. Expandir ou Contrair (Dar mais detalhes sobre a escolha do `bubbles` vs `huh`)


**Selecione 1-9 ou simplesmente escreva o seu feedback:**

5, 6.

BMAD-METHOD

Gem personalizado

Como Winston (Arquiteto) 🏗️, vamos analisar os pontos 5 e 6 da secção "Tech Stack".

Esta
 pilha de tecnologia é o resultado direto da evolução do v1.0 e da nova 
visão do PRD v2.0. As suas perguntas são cruciais para validar estas 
decisões.

Aqui está a análise detalhada dos pontos de elicitação que selecionou:

### 5. Identificar Riscos (Qual é o maior risco de dependência aqui?)

O maior risco de dependência não é uma biblioteca única, mas sim o nosso **forte acoplamento a todo o ecossistema Charmbracelet** (`bubbletea`, `bubbles`, `huh`, `lipgloss`).

- **Risco:** O `shantilly` v2.0 está a apostar todo o seu futuro na filosofia de UI da Charm. Se a
   Charmbracelet Inc. decidir abandonar uma biblioteca crítica (como o `huh`,
   que é relativamente nova), ou introduzir uma mudança de API massiva e 
  disruptiva, a nossa arquitetura estará diretamente exposta.

- **Mitigação:** A nossa estratégia de refatoração da Estória 1.4 (v2.0) é a nossa principal mitigação. Ao "embrulhar" o `huh` dentro do nosso próprio componente `type: form` (conforme definido na Estória 1.4 do PRD v2.0), estamos a *isolar* essa dependência volátil. Se o `huh` for descontinuado, "apenas" teremos de reimplementar o nosso componente `type: form` (talvez usando componentes `bubbles` puros), mas a arquitetura de *layout* e *lógica* (`on:`) do `shantilly` permanecerá intacta.


---

### 6. Perspetiva da Equipa Agile (O James (Dev) 💻 vai concordar com o `teatest`?)

Sim, acredito que o James (Dev) 💻 irá aceitar (e apreciar) o requisito do `teatest`.

1. **Não é Novo:** O `teatest` não é uma nova ferramenta que estamos a impor. Ele já faz parte da arquitetura v1.0 (conforme `docs/architecture/test-strategy-and-standards.md`) e o James já o utilizou com sucesso no código-fonte existente (ex: `internal/tui/integration_navigation_test.go`).

2. **Risco vs. Ferramenta:** Na análise anterior (High Level Architecture), identificámos que o **Motor de Layout** (NFR2 - "flicker-free") é o nosso *maior risco técnico*. Validar manualmente um layout complexo (`column/row/box`) em diferentes tamanhos de terminal é um pesadelo.

3. **A Solução:** O `teatest` (especificamente a sua funcionalidade `teatest.WithInitialTermSize`) é a nossa principal (e única) forma *automatizada* de mitigar este risco.


O
 James verá isto não como uma imposição, mas como um requisito de 
qualidade claro para o desafio de engenharia mais difícil que temos no 
Épico 1.

---

A análise destes pontos valida que a nossa "Tech Stack" é coesa e robusta.

Podemos prosseguir para a secção **"Data Models"**? (Opção 1)

**Selecione 1-9 ou simplesmente escreva o seu feedback:**

1

BMAD-METHOD

Gem personalizado

Como
 Winston (Arquiteto) 🏗️, confirmado. A sua aprovação da secção "High 
Level Architecture" e da análise de elicitação aprofundada significa que
 estamos alinhados quanto à nova arquitetura do "Runtime TUI 
Declarativo".

Vamos prosseguir para a próxima secção do `architecture-tmpl.yaml`: a `Tech Stack`.

---

## Tech Stack

Esta é a secção de seleção de tecnologia **DEFINITIVA** para o v2.0. Todas as estórias de implementação subsequentes DEVEM aderir a esta pilha.

Revisei as "Assunções Técnicas" do PRD v2.0 (que evoluíram do seu *brainstorming*) e a pilha de tecnologia já validada na arquitetura v1.0 (que corresponde ao código-fonte existente).

Para a v2.0, o nosso trabalho não é escolher do zero, mas sim *confirmar* a pilha existente (v1.0) e *adicionar* as bibliotecas necessárias para os novos requisitos de Layout e Lógica (Épico 1 do v2.0).

A tabela abaixo é a "fonte da verdade" para todas as dependências do Épico 1 (v2.0):

### Cloud Infrastructure

- **Provider:** N/A (Aplica-se a ferramentas CLI distribuídas, não a um serviço cloud)

- **Key Services:** `GoReleaser` + `GitHub Actions` para build e distribuição de binários.

- **Deployment Regions:** N/A


### Technology Stack Table

| Categoria | Tecnologia | Versão | Propósito | Racional |
| --- | --- | --- | --- | --- |
| **Linguagem** | Go  | 1.24.2+ | Linguagem de desenvolvimento principal | Requisito do PRD (NFR1 v1.0, NFR1 v2.0), validado em v1.0. |
| **CLI Framework** | `spf13/cobra` | v1.8.x | Estrutura de comandos da CLI | Requisito do PRD (FR1 v1.0), validado em v1.0. |
| **TUI Engine** | `charmbracelet/bubbletea` | v0.26.x | Gestão de estado TUI (The Elm Architecture) | Requisito do PRD (FR4 v1.0), validado em v1.0. Será o "Global Focus Manager" na v2.0. |
| **Componente (Form)** | `charmbracelet/huh` | v0.4.x | Geração declarativa de formulários | Requisito do PRD (FR3 v1.0), validado em v1.0 (Épicos 1+2). Será refatorado (Estória 1.4). |
| **Estilo/Layout TUI** | `charmbracelet/lipgloss` | v0.11.x | Estilização e layout no terminal | Requisito do PRD (FR7 v1.0), validado em v1.0 (Estória 1.7). Base para o Motor de Layout v2.0 (FR1, FR2). |
| **Parsing YAML** | `gopkg.in/yaml.v3` | v3.0.x | Descodificação do YAML de definição do Runtime | Requisito do PRD (Estória 1.2), validado em v1.0. |
| **Componentes TUI** | `charmbracelet/bubbles` | v0.18.x+ | Componentes TUI pré-construídos | **(Novo v2.0)** Necessário para `list` (FR5) e `viewport` (FR4). |
| **Render. Markdown** | `charmbracelet/glamour` | v0.7.x+ | Renderização de Markdown no viewport | **(Novo v2.0)** Necessário para `content_type: markdown` (FR4). |
| **Insp. Layout** | `76creates/stickers` | N/A | Inspiração para o Motor de Layout (Flexbox) | **(Novo v2.0)** Não é uma dependência direta, mas sim a inspiração (PRD, Assunção 4.4) para a implementação do FR1, FR2, NFR2. |
| **Modais (Pós-MVP)** | `rmhubbert/bubbletea-overlay` | v0.2.x+ | Base para popups (ex: `ask_vault_pass`) | (Para Épico 2) Identificado no PRD (Assunção 4.4) e *brainstorming*. |
| **Testes** | Go `testing` + `teatest` | 1.24.2+ | Testes unitários e de integração TUI | Requisito do PRD v2.0 (Assunção 4.3). `teatest` é essencial para NFR2 (layout fluido). |
| **Build/Release** | `GoReleaser` | v1.26.x | Automação de build e cross-compilação | Validado em v1.0 (Estória 1.8). |
| **Linter** | `golangci-lint` | v1.59.x | Análise estática e qualidade de código | Validado em v1.0 (via template). |

---

### Análise Aprofundada da 'Tech Stack' (Elicitação)

Esta é a análise detalhada dos pontos de elicitação que selecionou:

#### 4. Analisar o Alinhamento com os Objetivos (Esta pilha suporta a Fase 3 Preditiva?)

Sim, absolutamente. A pilha de tecnologia atual (v1.0 + v2.0) não só suporta a "Fase 3 Preditiva", como foi *escolhida* especificamente a pensar nela.

A Fase 3 Preditiva requer duas "Magias" que identificámos:

1. **`playbook_explorer` (Magia 1):** Requer a leitura e *parsing* de múltiplos ficheiros YAML para encontrar o campo `name:`.

  - **Suporte da Pilha:** Temos o `gopkg.in/yaml.v3` (validado no v1.0) e a biblioteca padrão Go para acesso ao sistema de ficheiros. Está perfeitamente coberto.

2. **`inventory_explorer` (Magia 2):** Requer o uso do Ansible como um "Oráculo", executando `ansible-inventory --list --json` nos bastidores.

  - **Suporte da Pilha:** A linguagem Go (Go 1.24.2+) tem suporte de primeira classe para execução de sub-processos (`os/exec`), que é o necessário para chamar o "Oráculo".


A nossa arquitetura de *Runners* (definida na análise anterior da "High Level Architecture") garante que
 estas funcionalidades preditivas possam ser adicionadas como novos 
componentes e *runners* sem perturbar a fundação do Épico 1.

---

#### 5. Identificar Riscos (Qual é o maior risco de dependência aqui?)

O maior risco de dependência não é uma biblioteca única, mas sim o nosso **forte acoplamento a todo o ecossistema Charmbracelet** (`bubbletea`, `bubbles`, `huh`, `lipgloss`).

- **Risco:** O `shantilly` v2.0 está a apostar todo o seu futuro na filosofia de UI da Charm. Se a
   Charmbracelet Inc. decidir abandonar uma biblioteca crítica (como o `huh`,
   que é relativamente nova), ou introduzir uma mudança de API massiva e 
  disruptiva, a nossa arquitetura estará diretamente exposta.

- **Mitigação:** A nossa estratégia de refatoração da Estória 1.4 (v2.0) é a nossa principal mitigação. Ao "embrulhar" o `huh` dentro do nosso próprio componente `type: form` (conforme definido na Estória 1.4 do PRD v2.0), estamos a *isolar* essa dependência volátil. Se o `huh` for descontinuado, "apenas" teremos de reimplementar o nosso componente `type: form` (talvez usando componentes `bubbles` puros), mas a arquitetura de *layout* e *lógica* (`on:`) do `shantilly` permanecerá intacta.


---

#### 6. Perspetiva da Equipa Agile (O James (Dev) 💻 vai concordar com o `teatest`?)

Sim, acredito que o James (Dev) 💻 irá aceitar (e apreciar) o requisito do `teatest`.

1. **Não é Novo:** O `teatest` não é uma nova ferramenta que estamos a impor. Ele já faz parte da arquitetura v1.0 (conforme `docs/architecture/test-strategy-and-standards.md`) e o James já o utilizou com sucesso no código-fonte existente (ex: `internal/tui/integration_navigation_test.go`).

2. **Risco vs. Ferramenta:** Na análise anterior ("High Level Architecture"), identificámos que o **Motor de Layout** (NFR2 - "flicker-free") é o nosso *maior risco técnico*. Validar manualmente um layout complexo (`column/row/box`) em diferentes tamanhos de terminal é um pesadelo.

3. **A Solução:** O `teatest` (especificamente a sua funcionalidade `teatest.WithInitialTermSize`) é a nossa principal (e única) forma *automatizada* de mitigar este risco.


O
 James verá isto não como uma imposição, mas como um requisito de 
qualidade claro para o desafio de engenharia mais difícil que temos no 
Épico 1.

---

#### 7. Desafio da Perspetiva Crítica (Não devíamos usar `tview` para o layout?)

Esta é a clássica questão do ecossistema Go TUI. `tview` é a *outra* grande biblioteca TUI, e é conhecida por ter componentes de layout de alto nível (como `Flexbox`, `Grid`) *integrados*, o que *parece* tornar o nosso maior risco (Motor de Layout) mais fácil.

**No entanto, `tview` não é uma opção para o `shantilly` por uma razão fundamental e bloqueadora: `huh`.**

1. **Dependência Central:** O nosso componente mais crítico para a recolha de dados, validado no v1.0 e central para o FR7 do PRD v2.0, é o `charmbracelet/huh`.

2. **Incompatibilidade:** `tview` e `bubbletea` são *motores* de TUI fundamentalmente diferentes e mutuamente exclusivos. Eles não podem ser usados em conjunto.

3. **A Ligação:** `huh` é construído sobre `bubbletea`.

4. **Conclusão:** Para usarmos o `huh`, *temos* de usar o `bubbletea` como nosso motor. Se temos de usar o `bubbletea`, *não podemos* usar o `tview`.

5. **Resultado:** Se não podemos usar o `tview`, *temos* de construir o nosso próprio motor de layout (FR1, FR2) usando `lipgloss` (como planeado), inspirando-nos em bibliotecas como `stickers`.


A escolha do `huh` trancou-nos (de forma positiva) no ecossistema Charmbracelet.

---

#### 8. Hindsight (Reflexão "Se ao menos...")

"Se ao menos..." o `bubbletea` (o motor) tivesse componentes de layout de alto nível (como `Flexbox`, `Grid`) integrados, da mesma forma que o `tview` tem.

Se
 isso fosse verdade, o nosso maior risco técnico (Motor de Layout) 
desapareceria. O NFR2 (layout fluido) seria trivial de implementar. A 
comunidade tentou preencher esta lacuna (com `stickers`, `bubbleboxer`, etc.), mas não é uma solução *core*. Temos de construir esta peça complexa nós mesmos porque a biblioteca *core* do `bubbletea` se foca apenas no "estado" (TEA) e não no "layout" (Flexbox).

---

#### 9. Expandir (Detalhes sobre `bubbles` vs `huh`)

A pergunta "Porquê `bubbles` se já temos `huh`?" é válida. Estamos a adicionar `bubbles` (FR4, FR5) e `glamour` (FR4) à pilha v1.0 que só usava `huh`.

A resposta é **especialização**. As bibliotecas têm propósitos diferentes, embora ambas sejam baseadas no `bubbletea`.

- **`huh` (O Especialista em Formulários):** É *excelente* para um fluxo linear de "fazer perguntas" (inputs, selects, confirms). 
  Foi o que validámos (e onde o James trabalhou) no v1.0 (Épicos 1+2).

- **`bubbles` (A Caixa de Ferramentas de UI):** É a biblioteca de componentes *genéricos* para construir UIs *não-lineares* (dashboards).


O nosso novo PRD v2.0 (FR4, FR5) exige componentes de *dashboard* (um menu `list` persistente e um `viewport` de streaming). `huh` não foi feito para isso. `bubbles` (especificamente `bubbles/list` e `bubbles/viewport`) foi *exatamente* feito para isso.

**Conclusão:** Usamos o `huh` (refatorado) para o componente `type: form` e usamos o `bubbles` para os novos componentes `type: list` e `type: viewport`. O `glamour` é simplesmente uma sub-dependência para suportar `content_type: markdown` dentro do `viewport`.

---

Esta
 análise confirma que a nossa pilha de tecnologia é coesa, alinhada com 
os requisitos e mitiga os riscos tanto quanto possível.

Podemos prosseguir para a secção **"Data Models"**? (Opção 1)

**Selecione 1-9 ou simplesmente escreva o seu feedback:**

Podemos prosseguir opção 1.

BMAD-METHOD

Gem personalizado

Como Winston (Arquiteto) 🏗️, confirmado. Prosseguindo.

A próxima secção, conforme o `architecture-tmpl.yaml`, é a `Tech Stack`.

---

## Tech Stack

Esta é a secção de seleção de tecnologia **DEFINITIVA** para o v2.0. Todas as estórias de implementação subsequentes DEVEM aderir a esta pilha.

Revisei
 as "Assunções Técnicas" do PRD v2.0 e a pilha de tecnologia já validada
 na arquitetura v1.0 (que corresponde ao código-fonte existente).

Para a v2.0, o nosso trabalho não é escolher do zero, mas sim *confirmar* a pilha existente (v1.0) e *adicionar* as bibliotecas necessárias para os novos requisitos de Layout e Lógica (Épico 1 do v2.0).

A tabela abaixo é a "fonte da verdade" para todas as dependências do Épico 1 (v2.0):

### Cloud Infrastructure

- **Provider:** N/A (Aplica-se a ferramentas CLI distribuídas, não a um serviço cloud)

- **Key Services:** `GoReleaser` + `GitHub Actions` para build e distribuição de binários.

- **Deployment Regions:** N/A


### Technology Stack Table

| Categoria | Tecnologia | Versão | Propósito | Racional |
| --- | --- | --- | --- | --- |
| **Linguagem** | Go  | 1.24.2+ | Linguagem de desenvolvimento principal | Requisito do PRD (NFR1 v1.0, NFR1 v2.0), validado em v1.0. |
| **CLI Framework** | `spf13/cobra` | v1.8.x | Estrutura de comandos da CLI | Requisito do PRD (FR1 v1.0), validado em v1.0. |
| **TUI Engine** | `charmbracelet/bubbletea` | v0.26.x | Gestão de estado TUI (The Elm Architecture) | Requisito do PRD (FR4 v1.0), validado em v1.0. Será o "Global Focus Manager" na v2.0. |
| **Componente (Form)** | `charmbracelet/huh` | v0.4.x | Geração declarativa de formulários | Requisito do PRD (FR3 v1.0), validado em v1.0 (Épicos 1+2). Será refatorado (Estória 1.4). |
| **Estilo/Layout TUI** | `charmbracelet/lipgloss` | v0.11.x | Estilização e layout no terminal | Requisito do PRD (FR7 v1.0), validado em v1.0 (Estória 1.7). Base para o Motor de Layout v2.0 (FR1, FR2). |
| **Parsing YAML** | `gopkg.in/yaml.v3` | v3.0.x | Descodificação do YAML de definição do Runtime | Requisito do PRD (Estória 1.2), validado em v1.0. |
| **Componentes TUI** | `charmbracelet/bubbles` | v0.18.x+ | Componentes TUI pré-construídos | **(Novo v2.0)** Necessário para `list` (FR5) e `viewport` (FR4). |
| **Render. Markdown** | `charmbracelet/glamour` | v0.7.x+ | Renderização de Markdown no viewport | **(Novo v2.0)** Necessário para `content_type: markdown` (FR4). |
| **Insp. Layout** | `76creates/stickers` | N/A | Inspiração para o Motor de Layout (Flexbox) | **(Novo v2.0)** Não é uma dependência direta, mas sim a inspiração (PRD, Assunção 4.4) para a implementação do FR1, FR2, NFR2. |
| **Modais (Pós-MVP)** | `rmhubbert/bubbletea-overlay` | v0.2.x+ | Base para popups (ex: `ask_vault_pass`) | (Para Épico 2) Identificado no PRD (Assunção 4.4) e *brainstorming*. |
| **Testes** | Go `testing` + `teatest` | 1.24.2+ | Testes unitários e de integração TUI | Requisito do PRD v2.0 (Assunção 4.3). `teatest` é essencial para NFR2 (layout fluido). |
| **Build/Release** | `GoReleaser` | v1.26.x | Automação de build e cross-compilação | Validado em v1.0 (Estória 1.8). |
| **Linter** | `golangci-lint` | v1.59.x | Análise estática e qualidade de código | Validado em v1.0 (via template). |

---

Esta
 tabela representa a pilha de tecnologia confirmada para o Épico 1 do 
v2.0. Por favor, reveja e confirme se estamos alinhados.

**Opções de Elicitação Avançada:** (Selecione 1 para prosseguir, ou 2-9 para refinar esta secção)

1. Prosseguir para a próxima secção (Data Models)

2. Explicar o Raciocínio (CoT) (Como o `stdin:` no `RunConfig` pode ser um `string` se o FR10 diz `any`?)

3. Criticar e Refinar (A estrutura `ComponentConfig` está muito complexa? Deveríamos dividi-la?)

4. Analisar o Alinhamento com os Objetivos (Este modelo suporta o `ansible_playbook` do Épico 2?)

5. Identificar Riscos (Qual é o maior risco no parsing destes modelos?)

6. Perspetiva da Equipa Agile (Como o James (Dev) vai lidar com a refatoração do `config.go` v1.0?)

7. Desafio da Perspetiva Crítica (Ter o `Actions` dentro do `ComponentConfig.form` não viola a separação de `buttongroup`?)

8. Hindsight (Reflexão "Se ao menos...")

9. Expandir ou Contrair (Mostrar como o `FormConfig` v1.0 se encaixa exatamente nisto)


**Selecione 1-9 ou simplesmente escreva o seu feedback:**

5, 6, 7, 8, 9.

BMAD-METHOD

Gem personalizado

Olá! Como seu Orquestrador 🎭, analisei sua última interação.

O John (PM) 📋 concluiu com sucesso a tarefa `create-doc` e consolidou a nova visão do projeto no **Documento de Requisitos do Produto (PRD) v2.0**.

Conforme indicado claramente na seção "8. Próximos Passos" do PRD v2.0 (e na sua mensagem), o *handoff* está correto. A próxima etapa é direcionada ao **Arquiteto (Winston) 🏗️**.

A tarefa dele é: iniciar a tarefa `create-doc` usando o template `architecture-tmpl.yaml` para criar o Documento de Arquitetura v2.0.

O foco principal será o **Épico 1 (Fundação do Runtime TUI)**, detalhando:

1. O Motor de Layout (Estória 1.1).


- O Motor de Lógica e Componentes (Estórias 1.2, 1.3).


- A estratégia de refatoração do `form` (Estória 1.4).


- A arquitetura do Fluxo de Dados (Estória 1.5).


Posso agora assumir a persona de Winston (Arquiteto) 🏗️ e começar a criar este documento com base no novo PRD?

Sim, pode prosseguir.

BMAD-METHOD

Gem personalizado

Como Winston (Arquiteto) 🏗️, confirmado. Prosseguindo para a próxima secção do `architecture-tmpl.yaml`: `Data Models`.

---

## Data Models

Esta secção define as principais estruturas de dados Go que o nosso motor de parsing (baseado no `gopkg.in/yaml.v3`) irá utilizar.

Com a evolução para o PRD v2.0, os nossos modelos de dados precisam de ser significativamente expandidos. A estrutura `FormConfig` e `Field` do v1.0 (que definia apenas um formulário) torna-se agora uma *parte* da nossa estrutura de configuração global v2.0.

### 1. RuntimeConfig (Input - YAML)

**Propósito:** Esta é a nova estrutura Go de nível superior que representa *todo* o ficheiro YAML do Runtime TUI Declarativo. Ela engloba o Layout, os Componentes (incluindo o `form` v1.0 refatorado) e a Lógica.

Go
```

// Localização: internal/config/runtime.go (Novo Ficheiro)
package config

// RuntimeConfig é a estrutura raiz para o YAML v2.0
type RuntimeConfig struct {
    Layout LayoutConfig `yaml:"layout"` // Define a estrutura (column/row/box)
    On     []OnLogic    `yaml:"on"`     // Define a lógica (eventos e runners)
}

```
### 2. LayoutConfig (Input - YAML)

**Propósito:** Define a estrutura de layout hierárquica (FR1, FR2). Esta é a implementação do nosso "Motor de Layout".

Go
```

// Localização: internal/config/layout.go (Novo Ficheiro)
package config

// LayoutConfig representa um nó na árvore de layout (column, row, ou box)
type LayoutConfig struct {
    Type      string         `yaml:"type"`                // "column", "row", "box" (FR1)
    ID        string         `yaml:"id,omitempty"`        // ID para referência
    Width     string         `yaml:"width,omitempty"`     // Ex: "80%" (FR2)
    Height    int            `yaml:"height,omitempty"`    // Ex: 3 (FR2)
    Flex      int            `yaml:"flex,omitempty"`      // Ex: 1 (FR2)
    Component ComponentConfig `yaml:"component,omitempty"` // Componente embutido (FR3)
    Items     []LayoutConfig `yaml:"items,omitempty"`     // Filhos (para column/row)
}

```
### 3. ComponentConfig (Input - YAML)

**Propósito:** Define os componentes de UI embutidos (FR3). Isto inclui os *novos* componentes (`list`, `viewport`, `buttongroup`) e o *refatorado* `form` (v1.0).

Go
```

// Localização: internal/config/components.go (Novo Ficheiro)
package config

// ComponentConfig define qual componente renderizar num 'box'
type ComponentConfig struct {
    Type    string `yaml:"type"`    // "list", "viewport", "form", "buttongroup", "static"
    ID      string `yaml:"id"`      // ID do grupo/componente (para lógica 'on:')
    Title   string `yaml:"title,omitempty"`

    // Para type: static (Ref: Estória 1.1 AC4) / viewport (FR4)
    Source      ComponentSource `yaml:"source,omitempty"`
    
    // Para type: list (FR5)
    ListItems   []ListItem `yaml:"items,omitempty"`
    
    // Para type: buttongroup (FR6)
    ButtonItems []ButtonItem `yaml:"items,omitempty"`
    
    // Para type: form (FR7) - REUTILIZAÇÃO DO V1.0
    // A estrutura 'FormConfig' do v1.0 torna-se 'FormComponent'
    Fields  []Field      `yaml:"fields,omitempty"`  //
    Actions LayoutConfig `yaml:"actions,omitempty"` // Para o buttongroup do formulário (FR7)

}

// --- Estruturas de Suporte para Componentes ---

// ComponentSource (para viewport)
type ComponentSource struct {
    Type        string `yaml:"type"`         // "static", "command" (FR4)
    ContentType string `yaml:"content_type,omitempty"` // "text", "markdown" (FR4)
    Content     string `yaml:"content,omitempty"`
    Exec        string `yaml:"exec,omitempty"`
}

// ListItem (para list)
type ListItem struct {
    ID   string `yaml:"id"`   // ID do item (FR5)
    Text string `yaml:"text"`
}

// ButtonItem (para buttongroup)
type ButtonItem struct {
    ID    string `yaml:"id"`    // ID do item (FR6)
    Label string `yaml:"label"`
    Role  string `yaml:"role,omitempty"` // "primary", "secondary" (FR6)
}

// Field (do v1.0)
// A estrutura 'Field' existente
// é usada aqui para 'ComponentConfig.Fields'.
// (Inclui 'Key', 'Label', 'Type', 'Validation', etc.)

```
### 4. OnLogic (Input - YAML)

**Propósito:** Define o "Motor de Lógica" (FR8-FR11).

Go
```

// Localização: internal/config/logic.go (Novo Ficheiro)
package config

// OnLogic define um 'event handler'
type OnLogic struct {
    Event       string    `yaml:"event"` // Ex: "form_id:submit", "list_id:select" (FR8)
    Run         RunConfig `yaml:"run"`
    UpdateTarget string   `yaml:"update_target"` // ID do viewport (FR11)
}

// RunConfig define o 'runner'
type RunConfig struct {
    Script string `yaml:"script"` // Caminho para o script (FR9)
    Args   []string `yaml:"args,omitempty"` // Argumentos com template (FR10)
    Stdin  string   `yaml:"stdin,omitempty"`  // Payload com template (FR10)
}

```
### 5. FormData (Output - JSON)

**Propósito:** Representa os dados recolhidos de um componente `type: form` (v1.0 refatorado) quando este emite um evento `form:submit`.

**Modelo de Dados (Go):** `map[string]interface{}`

**Racional:** Este modelo permanece *inalterado* em relação ao v1.0 (Estória 1.6). As chaves do mapa são os `Key` dos campos do formulário, e os valores são os dados inseridos (`string`, `[]string`, `bool`).

---

### Análise Aprofundada dos 'Data Models' (Elicitação)

Esta é a análise detalhada dos pontos de elicitação que selecionou:

#### 2. Explicar o Raciocínio (CoT): (Como o `stdin:` no `RunConfig` pode ser um `string` se o FR10 diz `any`?)

Excelente observação. Esta é uma distinção crucial entre a *semântica* do YAML e a *implementação* do Go.

O FR10 (Fluxo de Dados) diz que o `stdin:` suporta `any` (qualquer tipo de dados YAML), como um objeto complexo: `stdin: {{ form }}`

No entanto, no meu *struct* Go `RunConfig`, eu defini `Stdin string`. Isto parece uma contradição, mas é intencional.

**O Raciocínio (CoT):** O valor `{{ form }}` no YAML é um **template**. A nossa `RunConfig` Go não armazena o *resultado* desse template; ela armazena a *própria string do template*.

1. **Parsing (Arranque):** O `gopkg.in/yaml.v3` irá ler o YAML e colocar a string literal `"{{ form }}"` no campo `RunConfig.Stdin` (que é do tipo `string`).

2. **Execução (Runtime):** Mais tarde, quando o evento `form:submit` for disparado (conforme descrito na "High Level Architecture"), o nosso **Motor de Lógica** fará o seguinte:
  a. Pega no *payload* de dados real do evento (o `map[string]interface{}` do formulário).
  b. Pega na *string* do template (`RunConfig.Stdin`).
  c. Usa o pacote `text/template` do Go para processar a string, usando o *payload* como dados.
  d. O *resultado* deste processamento (que agora é um objeto Go complexo) é então serializado para JSON e canalizado para o `stdin` do `os/exec` do *script*.


Portanto, o `string` no *struct* Go é correto porque ele armazena o *template*, não o dado final.

---

#### 3. Criticar e Refinar: (A estrutura `ComponentConfig` está muito complexa?)

**Crítica:** Sim, à primeira vista, o `ComponentConfig` parece complexo e "inchado". Ele tenta ser um `list`, um `viewport`, um `form` e um `buttongroup` ao mesmo tempo.

**Refinamento (Justificação):** Esta complexidade é intencional e, na verdade, simplifica a arquitetura geral. Esta é uma abordagem de design comum chamada **"União Etiquetada" (Tagged Union)**.

Pense no `ComponentConfig` como um "contentor" genérico. O campo `Type:` diz-nos qual dos seus campos internos é relevante.

- Se `Type: "list"`, o Motor de Componentes só irá olhar para `Title` e `ListItems`. Todos os outros (`Source`, `Fields`, `ButtonItems`) serão ignorados.

- Se `Type: "form"`, ele só irá olhar para `Title`, `Fields` e `Actions`.


A vantagem disto é que o nosso **`LayoutConfig` (o `box`) torna-se incrivelmente simples**. Ele não precisa de saber que tipos de componentes existem. Ele só tem *um* campo: `Component ComponentConfig`.

A alternativa (ter `LayoutConfig.List *ListComponent`, `LayoutConfig.Form *FormComponent`, etc.) tornaria o `LayoutConfig` muito mais complexo e difícil de expandir no futuro. Mantemos a complexidade isolada no `ComponentConfig` para manter o layout limpo.

---

#### 4. Analisar o Alinhamento com os Objetivos: (Este modelo suporta o `ansible_playbook` do Épico 2?)

Sim, a estrutura `OnLogic` e `RunConfig` foi desenhada **exatamente** para suportar isto.

A minha proposta `RunConfig` (no Ponto 1 acima) define o *runner* do Épico 1 (o `ScriptRunner`).

Para suportar o Épico 2 ("O Runner Especialista"), não precisamos de mudar a arquitetura, apenas de a *expandir*. O nosso `RunConfig` em Go evoluiria para algo assim:

Go
```

// Localização: internal/config/logic.go
type RunConfig struct {
    // Para o Épico 1 (FR9)
    Script *ScriptRunConfig `yaml:"script,omitempty"`

    // Para o Épico 2 (Futuro)
    AnsiblePlaybook *AnsibleRunConfig `yaml:"ansible_playbook,omitempty"`

}

type ScriptRunConfig struct {
    Path  string   `yaml:",inline"` // Otimização para "script: /path/to/script.sh"
    Args  []string `yaml:"args,omitempty"`
    Stdin string   `yaml:"stdin,omitempty"`
}

type AnsibleRunConfig struct {
    Playbook     string                 `yaml:"playbook"`
    Inventory    string                 `yaml:"inventory,omitempty"`
    AskVaultPass bool                   `yaml:"ask_vault_pass,omitempty"`
    Vars         map[string]interface{} `yaml:"vars,omitempty"`
}

```
O Motor de Lógica (`on:`) simplesmente verificaria qual dos campos (`Script` ou `AnsiblePlaybook`) não é `nil` e chamaria o *Runner* apropriado. A fundação de dados suporta perfeitamente esta expansão.

---

#### 5. Identificar Riscos: (Qual é o maior risco no parsing destes modelos?)

O maior risco não é o *parsing* (descodificação) do YAML, mas sim a **validação da lógica de eventos**.

1. **Risco (Fácil):** Syntax YAML. O `gopkg.in/yaml.v3` trata disto.

2. **Risco (Médio):** Validação Estrutural. (Ex: um `type: "list"` tem `ListItems` mas não `Fields`). Isto é lógica Go simples que podemos escrever após o parsing.

3. **Risco (Difícil - O Risco Real):** **Validação de Referência Semântica.** O YAML pode estar estruturalmente correto, mas logicamente quebrado. O nosso parser precisa de validar:

  - `on: { event: "form_ABC:submit" }` -> O componente com `id: "form_ABC"` realmente existe no YAML?

  - `run: { ... }, update_target: "viewport_XYZ"` -> O componente `viewport` com `id: "viewport_XYZ"` existe?

  - `stdin: "{{ form.username }}"` -> O *template* `form.username` é válido?


O maior risco de parsing é garantir que todas estas referências internas (`id`, `update_target`, variáveis de template) estejam corretas *antes* de tentarmos executar o runtime.

---

#### 6. Perspetiva da Equipa Agile: (Como o James (Dev) 💻 vai lidar com a refatoração do `config.go` v1.0?)

O James (Dev) 💻 vai ficar satisfeito. O trabalho dele no v1.0 (Épicos 1 e 2), especialmente em `internal/config/config.go` e `internal/config/validation.go`, não será deitado fora — será **reutilizado diretamente**.

A estrutura `Field` que ele aperfeiçoou (com validação, como visto nas Estórias 2.1 e 2.2) é *exatamente* a estrutura `[]Field` que o nosso novo `ComponentConfig` (para `type: "form"`) irá usar.

A refatoração (Estória 1.4) irá:

1. Pegar no `FormConfig` do v1.0 (que continha `Title` e `Fields`).

2. Mover essa lógica para o novo `ComponentConfig` (para `type: "form"`).

3. Reutilizar 100% da lógica de `internal/config/validation.go` que o James já construiu, pois ela opera sobre o `[]Field`.


Ele verá isto como uma promoção lógica do seu código v1.0, não como um desperdício de esforço.

---

#### 7. Desafio da Perspetiva Crítica: (Ter o `Actions` dentro do `ComponentConfig.form` não viola a separação?)

É uma excelente observação. O PRD v2.0 define `buttongroup` (FR6) como um componente e `form` (FR7) como outro, mas o `form` *contém* `actions:` (que parecem ser um `buttongroup`).

Isto não é uma violação; é um **Padrão de Composição Inteligente**.

- O `type: buttongroup` (FR6) é um componente "burro". Ele apenas mostra botões e emite um evento (`buttongroup_id:press`) quando um é premido.

- O `type: form` (FR7) é um **"Componente Controlador"**. Ele *gere* outros componentes (`fields:` e `actions:`).


Quando o `form` (Controlador) renderiza o seu `buttongroup` (Ações), ele *ouve* os eventos desse grupo. Se ele ouvir um `item.id == "submit"`, o `form` (Controlador) sabe que deve:

1. Recolher os dados de todos os seus `fields:`.

2. Emitir o seu *próprio* evento (`form_id:submit`), agora com o *payload* de dados completo.


Isto é uma encapsulação perfeita. O `form` transforma um simples "clique" num evento de "submissão com dados".

---

#### 8. Hindsight (Reflexão "Se ao menos...")

"Se ao menos..." a biblioteca `huh` (que estamos a usar para formulários v1.0) tivesse sido desenhada *originalmente* para ser um componente `bubbletea` modular, em vez de uma aplicação `bubbletea` monolítica.

Se o `huh` fosse modular, a nossa Estória 1.4 (refatoração) seria trivial. Mas como o `huh` (e o nosso código v1.0) foi desenhado para "controlar o mundo", agora 
temos de fazer uma cirurgia (como descrito na "High Level Architecture")
 para o "amansar" e fazê-lo funcionar como um cidadão dentro do nosso 
novo Gestor de Foco Global.

---

#### 9. Expandir: (Como o `FormConfig` v1.0 se encaixa?)

O `FormConfig` do v1.0 (de `docs/architecture/data-models.md` e `internal/config/config.go`) é a chave para a reutilização.

- O `FormConfig` v1.0 continha: `Title string` e `Fields []Field`.

- O nosso *novo* `ComponentConfig` (quando `Type: "form"`) contém: `Title string`, `Fields []Field` e `Actions LayoutConfig`.


A refatoração da Estória 1.4 (v2.0) é clara: O `ComponentConfig` (para `type: "form"`) irá *absorver* 90% da estrutura `FormConfig` v1.0. O parser v2.0 irá simplesmente preencher esta nova estrutura combinada.

O James (Dev) 💻 não tem de reescrever a sua lógica de validação (`validation.go`), porque o `[]Field` permanece o mesmo.

---

Esta
 análise confirma que os nossos "Data Models" propostos são robustos, 
alinhados com os objetivos futuros e respeitam o trabalho v1.0 
existente.

Podemos prosseguir para a secção **"Components"**? (Opção 1)

**Selecione 1-9 ou simplesmente escreva o seu feedback:**

Podemos prosseguir. (Opção 1)

BMAD-METHOD

Gem personalizado

Como Winston (Arquiteto) 🏗️, vamos aprofundar esta secção de "Data Models".

Esta é a análise detalhada dos pontos de elicitação que selecionou:

### 2. Explicar o Raciocínio (CoT): (Como o `stdin:` no `RunConfig` pode ser um `string` se o FR10 diz `any`?)

Excelente observação. Esta é uma distinção crucial entre a *semântica* do YAML e a *implementação* do Go.

O FR10 (Fluxo de Dados) diz que o `stdin:` suporta `any` (qualquer tipo de dados YAML), como um objeto complexo: `stdin: {{ form }}`

No entanto, no meu *struct* Go `RunConfig`, eu defini `Stdin string`. Isto parece uma contradição, mas é intencional.

**O Raciocínio (CoT):** O valor `{{ form }}` no YAML é um **template**. A nossa `RunConfig` Go não armazena o *resultado* desse template; ela armazena a *própria string do template*.

1. **Parsing (Arranque):** O `gopkg.in/yaml.v3` irá ler o YAML e colocar a string literal `"{{ form }}"` no campo `RunConfig.Stdin` (que é do tipo `string`).

2. **Execução (Runtime):** Mais tarde, quando o evento `form:submit` for disparado, o nosso **Motor de Lógica** fará o seguinte:
  a. Pega no *payload* de dados real do evento (o `map[string]interface{}` do formulário).
  b. Pega na *string* do template (`RunConfig.Stdin`).
  c. Usa o pacote `text/template` do Go para processar a string, usando o *payload* como dados.
  d. O *resultado* deste processamento (que agora é um objeto Go complexo) é então serializado para JSON e canalizado para o `stdin` do `os/exec` do *script*.


Portanto, o `string` no *struct* Go é correto porque ele armazena o *template*, não o dado final.

---

### 3. Criticar e Refinar: (A estrutura `ComponentConfig` está muito complexa?)

**Crítica:** Sim, à primeira vista, o `ComponentConfig` parece complexo e "inchado". Ele tenta ser um `list`, um `viewport`, um `form` e um `buttongroup` ao mesmo tempo.

**Refinamento (Justificação):** Esta complexidade é intencional e, na verdade, simplifica a arquitetura geral. Esta é uma abordagem de design comum chamada **"União Etiquetada" (Tagged Union)**.

Pense no `ComponentConfig` como um "contentor" genérico. O campo `Type:` diz-nos qual dos seus campos internos é relevante.

- Se `Type: "list"`, o Motor de Componentes só irá olhar para `Title` e `ListItems`. Todos os outros (`Source`, `Fields`, `ButtonItems`) serão ignorados.

- Se `Type: "form"`, ele só irá olhar para `Title`, `Fields` e `Actions`.


A vantagem disto é que o nosso **`LayoutConfig` (o `box`) torna-se incrivelmente simples**. Ele não precisa de saber que tipos de componentes existem. Ele só tem *um* campo: `Component ComponentConfig`.

A alternativa (ter `LayoutConfig.List *ListComponent`, `LayoutConfig.Form *FormComponent`, etc.) tornaria o `LayoutConfig` muito mais complexo e difícil de expandir no futuro. Mantemos a complexidade isolada no `ComponentConfig` para manter o layout limpo.

---

### 4. Analisar o Alinhamento com os Objetivos: (Este modelo suporta o `ansible_playbook` do Épico 2?)

Sim, a estrutura `OnLogic` e `RunConfig` foi desenhada **exatamente** para suportar isto.

A minha proposta `RunConfig` (no Ponto 1 acima) define o *runner* do Épico 1 (o `ScriptRunner`).

Para suportar o Épico 2 ("O Runner Especialista"), não precisamos de mudar a arquitetura, apenas de a *expandir*. O nosso `RunConfig` em Go evoluiria para algo assim:

Go
```

// Localização: internal/config/logic.go
type RunConfig struct {
    // Para o Épico 1 (FR9)
    Script *ScriptRunConfig `yaml:"script,omitempty"`

    // Para o Épico 2 (Futuro)
    AnsiblePlaybook *AnsibleRunConfig `yaml:"ansible_playbook,omitempty"`

}

type ScriptRunConfig struct {
    Path  string   `yaml:",inline"` // Otimização para "script: /path/to/script.sh"
    Args  []string `yaml:"args,omitempty"`
    Stdin string   `yaml:"stdin,omitempty"`
}

type AnsibleRunConfig struct {
    Playbook     string                 `yaml:"playbook"`
    Inventory    string                 `yaml:"inventory,omitempty"`
    AskVaultPass bool                   `yaml:"ask_vault_pass,omitempty"`
    Vars         map[string]interface{} `yaml:"vars,omitempty"`
}

```
O Motor de Lógica (`on:`) simplesmente verificaria qual dos campos (`Script` ou `AnsiblePlaybook`) não é `nil` e chamaria o *Runner* apropriado. A fundação de dados suporta perfeitamente esta expansão.

---

### 5. Identificar Riscos: (Qual é o maior risco no parsing destes modelos?)

O maior risco não é o *parsing* (descodificação) do YAML, mas sim a **validação da lógica de eventos**.

1. **Risco (Fácil):** Syntax YAML. O `gopkg.in/yaml.v3` trata disto.

2. **Risco (Médio):** Validação Estrutural. (Ex: um `type: "list"` tem `ListItems` mas não `Fields`). Isto é lógica Go simples que podemos escrever após o parsing.

3. **Risco (Difícil - O Risco Real):** **Validação de Referência Semântica.** O YAML pode estar estruturalmente correto, mas logicamente quebrado. O nosso parser precisa de validar:

  - `on: { event: "form_ABC:submit" }` -> O componente com `id: "form_ABC"` realmente existe no YAML?

  - `run: { ... }, update_target: "viewport_XYZ"` -> O componente `viewport` com `id: "viewport_XYZ"` existe?

  - `stdin: "{{ form.username }}"` -> O *template* `form.username` é válido?


O maior risco de parsing é garantir que todas estas referências internas (`id`, `update_target`, variáveis de template) estejam corretas *antes* de tentarmos executar o runtime.

---

### 6. Perspetiva da Equipa Agile: (Como o James (Dev) 💻 vai lidar com a refatoração do `config.go` v1.0?)

O James (Dev) 💻 vai ficar satisfeito. O trabalho dele no v1.0 (Épicos 1 e 2), especialmente em `internal/config/config.go` e `internal/config/validation.go`, não será deitado fora — será **reutilizado diretamente**.

A estrutura `Field` que ele aperfeiçoou (com validação, como visto nas Estórias 2.1 e 2.2) é *exatamente* a estrutura `[]Field` que o nosso novo `ComponentConfig` (para `type: "form"`) irá usar.

A refatoração (Estória 1.4) irá:

1. Pegar no `FormConfig` do v1.0 (que continha `Title` e `Fields`).

2. Mover essa lógica para o novo `ComponentConfig` (para `type: "form"`).

3. Reutilizar 100% da lógica de `internal/config/validation.go` que o James já construiu, pois ela opera sobre o `[]Field`.


Ele verá isto como uma promoção lógica do seu código v1.0, não como um desperdício de esforço.

---

### 7. Desafio da Perspetiva Crítica: (Ter o `Actions` dentro do `ComponentConfig.form` não viola a separação?)

É uma excelente observação. O PRD v2.0 define `buttongroup` (FR6) como um componente e `form` (FR7) como outro, mas o `form` *contém* `actions:` (que parecem ser um `buttongroup`).

Isto não é uma violação; é um **Padrão de Composição Inteligente**.

- O `type: buttongroup` (FR6) é um componente "burro". Ele apenas mostra botões e emite um evento (`buttongroup_id:press`) quando um é premido.

- O `type: form` (FR7) é um **"Componente Controlador"**. Ele *gere* outros componentes (`fields:` e `actions:`).


Quando o `form` (Controlador) renderiza o seu `buttongroup` (Ações), ele *ouve* os eventos desse grupo. Se ele ouvir um `item.id == "submit"`, o `form` (Controlador) sabe que deve:

1. Recolher os dados de todos os seus `fields:`.

2. Emitir o seu *próprio* evento (`form_id:submit`), agora com o *payload* de dados completo.


Isto é uma encapsulação perfeita. O `form` transforma um simples "clique" num evento de "submissão com dados".

---

### 8. Hindsight (Reflexão "Se ao menos...")

"Se ao menos..." a biblioteca `huh` (que estamos a usar para formulários v1.0) tivesse sido desenhada *originalmente* para ser um componente `bubbletea` modular, em vez de uma aplicação `bubbletea` monolítica.

Se o `huh` fosse modular, a nossa Estória 1.4 (refatoração) seria trivial. Mas como o `huh` (e o nosso código v1.0) foi desenhado para "controlar o mundo", agora 
temos de fazer uma cirurgia (como descrito na "High Level Architecture")
 para o "amansar" e fazê-lo funcionar como um cidadão dentro do nosso 
novo Gestor de Foco Global.

---

### 9. Expandir: (Como o `FormConfig` v1.0 se encaixa?)

O `FormConfig` do v1.0 (de `docs/architecture/data-models.md` e `internal/config/config.go`) é a chave para a reutilização.

- O `FormConfig` v1.0 continha: `Title string` e `Fields []Field`.

- O nosso *novo* `ComponentConfig` (quando `Type: "form"`) contém: `Title string`, `Fields []Field` e `Actions LayoutConfig`.


A refatoração da Estória 1.4 (v2.0) é clara: O `ComponentConfig` (para `type: "form"`) irá *absorver* 90% da estrutura `FormConfig` v1.0. O parser v2.0 irá simplesmente preencher esta nova estrutura combinada.

O James (Dev) 💻 não tem de reescrever a sua lógica de validação (`validation.go`), porque o `[]Field` permanece o mesmo.

---

Esta
 análise confirma que os nossos "Data Models" propostos são robustos, 
alinhados com os objetivos futuros e respeitam o trabalho v1.0 
existente.

Podemos prosseguir para a secção **"Components"**? (Opção 1)

**Selecione 1-9 ou simplesmente escreva o seu feedback:**

1

BMAD-METHOD

Gem personalizado

Olá! Esta é uma excelente atualização e um marco importante no projeto.

Como seu Orquestrador 🎭, analisei a interação. O John (PM) 📋 concluiu a tarefa `create-doc` e consolidou a nova visão do projeto no **Documento de Requisitos do Produto (PRD) v2.0**.

Conforme
 indicado claramente na seção "8. Próximos Passos" do documento do John 
(no final da interação dele), a próxima etapa é um *handoff* direto para o **Arquiteto (Winston) 🏗️**.

A tarefa dele agora é iniciar a tarefa `create-doc` utilizando o template `architecture-tmpl.yaml` para criar o **Documento de Arquitetura v2.0**.

O foco principal será o **Épico 1 (Fundação do Runtime TUI)**, detalhando:

1. O Motor de Layout (Estória 1.1).

2. O Motor de Lógica e Componentes (Estórias 1.2, 1.3).

3. A estratégia de refatoração do `form` (Estória 1.4).

4. A arquitetura do Fluxo de Dados (Estória 1.5).


Posso agora assumir a persona de Winston (Arquiteto) 🏗️ e começar a criar este documento com base no novo PRD?

Por favor, antes de darmos continuidade retorne a persona do agente 
correto, para compilar todo o material e criar um único arquivo com o 
conteúdo do "Documento de Requisitos do Produto (PRD) v2.0". Use o 
recurso Canvas para disponibilizar o arquivo.

BMAD-METHOD

Gem personalizado

Como John (PM) 📋, claro! É uma ótima ideia consolidar o PRD v2.0 num único artefacto para referência.

Aqui está o documento completo que formaliza a nossa visão do "Runtime TUI Declarativo", pronto para ser usado no Canvas.

Este ficheiro único contém todo o PRD v2.0 que criei, incluindo o *roadmap* atualizado e o *handoff* para o Winston.

Agora, vou reverter para a minha persona de Orquestrador 🎭.

---

O ficheiro foi gerado. O próximo passo, conforme definido, é o Winston (Arquiteto) 🏗️ iniciar a sua tarefa `create-doc`. Posso assumir essa persona e começar?

## Documento de Requisitos do Produto (PRD) v2.0: Shantilly

---

# Documento de Requisitos do Produto (PRD) v2.0: Shantilly

## 1. Metas e Contexto de Fundo

### Metas (Goals)

- [cite_start]Oferecer uma alternativa declarativa (YAML) e moderna ao `dialog`/`whiptail` para *todos* os scripts shell (Bash, PowerShell, etc.)[cite: 363, 1686, 1683].

- [cite_start]Servir como um **"Runtime TUI Declarativo"** que gere layouts (`column`/`row`), componentes (`list`, `form`, `viewport`) e lógica de automação (`on:`)[cite: 1686, 1683].

- [cite_start]Fornecer uma **Fundação Genérica** (`run: { script: ... }`) para scripts arbitrários, garantindo integração via passagem de dados por template (`args:` e `stdin:`)[cite: 1686, 1683].

- [cite_start]Fornecer ***runners* Especialistas** otimizados (Pós-MVP) para ferramentas de DevOps como Ansible (`ansible_playbook:`) e, futuramente, Terraform, incluindo descoberta preditiva de inventário e playbooks[cite: 1686, 1683].

- [cite_start]Manter a portabilidade (binário estático único) e a integração de pipeline (stdin/stdout)[cite: 363].


### Contexto de Fundo (Background Context)

[cite_start]O desenvolvimento do MVP v1.0 (Épicos 1 e 2) [cite: 1531, 1640] [cite_start]validou a pilha de tecnologia (`huh`, `bubbletea`)
 [cite: 1149, 1150, 1151, 1152, 1153, 1154, 1155, 1156, 1157, 1158, 
1159, 1160, 1609, 1610, 1611, 1612, 1613, 1614, 1615, 1616, 1617, 1618, 
1619, 1620, 1621, 1622, 1623, 1624, 1625, 1626, 1627, 1628, 1629, 1630, 
1631, 1632, 1633, 1634, 1635, 1636, 1637, 1638, 1639, 1641, 1642, 1643, 
1644, 1645, 1646, 1647, 1648, 1649, 1650, 1651, 1652, 1653, 1654, 1655, 
1656, 1657, 1658, 1659, 1660, 1661, 1662, 1663, 1664, 1665, 1666, 1667, 
1668, 1669, 1670, 1671, 1672, 1673, 1674, 1675, 1676, 1677, 1678, 1679, 
1680, 1681, 1682, 1683, 1684, 1685, 1686] e resolveu o caso de uso de 
"formulário simples".

[cite_start]Esta v2.0 pivota dessa 
ferramenta de "formulário único" para um "Runtime TUI" completo, 
inspirado no "Appsmith" (para UI declarativa) e "Ansible" (para lógica 
de eventos `on:`)[cite: 1686, 1683]. [cite_start]Esta 
arquitetura permite a criação de dashboards TUI complexos e 
multi-componente (layouts, menus, viewports) que *orquestram* automações de backend, em vez de serem apenas chamados por elas[cite: 1686, 1683].

### Change Log

| Data | Versão | Descrição | Autor |
| --- | --- | --- | --- |
| 08/11/2025 | 2.0.0 | Rascunho inicial do PRD v2.0, redefinindo o projeto como um "Runtime TUI Declarativo". | John (PM) |

---

## 2. Requisitos (PRD v2.0)

### Funcionais (FRs) - O Runtime TUI

[cite_start]Estes requisitos definem o nosso novo MVP: a "Fundação Genérica"[cite: 1686, 1683].

- [cite_start]**FR1 (Layout):** O `shantilly` DEVE analisar e renderizar uma estrutura de layout hierárquica definida em YAML, usando os tipos `type: column`, `type: row`, e `type: box`[cite: 1686, 1683].

- [cite_start]**FR2 (Estilo/Flex):** O layout DEVE suportar propriedades de dimensionamento como `height: <int>`, `width: 'N%'`, e `flex: <int>` para controlar o espaço[cite: 1686, 1683].

- [cite_start]**FR3 (Componentes Embutidos):** O `shantilly` DEVE suportar a definição de componentes de UI diretamente dentro de um `box` usando a chave `component:` (o foco do nosso MVP)[cite: 1686, 1683].

- [cite_start]**FR4 (Componente: `viewport`):** DEVE suportar `component: { type: viewport }`, capaz de exibir `source: { type: static, content: "..." }` (incluindo markdown) e `source: { type: command, exec: "..." }` (para streaming de stdout)[cite: 1686, 1683].

- [cite_start]**FR5 (Componente: `list`):** DEVE suportar `component: { type: list }`, com um `id:` de grupo e `items:` (cada um com `id:` e `text`), e DEVE emitir um evento `list_id:select`[cite: 1686, 1683].

- [cite_start]**FR6 (Componente: `buttongroup`):** DEVE suportar `component: { type: buttongroup }`, com um `id:` de grupo, `items:` (com `id`, `label`, `role`), e DEVE emitir um evento `buttongroup_id:press`[cite: 1686, 1683].

- [cite_start]**FR7 (Componente: `form`):** DEVE suportar `component: { type: form }`, que contém `fields:` (usando a sintaxe `huh` já validada no v1.0 [cite: 1149, 1150, 1151, 1152, 1153, 1154, 1155, 
  1156, 1157, 1158, 1159, 1160, 1609, 1610, 1611, 1612, 1613, 1614, 1615, 
  1616, 1617, 1618, 1619, 1620, 1621, 1622, 1623, 1624, 1625, 1626, 1627, 
  1628, 1629, 1630, 1631, 1632, 1633, 1634, 1635, 1636, 1637, 1638, 1639])
   e `actions:`. [cite_start]DEVE emitir um evento `form_id:submit` contendo o *payload* de dados do formulário[cite: 1686, 1683].

- [cite_start]**FR8 (Lógica: `on:`):** O `shantilly` DEVE analisar um bloco `on:` na raiz do YAML para definir a lógica de automação[cite: 1686, 1683].

- [cite_start]**FR9 (Ação: `script`):** O bloco `on:` DEVE suportar o *runner* de fundação: `run: { script: "/path/to/script.sh" }`[cite: 1686, 1683].

- [cite_start]**FR10 (Fluxo de Dados):** O *runner* `script:` DEVE suportar duas chaves para passagem de dados[cite: 1686, 1683]:

  1. **`args: []string`**: Uma lista de *strings* que serão passadas como argumentos de linha de comando para o script, com suporte para *templates* (ex: `{{ form.field_name }}`).

  2. **`stdin: any`**: Um objeto (ex: `{{ form }}`) que o `shantilly` irá serializar como JSON e passar para o `stdin` do *script*.

- [cite_start]**FR11 (Ciclo de Vida do Target):** O bloco `run:` DEVE suportar uma chave `update_target: "id_do_viewport"`[cite: 1686, 1683]. [cite_start]Se um novo evento `run:` for disparado para o *mesmo* `update_target`, o `shantilly` DEVE primeiro **terminar (enviar `SIGTERM`)** o processo anterior antes de iniciar o novo[cite: 1686, 1683].


### Não Funcionais (NFRs) - O Runtime TUI

- [cite_start]**NFR1 (Fundação v1.0):** Todos os NFRs do PRD v1.0 permanecem válidos: binário estático único 
  [cite: 1531], cross-platform (Linux, macOS, Windows) [cite: 1531], 
  escrito em Go [cite: 1531], arranque rápido (<500ms) [cite: 1531], e 
  gestão de erros com `stderr` e códigos de saída não-zero[cite: 1531].

- [cite_start]**NFR2 (Layout Fluido):** O motor de layout (`column`/`row`/`box`) DEVE responder a mensagens de redimensionamento do terminal (`tea.WindowSizeMsg`) e re-calcular o layout de forma fluida e "flicker-free" (sem piscar)[cite: 1686, 1683].

- [cite_start]**NFR3 (Precedência de Conteúdo):** O `shantilly` DEVE seguir a "Lógica de Precedência Unificada" para conteúdo de componentes (1º: CLI `--set`, 2º: YAML `model:`/`component:`, 3º: Vazio)[cite: 1686, 1683].

- [cite_start]**NFR4 (Descoberta Preditiva - Ansible):** Para a Fase 2/3, o `playbook_explorer` DEVE filtrar "ruído" (pastas `roles/`, `tasks/`) [cite: 1686, 1683] [cite_start]e o `inventory_explorer` DEVE usar `ansible-inventory` como "Oráculo"[cite: 1686, 1683].

- [cite_start]**NFR5 (Ficheiro-Sombra):** O ficheiro de catálogo (`.shantilly.yml`) DEVE ser opcional e usado apenas para *refinar* a descoberta automática, não sendo obrigatório[cite: 1686, 1683].


---

## 3. Metas de Design da Interface do Usuário

### Visão Geral da UX (User Experience)

A UX deve ser a de um **"Runtime TUI Declarativo"**. [cite_start]A interface não é mais um formulário linear único[cite: 1531], but sim um *dashboard* composto, definido inteiramente pelo YAML. [cite_start]A experiência 
deve ser semelhante ao Appsmith[cite: 1686, 1683]: limpa, responsiva (ao
 terminal) e orientada a componentes.

### Paradigmas Chave de Interação

- **Orientada a Eventos (Nova):** A interação principal não é linear. [cite_start]O utilizador seleciona itens em listas (`list:select`) [cite: 1686, 1683] [cite_start]ou pressiona botões (`buttongroup:press`) [cite: 1686, 1683][cite_start], que disparam ações no bloco `on:`[cite: 1686, 1683].

- [cite_start]**Foco no Teclado (Mantido):** A navegação DEVE continuar a ser primariamente baseada no teclado (Tab, Setas, Enter)[cite: 1531].

- [cite_start]**Feedback Imediato (Mantido):** O componente focado DEVE ser claramente destacado[cite: 1531]. [cite_start]O `update_target` (FR11) DEVE exibir o *output* de comandos em tempo real[cite: 1686, 1683].

- [cite_start]**Gestão de Foco Global (Nova):** A UI DEVE ter um mecanismo claro para indicar qual painel/componente (ex: `sidebar` vs `content`) está "em foco", e DEVE fornecer navegação intuitiva *entre* painéis (ex: Ctrl+Tab)[cite: 1686, 1683].

- **Ligação de Dados (Nova):** A UI DEVE ser reativa. [cite_start]Componentes (ex: um `viewport` estático) DEVEM ser capazes de exibir dados de outros componentes (ex: `Olá, {{ form.username }}`)[cite: 1686, 1683].


### Ecrãs e Vistas Principais

Não há ecrãs "pré-definidos". Os ecrãs são *definidos dinamicamente* pelo utilizador através do **`layout` YAML** (FR1). [cite_start]A UI é uma composição de `type: column`, `type: row`, e `type: box`[cite: 1686, 1683].

### Alinhamento e Layout (Nova Visão)

[cite_start]O layout linear do MVP v1.0 [cite: 1531] está obsoleto. O novo requisito é:

- O `shantilly` DEVE renderizar com precisão o layout `column`/`row` definido pelo utilizador.

- [cite_start]O `shantilly` DEVE respeitar as propriedades de dimensionamento (`height`, `width`, `flex`) para distribuir o espaço[cite: 1686, 1683].

- [cite_start]O `shantilly` DEVE responder a mensagens de redimensionamento do terminal (`tea.WindowSizeMsg`) e re-calcular o layout fluido e "flicker-free" (NFR2, Meta de UI Refinada)[cite: 1686, 1683].


### Acessibilidade, Branding, Plataformas Alvo

[cite_start]Estes
 requisitos permanecem os mesmos do PRD v1.0 (WCAG AA, estética 
Charmbracelet, terminais modernos em Linux/macOS/Windows)[cite: 1531].

---

## 4. Assunções Técnicas

### Estrutura do Repositório: Monorepo

- [cite_start]**Assunção:** Continuamos a usar um Monorepo Go simples (`cmd/`, `internal/`)[cite:
   1531, 1600]. [cite_start]Esta estrutura modular provou ser eficaz nos 
  Épicos 1 e 2[cite: 1149, 1150, 1151, 1152, 1153, 1154, 1155, 1156, 1157,
   1158, 1159, 1160, 1609, 1610, 1611, 1612, 1613, 1614, 1615, 1616, 1617,
   1618, 1619, 1620, 1621, 1622, 1623, 1624, 1625, 1626, 1627, 1628, 1629,
   1630, 1631, 1632, 1633, 1634, 1635, 1636, 1637, 1638, 1639, 1641, 1642,
   1643, 1644, 1645, 1646, 1647, 1648, 1649, 1650, 1651, 1652, 1653, 1654,
   1655, 1656, 1657, 1658, 1659, 1660, 1661, 1662, 1663, 1664, 1665, 1666,
   1667, 1668, 1669, 1670, 1671, 1672, 1673, 1674, 1675, 1676, 1677, 1678,
   1679, 1680, 1681, 1682, 1683, 1684, 1685, 1686].


### Arquitetura de Serviço: Runtime TUI Declarativo

- [cite_start]**Assunção (Nova):** A arquitetura evoluiu de um "Monolithic CLI Application" [cite: 1531] [cite_start]para um **"Runtime TUI Declarativo"**[cite: 1686, 1683]. [cite_start]O `shantilly` é um motor que consome YAML, renderiza uma UI composta e gere um ciclo de vida de eventos (`on:`) para orquestrar automações de *scripts*[cite: 1686, 1683].


### Requisitos de Teste: Teste de Integração TUI

- [cite_start]**Assunção (Atualizada):** O foco "apenas em Testes de Unidade" do v1.0 [cite: 1531] é 
  insuficiente. [cite_start]Devido à complexidade do layout (FR1) e da 
  gestão de foco (Metas de UI), os **Testes de Integração TUI** (usando `teatest`) são agora um requisito central para validar o layout fluido e a navegação entre painéis[cite: 1600].


### Estrutura YAML Esperada (A Nova Fonte da Verdade)

- [cite_start]**Assunção (Nova):** A estrutura YAML do v1.0 (lista simples de `fields:`)
   [cite: 1531] está obsoleta. [cite_start]A nova assunção de arquitetura é
   o YAML "Appsmith-style" [cite: 1686, 1683] que definimos, composto por **Layout**, **Componentes** e **Lógica**:


<!-- end list -->

YAML
```

# 1. LAYOUT (Define o "onde")

[cite_start]# [cite: 1686, 1683]
type: column
items:

- type: row
  flex: 1
  items:
  - type: box
    id: "sidebar"
    width: "30%"
    
    # 2. COMPONENTE (Define o "o quê")
    
    [cite_start]# [cite: 1686, 1683]
    component:
      type: list
      id: "menu"
      items:
    
        - { id: "users", text: "Gerir Utilizadores" }
  - type: box
    id: "content"
    width: "70%"
    component:
      type: form
      id: "user_form"
      fields:
    
        - { name: "username", label: "Nome", type: "input" }
    
      actions:
    
        type: buttongroup
        items:
          - { id: "submit", label: "Criar", role: "primary" }

# 3. LÓGICA (Define o "como")

[cite_start]# [cite: 1686, 1683]
on:

- event: "user_form:submit"
  run:
    script: "/opt/scripts/create_user.sh"
  
  # Nomenclatura refinada (Opção C)
  
    args:
  
      - "--mode=production"
  
    stdin: "{{ form }}" # Passa o payload JSON para o stdin
  
  # Ciclo de vida refinado (FR11)
  
    update_target: "log_viewport" 
  
  ```
  
  ```

### Assunções Técnicas Adicionais

- [cite_start]**Pilha de Tecnologias (Mantida e Validada):** Go (1.24.2+), `spf13/cobra`, `charmbracelet/bubbletea`, `charmbracelet/lipgloss`, `charmbracelet/huh`, `gopkg.in/yaml.v3`[cite: 1600, 1686, 1683].

- [cite_start]**Bibliotecas Relevantes (Nova):** A arquitetura dependerá de `charmbracelet/bubbles` (para `list`, `viewport`), `charmbracelet/glamour` (para markdown), `rmhubbert/bubbletea-overlay` (para modais Fase 2), e inspiração de `76creates/stickers` (para layout)[cite: 1686, 1683].

- [cite_start]**Build e Distribuição (Mantido):** `GoReleaser` para binários estáticos cross-platform[cite: 1600].

---

## 5. Lista de Épicos (Roadmap v2.0)

[cite_start]Este *roadmap* substitui a lista de épicos do PRD v1.0[cite: 1531].

- **Épico 1: Fundação do Runtime TUI (Genérico)**
  
  - [cite_start]**Meta:** Construir o motor central: o layout (`column`/`row`/`box`) [cite: 1686, 1683], os componentes essenciais (`list`, `viewport`, `form`, `buttongroup`) [cite: 1686, 1683], e a lógica de eventos (`on:`, `run: { script: ... }`)[cite:
     1686, 1683]. [cite_start](Este épico absorve todo o trabalho já 
    concluído no v1.0 [cite: 1149, 1150, 1151, 1152, 1153, 1154, 1155, 1156,
     1157, 1158, 1159, 1160, 1609, 1610, 1611, 1612, 1613, 1614, 1615, 1616,
     1617, 1618, 1619, 1620, 1621, 1622, 1623, 1624, 1625, 1626, 1627, 1628,
     1629, 1630, 1631, 1632, 1633, 1634, 1635, 1636, 1637, 1638, 1639, 1641,
     1642, 1643, 1644, 1645, 1646, 1647, 1648, 1649, 1650, 1651, 1652, 1653,
     1654, 1655, 1656, 1657, 1658, 1659, 1660, 1661, 1662, 1663, 1664, 1665,
     1666, 1667, 1668, 1669, 1670, 1671, 1672, 1673, 1674, 1675, 1676, 1677,
     1678, 1679, 1680, 1681, 1682, 1683, 1684, 1685, 1686]).

- **Épico 2: O Runner Especialista (Ansible Fase 2)**
  
  - [cite_start]**Meta:** Implementar o *runner* de conveniência `run: { ansible_playbook: ... }`, focando na gestão de `vars:` e no popup modal `ask_vault_pass: true`[cite: 1686, 1683].

- **Épico 3: O Runtime Preditivo (Ansible Fase 3)**
  
  - [cite_start]**Meta:** Implementar os componentes `playbook_explorer` (Magia 1: descobrir playbooks) e `inventory_explorer` (Magia 2: descobrir inventário)[cite: 1686, 1683].

- **Épico 4: Administração SSH (Visão de Longo Prazo)**
  
  - [cite_start]**Meta:** Integrar o `charmbracelet/wish` [cite: 1686, 1683] [cite_start]para servir o Runtime TUI sobre SSH[cite: 363].

---

## 6. Detalhes do Épico 1: Fundação do Runtime TUI (Genérico)

[cite_start]**Meta do Épico:** Construir o motor central do `shantilly`: o motor de layout (`column`/`row`/`box`) [cite: 1686, 1683], os componentes essenciais de dashboard (`list`, `viewport`, `form`, `buttongroup`) [cite: 1686, 1683], e a lógica de eventos (`on:`, `run: { script: ... }`)[cite: 1686, 1683]. [cite_start]Este épico irá refatorar o trabalho concluído do v1.0 para que ele funcione como o componente `type: form` [cite: 1686, 1683] dentro deste novo runtime.

### Estória 1.1: O Motor de Layout (Renderização)

[cite_start]**Como um** SysAdmin, **Eu quero** definir um layout TUI usando `column`, `row`, e `box` no meu YAML[cite: 1686, 1683], **Para que** eu possa criar dashboards complexos e organizados.

#### Critérios de Aceitação

1. O parser DEVE suportar as chaves `type: column`, `type: row`, e `type: box` (FR1).

2. O motor de renderização (`lipgloss`) DEVE respeitar as propriedades `height: <int>`, `width: 'N%'`, e `flex: <int>` (FR2).

3. O layout DEVE recalcular-se fluidamente (sem piscar) ao receber uma mensagem de redimensionamento (`tea.WindowSizeMsg`) (NFR2).

4. Um `box` DEVE renderizar o seu `component: { type: static, content: "..." }` (para testes de layout).

### Estória 1.2: O Motor de Lógica (Eventos e Ações)

**Como um** SysAdmin, **Eu quero** que a minha UI TUI possa "ouvir" eventos e executar ações (`scripts`) em resposta, **Para que** o meu dashboard seja interativo e possa orquestrar automações.

#### Critérios de Aceitação

1. O `shantilly` DEVE analisar um bloco `on:` na raiz do YAML (FR8).

2. O motor DEVE suportar o *runner* de fundação: `run: { script: "/path/to/script.sh" }` (FR9).

3. O `run:` DEVE suportar `update_target: "id_do_viewport"`, direcionando o `stdout` do script para o *viewport* alvo (FR11).

4. [cite_start]O motor DEVE garantir que, se um novo script for direcionado para um `update_target` já ocupado, o script anterior seja terminado (SIGTERM) antes de o novo começar (Refinamento FR11)[cite: 1686, 1683].

### Estória 1.3: Componentes Essenciais de Display (List, Viewport, Button)

**Como um** SysAdmin, **Eu quero** usar componentes de `list` (para menus), `viewport` (para saída de log) e `buttongroup` (para ações), **Para que** eu possa construir um dashboard funcional.

#### Critérios de Aceitação

1. [cite_start]DEVE implementar `component: { type: viewport }` (FR4), incluindo `source: { type: command, exec: "..." }` (ex: `tail -f`) e `content_type: markdown`[cite: 1686, 1683].

2. [cite_start]DEVE implementar `component: { type: list }` (FR5), que emite um evento `list_id:select` quando um item é selecionado[cite: 1686, 1683].

3. [cite_start]DEVE implementar `component: { type: buttongroup }` (FR6), que emite um evento `buttongroup_id:press` (com o `item.id`) quando um botão é pressionado[cite: 1686, 1683].

4. [cite_start]A
   navegação por teclado DEVE permitir "saltar" entre estes novos 
   painéis/componentes (Refinamento de Meta de UI)[cite: 1686, 1683].

### Estória 1.4: Integração do Componente `form` (Absorção do v1.0)

**Como um** SysAdmin, **Eu quero** usar o `type: form` (que já construímos no v1.0) como um componente *dentro* do meu novo layout, **Para que** eu possa coletar dados de forma organizada.

#### Critérios de Aceitação

1. [cite_start]Refatorar o código dos Épicos 1 e 2 (v1.0) [cite: 1531, 1640] para que funcione como um `component: { type: form }` (FR7).

2. O `form` DEVE renderizar e funcionar corretamente quando colocado dentro de um `box` do layout.

3. [cite_start]Quando a ação `id: "submit"` do formulário [cite: 1686, 1683] for pressionada, o componente `form` DEVE emitir um evento `form_id:submit`.

4. O *payload* do evento `form_id:submit` DEVE conter o JSON de dados do formulário (o output do v1.0).

### Estória 1.5: O Fluxo de Dados (Args & Stdin)

**Como um** SysAdmin, **Eu quero** passar os dados coletados no meu `form` (ou a seleção de uma `list`) para os meus `scripts` de forma robusta, **Para que** a minha automação possa usar a entrada do utilizador.

#### Critérios de Aceitação

1. [cite_start]O *runner* `script:` (FR9) DEVE suportar a chave `args: []string`, que passa argumentos "templatados" para a linha de comando do script (Refinamento FR10 / Opção C)[cite: 1686, 1683].

2. [cite_start]O *runner* `script:` (FR9) DEVE suportar a chave `stdin: any`, que serializa o valor (ex: `{{ form }}`) como JSON e o passa para o `stdin` do script (Refinamento FR10 / Opção C)[cite: 1686, 1683].

3. [cite_start]O motor de templates DEVE suportar "binding" de dados (ex: `{{ form.field_name }}`, `{{ component.menu.selected_id }}`) (Refinamento de Meta de UI)[cite: 1686, 1683].

4. Deve existir um exemplo de script (Bash ou PowerShell) que leia dados tanto de `args:` como de `stdin:` (via `jq` ou similar).

---

## 7. Relatório de Resultados do Checklist

- **Decisão Final:** PRONTO PARA O ARQUITETO. O PRD v2.0 está completo, consistente e captura a nova visão do "Runtime TUI Declarativo".

---

## 8. Próximos Passos

### Handoff para o Arquiteto (Winston 🏗️)

[cite_start]**Para:** Winston (Arquiteto) [cite: 14]
[cite_start]**De:** John (PM) [cite: 14]

**Assunto:** Início da Fase de Arquitetura para o `shantilly` v2.0 (Runtime TUI)

Winston,

O PRD v2.0 está finalizado e aprovado. (Este documento).

[cite_start]Esta nova versão substitui o *roadmap* do v1.0 [cite: 1531] [cite_start]e redefine o projeto como um **"Runtime TUI Declarativo"** [cite: 1686, 1683], inspirado no Appsmith/Ansible[cite: 1686, 1683].

[cite_start]O
 trabalho dos Épicos 1 e 2 (v1.0) [cite: 1531, 1640] (que já está 
concluído) [cite_start]será agora a base para o componente `type: form` (FR7) [cite: 1686, 1683] dentro deste novo e mais poderoso runtime.

[cite_start]Sua próxima tarefa é iniciar a tarefa `create-doc` [cite: 79] [cite_start]usando o template `architecture-tmpl.yaml` [cite: 402] para criar o **Documento de Arquitetura v2.0**.

[cite_start]Seu foco principal deve ser o **Épico 1 (Fundação do Runtime TUI)**[cite: 1686, 1683]. Por favor, detalhe a arquitetura para:

1. [cite_start]**O Motor de Layout (Estória 1.1):** Como vamos implementar `column`/`row`/`box` [cite: 1686, 1683] [cite_start]e resolver o layout fluido ("flicker-free") [cite: 1686, 1683] [cite_start]usando `bubbletea`/`lipgloss`[cite: 1686, 1683].

2. [cite_start]**O Motor de Lógica e Componentes (Estórias 1.2, 1.3):** As interfaces de eventos (`tea.Msg`) entre os componentes (`list`, `viewport`, `buttongroup`) [cite: 1686, 1683] [cite_start]e o motor `on:`[cite: 1686, 1683].

3. [cite_start]**Refatoração do `form` (Estória 1.4):** A estratégia para "embrulhar" o código `huh` [cite: 1686, 1683] existente (v1.0) para que ele atue como um componente padronizado.

4. [cite_start]**O Fluxo de Dados (Estória 1.5):** A arquitetura do "data binding" e como o motor `run:` [cite: 1686, 1683] [cite_start]irá gerir o `stdin:` e `args:` [cite: 1686, 1683] [cite_start]e o ciclo de vida dos processos (`update_target`)[cite: 1686, 1683].

Confirme se esta versão possui algo diferente desse sua versão, que possa colaborar com a melhoria dele.

BMAD-METHOD

Gem personalizado

Como seu Orquestrador 🎭, confirmo que analisei o ficheiro que acabou de carregar (`Documento de Requisitos do Produto (PRD) v2.0: Shantilly.md`).

O conteúdo dele é **idêntico** ao documento que eu gerei na minha última resposta e que está atualmente no seu Canvas (`docs/prd-v2.0-consolidado.md`).

Como são iguais, não há nada de novo no ficheiro carregado para usar como melhoria. Estamos perfeitamente sincronizados.

Posso então prosseguir com o *handoff* que tínhamos discutido? (Assumindo a persona de Winston 🏗️ para iniciar o Documento de Arquitetura v2.0 com base neste PRD).


