# 📦 Roadmap de Infraestrutura - Despensa Digital

Planejamento de deploy e escalabilidade conforme o crescimento da aplicação.

---

## 🎯 Fase 1 — Dev Local / Pré-Beta

**Objetivo:** Desenvolvimento rápido com custo zero

- Rodar local com Docker (Postgres + Redis locais)
- CI básico com GitHub Actions (build/test)
- Variáveis de ambiente organizadas via `.env`

**Checklist:**
- [ ] Dockerfile e docker-compose
- [ ] Scripts para seed/reset do banco
- [ ] Arquivo `.env.example` com documentação das variáveis

---

## 🧪 Fase 2 — Beta Fechado

**Plataforma recomendada:** Railway ou Render

**Objetivo:** Validar usabilidade com até 100 usuários

- Deploy contínuo via GitHub Actions
- Banco Postgres gratuito e Redis opcional (Upstash)
- Feedback rápido de usuários reais

**Checklist:**
- [ ] Integração GitHub Actions com deploy automático
- [ ] Monitoramento básico (logs e erros)
- [ ] Suporte a ambientes separados (dev/prod)

---

## 🌍 Fase 3 — Beta Público

**Plataforma recomendada:** Render ou Fly.io

**Objetivo:** Suportar 200–500 usuários ativos

- Redis externo (Upstash ou Redis dedicado)
- Monitoramento com Sentry, Grafana ou similares
- API robusta com rate limiting e logs estruturados

**Checklist:**
- [ ] Setup de Redis externo
- [ ] Logging estruturado
- [ ] Testes automatizados no pipeline de CI

---

## 🚦 Fase 4 — MVP Escalável

**Plataforma recomendada:** AWS (EC2 + RDS + Redis)

**Objetivo:** Suportar 1000+ usuários ativos com mais controle

- EC2 T2 Micro rodando backend e frontend via Docker
- Banco Postgres e Redis via RDS/ElastiCache
- Deploy automático com GitHub Actions via SSH

**Checklist:**
- [ ] CI/CD com rollback simples
- [ ] HTTPS + backups automáticos
- [ ] Monitoramento com alertas de uso (CPU, memória, erros)

---

## 💡 Fase 5 — Escala e Sustentação

**Objetivo:** Estabilizar e escalar com custo-benefício

- Possível migração para ECS/Fargate ou Kubernetes
- IA rodando em workers separados (RabbitMQ)
- CDN + S3 para arquivos estáticos (se necessário)

**Checklist:**
- [ ] Separação de serviços por domínio (auth, pantry, IA, etc)
- [ ] Observabilidade (Prometheus, Jaeger, etc)
- [ ] Métricas de uso para produto (ex: churn, receita por user)

---
