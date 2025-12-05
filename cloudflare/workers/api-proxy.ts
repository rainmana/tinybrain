/**
 * Cloudflare Worker: API Proxy
 * 
 * This worker provides:
 * - API request routing to Railway backend
 * - Edge caching for read operations
 * - Rate limiting
 * - Security headers
 * - Request/response transformation
 */

interface Env {
  API_URL: string;
  SUPABASE_URL: string;
  SUPABASE_ANON_KEY: string;
  CACHE: KVNamespace;
  ENVIRONMENT: string;
}

// CORS configuration
const CORS_HEADERS = {
  'Access-Control-Allow-Origin': '*',
  'Access-Control-Allow-Methods': 'GET, POST, PUT, DELETE, OPTIONS',
  'Access-Control-Allow-Headers': 'Content-Type, Authorization, X-Request-ID',
  'Access-Control-Max-Age': '86400',
};

// Security headers
const SECURITY_HEADERS = {
  'X-Content-Type-Options': 'nosniff',
  'X-Frame-Options': 'DENY',
  'X-XSS-Protection': '1; mode=block',
  'Referrer-Policy': 'strict-origin-when-cross-origin',
  'Permissions-Policy': 'accelerometer=(), camera=(), geolocation=(), microphone=()',
};

export default {
  async fetch(request: Request, env: Env, ctx: ExecutionContext): Promise<Response> {
    const url = new URL(request.url);

    // Handle OPTIONS preflight
    if (request.method === 'OPTIONS') {
      return handleOptions();
    }

    // Rate limiting
    const rateLimitResult = await checkRateLimit(request, env);
    if (!rateLimitResult.allowed) {
      return new Response('Rate limit exceeded', {
        status: 429,
        headers: {
          ...CORS_HEADERS,
          'Retry-After': rateLimitResult.retryAfter.toString(),
        },
      });
    }

    // Health check endpoint
    if (url.pathname === '/health') {
      return new Response(
        JSON.stringify({
          status: 'ok',
          timestamp: new Date().toISOString(),
          environment: env.ENVIRONMENT,
        }),
        {
          headers: {
            'Content-Type': 'application/json',
            ...CORS_HEADERS,
          },
        }
      );
    }

    // Route API requests to Railway backend
    if (url.pathname.startsWith('/api/')) {
      return handleAPIRequest(request, env, ctx);
    }

    // Default: 404
    return new Response('Not Found', {
      status: 404,
      headers: CORS_HEADERS,
    });
  },
};

/**
 * Handle OPTIONS preflight requests
 */
function handleOptions(): Response {
  return new Response(null, {
    status: 204,
    headers: CORS_HEADERS,
  });
}

/**
 * Check rate limit for the request
 */
async function checkRateLimit(
  request: Request,
  env: Env
): Promise<{ allowed: boolean; retryAfter: number }> {
  // Get client identifier (IP or API key)
  const clientId = request.headers.get('CF-Connecting-IP') || 'anonymous';
  const rateLimitKey = `ratelimit:${clientId}`;

  // Check if client has exceeded rate limit
  const count = await env.CACHE.get(rateLimitKey);
  const limit = 100; // requests per minute
  const window = 60; // seconds

  if (count && parseInt(count) >= limit) {
    return { allowed: false, retryAfter: window };
  }

  // Increment counter
  const newCount = count ? parseInt(count) + 1 : 1;
  await env.CACHE.put(rateLimitKey, newCount.toString(), { expirationTtl: window });

  return { allowed: true, retryAfter: 0 };
}

/**
 * Handle API requests by proxying to Railway backend
 */
async function handleAPIRequest(
  request: Request,
  env: Env,
  ctx: ExecutionContext
): Promise<Response> {
  const url = new URL(request.url);
  
  // Build backend URL
  const backendURL = new URL(url.pathname + url.search, env.API_URL);

  // Check cache for GET requests
  if (request.method === 'GET') {
    const cacheKey = `cache:${url.pathname}${url.search}`;
    const cached = await env.CACHE.get(cacheKey);
    
    if (cached) {
      return new Response(cached, {
        headers: {
          'Content-Type': 'application/json',
          'X-Cache': 'HIT',
          ...CORS_HEADERS,
          ...SECURITY_HEADERS,
        },
      });
    }
  }

  // Forward request to backend
  const backendRequest = new Request(backendURL.toString(), {
    method: request.method,
    headers: request.headers,
    body: request.method !== 'GET' && request.method !== 'HEAD' ? await request.blob() : undefined,
  });

  let response: Response;
  try {
    response = await fetch(backendRequest);
  } catch (error) {
    console.error('Backend request failed:', error);
    return new Response(
      JSON.stringify({
        error: 'Backend service unavailable',
        message: 'Unable to connect to API server',
      }),
      {
        status: 503,
        headers: {
          'Content-Type': 'application/json',
          ...CORS_HEADERS,
        },
      }
    );
  }

  // Clone response for caching
  const responseClone = response.clone();

  // Cache successful GET responses
  if (
    request.method === 'GET' &&
    response.ok &&
    response.headers.get('Content-Type')?.includes('application/json')
  ) {
    const body = await responseClone.text();
    const cacheKey = `cache:${url.pathname}${url.search}`;
    const cacheTTL = getCacheTTL(url.pathname);
    
    if (cacheTTL > 0) {
      ctx.waitUntil(
        env.CACHE.put(cacheKey, body, { expirationTtl: cacheTTL })
      );
    }
  }

  // Add custom headers
  const headers = new Headers(response.headers);
  Object.entries(CORS_HEADERS).forEach(([key, value]) => headers.set(key, value));
  Object.entries(SECURITY_HEADERS).forEach(([key, value]) => headers.set(key, value));
  headers.set('X-Cache', 'MISS');
  headers.set('X-Served-By', 'Cloudflare Workers');

  return new Response(response.body, {
    status: response.status,
    statusText: response.statusText,
    headers,
  });
}

/**
 * Determine cache TTL based on endpoint
 */
function getCacheTTL(pathname: string): number {
  // Don't cache auth endpoints
  if (pathname.includes('/auth/')) return 0;
  
  // Cache security data for longer
  if (pathname.includes('/security/')) return 3600; // 1 hour
  
  // Cache session lists briefly
  if (pathname.includes('/sessions')) return 60; // 1 minute
  
  // Cache memory searches briefly
  if (pathname.includes('/memories/search')) return 30; // 30 seconds
  
  // Default: short cache
  return 10; // 10 seconds
}
