# Backend — Harém Brasil

Este ficheiro é o **ponto de entrada** da documentação backend do repositório. O detalhe completo está em:

|| Conteúdo |
|---|---|
|---|Ecossistema técnico, domínio PostgreSQL, segurança, tempo real, pagamentos, observabilidade, CI/CD, infraestrutura Hostinger, estrutura de repositório, roadmap por fases. |
|---| Contrato HTTP `/api/v1`, convenções, RBAC, catálogo de endpoints, WebSocket, OpenAPI, pacotes `internal/`, checklist Go. |


## 3. Arquitetura lógica (resumo)

- **Clientes:** SPA (e apps futuros) → CDN/WAF opcional → **Nginx** (TLS) → **`api` (Go)** e **`realtime` (Go WebSocket)**.
- **Dados:** **PostgreSQL**, **Redis** (cache, filas, sessões), armazenamento de **objetos** para mídia.
- **Externos:** PSPs, e-mail transacional, analytics/ads conforme política.

**Evolução:** MVP com monólito modular (**`api`**, **`realtime`**, **`worker`**) na mesma VPS; crescimento com separação de processos/BD; escala com réplicas de leitura, filas (Redis Streams / NATS / SQS) e CDN.

Diagrama Mermaid e matriz **Go vs .NET** opcional: ver arquitetura, secções 3 e apêndice .NET.

---

## 4. Dados (PostgreSQL)

- **IDs expostos:** UUID (v4/v7); sem IDs sequenciais na API.
- **Convenções:** `timestamptz` UTC, soft delete `deleted_at` onde aplicável, UTF-8.
- **Grupos principais:** identidade (`users`, social, verificação criador, sessões/refresh, `audit_log`), conteúdo (`posts`, mídia, likes, comentários), fórum (`forum_categories`, `forum_topics`, `forum_posts`), chat (`chat_rooms`, `chat_members`, `chat_messages`), monetização (planos, `subscriptions`, encomendas/transações, comissões, `payouts` futuros), `notifications`.
- **Migrações:** `goose`, `golang-migrate` ou `atlas`; CI com gate em staging; **queries sempre parametrizadas**.

Pormenores de índices, backup/PITR, particionamento e LGPD: documento de **arquitetura**, secção 4.

---

## 5. Segurança e API (regras transversais)

1. JSON `Content-Type: application/json; charset=utf-8`; limites de corpo no Nginx e **`http.MaxBytesReader`** por rota.
2. **Whitelist** em `POST`/`PATCH`; sem mass assignment nem colunas internas nas respostas.
3. **BOLA:** validar posse ou papel em todos os recursos com `{id}`.
4. **BFLA:** **`/api/v1/admin/*`** isolado, middleware e política próprios, papel `admin`.
5. **JWT:** access curto; refresh opaco com rotação; validar `iss`, `aud`, `exp`, assinatura (`golang-jwt/jwt/v5`).
6. **Rate limit** por IP e por utilizador; cabeçalhos `RateLimit-*` / `Retry-After`.
7. **Webhooks:** corpo cru para HMAC; `200` rápido + fila; idempotência por `event_id`; logs sem payload completo.
8. **Logging estruturado** (`slog` + OpenTelemetry); eventos de segurança **sem PII** desnecessária.

Modelo de erros (`application/problem+json`), paginação por **cursor**, cabeçalhos (`Authorization`, `X-Request-Id`, `Idempotency-Key`, `If-Match`), limites de payload por contexto e **catálogo completo de rotas**.

**Papéis RBAC:** `guest`, `user`, `creator`, `moderator`, `admin` — matriz e permissões exemplo no mesmo documento.

---

## 6. Tempo real

- **WebSocket:** autenticação preferencial via **`POST /realtime/ticket`** + upgrade (evitar JWT longo na query).
- **Envelope** e tipos de eventos (`chat.send`, `chat.message`, `notification`, etc.): secção 8.
- Fan-out entre instâncias: **Redis** Pub/Sub ou streams quando houver vários processos `realtime`.
- **Go:** goroutines com `context`, heartbeat, limite de frame, verificação de membro da sala na subscrição.

