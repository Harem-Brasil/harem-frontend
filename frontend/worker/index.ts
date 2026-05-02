interface WorkerEnv {
  API_URL: string
  APP_ENV?: string
  COMMIT_HASH?: string
  ASSETS?: { fetch: typeof fetch }
}

export default {
  async fetch(request, env: WorkerEnv) {
    const url = new URL(request.url);

    if (url.pathname.startsWith("/api/")) {
      if (!env.API_URL) {
        return new Response(JSON.stringify({ error: 'API_URL not configured' }), {
          status: 500,
          headers: { 'Content-Type': 'application/json' },
        });
      }
      const backendPath = url.pathname.replace(/^\/api/, '') + url.search;
      const apiUrl = new URL(backendPath, env.API_URL);
      const modifiedRequest = new Request(apiUrl, {
        method: request.method,
        headers: request.headers,
        body: request.body,
      });
      return fetch(modifiedRequest);
    }

    let res: Response;

    if (env.ASSETS) {
      // Production / wrangler dev: ASSETS binding serves static files with SPA fallback
      res = await env.ASSETS.fetch(request);
      if (res.status === 404) {
        const indexUrl = new URL('/', url.origin);
        res = await env.ASSETS.fetch(new Request(indexUrl.toString()));
      }
    } else {
      // Vite dev mode: ASSETS is not injected — proxy back to Vite's own server at '/'
      // which serves index.html, letting the React SPA handle routing client-side.
      const indexUrl = new URL('/', url.origin);
      res = await fetch(new Request(indexUrl.toString(), { headers: request.headers }));
    }

    // Clone response to add custom headers
    const modified = new Response(res.body, {
      status: res.status,
      statusText: res.statusText,
      headers: res.headers,
    });
    modified.headers.set('X-Environment', env.APP_ENV || 'development');
    if (env.COMMIT_HASH) {
      modified.headers.set('X-Commit-Hash', env.COMMIT_HASH);
    }
    return modified;
  },
} satisfies ExportedHandler<WorkerEnv>;

