# TinyBrain Web Implementation

Welcome to the TinyBrain web-based implementation! This version transforms TinyBrain from a local Go MCP server into a full-stack web application using modern cloud infrastructure.

## 🎯 Overview

This implementation provides:

- **Cloud-Native Architecture**: Distributed deployment across Supabase, Railway, and Cloudflare
- **Multi-User Support**: Team collaboration and user management
- **Web Dashboard**: Rich web interface for managing memories and sessions
- **Real-Time Features**: Live updates and notifications
- **Scalability**: Auto-scaling infrastructure for growing needs
- **Enhanced Security**: Row-level security, authentication, and authorization

## 📁 Repository Structure

```
tinybrain/
├── docs/                           # Documentation
│   ├── WEB_ARCHITECTURE.md        # Architecture overview
│   ├── DEPLOYMENT_GUIDE.md        # Step-by-step deployment
│   └── WEB_IMPLEMENTATION_README.md # This file
│
├── supabase/                       # Supabase configuration
│   └── migrations/                # Database migrations
│       ├── 001_initial_schema.sql     # Core tables
│       └── 002_row_level_security.sql # RLS policies
│
├── railway/                        # Railway configuration
│   └── Dockerfile                 # Railway-optimized Docker image
│
├── cloudflare/                     # Cloudflare configuration
│   ├── wrangler.toml              # Workers configuration
│   └── workers/                   # Edge functions
│       └── api-proxy.ts           # API proxy worker
│
├── web/                            # Frontend application (to be created)
│   ├── src/                       # Source code
│   ├── public/                    # Static assets
│   ├── package.json               # Dependencies
│   └── next.config.js             # Next.js configuration
│
├── cmd/tinybrain/                  # Backend application
│   └── main.go                    # Main server code
│
├── .env.example                    # Environment variables template
├── railway.toml                    # Railway deployment config
└── README.md                       # Main project README
```

## 🚀 Quick Start

### Prerequisites

