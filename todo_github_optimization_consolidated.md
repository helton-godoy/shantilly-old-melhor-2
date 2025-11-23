# TODO: GitHub Optimization – Consolidado

> **Fonte única de verdade** para o controle das tarefas de otimização do GitHub no projeto Shantilly.
> Consolidado a partir de: `todo_github_optimization.md`, `todo_github_optimization_progress.md`,
> `todo_github_finalization.md`, `todo_progress_updated.md`, `todo_remaining_tasks.md`,
> `todo_final_progress.md`.
>
> Status inicial baseado principalmente em `todo_final_progress.md` (snapshot mais recente);
> detalhes complementares vindos dos demais arquivos.

---

## FASE 1 – Registro e Orquestração no GitHub

- [x] **1.1 Atualizar repositório principal Shantilly**  
  README otimizado, `.gitignore`, `SECURITY.md`.

- [x] **1.2 Criar 20+ issues iniciais para roadmap evolutivo**  
  25 issues estruturadas com épicos/waves e organização por sprints.

- [x] **1.3 Configurar templates personalizados de Issues**  
  `.github/ISSUE_TEMPLATE/` configurado.

- [x] **1.4 Criar Project principal "Shantilly Roadmap" com automações**  
  Estrutura e configuração inicial definidas (colunas, automações básicas, YAML).

- [x] **1.5 Definir 5+ milestones estratégicos**  
  Ex.: `v1.0 Alpha`, `v1.0 Beta`, `v2.0 Features`, `v2.0 Extensions`/`v2.0 Extensions`.

- [x] **1.6 Criar 15+ labels hierárquicas**  
  Families principais: `priority::*`, `type::*`, `area::*` (e possivelmente `status::*`, `epic::*`).

- [x] **1.7 Configurar Discussions reais no repositório GitHub**  
  - Categorias: **Ideas**, **Q&A**, **Polls**, **General**.  
  - Arquivo de config e documentação já existem; falta garantir/ajustar a configuração final na UI do GitHub.

- [x] **1.8 Adicionar Topics, configurar Releases e ativar Codespaces**  
  Topics configurados, Releases básicas criadas, `devcontainer.json` adicionado.

---

## FASE 2 – Documentação e Artefatos

- [x] **2.1 Atualizar README.md com badges e quickstart melhorado**

- [x] **2.2 Criar Wiki expandida com recursos GitHub integrados**  
  Documentação do workflow GitHub, incluindo colaboração IA-humanos.

- [x] **2.3 Documentar fluxo completo com diagramas Mermaid**  
  Diagrama(s) de fluxo cobrindo core workflows e integrações.

- [x] **2.4 Criar/atualizar documentos essenciais**  
  `CONTRIBUTING.md`, templates de PR, código de conduta, etc.

- [x] **2.5 Configurar GitHub Pages para site de documentação**  
  Workflow `pages.yml` criado e integrado ao conteúdo de docs.

- [x] **2.6 Criar Actions avançadas (validar PRs, auto-relatórios, IA)**  
  Workflows de auto-label, weekly report, integração de IA e demais automações.

---

## FASE 3 – Manutenção, Segurança e Escalabilidade

 - [x] **3.1 Configurar branches protegidas (`main` / `develop`)**  
  - Definir regras de proteção (mínimo de approvals, status checks obrigatórios, squash/rebase/merge policy).  
  - Restringir `force-push` e pushes diretos em `main` conforme política.

 - [x] **3.2 Ativar e refinar segurança avançada (Dependabot, CodeQL, Secret Scanning)**  
  - Dependabot já ativado: revisar intervalos, escopo e regras de auto-merge (se desejado).  
  - CodeQL: garantir configuração apropriada para Go + workflows relevantes.  
  - Secret scanning: confirmar escopo e alertas.

- [x] **3.3 Configurar monitoramento via Insights e dashboards internos**  
  - Selecionar gráficos principais (issues abertas/fechadas, tempo de review, etc.).  
  - Padronizar o uso de Insights como fonte de métricas de saúde do projeto.

- [ ] **3.4 Criar integrações externas (Slack / Discord / Email)**  
  - Webhook/integração para: issues críticas, falhas de CI, novas releases.  
  - Integração com Discussions (se fizer sentido, ex.: canal #discussions).  
  - Notificações por email para milestones importantes (se aplicável).

- [x] **3.5 Implementar acessibilidade e boas práticas**  
  - Diretrizes de revisão (checklist de acessibilidade e qualidade em PRs).  
  - Ajustes em templates de issue/PR para reforçar boas práticas.  
  - Política de triagem e rotulagem consistente.

---

## FASE 4 – Lançamento Inicial e Iteração Contínua

 - [x] **4.1 Criar Release v0.1 "GitHub Hub Initialized"**  
  - Criar tag `v0.1.0`.  
  - Gerar changelog consolidando tudo que esta iniciativa entregou.  
  - Incluir assets: documentação, diagramas, guias de setup/colaboração IA.

- [x] **4.2 Configurar monitoramento de adoção**  
  - Acompanhar crescimento de stars, forks e contribuidores.  
  - Definir cadência para revisão dessas métricas (ex.: quinzenal).  
  - Criar visão/relatório simples (manual ou automatizado) para essas métricas.

- [x] **4.3 Estabelecer ciclos de iteração semanais**  
  - Definir cadência de sprints (ex.: semanal ou quinzenal) usando o Project "Shantilly Roadmap".  
  - Criar rotina de grooming/planning/review vinculada a milestones e ao Project.  
  - Atualizar roadmap em datas pré-definidas (ex.: ao final de cada sprint).

---

## EXTRAS – Automação com IA e Colaboração Híbrida

> Estes itens já estão em grande parte implementados; ficam aqui como referência
> para manutenção e evolução da camada de IA.

- [x] **E1. Workflows de IA no GitHub Actions**  
  - Análise automática de PRs (AI Code Analysis).  
  - Auto-labeling inteligente de issues.  
  - Security scan guiado por IA.  
  - Geração automática de documentação.

- [x] **E2. Comandos IA de colaboração (`@ai-*`)**  
  - `@ai-review` para review automático.  
  - `@ai-test` para execução de testes.  
  - `@ai-docs` para docs.  
  - `@ai-security` para scans de segurança.

- [x] **E3. Evolução dos fluxos IA-humanos (contínuo)**  
  - Refinar prompts/fluxos com base em uso real.  
  - Ajustar automações para evitar ruído (ex.: comentários excessivos em PRs).  
  - Revisar periodicamente métricas de eficácia da automação.

---

## COMO USAR ESTE ARQUIVO

- Use **apenas este arquivo** como controle de tarefas de GitHub Optimization.  
- Quando concluir uma atividade, marque o checkbox correspondente aqui.  
- Os arquivos antigos (`todo_*`) passam a ser apenas histórico de design/planejamento; o status
  oficial deve ser atualizado somente neste arquivo.
