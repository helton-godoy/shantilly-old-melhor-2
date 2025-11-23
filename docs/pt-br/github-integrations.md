# Integrações Externas – Slack / Discord / Email (pt-BR)

Este documento registra a **arquitetura de integrações externas** planejada para o repositório Shantilly e como avançar rapidamente quando você decidir ativar cada canal.

## 1. Objetivos (3.4)

- Receber notificações sobre:
  - Issues críticas (ex.: `priority::critical`).
  - Falhas de CI importantes.
  - Novas releases publicadas.
- Suportar múltiplos canais:
  - **Email** (via notificações nativas do GitHub).
  - **Slack** (via Webhook).
  - **Discord** (via Webhook).

> Por enquanto **apenas email** está realmente em uso (via sistema nativo do GitHub). 
> Slack e Discord têm workflows esqueleto preparados, mas **não enviam nada** até que você configure os webhooks e ajuste os passos.

---

## 2. Email (via notificações do GitHub)

### 2.1 Como funciona

Em vez de configurar SMTP próprio, a arquitetura assume o uso de:

- **Notificações nativas do GitHub** por email
  - Para issues/PRs com label `priority::critical`.
  - Para releases publicadas.
  - Para falhas de CI que resultem em issues/PRs com comentários.

### 2.2 O que já está pronto

- Estrutura de labels, milestones e releases já faz com que:
  - Qualquer pessoa que esteja **watching** o repositório receba email de:
    - Novas issues e PRs.
    - Releases `v0.1.0` em diante.

### 2.3 Como ativar para você

1. No GitHub, na página do repositório:
   - Clique em **Watch** → selecione o nível desejado (ex.: `All activity`).
2. Em **Settings → Notifications** (no seu perfil GitHub):
   - Garanta que **Email** está habilitado para issues/PRs/releases.

> Futuro: se quiser, podemos adicionar workflows que criem issues/resumos específicos 
> para certos eventos (por exemplo, falhas críticas de CI), e você receberá o email 
> desses eventos automaticamente via GitHub.

---

## 3. Slack – Workflow esqueleto

### 3.1 Conceito

- Usar um Webhook do Slack (configurado como segredo `SLACK_WEBHOOK_URL`).
- Workflow `.github/workflows/external-notifications.yml` já preparado com um job
  `notify_slack` que hoje **apenas faz `echo`**, sem envio real.

### 3.2 Como ativar no futuro

1. Criar um **Incoming Webhook** no workspace Slack.
2. Adicionar o segredo no repositório:
   - `Settings → Secrets and variables → Actions → New repository secret`.
   - Nome sugerido: `SLACK_WEBHOOK_URL`.
3. Editar o job `notify_slack` em
   `.github/workflows/external-notifications.yml` para usar uma ação de Slack, como:

   ```yaml
   - name: Send Slack notification
     if: ${{ secrets.SLACK_WEBHOOK_URL != '' }}
     uses: slackapi/slack-github-action@v1
     with:
       payload: '{"text":"..."}'
     env:
       SLACK_WEBHOOK_URL: ${{ secrets.SLACK_WEBHOOK_URL }}
   ```

### 3.3 Eventos planejados

O workflow esqueleto foi desenhado para reagir a:

- `workflow_run` (para workflows como **Build**, **Lint**, **Release**).
- `release` (tipo `published`).
- `issues` (quando labels como `priority::critical` forem aplicadas).

Você poderá ajustar quais eventos realmente devem mandar mensagem.

---

## 4. Discord – Workflow esqueleto

### 4.1 Conceito

- Similar ao Slack, mas usando um **Discord Webhook**.
- Workflow `.github/workflows/external-notifications.yml` tem um job `notify_discord`
  que hoje apenas faz `echo`.

### 4.2 Como ativar no futuro

1. Criar um Webhook em um canal do Discord.
2. Adicionar o segredo no repositório:
   - Nome sugerido: `DISCORD_WEBHOOK_URL`.
3. Ajustar o job `notify_discord` para usar uma action de Discord ou um `curl` simples:

   ```yaml
   - name: Send Discord notification
     if: ${{ secrets.DISCORD_WEBHOOK_URL != '' }}
     run: |
       curl -H "Content-Type: application/json" \
            -d '{"content":"..."}' \
            ${{ secrets.DISCORD_WEBHOOK_URL }}
   ```

---

## 5. Resumo da arquitetura de integrações

- **Email**
  - Usa apenas infraestrutura nativa do GitHub (watch + notifications).
  - Nenhum segredo adicional necessário.

- **Slack / Discord**
  - Workflows esqueleto prontos em `.github/workflows/external-notifications.yml`.
  - Dependem exclusivamente de você configurar
    - webhooks nos serviços, e
    - secrets `SLACK_WEBHOOK_URL` / `DISCORD_WEBHOOK_URL` no repositório.

Assim que você tiver clareza sobre como quer usar Slack/Discord, basta:

1. Configurar os webhooks e secrets.
2. Substituir os passos `echo` do workflow por chamadas reais às actions/URLs.
3. Opcionalmente ajustar os eventos (`on:`) para o que fizer mais sentido.
