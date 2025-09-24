# 📋 Guia da Coleção Postman - Despensa Digital LLM API

## 🎯 Visão Geral

Esta coleção foi atualizada para incluir a nova funcionalidade de **seleção de provedor LLM** nas requisições. Agora você pode escolher entre OpenAI, Gemini ou outros provedores diretamente nas requisições, sem precisar alterar o provedor ativo globalmente.

## ⚙️ Configuração Inicial

### Variáveis da Coleção

Configure estas variáveis na coleção:

| Variável | Valor Padrão | Descrição |
|----------|-------------|-----------|
| `baseUrl` | `http://localhost:8080` | URL base da API |
| `jwt_token` | `your_jwt_token_here` | Token JWT obtido no login |
| `pantry_id` | `123e4567-e89b-12d3-a456-426614174000` | ID da despensa |
| `llm_provider` | `openai` | Provedor LLM padrão |

### Como Configurar

1. **Importe a coleção** no Postman
2. **Configure as variáveis**:
   - Clique com botão direito na coleção → "Edit"
   - Vá para aba "Variables"
   - Configure os valores necessários
3. **Obtenha o JWT Token**:
   - Execute o endpoint "Authentication → Login"
   - Copie o token da resposta
   - Cole na variável `jwt_token`

## 🚀 Funcionalidades Principais

### 1. Seleção de Provedor por Requisição

#### Chat Simples com Provedor Específico
```json
POST /api/v1/llm/chat
{
  "message": "Como fazer risoto de camarão?",
  "provider": "gemini",  // ou "openai"
  "context": "recipe_request"
}
```

#### Requisição LLM Avançada com Provedor
```json
POST /api/v1/llm/process
{
  "messages": [...],
  "provider": "openai",
  "temperature": 0.7,
  "max_tokens": 1000
}
```

### 2. Configuração de Provedores

#### Configurar OpenAI
- Endpoint: `Provider Configuration → Configure OpenAI Provider`
- Necessário: API key da OpenAI

#### Configurar Gemini
- Endpoint: `Provider Configuration → Configure Gemini Provider`
- Necessário: API key do Google AI Studio

### 3. Testes e Comparações

#### Comparar Respostas
Use os endpoints em "Quick Tests" para comparar respostas entre provedores:
- "Simple Chat - OpenAI" vs "Simple Chat - Gemini"
- Mesmo prompt, provedores diferentes

## 📁 Estrutura da Coleção

### 1. **Authentication**
- Login para obter JWT token

### 2. **Recipe Generation**
- Gerar receitas com provedor específico
- Buscar receitas por ingredientes

### 3. **LLM Direct**
- Chat simples com seleção de provedor
- Requisições LLM avançadas
- Construção de prompts
- Status dos provedores

### 4. **Provider Configuration**
- ✨ **NOVA SEÇÃO**
- Configurar OpenAI e Gemini
- Alternar entre provedores
- Testar conectividade

### 5. **Quick Tests**
- ✨ **ATUALIZADA**
- Testes rápidos com comparação
- OpenAI vs Gemini lado a lado

## 🔧 Exemplos de Uso

### Cenário 1: Desenvolvimento com Gemini (Gratuito)

1. Configure o Gemini:
   ```
   Provider Configuration → Configure Gemini Provider
   ```

2. Altere a variável `llm_provider` para `gemini`

3. Teste o chat:
   ```
   Quick Tests → Simple Chat - Gemini
   ```

### Cenário 2: Produção com OpenAI

1. Configure o OpenAI:
   ```
   Provider Configuration → Configure OpenAI Provider
   ```

2. Altere a variável `llm_provider` para `openai`

3. Teste funcionalidades avançadas:
   ```
   LLM Direct → Advanced LLM Request with Provider
   ```

### Cenário 3: Comparação A/B

1. Execute o mesmo prompt nos dois provedores:
   ```
   Quick Tests → Simple Chat - OpenAI
   Quick Tests → Simple Chat - Gemini
   ```

2. Compare as respostas para qualidade e estilo

## 🎨 Dicas de Uso

### Otimização de Custos
- Use **Gemini** para desenvolvimento (gratuito)
- Use **OpenAI** para produção quando necessário
- Teste ambos para encontrar o melhor para cada caso

### Testes de Performance
- Compare tempos de resposta
- Avalie qualidade das respostas
- Teste diferentes temperaturas

### Desenvolvimento Eficiente
1. **Comece com Gemini** para prototipagem
2. **Teste com OpenAI** para validação
3. **Configure fallback** no frontend
4. **Monitor custos** em produção

## 🔍 Troubleshooting

### Erro: "provider não está configurado"
**Solução**: Execute o endpoint de configuração do provedor antes de usar

### Erro: "API key inválida"
**Solução**: Verifique se as chaves foram copiadas corretamente

### Erro: 401 Unauthorized
**Solução**: Atualize o `jwt_token` fazendo login novamente

### Erro: "nenhum provedor ativo configurado"
**Solução**: Configure pelo menos um provedor ou especifique `provider` na requisição

## 📊 Monitoramento

### Variáveis para Debug
- Adicione `console.log` nos scripts de teste
- Use variáveis para rastrear qual provedor foi usado
- Monitor tempos de resposta

### Métricas Importantes
- Taxa de sucesso por provedor
- Tempo médio de resposta
- Qualidade das respostas
- Custos por requisição

## 🆕 O que Mudou

### ✨ Novidades
- ✅ Campo `provider` em requisições de chat
- ✅ Campo `provider` em geração de receitas
- ✅ Seção "Provider Configuration" completa
- ✅ Testes comparativos entre provedores
- ✅ Variável `llm_provider` para facilitar testes

### 🔄 Atualizações
- ✅ Todos os endpoints de chat agora suportam seleção de provedor
- ✅ Geração de receitas com provedor específico
- ✅ Testes rápidos divididos por provedor

### 📈 Benefícios
- **Flexibilidade**: Escolha o provedor por requisição
- **Economia**: Use gratuito para desenvolvimento
- **Qualidade**: Compare e escolha a melhor resposta
- **Confiabilidade**: Implemente fallback entre provedores

---

**Pronto para usar!** 🚀 A coleção está completamente atualizada com a nova funcionalidade de seleção de provedores LLM.