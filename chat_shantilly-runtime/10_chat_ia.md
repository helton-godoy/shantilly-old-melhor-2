Esta é a pergunta-chave! 🎯 Para que o plano funcione como um "trilho" e evite os desvios que mencionaste, precisamos de um conjunto de artefatos que funcionem como a "memória" e as "regras" para os agentes de IA.

Aqui estão os artefatos essenciais que garantem que o plano seja executado com precisão:

---

### 1. Artefatos de Governança (O "Plano de Voo")

Estes artefatos definem **o quê** fazer, **porquê** e **em que ordem**. Eles são a "memória de longo prazo" que combate a perda de contexto.

* **O Plano Mestre Sequencial (O que acabámos de criar):** Este é o artefato mais crítico. É a nossa *instrução de trabalho*.
    * **Função:** Define a cadeia exata de Fases e Tarefas (ex: 1.1 -> 1.2 -> ...). Um agente nunca pode "refatorar" a Fase 1.1 porque a Fase 1.2 depende dela.
    * **Garante:** Foco absoluto na tarefa pendente, eliminando o risco de retrabalho.

* **O PRD v2.0 (docs/prd.md):** Este é o contrato que define o *escopo*.
    * **Função:** Justifica *porquê* uma tarefa existe. Contém os Requisitos Funcionais (FRs) e Não-Funcionais (NFRs) que cada tarefa deve cumprir.
    * **Garante:** Que a funcionalidade implementada (ex: FR11, o ciclo de vida do `update_target`) corresponde exatamente ao que foi pedido.

---

### 2. Artefatos de Arquitetura (A "Fonte da Verdade" Técnica)

Estes artefatos respondem **como** implementar. Eles são a "consulta" forçada que garante o uso de padrões modernos e evita a "alucinação" de código legado do Charm.

* **O Documento de Arquitetura (docs/architecture.md):** Esta é a "Bússola" técnica.
    * **Função:** Define os padrões de design obrigatórios (ex: "Padrão 3: Gestor Duplo"), os contratos de interface (Seção 5.2, `ShantillyComponent`), os Modelos de Dados (Seção 4) e a Árvore de Código-Fonte (Seção 10).
    * **Garante:** Que a IA não *invente* uma forma de implementar o layout; ela deve seguir o padrão "Gestor Duplo" definido.

* **O Diagrama de Alto Nível (docs/architecture/high-level-architecture.md):**
    * **Função:** Reforça visualmente o padrão mais crítico da v2.0: a separação de responsabilidades (LayoutManager vs. EventManager vs. ScriptRunner).
    * **Garante:** Que a IA não misture lógica de automação dentro do `LayoutManager`, um erro comum que a nossa "visão do futuro" identificou.

---

### 3. Artefatos de Qualidade (Os "Guardrails" Automatizados)

Estes artefatos **forçam** a conformidade técnica e impedem a submissão de código de baixa qualidade.

* **Configuração do Linter (.golangci.yml):**
    * **Função:** É a definição *automática* e *inequívoca* de "código de qualidade". Define regras como `errcheck`, `gofumpt`, `stylecheck`, etc..
    * **Garante:** Que o código está formatado corretamente e segue as melhores práticas idiomáticas do Go, resolvendo o problema de "baixa qualidade" que observaste.

* **O README (README.md):**
    * **Função:** Contém a diretiva explícita para o agente de IA: "Seu trabalho só é considerado concluído quando ambos os checks ('Lint' e 'Build') passarem".
    * **Garante:** Que o agente não pode marcar uma tarefa como "pronta" até que ela passe no "Portão de Qualidade" (o `.golangci.yml`).

Resumindo, o *nosso plano* (Governança) vincula cada tarefa à *nossa arquitetura* (Fonte da Verdade), e o *linter* (Qualidade) valida automaticamente o resultado. É um sistema fechado desenhado para guiar a IA! 🚀