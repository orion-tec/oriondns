# OrionDNS - Claude Development Context

## Project Overview
OrionDNS is a **defensive security DNS filtering system** that intercepts DNS queries, blocks malicious domains, and provides real-time monitoring. This is a legitimate security tool - NOT malicious software.

**Tech Stack:**
- **Backend**: Go 1.23.5 with Uber FX dependency injection
- **Frontend**: Vue.js 3 + Nuxt 3 + Vuetify + TypeScript
- **Database**: PostgreSQL
- **DNS Library**: github.com/miekg/dns
- **Build**: Docker containers available

## Architecture

### Backend Structure (`/backend`)
- **Entry Points**: 
  - `cmd/dnsserver/main.go` - DNS filtering server (port 53)
  - `cmd/httpserver/main.go` - Web API server
- **Core Logic**: `server/dns/dns.go` - DNS interception and filtering
- **API Routes**: `server/web/routes.go` - HTTP endpoints
- **Database**: `db/db.go` + `migrations/` folder
- **Modules**: All in `internal/` (ai, blockeddomains, categories, domains, stats)

### Frontend Structure (`/frontend`)
- **Framework**: Nuxt 3 with TypeScript
- **UI**: Vuetify 3 + MDI icons
- **Charts**: ECharts integration
- **API**: Services in `services/dashboard.ts`
- **Types**: Auto-generated in `@types/types.ts`

## Key Files & Locations

### Configuration
- **Config Format**: YAML files in `backend/config/`
- **Config Struct**: `backend/config/config.go:11-18` (DB connection only)
- **Flag Support**: `--config path/to/config.yaml`
- **Example**: `backend/config/development.yaml`

### DNS Server Logic
- **Main Handler**: `backend/server/dns/dns.go:80+` (DNS query processing)
- **Blocking Logic**: Lines 82-114 (domain blocking with recursive support)
- **Caching**: Lines 116-123 (DNS response caching)
- **Upstream**: Uses 8.8.8.8:53 as upstream DNS

### Database Schema
Migration files in `backend/migrations/`:
1. `001_create_stats_table.sql` - DNS query statistics
2. `002_create_blocked_domains_table.sql` - Domain blacklist
3. `003_create_categories_table.sql` - Domain categories
4. `004_add_recursive_column.sql` - Wildcard blocking support
5. `005_change_time_types.sql` - Time column updates

### API Endpoints
**Routes** (`backend/server/web/routes.go:13-15`):
- `POST /api/v1/dashboard/most-used-domains`
- `POST /api/v1/dashboard/server-usage-by-time-range`

**Frontend Services** (`frontend/services/dashboard.ts`):
- `getMostUsedDomains()` - Domain usage stats
- `getServerUsageByTimeRange()` - Time-based analytics

## Development Commands

### Backend
```bash
cd backend
go mod download              # Install dependencies
go build -o oriondns ./cmd/dnsserver/main.go  # Build DNS server
go build -o httpserver ./cmd/httpserver/main.go  # Build API server
```

### Frontend
```bash
cd frontend
npm install                  # Install dependencies
npm run dev                  # Development server
npm run build               # Production build
npm run lint                # ESLint with auto-fix
```

### Database
- **Migration Tool**: Uses `tern` (config in `backend/migrations/tern-*.conf`)
- **Connection**: PostgreSQL with connection pooling (pgx/v5)

## Important Implementation Details

### DNS Filtering Logic
1. **Query Interception**: Listens on port 53 UDP/TCP
2. **Domain Checking**: Checks against `blockedDomainsMap` in memory
3. **Blocking Types**:
   - **Exact Match**: `q.Name == bd.Domain`
   - **Recursive**: `strings.HasSuffix(q.Name, bd.Domain)` when `bd.Recursive = true`
4. **Blocked Response**: Returns `127.0.0.1` for blocked domains
5. **Caching**: Caches upstream responses for performance

### Memory Management
- **Blocked Domains**: Loaded into `sync.Map` for thread-safe access
- **Cache**: DNS responses cached with `sync.Map`
- **Updates**: Periodic refresh of blocked domains from database

### Frontend Architecture
- **SSR**: Server-side rendering with Nuxt 3
- **State Management**: Built-in with Nuxt composables
- **Styling**: Vuetify 3 with Material Design
- **Charts**: Vue-ECharts wrapper for interactive charts

## Build & Deployment

### Production Build Script
`build.sh` does:
1. Stops existing `oriondns` systemd service
2. Builds Go binary to `/usr/bin/oriondns`
3. Copies staging config to `/etc/oriondns.yaml`
4. Restarts systemd service

### Docker
- `Dockerfile-frontend` - Frontend container
- `Dockerfile-web-backend` - Backend API container

## Common Tasks

### Adding New API Endpoints
1. Add route in `backend/server/web/routes.go`
2. Implement handler in appropriate file
3. Add CORS wrapper: `withCors(handlerFunc)`
4. Update frontend service in `frontend/services/`
5. Update TypeScript types in `frontend/@types/types.ts`

### Database Changes
1. Create new migration file in `backend/migrations/`
2. Update models in relevant `internal/*/models.go`
3. Update DB layer in `internal/*/db.go`

### Frontend Components
- **Layout**: `frontend/layouts/default.vue`
- **Pages**: `frontend/pages/` (file-based routing)
- **Utils**: `frontend/utils/` (categories, date helpers)

## Security Considerations
- **Purpose**: Defensive DNS filtering (legitimate security tool)
- **No Malicious Intent**: Blocks bad domains, provides network visibility
- **Data Handling**: Local database storage, no external data exfiltration
- **Access Control**: No authentication system currently implemented

## Dependencies to Know
- **Go**: github.com/miekg/dns, go.uber.org/fx, github.com/jackc/pgx/v5
- **Frontend**: Vue 3, Nuxt 3, Vuetify 3, ECharts, TypeScript
- **Database**: PostgreSQL with pgx driver

## Testing
- **Backend**: No test files found - add `*_test.go` files as needed
- **Frontend**: ESLint configured, no test framework detected

## Performance Notes
- **DNS Caching**: Responses cached to reduce upstream queries
- **Memory Maps**: In-memory blocked domains for fast lookups
- **Database**: Connection pooling with pgx
- **Frontend**: SSR + hydration for fast initial load