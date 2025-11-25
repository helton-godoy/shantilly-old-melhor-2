# 🌍 Estratégia de Internacionalização - Shantilly

## 📋 VISÃO GERAL

O projeto Shantilly adota uma estratégia **bilíngue** completa (inglês + português brasileiro) para:
- Valorizar a origem brasileira do projeto
- Acessibilidade para desenvolvedores brasileiros
- Expansão da comunidade global

## 🎯 PRINCÍPIOS FUNDAMENTAIS

### 1. **Suporte Bilíngue Padrão**
- ✅ **Inglês**: Idioma primário (padrão global)
- ✅ **Português Brasileiro**: Idioma secundário (valorização da origem)

### 2. **Paridade de Conteúdo**
- ✅ Todo documento principal deve ter tradução
- ✅ Discussões podem usar ambos idiomas
- ✅ Commits podem ser em qualquer idioma

### 3. **Automação Inteligente**
- ✅ GitHub Actions para auto-tradução
- ✅ Detecção de mudanças em documentação
- ✅ Geração automática de PRs de tradução

## 📂 ESTRUTURA DE ARQUIVOS

### Padrão de Documentação:
```
📁 docs/
├── 📄 README.md (inglês)
├── 📄 README_pt.md (português)
├── 📁 en-US/ (documentação em inglês)
└── 📁 pt-BR/ (documentação em português)
```

### Discussions Bilíngues:
```
🔗 GitHub Discussions
├── 💡 Ideas / Ideias
├── ❓ Q&A / Perguntas & Respostas  
├── 📊 Polls / Enquetes
└── 💬 General / Geral
```

## 🤖 AUTOMATIZAÇÃO COM IA

### Workflow de Auto-Tradução:
```yaml
# .github/workflows/auto-translate.yml
Triggers:
  - Push em docs/**/*.md
  - Pull Request com documentação

Actions:
  1. Detectar arquivos modificados
  2. Traduzir para português brasileiro
  3. Criar branch de tradução
  4. Gerar PR para revisão
```

### Regras para Agentes IA:
1. **SEMPRE** criar duas versões
2. **Usar** categorias apropriadas (en-US/pt-BR)
3. **Manter** qualidade técnica em ambos idiomas
4. **Priorizar** conteúdo brasileiro quando relevante

## 📊 IMPLEMENTAÇÃO ATUAL

### ✅ CONCLUÍDO:
- Discussions.yml atualizado com categorias bilíngues
- Auto-translate workflow criado
- AI Contract atualizado com regras i18n
- Estrutura para português brasileiro

### ⏳ EM PROGRESSO:
- Configuração de Discussions via interface web
- Importação de issues bilíngues
- Documentação em português

### 📋 PRÓXIMOS PASSOS:
- Ativar Discussions bilíngues
- Traduzir documentação principal
- Configurar GitHub Pages bilíngue
- Importar issues estruturadas

## 🎯 BENEFÍCIOS

### Para Comunidade Brasileira:
- ✅ Acesso facilitado à documentação
- ✅ Participação ativa nas Discussions
- ✅ Valorização da origem do projeto

### Para Expansão Global:
- ✅ Inclusão da maior comunidade de desenvolvedores hispanófona
- ✅ Redução de barreiras de entrada
- ✅ Maior adoção do projeto

### Para Desenvolvedores:
- ✅ Documentação completa em idioma nativo
- ✅ Auto-tradução via IA
- ✅ Participação bilíngue nas Discussions

## 📞 SUPORTE

**Contato para Dúvidas sobre i18n:**
- Discussions: Usar categoria apropriada (Q&A/pt-BR)
- Issues: Tag `i18n` para questões de tradução
- PRs: Sempre incluir ambas versões de documentação

---

**Última Atualização**: 2025-11-19  
**Versão**: 1.0  
**Responsável**: Equipe Shantilly
