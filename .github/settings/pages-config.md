# Configuração do GitHub Pages

Este documento descreve as configurações necessárias no GitHub para habilitar o GitHub Pages com deploy via GitHub Actions.

## Passos de Configuração

### 1. Acessar Configurações do GitHub Pages

1. Acesse o repositório: `https://github.com/helton-godoy/shantilly`
2. Clique em **Settings** (Configurações)
3. No menu lateral esquerdo, clique em **Pages**

### 2. Configurar Source

Na seção **Build and deployment**:

- **Source**: Selecione **GitHub Actions** no dropdown
  - ⚠️ **Importante:** NÃO selecione "Deploy from a branch"
  - ✅ Deve estar: "GitHub Actions"

Essa configuração permite que o workflow `.github/workflows/deploy-pages.yml` faça o deploy automático.

### 3. Custom Domain (Opcional)

Se você tiver um domínio personalizado:

1. Em **Custom domain**, adicione seu domínio (ex: `shantilly.dev`)
2. Marque **Enforce HTTPS**

**Nota:** Para o site padrão do GitHub Pages, deixe em branco.

### 4. Verificação

Após o primeiro deploy bem-sucedido:

- A URL do site será exibida na seção Pages
- Deve ser: `https://helton-godoy.github.io/shantilly/`

## Permissões Necessárias

As permissões já estão configuradas no workflow YAML:

```yaml
permissions:
  contents: read
  pages: write
  id-token: write
```

Essas permissões permitem:

- `contents: read` - Ler código do repositório
- `pages: write` - Escrever no GitHub Pages
- `id-token: write` - Gerar tokens de deployment

## Workflow de Deploy

O arquivo `.github/workflows/deploy-pages.yml` é acionado:

- **Automaticamente:** Quando há push para `main` com mudanças em `site/**`
- **Manualmente:** Via página Actions → "Deploy GitHub Pages" → "Run workflow"

## Troubleshooting

### Erro: Actions not enabled

**Problema:** GitHub Actions não está habilitado no repositório

**Solução:**

1. Settings → Actions → General
2. Em "Actions permissions", selecione "Allow all actions and reusable workflows"
3. Salve

### Erro: Deployment failed with permission denied

**Problema:** Workflow sem permissão para deploy

**Solução:**

1. Settings → Actions → General
2. Em "Workflow permissions", selecione "Read and write permissions"
3. Marque "Allow GitHub Actions to create and approve pull requests"
4. Salve

### Site não carrega (404)

**Problema:** Base path incorreto ou GitHub Pages desabilitado

**Solução:**

1. Verificar `vite.config.js` tem `base: '/shantilly/'`
2. Verificar que GitHub Pages está com source "GitHub Actions"
3. Aguardar 2-3 minutos após deploy bem-sucedido

## Verificação de Deploy

Após commit e push:

```bash
# Via linha de comando
gh run watch

# Ou via browser
# Acesse: https://github.com/helton-godoy/shantilly/actions
```

Deploy bem-sucedido mostrará:

- ✅ build job completo
- ✅ deploy job completo
- 🌐 Site acessível em `https://helton-godoy.github.io/shantilly/`

## Documentação Oficial

- [GitHub Pages Quick Start](https://docs.github.com/en/pages/quickstart)
- [GitHub Actions for Pages](https://docs.github.com/en/pages/getting-started-with-github-pages/configuring-a-publishing-source-for-your-github-pages-site#publishing-with-a-custom-github-actions-workflow)
