# Sistema de Lista de Compras com IA - Documentação da API

## Visão Geral

O sistema de lista de compras implementa funcionalidades completas para criar, gerenciar e gerar listas de compras inteligentes usando IA. O sistema considera o perfil do usuário, histórico da despensa e preços médios de produtos brasileiros.

## Módulos Implementados

### 1. Módulo Profile

Gerencia o perfil financeiro e de preferências do usuário.

#### Endpoints

- `POST /api/v1/profile` - Criar perfil
- `GET /api/v1/profile` - Obter perfil
- `PUT /api/v1/profile` - Atualizar perfil
- `DELETE /api/v1/profile` - Deletar perfil

#### Exemplo de Criação de Perfil

```bash
curl -X POST http://localhost:8080/api/v1/profile \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "monthly_income": 3500.00,
    "preferred_budget": 800.00,
    "household_size": 3,
    "dietary_restrictions": ["vegetariano", "sem lactose"],
    "preferred_brands": ["Nestlé", "Sadia", "Tio João"],
    "shopping_frequency": "weekly"
  }'
```

**Resposta:**
```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440001",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "monthly_income": 3500.00,
    "preferred_budget": 800.00,
    "household_size": 3,
    "dietary_restrictions": ["vegetariano", "sem lactose"],
    "preferred_brands": ["Nestlé", "Sadia", "Tio João"],
    "shopping_frequency": "weekly",
    "created_at": "2024-09-24T10:30:00Z",
    "updated_at": "2024-09-24T10:30:00Z"
  }
}
```

### 2. Módulo Shopping List

Gerencia listas de compras manuais e geradas por IA.

#### Endpoints

- `POST /api/v1/shopping-lists` - Criar lista manual
- `GET /api/v1/shopping-lists` - Listar listas do usuário
- `GET /api/v1/shopping-lists/{id}` - Obter lista específica
- `PUT /api/v1/shopping-lists/{id}` - Atualizar lista
- `DELETE /api/v1/shopping-lists/{id}` - Deletar lista
- `PUT /api/v1/shopping-lists/{id}/items/{itemId}` - Atualizar item da lista
- `DELETE /api/v1/shopping-lists/{id}/items/{itemId}` - Deletar item da lista
- `POST /api/v1/shopping-lists/generate` - **Gerar lista com IA**

#### Exemplo de Lista Manual

```bash
curl -X POST http://localhost:8080/api/v1/shopping-lists \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "Compras da Semana",
    "total_budget": 200.00,
    "items": [
      {
        "name": "Arroz",
        "quantity": 1,
        "unit": "kg",
        "estimated_price": 5.50,
        "category": "Grãos",
        "brand": "Tio João",
        "priority": 1,
        "notes": "Arroz integral"
      },
      {
        "name": "Frango",
        "quantity": 1.5,
        "unit": "kg",
        "estimated_price": 12.00,
        "category": "Proteínas",
        "priority": 1
      }
    ]
  }'
```

#### Exemplo de Geração com IA

```bash
curl -X POST http://localhost:8080/api/v1/shopping-lists/generate \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "Lista IA - Semanal",
    "total_budget": 300.00,
    "shopping_type": "weekly",
    "include_basics": true,
    "exclude_items": ["carne vermelha", "leite"],
    "preferred_brands": ["Nestlé", "Sadia"],
    "notes": "Foco em produtos saudáveis e econômicos"
  }'
```

**Como a IA funciona:**

1. **Análise do Perfil**: Considera renda mensal, orçamento preferido, tamanho da família, restrições alimentares
2. **Histórico da Despensa**: Analisa itens frequentemente comprados, preços históricos, padrões de consumo
3. **Lista Básica Brasileira**: Inclui itens essenciais para o brasileiro médio
4. **Pesquisa de Preços**: Para itens sem histórico, busca preços atuais em mercados brasileiros
5. **Otimização**: Balanceia qualidade, necessidade e orçamento

**Resposta da IA:**
```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440002",
    "name": "Lista IA - Semanal",
    "status": "pending",
    "total_budget": 300.00,
    "estimated_cost": 287.50,
    "generated_by": "ai",
    "items": [
      {
        "name": "Arroz integral",
        "quantity": 2,
        "unit": "kg",
        "estimated_price": 6.50,
        "category": "Grãos",
        "brand": "Tio João",
        "priority": 1,
        "notes": "Item básico, consumo frequente na família",
        "source": "ai_suggestion"
      },
      {
        "name": "Frango peito",
        "quantity": 2,
        "unit": "kg",
        "estimated_price": 15.90,
        "category": "Proteínas",
        "brand": "Sadia",
        "priority": 1,
        "notes": "Proteína principal, boa relação custo-benefício",
        "source": "ai_suggestion"
      },
      {
        "name": "Feijão preto",
        "quantity": 1,
        "unit": "kg",
        "estimated_price": 7.20,
        "category": "Grãos",
        "priority": 1,
        "notes": "Complemento proteico, essencial na dieta brasileira",
        "source": "ai_suggestion"
      }
    ]
  }
}
```

#### Tipos de Compra Suportados

