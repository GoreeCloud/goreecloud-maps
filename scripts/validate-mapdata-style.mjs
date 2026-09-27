import { readFile } from 'node:fs/promises';
import process from 'node:process';

const fail = (message) => {
  throw new Error(message);
};

const record = (value) => value !== null && typeof value === 'object' && !Array.isArray(value);
const spriteNamePattern = /^[A-Za-z0-9][A-Za-z0-9._@-]{0,119}$/;
const spriteIDPattern = /^[A-Za-z0-9][A-Za-z0-9._@:-]{0,119}$/;

const validateRelativeResource = (value, pattern, label) => {
  if (typeof value !== 'string' || !pattern.test(value)) fail(`${label} must use the release-local resource path.`);
};

const validateStyle = (manifest, style) => {
  if (!record(manifest) || typeof manifest.releaseId !== 'string') fail('Manifest is invalid.');
  if (!record(style) || style.version !== 8 || !record(style.sources) || !Array.isArray(style.layers)) {
    fail('Style must be a MapLibre v8 style with sources and layers.');
  }
  if ('imports' in style) fail('Style imports are not allowed in map-data schema v1.');
  if ('font-faces' in style) fail('External font-face resources are not allowed in map-data schema v1.');

  const sourceIDs = new Set(Object.keys(style.sources));
  if (sourceIDs.size === 0) fail('Style must contain at least one vector source.');
  for (const [sourceID, source] of Object.entries(style.sources)) {
    if (!record(source) || source.type !== 'vector') fail(`Source ${sourceID} must be vector in schema v1.`);
    if ('url' in source) fail(`Source ${sourceID} must not use TileJSON indirection.`);
    if (!Array.isArray(source.tiles) || source.tiles.length !== 1 || source.tiles[0] !== 'tiles/{z}/{x}/{y}.pbf') {
      fail(`Source ${sourceID} must use tiles/{z}/{x}/{y}.pbf.`);
    }
    if (source.scheme !== undefined && source.scheme !== 'xyz') fail(`Source ${sourceID} must use xyz tile addressing.`);
    if (source.encoding !== undefined && source.encoding !== 'mvt') fail(`Source ${sourceID} must use MVT encoding in schema v1.`);
  }

  if (style.glyphs !== undefined) {
    validateRelativeResource(style.glyphs, /^glyphs\/\{fontstack\}\/\{range\}\.pbf$/, 'glyphs');
  }
  if (style.sprite !== undefined) {
    if (typeof style.sprite === 'string') {
      const match = style.sprite.match(/^sprites\/([A-Za-z0-9][A-Za-z0-9._@-]{0,119})$/);
      if (!match || !match[1] || !spriteNamePattern.test(match[1])) fail('sprite must use a release-local sprite name.');
    } else if (Array.isArray(style.sprite)) {
      for (const [index, entry] of style.sprite.entries()) {
        if (!record(entry) || typeof entry.id !== 'string' || !spriteIDPattern.test(entry.id) || typeof entry.url !== 'string') {
          fail(`sprite[${index}] is invalid.`);
        }
        const match = entry.url.match(/^sprites\/([A-Za-z0-9][A-Za-z0-9._@-]{0,119})$/);
        if (!match || !match[1] || !spriteNamePattern.test(match[1])) fail(`sprite[${index}].url must remain release-local.`);
      }
    } else {
      fail('sprite must be a string or array.');
    }
  }

  const layerIDs = new Set();
  for (const [index, layer] of style.layers.entries()) {
    if (!record(layer) || typeof layer.id !== 'string' || layer.id.trim() === '') fail(`layers[${index}] requires an id.`);
    if (layerIDs.has(layer.id)) fail(`Duplicate layer id: ${layer.id}`);
    layerIDs.add(layer.id);
    if (layer.type === 'background') {
      if (layer.source !== undefined) fail(`Background layer ${layer.id} must not declare a source.`);
      continue;
    }
    if (typeof layer.source !== 'string' || !sourceIDs.has(layer.source)) fail(`Layer ${layer.id} references an unknown source.`);
    if (typeof layer['source-layer'] !== 'string' || layer['source-layer'].trim() === '') {
      fail(`Vector layer ${layer.id} requires source-layer.`);
    }
  }

  const expectedStyle = `releases/${manifest.releaseId}/style.json`;
  const expectedTiles = `releases/${manifest.releaseId}/tiles/{z}/{x}/{y}.pbf`;
  if (manifest.stylePath !== expectedStyle || manifest.tileTemplate !== expectedTiles) {
    fail('Manifest/style paths do not match the release identifier.');
  }
};

const [manifestPath, stylePath, ...extra] = process.argv.slice(2);
if (!manifestPath || !stylePath || extra.length > 0) {
  fail('Usage: node scripts/validate-mapdata-style.mjs <release-manifest.json> <style.json>');
}

const manifest = JSON.parse(await readFile(manifestPath, 'utf8'));
const style = JSON.parse(await readFile(stylePath, 'utf8'));
validateStyle(manifest, style);
process.stdout.write(`validated ${stylePath} against ${manifestPath}\n`);
