# PROJECT KNOWLEDGE BASE

**Generated:** 2026-02-13
**Commit:** d7385a2
**Branch:** main

## OVERVIEW

Gulrot-tracker — carrot tracking application. Log when you eat a carrot, view stats.

**Stack:** TypeScript monorepo (pnpm workspaces)
- **api/** — Hono REST API + Drizzle ORM + Turso (libSQL)
- **web/** — Astro static frontend

## STRUCTURE

```
vetles-roetter/
├── api/
│   ├── src/
│   │   ├── index.ts          # Hono server entry (Node.js)
│   │   ├── lambda.ts         # AWS Lambda entry
│   │   ├── routes/carrots.ts # API endpoints
│   │   └── db/
│   │       ├── schema.ts     # Drizzle schema (carrots table)
│   │       └── index.ts      # DB connection
│   ├── drizzle.config.ts
│   ├── Dockerfile
│   └── package.json
├── web/
│   ├── src/
│   │   ├── pages/index.astro # Main UI
│   │   └── layouts/Layout.astro
│   ├── astro.config.mjs
│   └── package.json
├── infra/                    # Terraform IaC
│   ├── main.tf               # Lambda + API Gateway
│   ├── variables.tf
│   └── outputs.tf
├── package.json              # Root workspace
└── pnpm-workspace.yaml
```

## API ENDPOINTS

| Method | Path | Description |
|--------|------|-------------|
| POST | /api/carrots | Log a carrot (auto-timestamp) |
| GET | /api/carrots | List entries (supports ?limit, ?offset) |
| GET | /api/carrots/stats | Stats: today, thisWeek, total |
| DELETE | /api/carrots/:id | Remove entry |

## COMMANDS

```bash
pnpm dev          # Run api + web in parallel
pnpm build        # Build all
pnpm typecheck    # Typecheck all

# API specific
pnpm --filter @vetles-roetter/api dev       # API dev server (port 3000)
pnpm --filter @vetles-roetter/api db:generate  # Generate Drizzle migrations
pnpm --filter @vetles-roetter/api db:migrate   # Run migrations
pnpm --filter @vetles-roetter/api db:studio    # Drizzle Studio

# Web specific
pnpm --filter @vetles-roetter/web dev       # Astro dev server
pnpm --filter @vetles-roetter/web build     # Build static site
```

## ENVIRONMENT

**api/.env** (copy from .env.example):
```
TURSO_DATABASE_URL=libsql://...
TURSO_AUTH_TOKEN=...
```

**web/.env** (copy from .env.example):
```
PUBLIC_API_URL=http://localhost:3000
```

For local dev without Turso: leave TURSO_DATABASE_URL empty → uses local.db (SQLite file).

## CONVENTIONS

- ESM throughout (`"type": "module"`)
- Strict TypeScript
- No comments in code (self-documenting)
- Hono's chained route pattern

## DEPLOYMENT

**AWS (Terraform):**
```bash
# Build Lambda bundle
pnpm --filter @vetles-roetter/api build:lambda

# Deploy infrastructure
cd infra
cp terraform.tfvars.example terraform.tfvars  # Add Turso credentials
terraform init
terraform plan
terraform apply
```

**API (Docker):** `docker build -t gulrot-api ./api && docker run -p 3000:3000 gulrot-api`
**Web:** Static deploy anywhere (Netlify, Vercel, Cloudflare Pages).

## NOTES

- Owner: Vetle Mangrud Refsnes
- License: MIT
