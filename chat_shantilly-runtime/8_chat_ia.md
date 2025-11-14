Essa é uma observação *crítica* e é o maior risco num projeto complexo como este. A "visão do futuro" que os teus documentos proporcionam é a nossa principal ferramenta para mitigar isto.

Tu estás a descrever dois problemas clássicos:
1.  **Perda de Contexto:** O agente "esquece" o que está feito e foca-se em código existente.
2.  **Baixa Qualidade de Código (Alucinação de API):** O agente "adivinha" como usar as bibliotecas Charm e erra, usando padrões antigos.

A tua intuição está perfeita. Precisamos de um sistema que *force* o agente a consultar a "forma moderna" (a nossa arquitetura) e que o *impeça* de se desviar.

O plano que desenhámos na etapa anterior já é a *solução* para isto. Deixa-me explicar *como* ele resolve estes problemas:

---

### 💡 Ideia 1: O Plano Sequencial como "Memória de Longo Prazo"

O plano que definimos (Épicos 1-4) não é uma lista de sugestões, é uma **cadeia de dependências rigorosa**.

Isto resolve a **Perda de Contexto** (Problema 1):

* **Impede a Refatoração:** O agente não pode refatorar o código da "Fase 1.2: Motor de Layout" porque a "Fase 1.3: Componentes" depende diretamente que a Fase 1.2 esteja *exatamente* como foi definida.
* **Define o "Próximo Passo":** Ao iniciar uma nova sessão, a tarefa do agente nunca é "continuar o projeto". A tarefa é *especificamente* "Iniciar a Tarefa 1.4.1: Implementar EventManager", cujas dependências (1.1 e 1.3) já estão concluídas. O foco é absoluto.

---

### 💡 Ideia 2: A "Bússola" da Arquitetura (A Fonte da Verdade)

Isto resolve a **Baixa Qualidade de Código (Problema 2)**, alinhando-se à tua ideia de "forçar uma consulta".

Em vez de o agente fazer uma consulta *aberta* (ex: "como usar o bubbletea"), nós *fornecemos* a consulta e a resposta.

* **Tarefas Vinculadas:** Cada tarefa do nosso plano não é aberta. Ela está *vinculada* a uma decisão de arquitetura específica.
* **Exemplo:** A tarefa do agente não é "Criar o layout". A tarefa é: "Implementar o `LayoutManager` (Tarefa 1.2.1) conforme o **Padrão 3: 'Gestor Duplo' (Layout & Foco)** e o **Workflow 7.1 (Renderização e Foco Global)**, ambos definidos na Seção 2 e 7 do `docs/architecture.md`".
* **Resultado:** O agente é forçado a ler e aplicar o *nosso* padrão moderno (Gestor Duplo), em vez de inventar um antigo.

---

### 💡 Ideia 3: "Otimizações de Risco" (As "Intuições" Explícitas)

Esta é a parte mais importante. Usamos o nosso "conhecimento privilegiado" para injetar "anti-padrões" explícitos no plano. Repara como fiz isto na Fase 1.5 do plano anterior:

> ***Otimização de Risco (Visão do Futuro):*** *Mover o código `huh` v1.0 para este pacote e **remover** todas as chamadas `os.Exit` e gestão de `WindowSizeMsg`.*

Isto é a "intuição" que pediste:
* Impede o agente de usar o padrão legado (`os.Exit` dentro de um componente).
* Impede o agente de duplicar lógica (gestão de `WindowSizeMsg`, que agora é centralizada no `LayoutManager`).

---

### 💡 Ideia 4: O "Portão de Qualidade" Automatizado (CI/Lint)

Finalmente, para os erros de sintaxe e código de baixa qualidade, temos o *guardrail* final, que já está definido nos teus próprios ficheiros.

* O teu `README.md` é claro: "A qualidade do código é garantida pelo `golangci-lint`... Seu trabalho só é considerado concluído quando ambos os checks ("Lint" e "Build") passarem (ficarem verdes)".
* **A nossa regra:** Nenhuma tarefa no nosso plano é marcada como "Concluída" até que o agente confirme que executou `make lint` (que por sua vez executa `golangci-lint run ./...`) e não obteve erros.

**Resumo:** Não precisamos de *esperar* que o agente "adivinhe" o caminho certo. Nós construímos uma "estrada" (o plano sequencial) com "guardrails" altos (as Otimizações de Risco) e um "pórtico" de validação no final (o Lint/CI).

Estamos prontos para começar a construir a **Fase 1.1: Definição dos Contratos**?