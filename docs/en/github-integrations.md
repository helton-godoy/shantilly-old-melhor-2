# External Integrations – Slack / Discord / Email (en)

This document records the planned **external integrations architecture** for the Shantilly repository and how to move fast once you decide to enable each channel.

## 1. Goals (3.4)

- Receive notifications about:
  - Critical issues (for example `priority::critical`).
  - Important CI failures.
  - Newly published releases.
- Support multiple channels:
  - **Email** (via GitHub's native notifications).
  - **Slack** (via Webhook).
  - **Discord** (via Webhook).

> For now, **only email** is effectively used (via GitHub's native notifications).
> Slack and Discord have skeleton workflows prepared, but they **do not send anything**
> until you configure webhooks and adjust the steps.

---

## 2. Email (via GitHub notifications)

### 2.1 How it works

Instead of configuring a custom SMTP server, the architecture relies on:

- **GitHub native notifications** via email
  - For issues/PRs with label `priority::critical`.
  - For published releases.
  - For CI failures that lead to issues/PRs or comments.

### 2.2 What is already in place

- The existing structure of labels, milestones, and releases ensures that:
  - Anyone **watching** the repository will receive email for:
    - New issues and PRs.
    - Releases `v0.1.0` and later.

### 2.3 How to enable it for yourself

1. On the GitHub repository page:
   - Click **Watch** → choose the desired level (for example, `All activity`).
2. In **Settings → Notifications** (your GitHub profile):
   - Ensure **Email** is enabled for issues/PRs/releases.

> In the future, we can add workflows that create specific issues/summaries for
> certain events (for example, critical CI failures), and you'll receive those
> via email using GitHub's native notification system.

---

## 3. Slack – Skeleton workflow

### 3.1 Concept

- Use a Slack Incoming Webhook (stored as a secret `SLACK_WEBHOOK_URL`).
- The workflow `.github/workflows/external-notifications.yml` already contains a
  `notify_slack` job that currently **only echoes information**, without sending
  real messages.

### 3.2 How to enable it later

1. Create an **Incoming Webhook** in your Slack workspace.
2. Add the secret to the repository:
   - `Settings → Secrets and variables → Actions → New repository secret`.
   - Suggested name: `SLACK_WEBHOOK_URL`.
3. Edit the `notify_slack` job in
   `.github/workflows/external-notifications.yml` to use a Slack action, for example:

   ```yaml
   - name: Send Slack notification
     if: ${{ secrets.SLACK_WEBHOOK_URL != '' }}
     uses: slackapi/slack-github-action@v1
     with:
       payload: '{"text":"..."}'
     env:
       SLACK_WEBHOOK_URL: ${{ secrets.SLACK_WEBHOOK_URL }}
   ```

### 3.3 Planned events

The skeleton workflow is designed to react to:

- `workflow_run` (for workflows such as **Build**, **Lint**, **Release**).
- `release` (type `published`).
- `issues` (when labels like `priority::critical` are applied).

You can later tune which events should actually send messages.

---

## 4. Discord – Skeleton workflow

### 4.1 Concept

- Similar to Slack, but using a **Discord Webhook**.
- The workflow `.github/workflows/external-notifications.yml` contains a
  `notify_discord` job that currently only does `echo`.

### 4.2 How to enable it later

1. Create a Webhook in a Discord channel.
2. Add the secret to the repository:
   - Suggested name: `DISCORD_WEBHOOK_URL`.
3. Adjust the `notify_discord` job to use a Discord action or raw `curl`, for example:

   ```yaml
   - name: Send Discord notification
     if: ${{ secrets.DISCORD_WEBHOOK_URL != '' }}
     run: |
       curl -H "Content-Type: application/json" \
            -d '{"content":"..."}' \
            ${{ secrets.DISCORD_WEBHOOK_URL }}
   ```

---

## 5. Integration architecture summary

- **Email**
  - Relies only on GitHub's native notifications (watch + notification settings).
  - No additional secrets required.

- **Slack / Discord**
  - Skeleton workflows ready in `.github/workflows/external-notifications.yml`.
  - They only become active when you configure:
    - webhooks in the external services, and
    - repository secrets `SLACK_WEBHOOK_URL` / `DISCORD_WEBHOOK_URL`.

Once you have a clearer idea of how you want to use Slack/Discord, you only need to:

1. Configure webhooks and secrets.
2. Replace the `echo` steps in the workflow with real Slack/Discord calls.
3. Optionally adjust the `on:` events to match your desired notification patterns.
