import { readFile } from 'node:fs/promises';
import process from 'node:process';

const releaseIDPattern = /^[A-Za-z0-9][A-Za-z0-9._-]{0,79}$/;
const allowedKeys = new Set([
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

const fail = (message) => {
  throw new Error(message);
};

const nonEmptyString = (value, label, max = 1000) => {
  if (typeof value !== 'string' || value.trim() === '' || value.length > max) fail(`${label} must be a non-empty string.`);
};

const validateManifest = (manifest) => {
  if (!manifest || typeof manifest !== 'object' || Array.isArray(manifest)) fail('Manifest must be a JSON object.');
  for (const key of Object.keys(manifest)) {
    if (!allowedKeys.has(key)) fail(`Unsupported manifest field: ${key}`);
  }
  if (manifest.schemaVersion !== 1) fail('schemaVersion must equal 1.');
  if (typeof manifest.releaseId !== 'string' || !releaseIDPattern.test(manifest.releaseId)) fail('releaseId is invalid.');
  nonEmptyString(manifest.generatedAt, 'generatedAt', 64);
  if (!Number.isFinite(Date.parse(manifest.generatedAt))) fail('generatedAt must be a valid date-time.');

  const expectedStyle = `releases/${manifest.releaseId}/style.json`;
  const expectedTiles = `releases/${manifest.releaseId}/tiles/{z}/{x}/{y}.pbf`;
  if (manifest.stylePath !== expectedStyle) fail(`stylePath must equal ${expectedStyle}.`);
  if (manifest.tileTemplate !== expectedTiles) fail(`tileTemplate must equal ${expectedTiles}.`);

  if (!Number.isInteger(manifest.minZoom) || !Number.isInteger(manifest.maxZoom)) fail('minZoom and maxZoom must be integers.');
  if (manifest.minZoom < 0 || manifest.maxZoom > 24 || manifest.minZoom > manifest.maxZoom) fail('Zoom range must be ordered and within 0..24.');

  if (!Array.isArray(manifest.bounds) || manifest.bounds.length !== 4 || manifest.bounds.some((value) => typeof value !== 'number' || !Number.isFinite(value))) {
    fail('bounds must contain four finite numbers.');
  }
  const [west, south, east, north] = manifest.bounds;
  if (west < -180 || west > 180 || east < -180 || east > 180 || south < -90 || south > 90 || north < -90 || north > 90) {
    fail('bounds coordinates are outside valid longitude/latitude ranges.');
  }
  if (west > east || south > north) fail('bounds must be ordered west/south/east/north.');

  if (!Array.isArray(manifest.attribution) || manifest.attribution.length === 0) fail('attribution must contain at least one entry.');
  for (const [index, entry] of manifest.attribution.entries()) {
    if (!entry || typeof entry !== 'object' || Array.isArray(entry)) fail(`attribution[${index}] must be an object.`);
    const keys = Object.keys(entry);
    if (keys.some((key) => key !== 'text' && key !== 'url')) fail(`attribution[${index}] contains an unsupported field.`);
    nonEmptyString(entry.text, `attribution[${index}].text`, 500);
    if (entry.url !== undefined) {
      nonEmptyString(entry.url, `attribution[${index}].url`, 1000);
      const parsed = new URL(entry.url);
      if (parsed.protocol !== 'https:') fail(`attribution[${index}].url must use HTTPS.`);
    }
  }

  if (!Array.isArray(manifest.sources) || manifest.sources.length === 0) fail('sources must contain at least one entry.');
  for (const [index, source] of manifest.sources.entries()) {
    if (!source || typeof source !== 'object' || Array.isArray(source)) fail(`sources[${index}] must be an object.`);
    const keys = Object.keys(source);
    if (keys.some((key) => !['name', 'datasetVersion', 'license', 'provenance'].includes(key))) fail(`sources[${index}] contains an unsupported field.`);
    nonEmptyString(source.name, `sources[${index}].name`, 160);
    nonEmptyString(source.datasetVersion, `sources[${index}].datasetVersion`, 160);
    nonEmptyString(source.license, `sources[${index}].license`, 160);
    if (source.provenance !== undefined) nonEmptyString(source.provenance, `sources[${index}].provenance`, 1000);
  }

  if (manifest.publicGeographicDataOnly !== true) fail('publicGeographicDataOnly must be true for this public map-data delivery contract.');
};

const paths = process.argv.slice(2);
if (paths.length === 0) fail('Provide at least one release-manifest path.');

for (const path of paths) {
  const raw = await readFile(path, 'utf8');
  const manifest = JSON.parse(raw);
  validateManifest(manifest);
  process.stdout.write(`validated ${path}\n`);
}
