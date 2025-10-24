# My Pantry Endpoint

## Overview
Este documento descreve o novo endpoint `GET /api/v1/pantries/my-pantry` que retorna a despensa principal do usuário autenticado.

## Contexto
Com base na premissa de que inicialmente o usuário terá apenas uma despensa, este endpoint simplifica o acesso à despensa do usuário sem necessidade de informar o ID específico.

## Endpoint

### GET /api/v1/pantries/my-pantry

Retorna a primeira despensa do usuário autenticado com todos os seus itens e contagem.

#### Headers
```
Authorization: Bearer <token>
```

#### Response Success (200 OK)
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "name": "Minha Despensa",
    "owner_id": "uuid",
    "item_count": 15,
    "created_at": "2025-10-22T10:00:00Z",
    "updated_at": "2025-10-22T10:00:00Z",
    "items": [
      {
        "id": "uuid",
        "pantry_id": "uuid",
        "name": "Arroz",
        "quantity": 2.0,
        "unit": "kg",
        "price_per_unit": 5.50,
        "total_price": 11.00,
        "added_by": "uuid",
        "category_id": "uuid",
        "expiration_date": "2025-12-31T23:59:59Z",
        "created_at": "2025-10-22T10:00:00Z",
        "updated_at": "2025-10-22T10:00:00Z"
      }
    ]
  }
}
```

#### Response Error (404 Not Found)
```json
{
  "success": false,
  "code": "NOT_FOUND",
  "message": "Pantry not found"
}
```

Retornado quando o usuário ainda não possui nenhuma despensa.

#### Response Error (500 Internal Server Error)
```json
{
  "success": false,
  "message": "Failed to list pantry items"
}
```

Retornado quando há erro ao buscar os itens da despensa.

## Implementação

### Fluxo de Execução

1. **Handler** (`pantry_handler.go`):
   - Extrai o `userID` do contexto (injetado pelo middleware de autenticação)
   - Chama o serviço `GetMyPantry`
   - Busca os itens da despensa retornada
   - Converte os dados para DTOs
   - Retorna resposta formatada

2. **Service** (`pantry_service.go`):
   - Busca todas as despensas do usuário usando `repo.GetByUser`
   - Retorna erro `ErrPantryNotFound` se não houver despensas
   - Pega a primeira despensa da lista
   - Busca a contagem de itens usando `itemRepo.CountByPantryID`
   - Retorna a despensa com a contagem de itens

3. **Repository**:
   - Utiliza os métodos existentes `GetByUser` e `CountByPantryID`
   - Não requer novas implementações no repositório

### Código Adicionado

#### Interface (pantry_port.go)
```go
type PantryService interface {
    // ... métodos existentes ...
    GetMyPantry(ctx context.Context, userID uuid.UUID) (*model.PantryWithItemCount, error)
}

type PantryHandler interface {
    // ... métodos existentes ...
    GetMyPantry(ctx *gin.Context)
}
```

#### Rota (router.go)
```go
pantryGroup.GET("/my-pantry", pantryHandlerInstance.GetMyPantry)
```

## Vantagens

1. **Simplicidade**: Usuário não precisa conhecer o ID da sua despensa
2. **Convenção**: Assume que o usuário tem apenas uma despensa inicialmente
3. **Reutilização**: Usa a mesma estrutura de resposta do endpoint `GET /pantries/:id`
4. **Performance**: Uma única query para buscar a despensa e os itens

## Considerações Futuras

Quando o sistema evoluir para suportar múltiplas despensas por usuário:

1. Este endpoint pode ser mantido retornando a "despensa principal" ou "despensa padrão"
2. Pode-se adicionar um campo `is_default` na tabela `pantries` para marcar a despensa principal
3. O comportamento atual (retornar a primeira) pode ser mantido como fallback

## Diferença dos Endpoints Existentes

### GET /api/v1/pantries
- Lista **todas** as despensas do usuário
- Retorna array de despensas com contagem de itens
- Não retorna os itens de cada despensa

### GET /api/v1/pantries/:id
- Busca uma despensa **específica** por ID
- Requer conhecimento prévio do ID da despensa
- Retorna a despensa com todos os seus itens

### GET /api/v1/pantries/my-pantry (NOVO)
- Busca a **primeira/principal** despensa do usuário
- Não requer ID na requisição
- Retorna a despensa com todos os seus itens
- Ideal para casos de uso com despensa única

## Exemplo de Uso no Frontend

```typescript
// Ao invés de fazer duas chamadas:
// 1. GET /pantries para listar
// 2. GET /pantries/:id para pegar os detalhes

// Agora pode fazer apenas uma:
const response = await fetch('/api/v1/pantries/my-pantry', {
  headers: {
    'Authorization': `Bearer ${token}`
  }
});

const { data } = await response.json();
// data já contém a despensa com todos os itens
```

## Testes Recomendados

1. **Usuário sem despensa**: Deve retornar 404
2. **Usuário com uma despensa**: Deve retornar a despensa com itens
3. **Usuário com múltiplas despensas**: Deve retornar a primeira despensa
4. **Despensa sem itens**: Deve retornar `item_count: 0` e `items: []`
5. **Erro ao buscar itens**: Deve retornar 500
