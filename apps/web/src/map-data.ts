import type { StyleSpecification } from 'maplibre-gl';

export type MapDataAttribution = {
  text: string;
  url?: string;
};

export type MapDataSource = {
  name: string;
  datasetVersion: string;
  license: string;
  provenance?: string;
};

export type MapDataManifest = {
  schemaVersion: 1;
  releaseId: string;
  generatedAt: string;
  stylePath: string;
  tileTemplate: string;
  minZoom: number;
  maxZoom: number;
  bounds: [number, number, number, number];
  attribution: MapDataAttribution[];
  sources: MapDataSource[];
  publicGeographicDataOnly: true;
};

export type MapDataRelease = {
  manifest: MapDataManifest;
  manifestURL: string;
  styleURL: string;
  style: StyleSpecification;
};

const maximumManifestBytes = 64 * 1024;
const maximumStyleBytes = 2 * 1024 * 1024;
const releaseIDPattern = /^[A-Za-z0-9][A-Za-z0-9._-]{0,79}$/;
const spriteNamePattern = /^[A-Za-z0-9][A-Za-z0-9._@-]{0,119}$/;
const spriteIDPattern = /^[A-Za-z0-9][A-Za-z0-9._@:-]{0,119}$/;
const allowedManifestKeys = new Set([
  'schemaVersion',
  'releaseId',
  'generatedAt',
  'stylePath',
  'tileTemplate',
  'minZoom',
  'maxZoom',
  'bounds',
  'attribution',
  'sources',
  'publicGeographicDataOnly',
]);

const isRecord = (value: unknown): value is Record<string, unknown> =>
  typeof value === 'object' && value !== null && !Array.isArray(value);

const nonEmptyString = (value: unknown, maximumLength: number): value is string =>
  typeof value === 'string' && value.trim() !== '' && value.length <= maximumLength;

const validateConfiguredURL = (raw: string): URL => {
  let parsed: URL;
  try {
    parsed = new URL(raw, window.location.origin);
  } catch {
    throw new Error('Map-data manifest URL is invalid.');
  }

  const localDevelopment =
    parsed.protocol === 'http:' && (parsed.hostname === 'localhost' || parsed.hostname === '127.0.0.1');
  if (parsed.protocol !== 'https:' && !localDevelopment) {
    throw new Error('Map-data manifest must use HTTPS outside local development.');
  }
  if (parsed.username || parsed.password || parsed.search || parsed.hash) {
    throw new Error('Map-data manifest URL must not contain credentials, query parameters, or fragments.');
  }
  return parsed;
};

const validateAttribution = (value: unknown): value is MapDataAttribution[] => {
  if (!Array.isArray(value) || value.length === 0 || value.length > 32) return false;
  return value.every((entry) => {
    if (!isRecord(entry)) return false;
    if (Object.keys(entry).some((key) => key !== 'text' && key !== 'url')) return false;
    if (!nonEmptyString(entry.text, 500)) return false;
    if (entry.url !== undefined) {
      if (!nonEmptyString(entry.url, 1000)) return false;
      try {
        const parsed = new URL(entry.url);
        if (parsed.protocol !== 'https:') return false;
      } catch {
        return false;
      }
    }
    return true;
  });
};

const validateSources = (value: unknown): value is MapDataSource[] => {
  if (!Array.isArray(value) || value.length === 0 || value.length > 64) return false;
  return value.every((entry) => {
    if (!isRecord(entry)) return false;
    if (Object.keys(entry).some((key) => !['name', 'datasetVersion', 'license', 'provenance'].includes(key))) {
      return false;
    }
    if (!nonEmptyString(entry.name, 160)) return false;
    if (!nonEmptyString(entry.datasetVersion, 160)) return false;
    if (!nonEmptyString(entry.license, 160)) return false;
    if (entry.provenance !== undefined && !nonEmptyString(entry.provenance, 1000)) return false;
    return true;
  });
};

const validateBounds = (value: unknown): value is [number, number, number, number] => {
  if (!Array.isArray(value) || value.length !== 4) return false;
  if (value.some((coordinate) => typeof coordinate !== 'number' || !Number.isFinite(coordinate))) return false;
  const [west, south, east, north] = value as number[];
  if (west === undefined || south === undefined || east === undefined || north === undefined) return false;
  if (west < -180 || west > 180 || east < -180 || east > 180) return false;
  if (south < -90 || south > 90 || north < -90 || north > 90) return false;
  return west <= east && south <= north;
};

