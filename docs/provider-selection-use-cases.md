# 🎯 Casos de Uso Práticos - Seleção de Provedores LLM

## 📋 Cenários Reais de Uso

### 1. 💰 Desenvolvimento Econômico

**Situação**: Desenvolvendo e testando funcionalidades sem custos

**Configuração**:
```json
{
  "llm_provider": "gemini"
}
```

**Endpoints a usar**:
- `POST /api/v1/llm/chat` com `"provider": "gemini"`
- `POST /api/v1/recipes/generate` com `"provider": "gemini"`

**Vantagens**:
- ✅ Gratuito até o limite
- ✅ Resposta rápida
- ✅ Boa qualidade para testes

---

### 2. 🎯 Produção com Alta Qualidade

**Situação**: Sistema em produção que precisa da melhor qualidade

**Configuração**:
```json
{
  "llm_provider": "openai"
}
```

**Exemplo de requisição**:
```json
POST /api/v1/llm/process
{
  "messages": [
    {
      "role": "system", 
      "content": "Você é um chef michelin especialista em culinária francesa"
    },
    {
      "role": "user",
      "content": "Crie um menu completo para um jantar romântico"
    }
  ],
  "provider": "openai",
  "temperature": 0.8,
  "max_tokens": 2000
}
```

---

### 3. ⚖️ Comparação A/B de Qualidade

**Situação**: Quer comparar a qualidade das respostas para escolher o melhor

**Teste 1 - OpenAI**:
```json
POST /api/v1/llm/chat
{
  "message": "Explique como fazer carbonara autêntica italiana",
  "provider": "openai",
  "context": "recipe_detailed"
}
```

**Teste 2 - Gemini**:
```json
POST /api/v1/llm/chat
{
  "message": "Explique como fazer carbonara autêntica italiana",
  "provider": "gemini", 
  "context": "recipe_detailed"
}
```

**Análise**:
- Compare detalhamento
- Avalie precisão técnica
- Verifique autenticidade cultural
- Analise clareza das instruções

---

### 4. 🔄 Fallback Automático

**Situação**: Implementar redundância no frontend

**Estratégia**:
1. **Primeira tentativa**: Gemini (mais barato)
2. **Fallback**: OpenAI (se Gemini falhar)

**Implementação no Frontend**:
```javascript
async function generateRecipe(recipeData) {
  // Tenta primeiro com Gemini
  try {
    const response = await fetch('/api/v1/recipes/generate', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        ...recipeData,
        provider: 'gemini'
      })
    });
    
    if (response.ok) {
      return await response.json();
    }
  } catch (error) {
    console.log('Gemini falhou, tentando OpenAI...');
  }
  
  // Fallback para OpenAI
  const response = await fetch('/api/v1/recipes/generate', {
    method: 'POST', 
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      ...recipeData,
      provider: 'openai'
    })
  });
  
  return await response.json();
}
```

---

### 5. 🎨 Personalização por Tipo de Conteúdo

**Situação**: Usar diferentes provedores para diferentes tipos de tarefa

#### Para Receitas Criativas:
```json
POST /api/v1/llm/chat
{
  "message": "Crie uma receita fusion japonês-brasileiro inovadora",
  "provider": "openai",
  "context": "creative_cooking"
}
```

#### Para Receitas Tradicionais:
```json
POST /api/v1/llm/chat
{
  "message": "Como fazer feijoada tradicional brasileira",
  "provider": "gemini",
  "context": "traditional_recipe"
}
```

#### Para Análise Nutricional:
```json
POST /api/v1/llm/process
{
  "messages": [
    {
      "role": "system",
      "content": "Você é um nutricionista especializado"
    },
    {
      "role": "user", 
      "content": "Analise o valor nutricional desta receita: [receita]"
    }
  ],
  "provider": "openai",
  "temperature": 0.3
}
```

---

### 6. 📊 Testes de Performance

**Situação**: Medir e comparar performance entre provedores

**Teste de Latência**:
```json
// Medir tempo de resposta
const startTime = Date.now();

await fetch('/api/v1/llm/chat', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    message: "Receita rápida para o jantar",
    provider: "gemini"
  })
});

const geminiTime = Date.now() - startTime;
```

**Teste de Qualidade**:
```json
// Mesmo prompt, provedores diferentes
const prompts = [
  {
    provider: "openai",
    message: "Crie uma receita de bolo de chocolate para diabéticos"
  },
  {
    provider: "gemini", 
    message: "Crie uma receita de bolo de chocolate para diabéticos"
  }
];
```

---

### 7. 🌍 Otimização por Região/Idioma

**Situação**: Alguns provedores podem ser melhores para conteúdo específico

#### Culinária Brasileira:
```json
POST /api/v1/llm/chat
{
  "message": "Receitas tradicionais do Nordeste brasileiro",
  "provider": "gemini",
  "context": "regional_cuisine_brazil"
}
```

#### Culinária Internacional:
```json
POST /api/v1/llm/chat
{
  "message": "Authentic French coq au vin recipe",
  "provider": "openai",
  "context": "international_cuisine"
}
```

---

### 8. 💡 Desenvolvimento Iterativo

**Situação**: Refinamento progressivo de receitas

**Fase 1 - Brainstorm** (Gemini - gratuito):
```json
POST /api/v1/llm/chat
{
  "message": "5 ideias de receitas com frango e legumes",
  "provider": "gemini"
}
```

**Fase 2 - Refinamento** (OpenAI - precisão):
```json
POST /api/v1/llm/process
{
  "messages": [
    {
      "role": "system",
      "content": "Você é um chef profissional"
    },
    {
      "role": "user",
      "content": "Refine esta receita: [receita escolhida]"
    }
  ],
  "provider": "openai",
  "temperature": 0.7
}
```

---

## 🔧 Configuração no Postman

### Ambiente de Desenvolvimento:
```json
{
  "llm_provider": "gemini",
  "environment": "development"
}
```

### Ambiente de Produção:
```json
{
  "llm_provider": "openai",
  "environment": "production"
}
```

### Ambiente de Teste:
```json
{
  "llm_provider": "{{provider_variable}}",
  "environment": "testing"
}
```

---

## 📈 Métricas para Monitorar

### Por Provedor:
- **Taxa de sucesso** (%)
- **Tempo médio de resposta** (ms)
- **Custo por requisição** ($)
- **Satisfação do usuário** (rating)

### Por Tipo de Conteúdo:
- **Receitas simples**: Qual provedor é mais eficiente?
- **Receitas complexas**: Qual tem melhor qualidade?
- **Análises nutricionais**: Qual é mais preciso?

---

## 🎯 Recomendações Finais

1. **Comece com Gemini** para tudo
2. **Identifique casos** onde OpenAI é superior
3. **Implemente fallback** para alta disponibilidade
4. **Monitor custos** constantemente
5. **Colete feedback** dos usuários sobre qualidade

**Resultado**: Sistema flexível, econômico e confiável! 🚀