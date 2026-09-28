interface Env {
  MAPS_DATA: R2Bucket;
}

type ObjectRoute = {
  key: string;
  immutable: boolean;
  contentType: string;
};

const apiPrefix = '/map-data/v1/';
const currentManifestPath = `${apiPrefix}manifest.json`;
const currentManifestKey = 'manifests/current.json';
const releaseIDPattern = /^[A-Za-z0-9][A-Za-z0-9._-]{0,79}$/;
const spritePattern = /^[A-Za-z0-9][A-Za-z0-9._@-]{0,119}\.(?:json|png)$/;
const rangePattern = /^\d{1,5}-\d{1,5}\.pbf$/;

const publicHeaders = (): Headers => {
  const headers = new Headers();
  headers.set('Access-Control-Allow-Origin', '*');
  headers.set('Cross-Origin-Resource-Policy', 'cross-origin');
  headers.set('Referrer-Policy', 'no-referrer');
  headers.set('X-Content-Type-Options', 'nosniff');
  return headers;
};

const jsonResponse = (status: number, body: Record<string, unknown>, cacheControl = 'no-store'): Response => {
  const headers = publicHeaders();
  headers.set('Content-Type', 'application/json; charset=utf-8');
  headers.set('Cache-Control', cacheControl);
  return new Response(JSON.stringify(body), { status, headers });
};

const safeDecode = (value: string): string | null => {
  try {
    return decodeURIComponent(value);
  } catch {
    return null;
  }
};

const validGlyphFontstack = (raw: string): boolean => {
  const decoded = safeDecode(raw);
  if (!decoded || decoded.length > 180) return false;
  if (/[\u0000-\u001f\u007f\\/]/u.test(decoded)) return false;
  return true;
};

const tileCoordinatesValid = (zRaw: string, xRaw: string, yFilename: string): boolean => {
  if (!/^\d{1,2}$/u.test(zRaw) || !/^\d{1,10}$/u.test(xRaw) || !/^\d{1,10}\.pbf$/u.test(yFilename)) {
    return false;
  }
  const z = Number(zRaw);
  const x = Number(xRaw);
  const y = Number(yFilename.slice(0, -4));
  if (!Number.isSafeInteger(z) || !Number.isSafeInteger(x) || !Number.isSafeInteger(y) || z < 0 || z > 24) {
    return false;
  }
  const maximum = 2 ** z - 1;
  return x >= 0 && x <= maximum && y >= 0 && y <= maximum;
};

const resolveObjectRoute = (pathname: string): ObjectRoute | null => {
  if (pathname === currentManifestPath) {
    return { key: currentManifestKey, immutable: false, contentType: 'application/json; charset=utf-8' };
  }
  if (!pathname.startsWith(apiPrefix)) return null;

  const key = pathname.slice(apiPrefix.length);
  if (key.length === 0 || key.length > 1024 || key.includes('//')) return null;
  const segments = key.split('/');
  if (segments.some((segment) => segment === '' || segment === '.' || segment === '..')) return null;
  if (segments[0] !== 'releases' || segments.length < 3) return null;

  const releaseID = segments[1];
  if (!releaseID || !releaseIDPattern.test(releaseID)) return null;

  if (segments.length === 3 && segments[2] === 'manifest.json') {
    return { key, immutable: true, contentType: 'application/json; charset=utf-8' };
  }
  if (segments.length === 3 && segments[2] === 'style.json') {
    return { key, immutable: true, contentType: 'application/json; charset=utf-8' };
  }
  if (
    segments.length === 6 &&
    segments[2] === 'tiles' &&
    segments[3] !== undefined &&
    segments[4] !== undefined &&
    segments[5] !== undefined &&
    tileCoordinatesValid(segments[3], segments[4], segments[5])
  ) {
    return { key, immutable: true, contentType: 'application/vnd.mapbox-vector-tile' };
  }
  if (segments.length === 4 && segments[2] === 'sprites' && segments[3] && spritePattern.test(segments[3])) {
    return {
      key,
      immutable: true,
      contentType: segments[3].endsWith('.png') ? 'image/png' : 'application/json; charset=utf-8',
    };
  }
  if (
    segments.length === 5 &&
    segments[2] === 'glyphs' &&
    segments[3] !== undefined &&
    segments[4] !== undefined &&
    validGlyphFontstack(segments[3]) &&
    rangePattern.test(segments[4])
  ) {
    return { key, immutable: true, contentType: 'application/x-protobuf' };
  }
  return null;
};

const objectResponse = (request: Request, object: R2ObjectBody, route: ObjectRoute): Response => {
  const headers = publicHeaders();
  object.writeHttpMetadata(headers);
  headers.set('Content-Type', headers.get('Content-Type') || route.contentType);
  headers.set('ETag', object.httpEtag);
  const cacheControl = route.immutable
    ? 'public, max-age=31536000, immutable'
    : 'public, max-age=60, stale-while-revalidate=300';
  headers.set('Cache-Control', cacheControl);
  headers.set('Cloudflare-CDN-Cache-Control', cacheControl);

  return new Response(request.method === 'HEAD' ? null : object.body, {
    status: 200,
    headers,
  });
};

export default {
  async fetch(request: Request, env: Env): Promise<Response> {
    const url = new URL(request.url);

    if (request.method === 'OPTIONS') {
      const headers = publicHeaders();
      headers.set('Access-Control-Allow-Methods', 'GET, HEAD, OPTIONS');
      headers.set('Access-Control-Max-Age', '86400');
      headers.set('Cache-Control', 'public, max-age=86400');
      return new Response(null, { status: 204, headers });
    }

    if (request.method !== 'GET' && request.method !== 'HEAD') {
      const response = jsonResponse(405, { error: 'method_not_allowed' });
      response.headers.set('Allow', 'GET, HEAD, OPTIONS');
      return response;
    }

    if (url.pathname === '/healthz') {
      return jsonResponse(200, { status: 'ok', service: 'goreecloud-maps-data' });
    }

    if (url.search !== '') {
      return jsonResponse(400, { error: 'query_parameters_not_supported' });
    }

    const route = resolveObjectRoute(url.pathname);
    if (!route) return jsonResponse(404, { error: 'not_found' });

    const object = await env.MAPS_DATA.get(route.key);
    if (!object) return jsonResponse(404, { error: 'not_found' });
    return objectResponse(request, object, route);
  },
} satisfies ExportedHandler<Env>;
