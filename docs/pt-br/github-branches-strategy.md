# Estratégia de Branches - Shantilly

## Objetivos

- Manter o repositório **limpo e previsível**.
- Garantir que apenas código **estável** chegue em `main`.
- Permitir desenvolvimento contínuo em `develop` com segurança.
- Usar branches curtas e específicas para features e correções.

## Branches Permanentes

### `main`

- Representa o **estado de produção** / releases estáveis.
- Regras:
  - Apenas código **testado e aprovado**.
  - Entradas em `main` devem ocorrer via **PRs de release** (ex: `integration-v2.0`) ou hotfixes muito bem justificados.
  - Tags de versão (ex: `v0.2.0`) devem apontar para commits em `main`.

### `develop`

- É o **tronco de desenvolvimento contínuo**.
- Sempre parte de um ponto estável de `main` (por exemplo, logo após um release).
- Regras:
  - Novas features e correções partem de `develop`.
  - `develop` pode receber merges de branches `feat/*` e `fix/*` via PR.
  - Periodicamente, quando estável, o estado de `develop` é promovido para `main` como um novo release.

## Branches Temporárias

### `feat/*`

- Para **novas funcionalidades**.
- Criadas a partir de `develop`.
- Exemplo: `feat/runtime-logging`, `feat/new-component-select`.
- Fluxo:
  1. `git checkout develop`
  2. `git checkout -b feat/nome-da-feature`
  3. Implementação + testes locais.
  4. Abrir PR `feat/nome-da-feature` → `develop`.
  5. Após o merge, **apagar** a branch remota e local.

### `fix/*`

- Para **correções pontuais** (bugs, ajustes de CI, etc.).
- Em geral também partem de `develop`.
- Exemplo: `fix/runtime-panic-empty-config`, `fix/weekly-report-permissions`.
- Fluxo:
  1. `git checkout develop`
  2. `git checkout -b fix/descreve-o-bug`
  3. Corrigir, testar (`go test ./...`, smoke tests relevantes).
  4. Abrir PR `fix/...` → `develop` (ou → `main` em caso de hotfix crítico).
  5. Após o merge, **apagar** a branch.

## Política de Integração

1. **Code Review obrigatório** para merges em `develop` e `main`.
2. **Testes automáticos** devem rodar em PRs:
   - `go test ./...`
   - Smoke tests de runtime (ex.: `./shantilly runtime --file app.example.yaml`).
3. **`main` nunca recebe commits diretos** (apenas via PR/release ou hotfixes muito pontuais com justificativa).
4. **`develop` é a base para trabalho diário**, mas deve ser mantida em estado aceitavelmente estável.

Para um checklist detalhado de revisão, ver:

- [Checklist de Pull Requests (PT-BR)](github-pr-checklist.md)

## Limpeza de Branches

- Branches como `docs-i18n-and-advanced-site`, `feat/runtime-migration`, `fix/package-lock`, `fix/weekly-report-permissions` e `fix/package-lock2` foram **avaliadas** e seu conteúdo de valor foi **incorporado** em `main` e/ou `develop`.
- Após integração:
  - As branches antigas foram removidas do remoto para reduzir ruído.
  - O histórico foi preservado via commits e PRs (por exemplo, PR #68 para integração v2.0).

### Regra Prática

- Toda branch `feat/*` ou `fix/*` deve:
  - Nascer de `develop`.
  - Ser integrada via PR.
  - Ser apagada após o merge.

## Releases

- Releases formais devem ser feitos a partir de `main`, usando tags:
  - Exemplo: `v0.2.0` para o release que integra o runtime v2.0 + documentação bilíngue.
- Recomenda-se manter um arquivo de notas de release (ex.: `RELEASES.md`) apontando para as principais mudanças por versão.

## Long-running features

Em vez de manter branches gigantes e de longa duração (que acumulam conflitos e dificultam revisão), a estratégia preferida é:

- **Quebrar features grandes em incrementos menores**, cada um entregue em uma branch `feat/*` específica.
- Sempre que possível, usar **feature flags** ou configurações que permitam:
  - Mesclar código parcialmente implementado em `develop`.
  - Manter a funcionalidade desativada por padrão até estar pronta.
- Quando várias features grandes precisam ser coordenadas:
  - Criar, se necessário, uma **branch de integração temporária** a partir de `develop` (ex.: `integration/epic-1-runtime-v2`).
  - Mesclar nela branches `feat/*` relacionadas.
  - Após estabilizar, promover o resultado para `develop` (via PR) e **apagar** a branch de integração.

Regras importantes para long-running features:

- Evitar trabalhar por semanas em uma branch isolada sem atualizar de `develop`.
- Preferir ciclos curtos: pequenos PRs, revisões rápidas, feedback contínuo.
- Usar sempre PRs para visibilidade histórica e discussão de decisões.

## Diagrama do fluxo de branches

```mermaid
flowchart LR
    A[main\n(releases estáveis)] --> B[develop\n(desenvolvimento contínuo)]

    B --> C[feat/nova-feature]
    C --> B

    B --> D[fix/bug-específico]
    D --> B

    B --> E[integration/epic-X\n(branch temporária)]
    C --> E
    D --> E
    E --> B
```

## Resumo

- **`main`**: produção / releases estáveis.
- **`develop`**: desenvolvimento contínuo, sempre baseado em um ponto estável de `main`.
- **`feat/*` e `fix/*`**: branches curtas, específicas, sempre integradas via PR e removidas após o merge.
- Histórico preservado em **commits, PRs e tags**, não em branches long-lived desnecessárias.
