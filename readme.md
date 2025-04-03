# 🥫 Despensa Digital

A **Despensa Digital** é uma aplicação SaaS para organização e controle de itens em casa ou pequenos comércios, com foco em:

- Gerenciar sua despensa (itens que você tem)
- Criar listas de compras automáticas com base no que falta
- Descobrir receitas com base no que está disponível
- Autenticação segura com JWT + Refresh Token
- Documentação Swagger completa e API padronizada

---

## 🚀 Tecnologias

- **Backend:** Golang (Gin, GORM, JWT, Redis, PostgreSQL)
- **Frontend:** React (em desenvolvimento)
- **Mobile:** (a definir)
- **Infra:** Docker (em breve)

---

## 🧠 Principais Funcionalidades

- ✅ Cadastro e login com JWT + refresh token (cookies seguros)
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

| Método | Rota            | Descrição                   |
|--------|------------------|-----------------------------|
| POST   | `/auth/register` | Registro de usuário         |
| POST   | `/auth/login`    | Login e geração de tokens   |
| POST   | `/auth/logout`   | Logout e remoção do refresh |
| POST   | `/auth/refresh`  | Gera novo access token      |
| GET    | `/user/me`       | Usuário logado              |
| GET    | `/user/:id`      | Usuário por Id              |
| GET    | `/user/all`      | Lista todos os usuários     |

Mais rotas serão adicionadas conforme o projeto evolui 🚧

---

## 📦 Instalação do Swagger (para desenvolvimento)

```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g cmd/server/main.go -o cmd/server/docs
```

---

## 📄 Licença

MIT © [Nicolas Guadagno](https://github.com/nclsgg)