- Node.js 18+ and npm/yarn
- Go 1.24+
- Git
- Accounts on:
  - [Supabase](https://supabase.com)
  - [Railway](https://railway.app)
  - [Cloudflare](https://cloudflare.com)

### Local Development Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/rainmana/tinybrain.git
   cd tinybrain
   ```

2. **Set up Supabase**
   ```bash
   # Install Supabase CLI
   npm install -g supabase

   # Initialize Supabase (optional for local dev)
   npx supabase init
   npx supabase start

   # Or use your cloud Supabase project
   ```

3. **Configure environment variables**
   ```bash
   cp .env.example .env.local
   # Edit .env.local with your Supabase credentials
   ```

4. **Run database migrations**
   ```bash
   # Apply migrations to your Supabase project
   psql $DATABASE_URL -f supabase/migrations/001_initial_schema.sql
   psql $DATABASE_URL -f supabase/migrations/002_row_level_security.sql
   ```

5. **Start the backend**
   ```bash
   go run ./cmd/tinybrain serve
   ```

6. **Start the frontend** (once created)
   ```bash
   cd web
   npm install
   npm run dev
   ```

7. **Access the application**
   - Backend API: http://localhost:8090
   - Frontend: http://localhost:3000
   - Supabase Dashboard: https://app.supabase.com

## 📦 Deployment

### Production Deployment

Follow the comprehensive [Deployment Guide](./DEPLOYMENT_GUIDE.md) for step-by-step instructions.

**Quick Summary:**

1. **Supabase**: Create project, run migrations, configure auth
2. **Railway**: Connect repo, set environment variables, deploy backend
3. **Cloudflare Pages**: Connect repo, configure build, deploy frontend

### Environment-Specific Deployments

- **Development**: Local Supabase + local backend + local frontend
- **Staging**: Staging Supabase + Railway staging + Cloudflare preview
- **Production**: Production Supabase + Railway production + Cloudflare production

## 🏗️ Architecture

See [WEB_ARCHITECTURE.md](./WEB_ARCHITECTURE.md) for detailed architecture documentation.

### High-Level Flow

```
User Browser
    ↓
Cloudflare Pages (Frontend + Workers)
    ↓
Railway.app (Backend API)
    ↓
Supabase (Database + Auth + Storage)
```

### Key Components

1. **Supabase**: PostgreSQL database with auth, storage, and real-time
2. **Railway**: Go backend API with MCP compatibility
3. **Cloudflare Pages**: Static frontend with edge workers

## 🔐 Security

### Authentication

- Supabase Auth handles user authentication
- Support for email/password, OAuth (Google, GitHub)
- JWT tokens for API access

### Authorization

- Row-Level Security (RLS) in Supabase
- User-level and team-level data isolation
- API-level permission checks

### Data Protection

- TLS 1.3 for all connections
- At-rest encryption in Supabase
- Environment variable-based secrets

## 🎨 Frontend (To Be Implemented)

The frontend will be built with:

- **Framework**: Next.js 14+ with App Router
- **Language**: TypeScript
- **Styling**: Tailwind CSS
- **State**: React Context + Supabase client
- **Real-time**: Supabase Realtime subscriptions

### Planned Features

- [ ] Dashboard with session overview
- [ ] Memory browser with search and filters
- [ ] Session management
- [ ] Real-time collaboration
- [ ] MITRE ATT&CK visualization
- [ ] Security knowledge hub
- [ ] User settings and team management

## 🔧 Development

### Backend Development

```bash
# Run tests
go test -v ./...

# Build
go build -o server ./cmd/tinybrain

# Format code
go fmt ./...

# Lint
go vet ./...
```

### Frontend Development

```bash
cd web

# Install dependencies
npm install

# Development server
npm run dev

# Build for production
npm run build

# Run tests
npm run test

# Lint and format
npm run lint
npm run format
```

### Database Development

```bash
# Create new migration
# Create file: supabase/migrations/003_your_migration.sql

# Apply migration
psql $DATABASE_URL -f supabase/migrations/003_your_migration.sql

# Reset database (development only!)
npx supabase db reset
```

## 📊 Monitoring

### Production Monitoring

- **Railway**: Built-in metrics and logs
- **Supabase**: Database metrics and query performance
- **Cloudflare**: Analytics and error tracking

### Recommended Tools

- **Error Tracking**: Sentry
- **Logging**: Railway logs + external aggregator
- **APM**: Datadog or New Relic
- **Uptime**: UptimeRobot or Pingdom

## 🧪 Testing

### Backend Tests

```bash
# Unit tests
go test ./internal/...

# Integration tests
go test -tags=integration ./test/...

# Coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Frontend Tests

```bash
cd web

# Unit tests
npm run test

# E2E tests
npm run test:e2e

# Coverage
npm run test:coverage
```

## 📝 API Documentation

### REST API Endpoints

The backend provides RESTful APIs for all operations:

```
GET    /health                      - Health check
POST   /api/auth/login              - User login
POST   /api/auth/signup             - User signup
GET    /api/sessions                - List sessions
POST   /api/sessions                - Create session
GET    /api/sessions/:id            - Get session
PUT    /api/sessions/:id            - Update session
DELETE /api/sessions/:id            - Delete session
GET    /api/memories                - List memories
POST   /api/memories                - Create memory
GET    /api/memories/:id            - Get memory
PUT    /api/memories/:id            - Update memory
DELETE /api/memories/:id            - Delete memory
POST   /api/memories/search         - Search memories
GET    /api/security/nvd            - Query NVD data
GET    /api/security/mitre          - Query MITRE ATT&CK
GET    /api/security/owasp          - Query OWASP data
```

### MCP Protocol Compatibility

The backend maintains MCP protocol compatibility:

```
POST   /mcp                         - MCP JSON-RPC endpoint
```

## 🔄 Migration from Local Version

To migrate from the local TinyBrain version:

1. **Export existing data**
   ```bash
   ./tinybrain export --output=data.json
   ```

2. **Set up cloud infrastructure** (Supabase, Railway, Cloudflare)

3. **Import data to Supabase**
   ```bash
   # Use the migration script (to be created)
   go run scripts/migrate_to_postgres.go --input=data.json
   ```

4. **Verify data integrity**

5. **Update MCP client configuration** to point to new API

## 💡 Tips & Best Practices

### Development

- Use `.env.local` for local development, never commit it
- Test RLS policies thoroughly before deploying
- Use Supabase local development for faster iteration
- Implement feature flags for gradual rollouts

### Deployment

- Always test in staging before production
- Use Railway preview deployments for PRs
- Monitor logs after each deployment
- Keep secrets secure and rotate regularly

### Performance

- Use Cloudflare Workers for edge caching
- Implement pagination for large datasets
- Add database indexes for common queries
- Monitor slow queries in Supabase dashboard

## 🆘 Troubleshooting

See [DEPLOYMENT_GUIDE.md](./DEPLOYMENT_GUIDE.md) for detailed troubleshooting steps.

Common issues:

- **CORS errors**: Check `CORS_ALLOWED_ORIGINS` in Railway
- **Auth failures**: Verify JWT secrets match across services
- **Database connection**: Check `DATABASE_URL` and firewall rules
- **Build failures**: Ensure Go version 1.24+ and all dependencies

## 🤝 Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

See [CONTRIBUTING.md](../CONTRIBUTING.md) for detailed guidelines.

## 📄 License

MIT License - see [LICENSE](../LICENSE) file for details.

## 🔗 Links

- **Documentation**: https://rainmana.github.io/tinybrain/
- **GitHub**: https://github.com/rainmana/tinybrain
- **Issues**: https://github.com/rainmana/tinybrain/issues
- **Discussions**: https://github.com/rainmana/tinybrain/discussions

## 📞 Support

- GitHub Issues: For bugs and feature requests
- GitHub Discussions: For questions and community support
- Email: [Your support email]

---

**Status**: 🚧 In Development

This web implementation is currently under active development. The infrastructure configurations are ready, and the implementation is in progress. Check back for updates or contribute to help move it forward!