- `weekly`: Compras semanais básicas
- `monthly`: Estoque mensal
- `stock_up`: Reposição de estoque
- `emergency`: Compras de emergência

#### Sistema de Prioridades

- `1`: Essencial (itens básicos, acabando)
- `2`: Importante (itens úteis, em falta)
- `3`: Desejável (itens de conveniência)

### 3. Funcionalidades da IA

A IA considera múltiplos fatores para criar listas inteligentes:

#### Análise de Perfil
- **Renda mensal**: Ajusta sugestões ao poder aquisitivo
- **Orçamento preferido**: Respeita limites financeiros
- **Tamanho da família**: Ajusta quantidades
- **Restrições alimentares**: Filtra produtos inadequados
- **Frequência de compras**: Adapta quantidades ao período

#### Análise da Despensa
- **Itens frequentes**: Prioriza produtos consumidos regularmente
- **Preços históricos**: Usa dados de compras anteriores
- **Padrões de consumo**: Identifica necessidades típicas
- **Estoque baixo**: Sugere reposição de itens acabando

#### Inteligência de Preços
- **Histórico do usuário**: Usa preços pagos anteriormente
- **Pesquisa em tempo real**: Para produtos novos, consulta mercados brasileiros
- **Sazonalidade**: Considera variações de preço por época
- **Marcas alternativas**: Sugere opções mais econômicas quando apropriado

### 4. Fluxo de Uso Típico

1. **Configurar Perfil**: Usuário define perfil financeiro e preferências
2. **Usar Sistema Normalmente**: Registrar compras na despensa
3. **Gerar Lista IA**: Solicitar lista baseada em histórico
4. **Revisar e Ajustar**: Modificar itens conforme necessário
5. **Fazer Compras**: Marcar itens como comprados
6. **Registrar Preços Reais**: Atualizar custos para melhorar IA

### 5. Exemplos de Prompts da IA

A IA usa prompts elaborados que incluem:

```
Você é um assistente especializado em criar listas de compras inteligentes para brasileiros.

CONTEXTO DO USUÁRIO:
- Orçamento para compras: R$ 300,00
- Tipo de compra: weekly
- Incluir itens básicos: true

PERFIL DO USUÁRIO:
- Renda mensal: R$ 3.500,00
- Orçamento preferido: R$ 800,00
- Tamanho da família: 3 pessoas
- Frequência de compras: weekly
- Restrições alimentares: vegetariano, sem lactose
- Marcas preferidas: Nestlé, Sadia, Tio João

ANÁLISE DA DESPENSA:
- Total de itens cadastrados: 45
- Categorias mais comuns: Grãos, Proteínas, Laticínios
- Itens frequentemente comprados:
  * Arroz (1.0 kg) - R$ 5.50 em média
  * Feijão (1.0 kg) - R$ 7.20 em média

INSTRUÇÕES:
1. Crie uma lista balanceada e econômica
2. Considere preços médios de mercados brasileiros (usando dados de 2024/2025)
3. Priorize itens essenciais e de qualidade
4. Para produtos sem preço histórico, pesquise preços atuais no Brasil
5. Considere a proporção família/orçamento
6. Inclua uma breve explicação para cada item
```

### 6. Estrutura do Banco de Dados

#### Tabela Profiles
```sql
CREATE TABLE profiles (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL UNIQUE,
    monthly_income NUMERIC,
    preferred_budget NUMERIC,
    household_size INTEGER DEFAULT 1,
    dietary_restrictions TEXT[],
    preferred_brands TEXT[],
    shopping_frequency VARCHAR(20) DEFAULT 'weekly',
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
```

#### Tabela Shopping Lists
```sql
CREATE TABLE shopping_lists (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    name VARCHAR NOT NULL,
    status VARCHAR DEFAULT 'pending',
    total_budget NUMERIC,
    estimated_cost NUMERIC,
    actual_cost NUMERIC,
    generated_by VARCHAR DEFAULT 'manual',
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
```

#### Tabela Shopping List Items
```sql
CREATE TABLE shopping_list_items (
    id UUID PRIMARY KEY,
    shopping_list_id UUID NOT NULL,
    name VARCHAR NOT NULL,
    quantity NUMERIC NOT NULL,
    unit VARCHAR NOT NULL,
    estimated_price NUMERIC,
    actual_price NUMERIC,
    category VARCHAR,
    brand VARCHAR,
    priority INTEGER DEFAULT 3,
    purchased BOOLEAN DEFAULT false,
    notes TEXT,
    source VARCHAR,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
```

## Próximos Passos

Para usar o sistema em produção:

1. **Configure as variáveis de ambiente** para o LLM (OpenAI/Gemini)
2. **Execute as migrações** do banco de dados
3. **Configure autenticação JWT** 
4. **Teste os endpoints** com dados reais
5. **Ajuste os prompts da IA** conforme necessário
6. **Implemente cache** para melhorar performance das consultas de preço

O sistema está pronto para uso e pode ser estendido com funcionalidades adicionais como comparação de preços entre mercados, notificações de ofertas, e integração com APIs de supermercados.