const parseManifest = (value: unknown): MapDataManifest => {
  if (!isRecord(value)) throw new Error('Map-data manifest must be a JSON object.');
  if (Object.keys(value).some((key) => !allowedManifestKeys.has(key))) {
    throw new Error('Map-data manifest contains unsupported fields.');
  }
  if (value.schemaVersion !== 1) throw new Error('Map-data manifest schema version is unsupported.');
  if (typeof value.releaseId !== 'string' || !releaseIDPattern.test(value.releaseId)) {
    throw new Error('Map-data release identifier is invalid.');
  }
  if (!nonEmptyString(value.generatedAt, 64) || !Number.isFinite(Date.parse(value.generatedAt))) {
    throw new Error('Map-data generation timestamp is invalid.');
  }

  const expectedStyle = `releases/${value.releaseId}/style.json`;
  const expectedTiles = `releases/${value.releaseId}/tiles/{z}/{x}/{y}.pbf`;
  if (value.stylePath !== expectedStyle || value.tileTemplate !== expectedTiles) {
    throw new Error('Map-data release paths do not match the release identifier.');
  }
  if (!Number.isInteger(value.minZoom) || !Number.isInteger(value.maxZoom)) {
    throw new Error('Map-data zoom range is invalid.');
  }
  const minZoom = value.minZoom as number;
  const maxZoom = value.maxZoom as number;
  if (minZoom < 0 || maxZoom > 24 || minZoom > maxZoom) {
    throw new Error('Map-data zoom range is outside the supported bounds.');
  }
  if (!validateBounds(value.bounds)) throw new Error('Map-data geographic bounds are invalid.');
  if (!validateAttribution(value.attribution)) throw new Error('Map-data attribution contract is invalid.');
  if (!validateSources(value.sources)) throw new Error('Map-data source/provenance contract is invalid.');
  if (value.publicGeographicDataOnly !== true) {
    throw new Error('Map-data release is not marked as public geographic data only.');
  }

  return value as MapDataManifest;
};

const readBoundedJSON = async (url: URL, maximumBytes: number, label: string, cache: RequestCache): Promise<unknown> => {
  const response = await fetch(url, {
    method: 'GET',
    headers: { Accept: 'application/json' },
    credentials: 'omit',
    cache,
    redirect: 'error',
  });
  if (!response.ok) throw new Error(`${label} request failed with status ${response.status}.`);

  const contentLength = Number(response.headers.get('Content-Length'));
  if (Number.isFinite(contentLength) && contentLength > maximumBytes) {
    throw new Error(`${label} is too large.`);
  }

  const raw = await response.text();
  if (new TextEncoder().encode(raw).byteLength > maximumBytes) {
    throw new Error(`${label} is too large.`);
  }
  try {
    return JSON.parse(raw) as unknown;
  } catch {
    throw new Error(`${label} is not valid JSON.`);
  }
};

const escapeHTML = (value: string): string =>
  value.replace(/[&<>'"]/g, (character) => {
    const entities: Record<string, string> = {
      '&': '&amp;',
      '<': '&lt;',
      '>': '&gt;',
      "'": '&#39;',
      '"': '&quot;',
    };
    return entities[character] ?? character;
  });

const manifestAttribution = (manifest: MapDataManifest): string =>
  manifest.attribution
    .map((entry) =>
      entry.url
        ? `<a href="${escapeHTML(entry.url)}" target="_blank" rel="noopener noreferrer">${escapeHTML(entry.text)}</a>`
        : escapeHTML(entry.text),
    )
    .join(' · ');

const absoluteReleaseResource = (raw: string, styleURL: URL, allowedRelative: RegExp, label: string): string => {
  if (allowedRelative.test(raw)) return `${styleURL.toString().slice(0, styleURL.toString().lastIndexOf('/') + 1)}${raw}`;

  let parsed: URL;
  try {
    parsed = new URL(raw);
  } catch {
    throw new Error(`${label} must remain within the immutable map-data release.`);
  }
  if (parsed.origin !== styleURL.origin || parsed.username || parsed.password || parsed.search || parsed.hash) {
    throw new Error(`${label} escaped the configured map-data origin.`);
  }
  const releasePrefix = styleURL.pathname.slice(0, styleURL.pathname.lastIndexOf('/') + 1);
  if (!parsed.pathname.startsWith(releasePrefix)) {
    throw new Error(`${label} escaped the immutable map-data release path.`);
  }
  return parsed.toString();
};

