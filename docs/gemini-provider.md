# Provedor Gemini - Documentação

## Visão Geral
O provedor Gemini foi adicionado ao sistema de LLM da Despensa Digital, permitindo usar os modelos do Google AI (Gemini) para processamento de linguagem natural.

## Modelos Disponíveis
- `gemini-1.5-flash` (padrão) - Modelo rápido e eficiente
- `gemini-1.5-pro` - Modelo mais avançado para tarefas complexas
- `gemini-1.0-pro` - Modelo anterior, ainda disponível

## Configuração

### Variáveis de Ambiente
Adicione ao seu arquivo `.env`:

```env
# Google Gemini API Key - Obtenha em https://aistudio.google.com/app/apikey  
GEMINI_API_KEY=seu-api-key-aqui
```

### Configuração Automática
O sistema detecta automaticamente a presença da chave API e configura o provedor com as seguintes configurações padrão:

- **Model**: `gemini-1.5-flash`
- **Max Tokens**: `2000`
- **Temperature**: `0.7`
- **Timeout**: `30s`

### Configuração Manual via API

```bash
POST /api/v1/llm/provider/configure
Content-Type: application/json

{
  "name": "gemini",
  "provider": "gemini",
  "api_key": "seu-api-key-aqui",
  "model": "gemini-1.5-flash",
  "max_tokens": 2000,
  "temperature": 0.7,
  "timeout": "30s"
}
```

## Exemplos de Uso

### Chat Básico com Gemini

```bash
POST /api/v1/llm/chat
Content-Type: application/json

{
  "message": "Olá! Como você pode me ajudar com receitas?",
  "context": "user_interaction"
}
```

### Mudança de Provedor para Gemini

```bash
POST /api/v1/llm/provider/switch
Content-Type: application/json

{
  "provider_name": "gemini"
}
```

### Geração de Receitas com Gemini

```bash
POST /api/v1/recipes/generate
Content-Type: application/json

{
  "ingredients": ["arroz", "frango", "cebola"],
  "dietary_restrictions": [],
  "max_recipes": 3,
  "llm_provider": "gemini"
}
```

### Teste de Status do Provedor

```bash
GET /api/v1/llm/provider/status
```

Resposta esperada quando Gemini está configurado:
```json
{
  "success": true,
  "data": {
    "providers": [
      {
        "name": "openai", 
        "active": false,
        "model": "gpt-3.5-turbo",
        "status": "configured"
      },
      {
        "name": "gemini",
        "active": true,
        "model": "gemini-1.5-flash", 
        "status": "configured"
      }
    ],
    "active_provider": "gemini"
  }
}
```

## Características do Gemini

### Vantagens
- **Rápido**: Especialmente o modelo flash
- **Multimodal**: Suporte nativo para texto e imagens (em versões futuras)
- **Gratuito**: Tier gratuito generoso para desenvolvimento
- **Contexto**: Janela de contexto grande

### Considerações
- **Rate Limits**: Diferentes dos outros provedores
- **Formato de Resposta**: Ligeiramente diferente do OpenAI
- **Segurança**: Filtros de segurança mais rigorosos

## Migração de OpenAI para Gemini

Para alternar do OpenAI para Gemini como provedor padrão:

1. **Configure a API Key**:
   ```env
   GEMINI_API_KEY=seu-api-key-aqui
   ```

2. **Reinicie o servidor** para auto-configuração

3. **Alterne o provedor via API**:
   ```bash
   POST /api/v1/llm/provider/switch
   Content-Type: application/json
   
   {
     "provider_name": "gemini"
   }
   ```

## Troubleshooting

### Erro: "provedor 'gemini' não é suportado"
- Verifique se o Gemini foi registrado no ProviderFactory
- Confirme se o build incluiu os arquivos do provedor

### Erro: "API key é obrigatória"
- Configure a variável `GEMINI_API_KEY` no .env
- Reinicie o servidor
- Ou configure manualmente via API

### Erro: "invalid API key"
- Verifique se a chave foi copiada corretamente
- Confirme se a chave tem as permissões necessárias
- Teste a chave diretamente na API do Google

## Obtendo uma API Key do Gemini

1. Acesse: https://aistudio.google.com/app/apikey
2. Faça login com sua conta Google
3. Clique em "Create API Key"
4. Selecione um projeto Google Cloud (ou crie um novo)
5. Copie a chave gerada
6. Cole no arquivo `.env`: `GEMINI_API_KEY=sua-chave-aqui`

## Comparação de Custos (Estimativa)

| Provedor | Modelo | Input (1K tokens) | Output (1K tokens) |
|----------|---------|-------------------|-------------------|
| OpenAI | gpt-3.5-turbo | $0.0015 | $0.002 |
| OpenAI | gpt-4 | $0.03 | $0.06 |
| Gemini | gemini-1.5-flash | Gratuito* | Gratuito* |
| Gemini | gemini-1.5-pro | $0.00125 | $0.00375 |

*Dentro dos limites do tier gratuito

## Próximos Passos

Com o Gemini implementado, você pode:

1. **Testar Performance**: Compare respostas entre OpenAI e Gemini
2. **Otimizar Custos**: Use Gemini para reduzir custos de development
3. **Implementar Fallback**: Configure múltiplos provedores para alta disponibilidade
4. **Explorar Modelos**: Teste diferentes modelos Gemini para casos específicos

## Suporte Técnico

Para dúvidas sobre a implementação do Gemini:
- Verifique logs do servidor para erros detalhados
- Use o endpoint `/api/v1/llm/provider/test/gemini` para diagnósticos
- Consulte a documentação oficial do Google AI Studio