---

## 7. Pacotes Go sugeridos (`internal/`)

| Pacote | Função |
|--------|--------|
| `httpapi` | Router, `v1`, middleware |
| `auth` | JWT, refresh, OIDC, política de senha |
| `user` | Perfis, pesquisa |
| `feed` | Posts, likes, comments, feed |
| `forum` | Categorias, tópicos, respostas |
| `chat` | REST + autorização de salas; hub WS |
| `billing` | Planos, checkout, subscrição |
| `webhook` | HMAC, deduplicação, fila |
| `creator` | Verificação, catálogo, pedidos, ganhos |
| `notify` | Notificações, fan-out WS |
| `admin` | Handlers admin + auditoria |
| `middleware` | Request ID, auth, rate limit, recover, logging, `MaxBytesReader` |
| `realtime` | Hub WebSocket, tickets, quotas |

**Stack de referência:** router (`chi` / `echo` / `fiber`), **pgx** + **sqlc** (ou equivalente parametrizado), validação (`validator/v10` ou contrato gerado **OpenAPI** / `oapi-codegen`), WebSocket (`nhooyr.io/websocket` ou `gorilla/websocket`).

**OpenAPI:** manter `openapi/openapi.yaml` como fonte de verdade; CI com lint (`spectral` / `redocly`).

---

## 8. Estrutura de repositório (sugestão)

```
harem/
  cmd/
    api/main.go
    realtime/main.go
    worker/main.go
  internal/
    auth/
    forum/
    chat/
    billing/
    storage/
    middleware/
  pkg/            # só se reutilização real entre módulos
  migrations/
  deployments/
    nginx/
    systemd/
  docs/
```

---

## 9. Infraestrutura e operações (síntese)

- **Hostinger VPS:** Ubuntu LTS, Nginx TLS, binários Go com **systemd**, UFW, Fail2ban, SSL (Let’s Encrypt).
- **PostgreSQL:** mesma VPS (MVP) ou VPS dedicada; `pg_hba` SCRAM; role da app com privilégios mínimos; **PgBouncer** em escala.
- **Mídia:** object storage compatível S3 (Hostinger ou R2/B2); URLs pré-assinadas para upload.
- **Observabilidade:** logs JSON, métricas (Prometheus), tracing (OTel), alertas (picos 4xx/5xx, falhas de login, filas de webhook, disco PG).

CI/CD, ambientes e política de migrações em produção: **arquitetura**, secções 8–9.

---

## 10. Roadmap técnico (alinhado ao PDF de produto)

| Fase | Duração indicada | Foco |
|------|-------------------|------|
| 1 | ~2 semanas | Stack, repos, CI, primeira migração, ambientes base |
| 2 | ~3 semanas | Auth, JWT/sessões, roles, verificação criador |
| 3 | ~3 semanas | Posts, mídia, feed, likes, comentários |
| 4 | ~4 semanas | WebSocket, salas, histórico, Redis multi-instância |
| 5 | ~3 semanas | Fórum CRUD, moderação, RBAC |
| 6 | ~4 semanas | Planos, webhooks, acesso Premium/VIP |
| 7 | ~3 semanas | Área criador, uploads, dashboard de ganhos |

---

## 11. Próximos passos (equipa)

1. Validar **juridicamente** idade, termos, moderação e tratamento de dados (LGPD + conteúdo adulto).
2. Fixar versões: Go, PostgreSQL, política de retenção de chat e mídia.
3. Carga mínima (k6/vegeta) e plano de **disaster recovery** documentado.

---

## 12. Glossário rápido

| Sigla | Significado |
|--------|-------------|
| **BOLA** | Autorização ao nível do objeto (acesso a recurso alheio). |
| **BFLA** | Autorização ao nível de função (ex.: rotas admin). |
| **PSP** | Provedor de pagamento (Stripe, PagSeguro, Mercado Pago). |
| **PITR** | Recuperação PostgreSQL point-in-time. |

---
