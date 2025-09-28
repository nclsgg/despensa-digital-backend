cd despensa-digital/backend
go install github.com/swaggo/swag/cmd/swag@latest
# 🥫 Despensa Digital — Backend

API em Go que alimenta a plataforma **Despensa Digital**. O serviço organiza o dia a dia da despensa, sugere compras com apoio de IA e integra autenticação OAuth2 com provedores populares.

---

## 🧩 Arquitetura em alto nível

- **Camadas por domínio** em `internal/modules/<domínio>` divididas em `domain`, `service`, `handler`, `repository`, `dto` e `model`.
- **Handlers** traduzem HTTP ↔️ domínio utilizando `pkg/response` para payloads consistentes.
- **Services** concentram regras de negócio, expõem interfaces, usam sentinelas (`Err<Algo>`) e operam sempre com `context.Context`.
- **Repositories** encapsulam o GORM; retornos de erro são convertidos em sentinelas nas services, mantendo handlers desacoplados do banco.
- **Router** (`internal/router`) injeta dependências, middlewares de auth/profile e compõe o servidor Gin.
- **Integração com IA** via módulo `llm` que cria prompts, escolhe provedores e expõe serviços consumidos por receitas e listas automáticas.

---

## 📦 Módulos principais

| Módulo | Responsabilidade | Destaques |
| --- | --- | --- |
| `auth` | Fluxo OAuth (Google/GitHub), tokens JWT, conclusão de perfil | Serviços de OAuth com refresh token no Redis, sentinelas `auth:` |
| `user` | Consultas de usuário autenticado e operações administrativas | `ErrUserNotFound`, profile completion via service |
| `profile` | Preferências de compra do usuário | Conversão `StringArray`, deduplicação, sentinelas `ErrProfile*` |
| `pantry` | Gestão de despensas e membros | Regras de acesso e soft delete via GORM |
| `item` | Inventário de itens da despensa | DTOs com formatação ISO8601, filtros e validações |
| `shopping_list` | Listas manuais e geradas por IA | Prompt builder, domínio rico, sentinelas para autorização/IA |
| `recipe` | Sugestões de receitas a partir do estoque | Integra LLM com preferências do usuário |
| `llm` | Abstrações para provedores e prompts | Seleção de provider, builders e sessão |

Outros pacotes relevantes:

- `config`: carrega variáveis de ambiente e configurações (JWT, Redis, banco).
- `pkg/response`: padroniza envelopes de resposta para os handlers.
- `pkg/database`: inicialização de PostgreSQL e Redis.

---

## 🚀 Stack

- **Go 1.22+**
- **Gin** para HTTP routing
- **GORM** + PostgreSQL
- **Redis** para tokens
- **google/uuid**, **golang-jwt/jwt/v5**
- **swaggo/swag** para documentação Swagger
- **miniredis** + **testify** em testes

---

## � Estrutura do backend

```
backend/
├── cmd/server/               # main.go, configuração do Swagger
├── config/                   # Configuração e carregamento de .env
├── internal/
│   ├── modules/
│   │   ├── auth/
│   │   ├── item/
│   │   ├── llm/
│   │   ├── pantry/
│   │   ├── profile/
│   │   ├── recipe/
│   │   ├── shopping_list/
│   │   └── user/
│   ├── core/                 # Modelos e erros globais
│   ├── router/               # Inicialização das rotas e middlewares
│   └── utils/                # Validadores e helpers
├── pkg/
│   ├── database/             # Conexões PostgreSQL/Redis
│   └── response/             # Helpers para JSON/API
├── docs/                     # Swagger gerado
└── tmp/                      # Artefatos temporários (builds, logs)
```

---

## 🛠️ Guia rápido de desenvolvimento

### 1. Clonar e instalar dependências

```bash
git clone https://github.com/nclsgg/despensa-digital.git
cd despensa-digital/backend
go mod download
```

### 2. Configurar variáveis de ambiente

Copie o `.env.example` para `.env` e ajuste conforme seu ambiente:

```
PORT=5310
GIN_MODE=debug

DATABASE_URL=postgres://user:password@localhost:5432/despensa?sslmode=disable

JWT_SECRET=change-me
JWT_EXPIRATION=1h
JWT_ISSUER=despensa-digital
JWT_AUDIENCE=app

REDIS_URL=localhost:6379
REDIS_USERNAME=
REDIS_PASSWORD=
REDIS_DB=0
```

### 3. Subir infra (opcional)

O repositório possui `docker-compose.yaml` para PostgreSQL/Redis de desenvolvimento:

```bash
docker compose up -d
```

### 4. Executar o servidor

```bash
go run cmd/server/main.go
```

### 5. Acessar a documentação Swagger

> http://localhost:5310/swagger/index.html

---

## ✅ Fluxos suportados

### Autenticação OAuth

| Método | Rota | Descrição |
| --- | --- | --- |
| GET | `/oauth/login/:provider` | Inicia o fluxo OAuth (Google/GitHub) |
| GET | `/oauth/callback` | Callback do provedor OAuth |
| POST | `/oauth/logout` | Revoga o refresh token |
| POST | `/oauth/refresh` | Gera novo access token |

### Rotas principais

| Domínio | Exemplo de rotas | Observações |
| --- | --- | --- |
| Auth | `/auth/login`, `/auth/logout` | Fluxo tradicional (em revisão) |
| User | `/user/me`, `/user/:id`, `/user/all` | Sentinelas para not-found, rotas admin |
| Profile | `/profile` (CRUD) | Exige perfil único por usuário |
| Pantry | `/pantries`, `/pantries/{id}/users` | Controle de acesso por owner/membros |
| Item | `/items`, `/items/pantry/{id}` | Respostas ISO8601, filtros |
| Shopping List | `/shopping-lists`, `/shopping-lists/generate` | Geração manual e IA |
| Recipe | `/recipes/generate`, `/recipes/ingredients` | Depende de pantry + IA |

Consulte `docs/swagger.yaml` ou a Wiki para detalhes completos dos contratos.

---

## 🔍 Qualidade e testes

Execute estes passos antes de abrir PRs:

```bash
gofmt -w ./internal ./pkg ./cmd
go test ./...
```

- Os módulos principais possuem testes de serviço com `testify` e `miniredis` (auth, profile, user, shopping_list, etc.).
- Utilize sentinelas definidas em `internal/modules/<domínio>/domain/errors.go` para manter mapas de erro consistentes.
- Handlers nunca retornam mensagens cruas de GORM — sempre converta para códigos/labels de domínio.

---

## � Ferramentas auxiliares

- **Swagger:** `go install github.com/swaggo/swag/cmd/swag@latest` e `swag init -g cmd/server/main.go -o cmd/server/docs`.
- **Postman:** veja `postman_collection.json` e exemplos em `docs/postman-collection-guide.md`.
- **LLM providers:** orientações em `docs/provider-selection.md` e `docs/gemini-provider.md`.

---

## 📄 Licença

MIT © [Nicolas Guadagno](https://github.com/nclsgg)

