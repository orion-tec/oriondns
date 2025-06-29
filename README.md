# OrionDNS

A DNS filtering and monitoring system with AI-powered domain categorization and real-time statistics dashboard.

## Overview

OrionDNS is a defensive security tool that intercepts DNS queries, filters malicious or unwanted domains, and provides comprehensive monitoring through a web dashboard. It combines traditional DNS filtering capabilities with AI-powered domain analysis and categorization.

## Features

### DNS Server
- **DNS Query Interception**: Captures and processes all DNS requests
- **Domain Blocking**: Blocks access to domains based on configurable rules
- **Recursive Filtering**: Supports wildcard blocking for subdomains
- **DNS Caching**: Improves performance with intelligent caching
- **Upstream DNS**: Uses Google DNS (8.8.8.8) as upstream resolver

### Domain Management
- **Blocked Domains Database**: Persistent storage of blocked domains
- **Category-based Filtering**: Organize blocked domains by categories
- **AI-Powered Analysis**: Automatic domain categorization using AI
- **Real-time Updates**: Dynamic updates to blocking rules without restart

### Analytics & Monitoring
- **Real-time Statistics**: Track DNS queries, blocks, and performance
- **Time-based Analytics**: Historical data with configurable time ranges
- **Web Dashboard**: Modern Vue.js interface with interactive charts
- **Export Capabilities**: Data export functionality

## Architecture

### Backend (Go)
- **DNS Server** (`cmd/dnsserver`): Core DNS filtering service
- **HTTP API** (`cmd/httpserver`): REST API for management and statistics
- **Database Layer**: PostgreSQL integration for data persistence
- **AI Integration**: Domain classification and analysis
- **Modular Design**: Clean architecture with dependency injection (Uber FX)

### Frontend (Vue.js/Nuxt)
- **Dashboard Interface**: Real-time monitoring and statistics
- **Chart Visualization**: Interactive charts using ECharts
- **Responsive Design**: Built with Vuetify UI framework
- **TypeScript Support**: Type-safe development

### Database
- **PostgreSQL**: Primary data store
- **Migrations**: Version-controlled database schema
- **Tables**:
  - `stats`: DNS query statistics
  - `blocked_domains`: Domain blacklist with categories
  - `categories`: Domain categorization system

## Quick Start

### Prerequisites
- Go 1.23.5+
- Node.js 18+
- PostgreSQL 12+

### Backend Setup
```bash
cd backend
go mod download
go build -o oriondns ./cmd/dnsserver/main.go
```

### Frontend Setup
```bash
cd frontend
npm install
npm run dev
```

### Database Setup
```bash
# Run migrations
cd backend/migrations
# Configure your database connection in config files
```

### Configuration
Create configuration files based on the templates:
- `backend/config/development.yaml` - Development settings
- `backend/config/staging.yaml` - Production settings

Example configuration:
```yaml
db:
  host: localhost
  port: 5432
  user: postgres
  name: oriondns_dev
```

## Usage

### Running the DNS Server
```bash
./oriondns
```

### Running the Web Interface
```bash
cd frontend
npm run dev
```

### DNS Configuration
Configure your system or router to use OrionDNS as the primary DNS resolver.

## Development

### Project Structure
```
oriondns/
├── backend/
│   ├── cmd/                 # Application entry points
│   │   ├── dnsserver/      # DNS server binary
│   │   └── httpserver/     # HTTP API server
│   ├── internal/           # Internal packages
│   │   ├── ai/            # AI integration
│   │   ├── blockeddomains/ # Domain blocking logic
│   │   ├── categories/     # Domain categorization
│   │   ├── domains/        # Domain management
│   │   └── stats/          # Statistics collection
│   ├── server/             # Server implementations
│   │   ├── dns/           # DNS server logic
│   │   └── web/           # HTTP server and routes
│   ├── migrations/         # Database migrations
│   └── config/            # Configuration files
├── frontend/
│   ├── pages/             # Vue.js pages
│   ├── services/          # API client services
│   ├── utils/             # Utility functions
│   └── @types/           # TypeScript definitions
└── build.sh              # Build and deployment script
```

### Key Components

#### DNS Filtering (`backend/server/dns/dns.go`)
- Intercepts DNS queries on port 53
- Checks against blocked domains database
- Supports both exact and recursive (wildcard) matching
- Caches responses for performance
- Logs statistics for monitoring

#### Web API (`backend/server/web/`)
- RESTful API for domain management
- Statistics endpoints for dashboard
- CORS support for frontend integration

#### Dashboard (`frontend/pages/dashboard/`)
- Real-time DNS statistics
- Interactive charts and graphs
- Time range filtering
- Domain management interface

### Building for Production
```bash
# Build backend
cd backend
go build -o oriondns ./cmd/dnsserver/main.go

# Build frontend
cd frontend
npm run build

# Or use the provided build script
./build.sh
```

## Docker Support
Dockerfiles are provided for containerized deployment:
- `Dockerfile-frontend`: Frontend container
- `Dockerfile-web-backend`: Backend API container

## Contributing
1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## Security
OrionDNS is designed as a defensive security tool. It helps protect networks by:
- Blocking known malicious domains
- Monitoring DNS traffic for suspicious activity
- Providing visibility into network DNS patterns
- Categorizing domains for policy enforcement

## License
[Add your license information here]

## Support
[Add support information here]