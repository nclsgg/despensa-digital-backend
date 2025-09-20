# 🥫 Despensa Digital

A **Despensa Digital** é uma aplicação SaaS para organização e contro### 🔐 Autenticação OAuth

| Método | Rota                    | Descrição                         |
|--------|------------------------|-----------------------------------|
| GET    | `/oauth/login/:provider`| Iniciar fluxo OAuth (Google/GitHub)|
| GET    | `/oauth/callback`       | Callback do provedor OAuth         |
| POST   | `/oauth/logout`         | Logout e remoção do token          |
| POST   | `/oauth/refresh`        | Gera novo access token             |ens em casa ou pequenos comércios, com foco em:

- Gerenciar sua despensa (itens que você tem)
- Criar listas de compras automáticas com base no que falta
- Descobrir receitas com base no que está disponível
- Autenticação OAuth2 com provedores populares
- Documentação Swagger completa e API padronizada

---

## 🚀 Tecnologias

- **Backend:** Golang (Gin, GORM, OAuth2, JWT, Redis, PostgreSQL)
- **Frontend:** React + Next.js
- **Mobile:** (a definir)
- **Infra:** Docker (em breve)

---

## 🧠 Principais Funcionalidades

- ✅ Autenticação OAuth2 (Google, GitHub) com JWT refresh tokens
- ✅ Controle de usuários com roles (admin, user)
- ✅ Gestão de itens da despensa
- ✅ Listagem de usuários (admin)
- 🧠 Recomendações de receitas via IA (em breve)
- 🛒 Lista de compras automática com base na rotina (em breve)

---

## 📁 Estrutura de Pastas

```
backend/
│
├── cmd/server              # Ponto de entrada da aplicação
├── internal/
│   ├── modules/            # Módulos organizados por domínio (auth, user, pantry...)
│   ├── router/             # Rotas e middlewares
│   ├── utils/              # Funções utilitárias
│   └── core/               # Modelos compartilhados
├── pkg/                    # Pacotes utilitários (response, etc.)
├── config/                 # Carregamento de variáveis e configuração
├── docs/                   # Swagger gerado com swag init
```

---

## 🛠️ Como rodar localmente

### 1. Clone o repositório

```bash
git clone https://github.com/nclsgg/despensa-digital.git
cd despensa-digital/backend
```

### 2. Configure as variáveis de ambiente

Crie um `.env` com base no `.env.example`:

```
PORT=5310
GIN_MODE=debug

DATABASE_URL=

JWT_SECRET=
JWT_EXPIRATION=1h
JWT_ISSUER=
JWT_AUDIENCE=

REDIS_URL=
REDIS_USERNAME=
REDIS_PASSWORD=
REDIS_DB=0
```

### 3. Rode a aplicação

```bash
go run cmd/server/main.go
```

### 4. Acesse a documentação Swagger

> http://localhost:5310/swagger/index.html

---

## ✅ Rotas principais

### 🔐 Autenticação

| Método | Rota             | Descrição                    |
|--------|------------------|------------------------------|
| POST   | `/auth/register` | Registro de usuário          |
| POST   | `/auth/login`    | Login e geração de tokens    |
| POST   | `/auth/logout`   | Logout e remoção do token    |
| POST   | `/auth/refresh`  | Gera novo access token       |

### 👤 Usuário

| Método | Rota         | Descrição                         |
|--------|--------------|------------------------------------|
| GET    | `/user/me`   | Dados do usuário logado            |
| GET    | `/user/:id`  | Buscar usuário por ID              |
| GET    | `/user/all`  | Listar todos os usuários (admin)   |

### 🥫 Despensas

| Método | Rota                     | Descrição                            |
|--------|--------------------------|---------------------------------------|
| GET    | `/pantries`              | Listar despensas do usuário           |
| POST   | `/pantries`              | Criar nova despensa                   |
| GET    | `/pantries/{id}`         | Detalhes de uma despensa              |
| PUT    | `/pantries/{id}`         | Atualizar nome da despensa            |
| DELETE | `/pantries/{id}`         | Deletar despensa (soft delete)        |

### 👥 Usuários da Despensa

| Método | Rota                               | Descrição                          |
|--------|------------------------------------|-------------------------------------|
| GET    | `/pantries/{id}/users`            | Listar usuários da despensa         |
| POST   | `/pantries/{id}/users`            | Adicionar usuário à despensa        |
| DELETE | `/pantries/{id}/users`            | Remover usuário da despensa         |

### 💼 Itens da Despensa

| Método | Rota                         | Descrição                            |
|--------|------------------------------|---------------------------------------|
| POST   | `/items`                     | Criar item na despensa                |
| GET    | `/items/pantry/{id}`         | Listar itens da despensa              |
| GET    | `/items/{id}`                | Obter detalhes de um item             |
| PUT    | `/items/{id}`                | Atualizar item                        |
| DELETE | `/items/{id}`                | Deletar item                          |

### 🌍 Categorias de Itens

| Método | Rota                                                         | Descrição                                              |
|--------|--------------------------------------------------------------|-----------------------------------------------------------|
| POST   | `/item-categories`                                           | Criar nova categoria de item                              |
| POST   | `/item-categories/default`                                   | Criar categoria padrão (admin)                           |
| POST   | `/item-categories/from-default/{default_id}/pantry/{pantry_id}` | Clonar categoria padrão para despensa                    |
| GET    | `/item-categories/pantry/{id}`                               | Listar categorias de uma despensa                         |
| GET    | `/item-categories/{id}`                                      | Obter detalhes da categoria                               |
| PUT    | `/item-categories/{id}`                                      | Atualizar categoria                                       |
| DELETE | `/item-categories/{id}`                                      | Deletar categoria                                         |
| GET    | `/item-categories/user`                                      | Listar categorias do usuário                             |

---

## 📦 Instalação do Swagger (para desenvolvimento)

```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g cmd/server/main.go -o cmd/server/docs
```

---

## 📄 Licença

MIT © [Nicolas Guadagno](https://github.com/nclsgg)

