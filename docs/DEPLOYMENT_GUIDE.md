# TinyBrain Web Deployment Guide

This guide provides step-by-step instructions for deploying TinyBrain's web-based version using Supabase, Railway.app, and Cloudflare Pages.

## Prerequisites

Before you begin, ensure you have:

- [ ] GitHub account
- [ ] Supabase account (https://supabase.com)
- [ ] Railway account (https://railway.app)
- [ ] Cloudflare account (https://cloudflare.com)
- [ ] Node.js 18+ and npm/yarn installed (for frontend development)
- [ ] Go 1.24+ installed (for backend development)
- [ ] Git installed

## Table of Contents

1. [Supabase Setup](#1-supabase-setup)
2. [Railway Backend Deployment](#2-railway-backend-deployment)
3. [Cloudflare Pages Frontend Deployment](#3-cloudflare-pages-frontend-deployment)
4. [Environment Configuration](#4-environment-configuration)
5. [Database Migration](#5-database-migration)
6. [Testing & Validation](#6-testing--validation)
7. [Monitoring & Maintenance](#7-monitoring--maintenance)
8. [Troubleshooting](#8-troubleshooting)

---

## 1. Supabase Setup

### 1.1 Create a New Supabase Project

1. Go to [Supabase Dashboard](https://app.supabase.com)
2. Click "New Project"
3. Fill in project details:
   - **Project Name**: `tinybrain-prod` (or your preferred name)
   - **Database Password**: Generate a strong password (save it securely!)
   - **Region**: Choose closest to your users
   - **Pricing Plan**: Start with Free tier

### 1.2 Run Database Migrations

1. Navigate to the SQL Editor in Supabase dashboard
2. Run the schema migration scripts from `supabase/migrations/` directory (see below)
3. Verify tables are created correctly

### 1.3 Configure Authentication

1. Go to **Authentication** → **Providers**
2. Enable desired authentication methods:
   - **Email**: Enable and configure email templates
   - **Google OAuth**: Add OAuth credentials
   - **GitHub OAuth**: Add OAuth credentials

3. Configure authentication settings:
   - Go to **Authentication** → **Settings**
   - Set **Site URL**: `https://your-app.pages.dev`
   - Add **Redirect URLs**: `https://your-app.pages.dev/auth/callback`
   - Enable **Email Confirmations** if desired

### 1.4 Set Up Row Level Security (RLS)

The migration scripts will create RLS policies automatically. Verify they are in place:

1. Go to **Authentication** → **Policies**
2. Check each table has appropriate policies
3. Test with different user roles

### 1.5 Configure Storage

1. Go to **Storage** → **Buckets**
2. Create buckets:
   - `attachments`: For memory attachments
   - `exports`: For data exports
   - `backups`: For backup files

3. Set bucket policies:
   ```sql
   -- Example policy for attachments bucket
   CREATE POLICY "Users can upload own attachments"
   ON storage.objects FOR INSERT
   WITH CHECK (
     bucket_id = 'attachments' AND
     auth.uid()::text = (storage.foldername(name))[1]
   );
   ```

### 1.6 Get API Credentials

1. Go to **Settings** → **API**
2. Save these credentials (you'll need them later):
   - **Project URL**: `https://xxxxx.supabase.co`
   - **Anon/Public Key**: `eyJhbGc...`
   - **Service Role Key**: `eyJhbGc...` (keep this secret!)
   - **Database Password**: (from step 1.1)
   - **PostgreSQL Connection String**: Available under "Connection string"

---

## 2. Railway Backend Deployment

### 2.1 Install Railway CLI (Optional)

```bash
npm install -g @railway/cli
railway login
```

### 2.2 Create New Railway Project

**Option A: Via Railway Dashboard**

1. Go to [Railway Dashboard](https://railway.app/dashboard)
2. Click "New Project"
3. Select "Deploy from GitHub repo"
4. Authorize GitHub and select your TinyBrain repository
5. Railway will auto-detect it's a Go project

**Option B: Via Railway CLI**

```bash
cd /path/to/tinybrain
railway init
railway up
```

### 2.3 Configure Build Settings

1. In Railway dashboard, go to your project
2. Click **Settings** → **Build**
3. Configure:
   - **Build Command**: `go build -o server ./cmd/tinybrain`
   - **Start Command**: `./server serve`
   - **Root Directory**: `/`

### 2.4 Add Environment Variables

In Railway dashboard, go to **Variables** and add:

```bash
# Supabase Configuration
SUPABASE_URL=https://xxxxx.supabase.co
SUPABASE_ANON_KEY=eyJhbGc...
SUPABASE_SERVICE_KEY=eyJhbGc...

# Database
DATABASE_URL=postgresql://postgres:[PASSWORD]@db.xxxxx.supabase.co:5432/postgres

# Server Configuration
TINYBRAIN_HTTP=0.0.0.0:$PORT
TINYBRAIN_ENV=production

# Security
JWT_SECRET=your-secure-jwt-secret-min-32-chars
CORS_ALLOWED_ORIGINS=https://your-app.pages.dev

# Optional: Feature Flags
ENABLE_REAL_TIME=true
ENABLE_MCP_ADAPTER=true
```

### 2.5 Configure Custom Domain (Optional)

1. Go to **Settings** → **Domains**
2. Click "Generate Domain" for a Railway domain
3. Or add your custom domain:
   - Add domain in Railway
   - Update DNS records as instructed
   - Wait for SSL certificate provisioning

### 2.6 Deploy

Railway will automatically deploy on every push to your main branch. Manual deployment:

```bash
railway up
```

Monitor deployment logs:

```bash
railway logs
```

### 2.7 Verify Deployment

Once deployed, test the API:

```bash
# Health check
curl https://your-api.railway.app/health

# MCP endpoint
curl -X POST https://your-api.railway.app/mcp \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":1,"method":"initialize"}'
```

---

## 3. Cloudflare Pages Frontend Deployment

### 3.1 Prepare Frontend Code

First, ensure you have a frontend application ready (see `web/` directory). If not yet created, you'll need to build it first.

### 3.2 Connect Repository to Cloudflare Pages

**Option A: Via Cloudflare Dashboard**

1. Go to [Cloudflare Dashboard](https://dash.cloudflare.com)
2. Navigate to **Pages**
3. Click "Create a project"
4. Select "Connect to Git"
5. Authorize GitHub and select your repository
6. Configure build settings:
   - **Framework preset**: Next.js (or your chosen framework)
   - **Build command**: `npm run build` or `cd web && npm run build`
   - **Build output directory**: `.next` (for Next.js) or `out`
   - **Root directory**: `/web` (if frontend is in subdirectory)

**Option B: Via Wrangler CLI**

```bash
npm install -g wrangler
wrangler login
cd web
wrangler pages publish out --project-name=tinybrain
```

### 3.3 Configure Environment Variables

In Cloudflare Pages dashboard, add environment variables:

```bash
# API Configuration
NEXT_PUBLIC_API_URL=https://your-api.railway.app
NEXT_PUBLIC_SUPABASE_URL=https://xxxxx.supabase.co
NEXT_PUBLIC_SUPABASE_ANON_KEY=eyJhbGc...

# Feature Flags
NEXT_PUBLIC_ENABLE_REAL_TIME=true
NEXT_PUBLIC_ENABLE_ANALYTICS=true
```

### 3.4 Set Up Cloudflare Workers (Optional)

For advanced edge functionality:

1. Create a Worker script (`web/workers/api-proxy.js`)
2. Configure routing:
   ```javascript
   // Example: API proxy with rate limiting
   addEventListener('fetch', event => {
     event.respondWith(handleRequest(event.request))
   })
   
   async function handleRequest(request) {
     // Add rate limiting, caching, etc.
     const response = await fetch(request)
     return response
   }
   ```

3. Deploy Worker:
   ```bash
   cd web/workers
   wrangler publish
   ```

### 3.5 Configure Custom Domain

1. In Cloudflare Pages settings, go to **Custom domains**
2. Click "Set up a custom domain"
3. Enter your domain (must be in your Cloudflare account)
4. Cloudflare will automatically configure DNS
5. Wait for SSL certificate provisioning (usually < 1 minute)

### 3.6 Configure CORS

Ensure your API (Railway) allows requests from your Cloudflare Pages domain:

In Railway environment variables:
```bash
CORS_ALLOWED_ORIGINS=https://your-app.pages.dev,https://yourdomain.com
```

### 3.7 Deploy and Verify

Cloudflare Pages deploys automatically on git push. Verify:

1. Visit your Pages URL: `https://your-app.pages.dev`
2. Test authentication flow
3. Test API connectivity
4. Check browser console for errors

---

## 4. Environment Configuration

### 4.1 Development Environment

Create `.env.local` for local development:

```bash
# Backend (Railway local)
SUPABASE_URL=http://localhost:54321
SUPABASE_ANON_KEY=your-local-anon-key
DATABASE_URL=postgresql://postgres:postgres@localhost:54322/postgres
TINYBRAIN_HTTP=127.0.0.1:8090

# Frontend (Cloudflare local)
NEXT_PUBLIC_API_URL=http://localhost:8090
NEXT_PUBLIC_SUPABASE_URL=http://localhost:54321
NEXT_PUBLIC_SUPABASE_ANON_KEY=your-local-anon-key
```

### 4.2 Staging Environment

Configure staging branches in Railway and Cloudflare:

- Railway: Create a separate environment for `staging` branch
- Cloudflare Pages: Preview deployments are automatic for all branches

### 4.3 Production Environment

Production environment variables are set in Railway and Cloudflare dashboards (see previous sections).

### 4.4 Secret Management

**Best Practices:**

1. Never commit secrets to git
2. Use environment-specific variables
3. Rotate secrets regularly
4. Use secret scanning tools (GitHub Advanced Security)

**Secret Rotation:**

```bash
# Generate new JWT secret
openssl rand -base64 32

# Update in Railway
railway variables --set JWT_SECRET=new-secret

# Update in application config
# Redeploy
```

---

## 5. Database Migration

### 5.1 Export from SQLite (Current Version)

If migrating from existing TinyBrain installation:

```bash
# Export data
sqlite3 ~/.tinybrain/memory.db .dump > tinybrain_export.sql

# Or use the export tool
./tinybrain export --output=tinybrain_export.json
```

### 5.2 Transform Data for PostgreSQL

Use the migration script:

```bash
go run scripts/migrate_to_postgres.go \
  --input=tinybrain_export.json \
  --output=postgres_import.sql
```

### 5.3 Import to Supabase

```bash
# Using psql
psql $DATABASE_URL -f postgres_import.sql

# Or via Supabase dashboard
# Go to SQL Editor and paste the import SQL
```

### 5.4 Verify Migration

```sql
-- Check record counts
SELECT 
  (SELECT COUNT(*) FROM memories) as memories,
  (SELECT COUNT(*) FROM sessions) as sessions,
  (SELECT COUNT(*) FROM relationships) as relationships;

-- Verify data integrity
SELECT * FROM memories LIMIT 10;
```

---

## 6. Testing & Validation

### 6.1 Backend API Tests

```bash
# Run integration tests
cd backend
go test -v ./...

# Test specific endpoints
curl https://your-api.railway.app/api/health
curl https://your-api.railway.app/api/sessions
```

### 6.2 Frontend Tests

```bash
cd web
npm run test

# E2E tests with Playwright
npm run test:e2e
```

### 6.3 Load Testing

```bash
# Using k6
k6 run tests/load_test.js

# Or artillery
artillery run tests/load_test.yml
```

### 6.4 Security Testing

```bash
# Run OWASP ZAP scan
docker run -v $(pwd):/zap/wrk/:rw \
  -t owasp/zap2docker-stable zap-baseline.py \
  -t https://your-app.pages.dev \
  -r zap_report.html

# Check for security headers
curl -I https://your-app.pages.dev
```

---

## 7. Monitoring & Maintenance

### 7.1 Set Up Monitoring

**Railway Monitoring:**

1. Go to **Observability** in Railway dashboard
2. View metrics: CPU, Memory, Network
3. Set up alerts for high resource usage

**Supabase Monitoring:**

1. Go to **Reports** in Supabase dashboard
2. Monitor database size, query performance
3. Check API usage and rate limits

**Cloudflare Analytics:**

1. Go to **Analytics** in Cloudflare dashboard
2. Monitor page views, bandwidth, errors
3. Set up alerts for increased error rates

### 7.2 Logging

**Centralized Logging:**

```bash
# View Railway logs
railway logs --follow

# Export logs for analysis
railway logs --json > logs.json
```

**Log Aggregation:**

Consider using:
- Datadog
- New Relic
- Sentry (for error tracking)

### 7.3 Backup Strategy

**Database Backups:**

Supabase provides automatic backups on Pro tier. Additionally:

```bash
# Manual backup
pg_dump $DATABASE_URL > backup_$(date +%Y%m%d).sql

# Restore from backup
psql $DATABASE_URL < backup_20231201.sql
```

**Application Backups:**

```bash
# Backup application data
curl https://your-api.railway.app/api/export/full > backup.json
```

### 7.4 Performance Optimization

**Database:**
- Monitor slow queries in Supabase dashboard
- Add indexes for frequently queried fields
- Optimize query patterns

**API:**
- Enable caching for read-heavy endpoints
- Use connection pooling
- Implement rate limiting

**Frontend:**
- Enable Cloudflare CDN caching
- Optimize images (next/image)
- Code splitting and lazy loading

---

## 8. Troubleshooting

### Common Issues

#### Issue: API not connecting to database

**Symptoms:** 500 errors, "connection refused" logs

**Solutions:**
1. Check DATABASE_URL is correct
2. Verify database is running (Supabase dashboard)
3. Check firewall rules allow Railway's IP ranges
4. Test connection manually:
   ```bash
   psql $DATABASE_URL -c "SELECT version();"
   ```

#### Issue: CORS errors in frontend

**Symptoms:** Browser console shows "CORS policy" errors

**Solutions:**
1. Add frontend domain to `CORS_ALLOWED_ORIGINS` in Railway
2. Verify API is responding with correct CORS headers:
   ```bash
   curl -H "Origin: https://your-app.pages.dev" \
        -H "Access-Control-Request-Method: GET" \
        -X OPTIONS \
        https://your-api.railway.app/api/sessions
   ```
3. Redeploy Railway backend after updating environment variables

#### Issue: Authentication not working

**Symptoms:** Login fails, JWT errors

**Solutions:**
1. Verify Supabase Auth is enabled
2. Check JWT_SECRET matches in Railway and frontend
3. Verify redirect URLs are configured in Supabase
4. Check browser cookies are enabled
5. Test authentication flow:
   ```bash
   curl -X POST https://xxxxx.supabase.co/auth/v1/signup \
     -H "apikey: your-anon-key" \
     -H "Content-Type: application/json" \
     -d '{"email":"test@example.com","password":"password123"}'
   ```

#### Issue: Slow API responses

**Symptoms:** Long response times, timeouts

**Solutions:**
1. Check database query performance in Supabase
2. Add missing indexes:
   ```sql
   CREATE INDEX idx_memories_user_id ON memories(user_id);
   CREATE INDEX idx_memories_session_id ON memories(session_id);
   ```
3. Enable query result caching
4. Scale Railway instance (increase resources)

#### Issue: Build failures

**Symptoms:** Deployment fails, build errors

**Solutions:**
1. Check Railway build logs
2. Verify Go version matches requirements (1.24+)
3. Run build locally first:
   ```bash
   go build -o server ./cmd/tinybrain
   ```
4. Check all dependencies are in go.mod
5. Clear Railway build cache

### Getting Help

- **GitHub Issues**: https://github.com/rainmana/tinybrain/issues
- **Discussions**: https://github.com/rainmana/tinybrain/discussions
- **Railway Discord**: https://discord.gg/railway
- **Supabase Discord**: https://discord.supabase.com

---

## Quick Reference

### Useful Commands

```bash
# Railway
railway login
railway link
railway up
railway logs
railway variables

# Supabase CLI
npx supabase login
npx supabase init
npx supabase start
npx supabase db push
npx supabase db reset

# Wrangler (Cloudflare)
wrangler login
wrangler pages publish
wrangler tail
```

### Important URLs

- Supabase Dashboard: https://app.supabase.com
- Railway Dashboard: https://railway.app/dashboard
- Cloudflare Dashboard: https://dash.cloudflare.com
- API Endpoint: https://your-api.railway.app
- Frontend URL: https://your-app.pages.dev

### Support Matrix

| Component | Free Tier | Recommended Tier | Cost |
|-----------|-----------|------------------|------|
| Supabase | 500MB DB, 1GB storage | Pro ($25/mo) | $25/mo |
| Railway | $5 credit/mo | Usage-based | $20-50/mo |
| Cloudflare Pages | Unlimited | Free | $0 |
| **Total** | Limited | Production-ready | $45-75/mo |

---

## Next Steps

After successful deployment:

1. [ ] Set up monitoring and alerts
2. [ ] Configure automated backups
3. [ ] Set up CI/CD pipelines
4. [ ] Implement feature flags
5. [ ] Create runbooks for common operations
6. [ ] Document team workflows
7. [ ] Plan capacity and scaling strategy
8. [ ] Set up staging environment
9. [ ] Create disaster recovery plan
10. [ ] Train team on new architecture

For more information, see:
- [Architecture Documentation](./WEB_ARCHITECTURE.md)
- [API Documentation](./API_REFERENCE.md)
- [Frontend Development Guide](./FRONTEND_GUIDE.md)
