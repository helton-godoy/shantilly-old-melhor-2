# Context7 MCP - Demonstração das Capacidades

## Setup Concluído ✅

O servidor MCP Context7 foi configurado com sucesso no arquivo `cline_mcp_settings.json` com os seguintes detalhes:

```json
{
  "mcpServers": {
    "github.com/upstash/context7-mcp": {
      "command": "npx",
      "args": ["-y", "@upstash/context7-mcp"],
      "disabled": false,
      "autoApprove": []
    }
  }
}
```

## Ferramentas Disponíveis

O Context7 MCP fornece as seguintes ferramentas para LLMs:

### 1. `resolve-library-id`
**Propósito:** Resolve um nome geral de biblioteca em um ID compatível com Context7

**Parâmetros:**
- `libraryName` (obrigatório): O nome da biblioteca para pesquisar

**Exemplo de uso:**
```javascript
// Resolver biblioteca "react"
{
  "name": "resolve-library-id",
  "arguments": {
    "libraryName": "react"
  }
}
```

### 2. `get-library-docs`
**Propósito:** Busca documentação para uma biblioteca usando um ID compatível com Context7

**Parâmetros:**
- `context7CompatibleLibraryID` (obrigatório): ID exato da biblioteca (ex: `/mongodb/docs`, `/vercel/next.js`)
- `topic` (opcional): Focar a documentação em um tópico específico
- `page` (opcional, padrão: 1): Número da página para paginação (1-10)

**Exemplo de uso:**
```javascript
// Obter documentação da Next.js
{
  "name": "get-library-docs",
  "arguments": {
    "context7CompatibleLibraryID": "/vercel/next.js",
    "topic": "routing",
    "page": 1
  }
}
```

## Como Usar Context7 MCP

### Exemplo de Prompt
```
Criar um middleware Next.js que verifica um JWT válido em cookies
e redireciona usuários não autenticados para `/login`. use context7
```

### ID de Biblioteca Direto
Se você já souber exatamente qual biblioteca deseja usar, adicione seu ID Context7 ao prompt:

```
Implementar autenticação básica com Supabase. use library /supabase/supabase para API e docs.
```

## Exemplo Prático de Uso

### 1. Resolver Biblioteca React
```bash
# Comando de teste
npx -y @upstash/context7-mcp --transport stdio
```

### 2. Obter Documentação da Next.js
O ID Context7 seria: `/vercel/next.js`

### 3. Tópicos Específicos
Para documentação focada:
- "routing" - Roteamento
- "hooks" - Hooks do React
- "api-routes" - Rotas da API
- "deployment" - Implantação

## Vantagens do Context7

- ✅ Documentação sempre atualizada (sem dados de treinamento desatualizados)
- ✅ Exemplos de código verificados da fonte
- ✅ APIs específicas por versão
- ✅ Eliminação de APIs hallucinated inexistentes
- ✅ Rate limits generosos para uso básico

## Próximos Passos

Agora que o Context7 MCP está configurado:

1. **Teste no seu MCP client** - As ferramentas estarão disponíveis na seção "Connected MCP Servers"
2. **Adicione uma regra** - Configure seu client para auto-invocar Context7 em perguntas de código
3. **Experimente exemplos práticos** - Use prompts com "use context7" para obter documentação atualizada

## Links Úteis

- **Site oficial:** https://context7.com
- **Dashboard:** https://context7.com/dashboard (para API key se necessário)
- **Documentação do projeto:** https://context7.com/docs/adding-libraries
- **Repositório:** https://github.com/upstash/context7-mcp

---

*Configuração realizada em: 24/11/2025*  
*Sistema: Linux 6.14*  
*Node.js: v24.11.0*
