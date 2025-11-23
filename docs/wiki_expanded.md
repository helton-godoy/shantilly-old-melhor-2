# Wiki Expandida

Esta página contém a documentação ampliada do projeto Shantilly, integrando recursos do GitHub como Issues, Discussions e Actions. Abaixo estão descritos os fluxos de trabalho principais e como utilizá‑los.

## Integração com Issues

- Cada issue deve seguir o template padrão definido em `.github/ISSUE_TEMPLATE`.
- Labels hierárquicas (`priority::critical`, `type::feature`, `area::core`) são aplicadas automaticamente via o workflow `auto-label.yml`.

## Discussions

- As categorias configuradas são **Ideas**, **Q&A**, **Polls** e **General**.
- Use a tag `discussion` ao criar novas discussões para que o bot de triagem as classifique.

## Actions

- O workflow `deploy-pages.yml` publica o site estático.
- O workflow `auto-label.yml` rotula issues conforme padrões.

## Diagrama de Fluxo

Abaixo está o diagrama Mermaid que ilustra o fluxo completo de desenvolvimento, desde a criação de uma issue até a publicação no GitHub Pages.

```mermaid
flowchart TD
    A[Ideia] --> B[Issue criada]
    B --> C[Triagem automática]
    C --> D[Pull Request aberto]
    D --> E[Revisão + Aprovações]
    E --> F[Merge na branch develop]
    F --> G[Deploy de preview (Vercel/Netlify)]
    G --> H[Merge na branch main]
    H --> I[Deploy GitHub Pages]
    I --> J[Site publicado]
```

---
*Esta página será mantida atualizada conforme novas integrações e processos forem adicionados.*