const normalizeReleaseStyle = (value: unknown, manifest: MapDataManifest, styleURL: URL): StyleSpecification => {
  if (!isRecord(value) || value.version !== 8 || !isRecord(value.sources) || !Array.isArray(value.layers)) {
    throw new Error('Map-data style is not a valid MapLibre style foundation.');
  }
  if ('imports' in value || 'font-faces' in value) {
    throw new Error('Map-data style v1 does not permit imported styles or external font-face resources.');
  }

  const releaseBase = styleURL.toString().slice(0, styleURL.toString().lastIndexOf('/') + 1);
  const expectedTiles = `${releaseBase}tiles/{z}/{x}/{y}.pbf`;
  const attribution = manifestAttribution(manifest);
  let firstVectorSource = true;

  for (const [sourceID, rawSource] of Object.entries(value.sources)) {
    if (!isRecord(rawSource) || rawSource.type !== 'vector') {
      throw new Error(`Map-data style source ${sourceID} must be a vector source in schema v1.`);
    }
    if ('url' in rawSource) {
      throw new Error(`Map-data style source ${sourceID} must use the release tile template, not TileJSON.`);
    }
    if (!Array.isArray(rawSource.tiles) || rawSource.tiles.length !== 1 || typeof rawSource.tiles[0] !== 'string') {
      throw new Error(`Map-data style source ${sourceID} must contain exactly one tile template.`);
    }
    const tileTemplate = rawSource.tiles[0];
    if (tileTemplate !== 'tiles/{z}/{x}/{y}.pbf' && tileTemplate !== expectedTiles) {
      throw new Error(`Map-data style source ${sourceID} does not use the declared release tiles.`);
    }
    rawSource.tiles = [expectedTiles];
    rawSource.minzoom = manifest.minZoom;
    rawSource.maxzoom = manifest.maxZoom;
    rawSource.bounds = [...manifest.bounds];
    rawSource.scheme = 'xyz';
    rawSource.encoding = 'mvt';
    if (firstVectorSource) {
      rawSource.attribution = attribution;
      firstVectorSource = false;
    } else {
      delete rawSource.attribution;
    }
  }

  if (value.glyphs !== undefined) {
    if (typeof value.glyphs !== 'string') throw new Error('Map-data glyph URL is invalid.');
    value.glyphs = absoluteReleaseResource(
      value.glyphs,
      styleURL,
      /^glyphs\/\{fontstack\}\/\{range\}\.pbf$/u,
      'Map-data glyph URL',
    );
  }

  if (value.sprite !== undefined) {
    if (typeof value.sprite === 'string') {
      if (!/^sprites\/[A-Za-z0-9][A-Za-z0-9._@-]{0,119}$/u.test(value.sprite)) {
        value.sprite = absoluteReleaseResource(value.sprite, styleURL, /^sprites\/[A-Za-z0-9][A-Za-z0-9._@-]{0,119}$/u, 'Map-data sprite URL');
      } else {
        value.sprite = `${releaseBase}${value.sprite}`;
      }
    } else if (Array.isArray(value.sprite)) {
      value.sprite = value.sprite.map((entry) => {
        if (!isRecord(entry) || typeof entry.id !== 'string' || !spriteIDPattern.test(entry.id) || typeof entry.url !== 'string') {
          throw new Error('Map-data sprite entry is invalid.');
        }
        const relativeMatch = entry.url.match(/^sprites\/([A-Za-z0-9][A-Za-z0-9._@-]{0,119})$/u);
        const normalizedURL = relativeMatch && relativeMatch[1] && spriteNamePattern.test(relativeMatch[1])
          ? `${releaseBase}${entry.url}`
          : absoluteReleaseResource(entry.url, styleURL, /^sprites\/[A-Za-z0-9][A-Za-z0-9._@-]{0,119}$/u, 'Map-data sprite URL');
        return { id: entry.id, url: normalizedURL };
      });
    } else {
      throw new Error('Map-data sprite configuration is invalid.');
    }
  }

  return value as unknown as StyleSpecification;
};

export const configuredMapDataManifestURL = (): string =>
  import.meta.env.VITE_MAP_DATA_MANIFEST_URL?.trim() ?? '';

export const resolveMapDataRelease = async (): Promise<MapDataRelease | null> => {
  const configured = configuredMapDataManifestURL();
  if (!configured) return null;

  const requestedURL = validateConfiguredURL(configured);
  const manifest = parseManifest(await readBoundedJSON(requestedURL, maximumManifestBytes, 'Map-data manifest', 'no-store'));

  const manifestURL = requestedURL.toString();
  const styleURL = new URL(manifest.stylePath, manifestURL);
  if (styleURL.origin !== requestedURL.origin || styleURL.username || styleURL.password || styleURL.search || styleURL.hash) {
    throw new Error('Map-data style path escaped the configured release origin.');
  }

  const style = normalizeReleaseStyle(
    await readBoundedJSON(styleURL, maximumStyleBytes, 'Map-data style', 'force-cache'),
    manifest,
    styleURL,
  );

  return {
    manifest,
    manifestURL,
    styleURL: styleURL.toString(),
    style,
  };
};
