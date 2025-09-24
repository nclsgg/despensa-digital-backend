# 🎯 Seleção de Provedor LLM nas Requisições

## Visão Geral
Agora você pode escolher qual provedor LLM usar diretamente em cada requisição, sem precisar alterar o provedor ativo globalmente. Isso permite usar diferentes provedores para diferentes tipos de tarefas.

## 🚀 Como Usar

### 1. Chat Simples com Seleção de Provedor

**Endpoint**: `POST /api/v1/llm/chat`

**Exemplo usando OpenAI:**
```json
{
  "message": "Olá! Como você pode me ajudar com receitas?",
  "provider": "openai",
  "context": "user_interaction"
}
```

**Exemplo usando Gemini:**
```json
{
  "message": "Sugira uma receita com frango e arroz",
  "provider": "gemini",
  "context": "recipe_suggestion"
}
```

**Exemplo sem especificar provedor (usa o ativo):**
```json
{
  "message": "Como fazer um bolo de chocolate?",
  "context": "recipe_request"
}
```

### 2. Requisição LLM Avançada com Seleção de Provedor

**Endpoint**: `POST /api/v1/llm/process`

**Exemplo com configurações avançadas:**
```json
{
  "messages": [
    {
      "role": "system",
      "content": "Você é um chef especialista em culinária brasileira"
    },
    {
      "role": "user", 
      "content": "Crie uma receita de feijoada para 6 pessoas"
    }
  ],
  "provider": "gemini",
  "max_tokens": 1500,
  "temperature": 0.8,
  "top_p": 0.9
}
```

### 3. Geração de Receitas com Provedor Específico

**Endpoint**: `POST /api/v1/recipes/generate`

```json
{
  "pantry_id": "123e4567-e89b-12d3-a456-426614174000",
  "provider": "openai",
  "cooking_time": 45,
  "meal_type": "lunch",
  "difficulty": "medium",
  "dietary_restrictions": ["vegetarian"]
}
```

## 📋 Provedores Disponíveis

| Provedor | Valor | Status |
|----------|--------|--------|
| OpenAI | `openai` | ✅ Implementado |
| Google Gemini | `gemini` | ✅ Implementado |
| Anthropic Claude | `anthropic` | 🚧 Em desenvolvimento |
| Ollama (Local) | `ollama` | 🚧 Planejado |

## 🎛️ Comportamento do Sistema

### Quando `provider` é especificado:
1. ✅ Sistema usa o provedor especificado
2. ✅ Retorna erro se provedor não estiver configurado
3. ✅ Ignora o provedor ativo global
4. ✅ Inclui informação do provedor usado na resposta

### Quando `provider` é omitido:
1. ✅ Sistema usa o provedor ativo atual
2. ✅ Retorna erro se nenhum provedor estiver ativo
3. ✅ Comportamento compatível com versões anteriores

## 💡 Exemplos de Casos de Uso

### Comparação de Provedores
```bash
# Teste a mesma pergunta com diferentes provedores
curl -X POST http://localhost:3030/api/v1/llm/chat \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -d '{
    "message": "Explique como fazer risoto de camarão",
    "provider": "openai"
  }'

curl -X POST http://localhost:3030/api/v1/llm/chat \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -d '{
    "message": "Explique como fazer risoto de camarão", 
    "provider": "gemini"
  }'
```

### Otimização de Custos
```json
// Use Gemini para tarefas simples (gratuito)
{
  "message": "Liste 5 ingredientes básicos para uma salada",
  "provider": "gemini"
}

// Use OpenAI para tarefas complexas
{
  "messages": [
    {
      "role": "system",
      "content": "Analise nutricionalmente esta receita..."
    },
    {
      "role": "user",
      "content": "Receita complexa aqui..."
    }
  ],
  "provider": "openai",
  "temperature": 0.3
}
```

### Fallback Automático
```javascript
// Exemplo de implementação no frontend
async function askLLM(message) {
  try {
    // Tenta primeiro com Gemini (gratuito)
    return await fetch('/api/v1/llm/chat', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        message: message,
        provider: 'gemini'
      })
    });
  } catch (error) {
    // Fallback para OpenAI se Gemini falhar
    return await fetch('/api/v1/llm/chat', {
      method: 'POST', 
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        message: message,
        provider: 'openai'
      })
    });
  }
}
```

## 📊 Respostas do Sistema

### Resposta Chat Simples
```json
{
  "success": true,
  "data": {
    "response": "Para fazer risoto de camarão, você vai precisar...",
    "provider": "gemini",
    "model": "gemini-1.5-flash",
    "usage": {
      "prompt_tokens": 25,
      "completion_tokens": 150,
      "total_tokens": 175
    }
  }
}
```

### Resposta LLM Avançada
```json
{
  "success": true,
  "data": {
    "id": "llm-req-123456789",
    "response": "Receita detalhada aqui...",
    "model": "gpt-3.5-turbo",
    "usage": {
      "prompt_tokens": 45,
      "completion_tokens": 320,
      "total_tokens": 365
    },
    "metadata": {
      "used_provider": "openai"
    }
  }
}
```

## 🔧 Configuração de Provedores

Para usar múltiplos provedores, configure as chaves API:

```env
# .env
OPENAI_API_KEY=sk-proj-seu-openai-key
GEMINI_API_KEY=seu-gemini-key
```

O sistema automaticamente detectará e configurará ambos os provedores.

## 🎯 Vantagens

✅ **Flexibilidade**: Escolha o melhor provedor para cada tarefa  
✅ **Otimização de Custos**: Use provedores gratuitos quando possível  
✅ **Redundância**: Fallback automático se um provedor falhar  
✅ **Testes A/B**: Compare respostas de diferentes provedores  
✅ **Compatibilidade**: Funciona com código existente  

## 🚨 Tratamento de Erros

```json
// Provedor não configurado
{
  "success": false,
  "error": {
    "code": "INTERNAL_SERVER_ERROR",
    "message": "Erro ao processar requisição LLM: provedor 'anthropic' não está configurado"
  }
}

// Provedor inválido
{
  "success": false,
  "error": {
    "code": "BAD_REQUEST", 
    "message": "Dados de entrada inválidos: provider deve ser um dos: openai, gemini, anthropic, ollama"
  }
}
```

## 🔄 Migração do Código Existente

Seu código atual continuará funcionando sem mudanças:

```json
// ✅ Ainda funciona (usa provedor ativo)
{
  "message": "Como fazer um bolo?"
}

// ✨ Nova funcionalidade (escolhe provedor)
{
  "message": "Como fazer um bolo?",
  "provider": "gemini"
}
```