# Checklist de Pull Requests - Shantilly

Este checklist deve ser seguido **antes de aprovar um PR** para `develop` ou `main`.

## 1. Escopo e propósito

- [ ] O título do PR é **claro e descritivo**.
- [ ] A descrição explica **o problema** e **a solução proposta**.
- [ ] O PR está focado em **um escopo razoável** (evitar PRs gigantes e mistos).

## 2. Branch de origem

- [ ] A branch segue a convenção:
  - `feat/nome-da-feature` para novas funcionalidades.
  - `fix/descricao-do-bug` para correções.
- [ ] A branch foi criada a partir de `develop` (ou de `main` no caso de hotfix crítico).

## 3. Código e arquitetura

- [ ] Mudanças respeitam a **estratégia de branches** descrita em `github-branches-strategy.md`.
- [ ] Não há violações claras de padrões arquiteturais (ex.: evitar `os.Exit` fora de `main`, etc.).
- [ ] Nomes de funções, arquivos e pacotes são consistentes com o restante do projeto.
- [ ] Não há código comentado ou lixo que deveria ter sido removido.

## 4. Testes

- [ ] Foram executados os testes locais relevantes:
  - [ ] `go test ./...`
  - [ ] Smoke test de runtime, quando aplicável:
        `./shantilly runtime --file app.example.yaml`
- [ ] Novos comportamentos críticos possuem **testes automatizados** quando fizer sentido.
- [ ] O PR **não quebra** testes existentes.

## 5. Documentação e comunicação

- [ ] Documentação foi atualizada quando necessário:
  - [ ] `README.md`
  - [ ] Documentos em `docs/pt-br/` e `docs/en/` relevantes.
- [ ] Changelog / notas de release foram atualizados (quando o PR impacta release).
- [ ] Comentários no código explicam apenas o que é realmente não óbvio.

## 6. Impacto em CI/CD

- [ ] Workflows do GitHub Actions foram revisados se o PR mexe em `.github/workflows/`.
- [ ] Verificado se o PR não quebra build, lint ou deploy em ambientes automatizados.

## 7. Revisão de qualidade

- [ ] O diff foi lido por pelo menos **um revisor** (pode ser o próprio autor em projetos solo, mas com atenção).
- [ ] Não há "surpresas" no PR (arquivos grandes adicionados sem necessidade, secrets, etc.).
- [ ] Para long-running features, foi considerada a estratégia de:
  - Quebrar em vários PRs menores quando possível.
  - Usar feature flags para código parcialmente implementado.

## 8. Regras específicas para `main`

Antes de aprovar um PR com base `main` (release ou hotfix):

- [ ] O PR foi originado de um ponto **estável** de `develop` (exceto hotfix pontual).
- [ ] Todos os testes passaram em ambiente CI.
- [ ] A tag de versão planejada (ex.: `v0.2.1`) foi definida e documentada.
- [ ] Há uma descrição clara do impacto para usuários finais.

Seguir este checklist ajuda a manter o repositório **saudável**, previsível e com histórico claro de decisões.
