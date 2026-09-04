import * as maplibregl from 'maplibre-gl';
import 'maplibre-gl/dist/maplibre-gl.css';
import './styles.css';
import './integration.css';
import { MapsAPI, MapsApiError, type Collection, type Route, type SavedPlace, type SearchResult } from './api';
import { IdentityClient } from './auth';
import { configuredMapDataManifestURL, resolveMapDataRelease } from './map-data';

const app = document.querySelector<HTMLDivElement>('#app');
if (!app) throw new Error('Maps application root was not found.');

app.innerHTML = `
  <main class="maps-shell" data-panel="explore">
    <div id="map" class="map-canvas" aria-label="Interactive map"></div>

    <header class="top-island glaze glaze--soft" aria-label="Map search and account controls">
      <button class="brand-button" type="button" aria-label="GoreeCloud Maps home">
        <span class="brand-mark" aria-hidden="true">M</span>
        <span class="brand-name">Maps</span>
      </button>

      <form class="search-control" role="search">
        <label class="sr-only" for="map-search">Search places and addresses</label>
        <span class="search-icon" aria-hidden="true">⌕</span>
        <input id="map-search" type="search" autocomplete="off" placeholder="Search places and addresses" />
        <button class="search-clear" type="button" aria-label="Clear search" hidden>×</button>
      </form>

      <button class="account-button" type="button" aria-label="Sign in with GoreeCloud Identity">
        <span data-account-label aria-hidden="true">ID</span>
      </button>
    </header>

    <aside class="desktop-panel surface" aria-label="Explore Maps">
      <div class="panel-heading">
        <p class="eyebrow">Explore</p>
        <h1>Find somewhere worth going.</h1>
        <p class="panel-copy">Search places, plan routes, and work with shared collections through GoreeCloud-controlled service boundaries.</p>
      </div>

      <div class="category-row" aria-label="Place categories">
        <button type="button" class="chip" data-category="Food">Food</button>
        <button type="button" class="chip" data-category="Coffee">Coffee</button>
        <button type="button" class="chip" data-category="Shopping">Shopping</button>
        <button type="button" class="chip" data-category="Outdoors">Outdoors</button>
      </div>

      <section class="status-card surface-raised" aria-labelledby="provider-title">
        <div class="status-dot" aria-hidden="true"></div>
        <div>
          <h2 id="provider-title" data-provider-title>Checking Maps capabilities…</h2>
          <p data-provider-copy>The renderer will remain on its privacy-safe local fallback until approved geographic services are configured.</p>
        </div>
      </section>

      <section class="quick-actions" aria-label="Map actions">
        <button class="action-tile" type="button" data-action="directions">
          <span class="action-symbol" aria-hidden="true">↗</span>
          <span><strong>Directions</strong><small>Plan a route</small></span>
        </button>
        <button class="action-tile" type="button" data-action="saved">
          <span class="action-symbol" aria-hidden="true">◇</span>
          <span><strong>Saved</strong><small>Your places</small></span>
        </button>
        <button class="action-tile" type="button" data-action="shared">
          <span class="action-symbol" aria-hidden="true">◎</span>
          <span><strong>Shared</strong><small>Collaborative maps</small></span>
        </button>
      </section>

      <div data-integration-host></div>
    </aside>

    <div class="map-controls glaze" aria-label="Map view controls">
      <button type="button" data-map-action="zoom-in" aria-label="Zoom in">+</button>
      <button type="button" data-map-action="zoom-out" aria-label="Zoom out">−</button>
      <span class="control-divider" aria-hidden="true"></span>
      <button type="button" data-map-action="reset-bearing" aria-label="Reset map orientation">↑</button>
    </div>

    <button class="locate-control glaze" type="button" aria-describedby="location-note">
      <span aria-hidden="true">⌾</span>
      <span class="sr-only">Show my location</span>
    </button>
    <p id="location-note" class="sr-only">Current-location display will use GoreeCloud Location when that integration is available.</p>

    <section class="mobile-sheet surface" aria-label="Explore nearby">
      <div class="sheet-handle" aria-hidden="true"></div>
      <div class="sheet-summary">
        <div>
          <p class="eyebrow">Explore nearby</p>
          <h2>Maps service controls</h2>
        </div>
        <button class="directions-button" type="button" data-action="directions">Directions</button>
      </div>
      <div class="mobile-integration-content" data-integration-host></div>
    </section>

    <nav class="navigation-capsule glaze" aria-label="Primary navigation">
      <button class="nav-item is-selected" type="button" data-panel="explore" aria-current="page">
        <span aria-hidden="true">⌖</span><span>Explore</span>
      </button>
      <button class="nav-item" type="button" data-panel="saved">
        <span aria-hidden="true">◇</span><span>Saved</span>
      </button>
      <button class="nav-item" type="button" data-panel="shared">
        <span aria-hidden="true">◎</span><span>Shared</span>
      </button>
    </nav>

    <div class="toast glaze glaze--deep" role="status" aria-live="polite" hidden></div>
  </main>
`;

const configuredManifest = configuredMapDataManifestURL();
const configuredStyle = import.meta.env.VITE_MAP_STYLE_URL?.trim() ?? '';
const initialStyle = configuredManifest ? '/map-style.json' : configuredStyle || '/map-style.json';
const map = new maplibregl.Map({
  container: 'map',
  style: initialStyle,
  center: [0, 20],
  zoom: 1.7,
  attributionControl: false,
  pitchWithRotate: true,
  dragRotate: true,
  touchPitch: true,
});
map.addControl(new maplibregl.AttributionControl({ compact: true }), 'bottom-right');

const identity = new IdentityClient();
const api = new MapsAPI(() => identity.getAccessToken());
let capabilities = { geocoding: false, routing: false };
let apiReachable = false;
let selectedMarker: maplibregl.Marker | null = null;
type MapDataState = {
  mode: 'local' | 'direct' | 'resolving' | 'release' | 'error';
  releaseId?: string;
};
let mapDataState: MapDataState = configuredManifest
  ? { mode: 'resolving' }
  : configuredStyle
    ? { mode: 'direct' }
    : { mode: 'local' };

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

const showToast = (message: string): void => {
  const toast = document.querySelector<HTMLDivElement>('.toast');
  if (!toast) return;
  toast.textContent = message;
  toast.hidden = false;
  window.setTimeout(() => {
    toast.hidden = true;
  }, 3600);
};

const setIntegrationContent = (content: string): void => {
  document.querySelectorAll<HTMLElement>('[data-integration-host]').forEach((host) => {
    host.innerHTML = `<section class="integration-panel surface-raised">${content}</section>`;
  });
};

const mapDataStatusLabel = (): string => {
  switch (mapDataState.mode) {
    case 'release':
      return `release ${mapDataState.releaseId ?? 'active'}`;
    case 'direct':
      return 'direct style';
    case 'resolving':
      return 'checking release';
    case 'error':
      return 'release unavailable';
    case 'local':
      return 'local fallback';
  }
};

const renderExplore = (): void => {
  setIntegrationContent(`
    <div class="integration-heading">
      <div>
        <h2>Connected services</h2>
        <p>Maps exposes provider availability without revealing provider origins.</p>
      </div>
    </div>
    <div class="result-meta">
      <span class="status-pill">Map data: ${escapeHTML(mapDataStatusLabel())}</span>
      <span class="status-pill">Geocoding: ${capabilities.geocoding ? 'configured' : 'unavailable'}</span>
      <span class="status-pill">Routing: ${capabilities.routing ? 'configured' : 'unavailable'}</span>
      <span class="status-pill">Identity: ${identity.configured ? (identity.authenticated ? 'signed in' : 'available') : 'not registered'}</span>
    </div>
  `);
};

const renderSignInRequired = (purpose: string): void => {
  setIntegrationContent(`
    <div class="integration-heading">
      <div><h2>Sign in required</h2><p>${escapeHTML(purpose)}</p></div>
    </div>
    <p class="integration-note">${
      identity.configured
        ? 'Maps uses GoreeCloud Identity through Authorization Code + PKCE. Access tokens are kept in memory for this source foundation.'
        : 'This build has no GoreeCloud Identity issuer/client registration configured. No fallback account system is used.'
    }</p>
    ${identity.configured ? '<div class="integration-actions"><button class="integration-button primary" type="button" data-sign-in>Sign in</button></div>' : ''}
  `);
  document.querySelectorAll<HTMLButtonElement>('[data-sign-in]').forEach((button) => {
    button.addEventListener('click', () => void startSignIn());
  });
};

const updateAccountButton = (): void => {
  const button = document.querySelector<HTMLButtonElement>('.account-button');
  const label = button?.querySelector<HTMLElement>('[data-account-label]');
  if (!button || !label) return;
  if (identity.authenticated) {
    label.textContent = '✓';
    button.dataset.state = 'signed-in';
    button.setAttribute('aria-label', 'Disconnect Maps from the current GoreeCloud Identity session');
  } else if (identity.configured) {
    label.textContent = 'ID';
    delete button.dataset.state;
    button.setAttribute('aria-label', 'Sign in with GoreeCloud Identity');
  } else {
    label.textContent = 'ID';
    delete button.dataset.state;
    button.setAttribute('aria-label', 'GoreeCloud Identity is not configured for this Maps build');
  }
};

const updateProviderStatus = (): void => {
  const title = document.querySelector<HTMLElement>('[data-provider-title]');
  const copy = document.querySelector<HTMLElement>('[data-provider-copy]');
  if (!title || !copy) return;
  if (mapDataState.mode === 'error') {
    title.textContent = 'Map data release unavailable';
    copy.textContent = 'The configured release manifest was rejected or unavailable, so Maps kept the privacy-safe local empty style.';
    return;
  }
  if (mapDataState.mode === 'resolving') {
    title.textContent = 'Checking map data release…';
    copy.textContent = 'Maps is validating the configured public geographic-data manifest before loading its immutable style release.';
    return;
  }
  if (!apiReachable) {
    title.textContent = 'Maps API unavailable';
    copy.textContent = 'The renderer remains available, but search, routing, and collaboration require the same-origin Maps API.';
    return;
  }
  if (!capabilities.geocoding && !capabilities.routing) {
    title.textContent = 'Geographic providers not configured';
    copy.textContent = 'The Maps API is reachable, but live geocoding and routing remain disabled. No public third-party provider is used automatically.';
    return;
  }
  title.textContent = 'Maps service capabilities available';
  copy.textContent = `Geocoding ${capabilities.geocoding ? 'is configured' : 'is unavailable'}; routing ${capabilities.routing ? 'is configured' : 'is unavailable'}. Provider origins remain server-side.`;
};

const applyConfiguredMapData = async (): Promise<void> => {
  if (!configuredManifest) return;
  try {
    const release = await resolveMapDataRelease();
    if (!release) {
      mapDataState = { mode: 'local' };
      return;
    }
    mapDataState = { mode: 'release', releaseId: release.manifest.releaseId };
    const apply = (): void => {
      map.setStyle(release.style);
      const [west, south, east, north] = release.manifest.bounds;
      const globalCoverage = west <= -179.9 && east >= 179.9 && south <= -84 && north >= 84;
      if (!globalCoverage) {
        map.fitBounds(
          [
            [west, south],
            [east, north],
          ],
          { padding: 48, maxZoom: Math.min(release.manifest.maxZoom, 12), duration: 0 },
        );
      }
    };
    if (map.isStyleLoaded()) apply();
    else map.once('load', apply);
  } catch (error) {
    mapDataState = { mode: 'error' };
    showToast(error instanceof Error ? error.message : 'The configured map-data release could not be loaded.');
  }
};

const startSignIn = async (): Promise<void> => {
  try {
    await identity.signIn();
  } catch (error) {
    showToast(error instanceof Error ? error.message : 'GoreeCloud Identity sign-in could not start.');
  }
};

const requireAuthentication = (purpose: string): boolean => {
  if (identity.authenticated) return true;
  renderSignInRequired(purpose);
  return false;
};

const handleApiError = (error: unknown): void => {
  if (error instanceof MapsApiError) {
    if (error.status === 401) identity.disconnect();
    updateAccountButton();
    showToast(error.message);
    return;
  }
  showToast(error instanceof Error ? error.message : 'Maps could not complete that request.');
};

const renderSearchResults = (query: string, results: SearchResult[]): void => {
  if (results.length === 0) {
    setIntegrationContent(`<div class="integration-heading"><div><h2>No results</h2><p>No places matched “${escapeHTML(query)}”.</p></div></div>`);
    return;
  }
  setIntegrationContent(`
    <div class="integration-heading"><div><h2>Search results</h2><p>${results.length} result${results.length === 1 ? '' : 's'} for “${escapeHTML(query)}”.</p></div></div>
    <div class="integration-list">
      ${results
        .map(
          (result, index) => `
            <button type="button" class="result-card" data-result-index="${index}">
              <strong>${escapeHTML(result.name || result.label)}</strong>
              <p>${escapeHTML(result.label)}</p>
              <div class="result-meta">${result.category ? `<span class="status-pill">${escapeHTML(result.category)}</span>` : ''}${result.type ? `<span class="status-pill">${escapeHTML(result.type)}</span>` : ''}</div>
            </button>`,
        )
        .join('')}
    </div>
  `);
  document.querySelectorAll<HTMLButtonElement>('[data-result-index]').forEach((button) => {
    button.addEventListener('click', () => {
      const index = Number(button.dataset.resultIndex);
      const result = results[index];
      if (!result) return;
      map.flyTo({ center: [result.longitude, result.latitude], zoom: Math.max(map.getZoom(), 14) });
      selectedMarker?.remove();
      selectedMarker = new maplibregl.Marker().setLngLat([result.longitude, result.latitude]).addTo(map);
    });
  });
};

const performSearch = async (query: string): Promise<void> => {
  if (!capabilities.geocoding) {
    showToast(apiReachable ? 'Place search is not configured.' : 'The Maps API is unavailable.');
    return;
  }
  if (!requireAuthentication('Sign in to search live place data through the configured Maps geocoder.')) return;
  setIntegrationContent(`<div class="integration-heading"><div><h2>Searching…</h2><p>Looking for “${escapeHTML(query)}”.</p></div></div>`);
  try {
    renderSearchResults(query, await api.search(query));
  } catch (error) {
    handleApiError(error);
  }
};

const decodeValhallaShape = (encoded: string): [number, number][] => {
  const coordinates: [number, number][] = [];
  let index = 0;
  let latitude = 0;
  let longitude = 0;
  const decodeValue = (): number => {
    let result = 0;
    let shift = 0;
    while (index < encoded.length) {
      const value = encoded.charCodeAt(index++) - 63;
      result |= (value & 0x1f) << shift;
      shift += 5;
      if (value < 0x20) return result & 1 ? ~(result >> 1) : result >> 1;
      if (shift > 30) break;
    }
    throw new Error('Route geometry is malformed.');
  };
  while (index < encoded.length) {
    latitude += decodeValue();
    longitude += decodeValue();
    coordinates.push([longitude / 1e6, latitude / 1e6]);
  }
  return coordinates;
};

const drawRoute = (route: Route): void => {
  const coordinates = route.legs.flatMap((leg) => (leg.shape ? decodeValhallaShape(leg.shape) : []));
  if (coordinates.length < 2) return;
  const feature = {
    type: 'Feature' as const,
    properties: {},
    geometry: { type: 'LineString' as const, coordinates },
  };
  const source = map.getSource('active-route');
  if (source) {
    (source as maplibregl.GeoJSONSource).setData(feature);
  } else {
    map.addSource('active-route', { type: 'geojson', data: feature });
    map.addLayer({
      id: 'active-route-line',
      type: 'line',
      source: 'active-route',
      layout: { 'line-cap': 'round', 'line-join': 'round' },
      paint: { 'line-color': '#2768d7', 'line-width': 6, 'line-opacity': 0.9 },
    });
  }
  const bounds = new maplibregl.LngLatBounds();
  coordinates.forEach((coordinate) => bounds.extend(coordinate));
  map.fitBounds(bounds, { padding: 72, maxZoom: 15, duration: 500 });
};

const renderRoute = (origin: SearchResult, destination: SearchResult, route: Route): void => {
  const minutes = Math.max(1, Math.round(route.durationSeconds / 60));
  const kilometers = route.distanceMeters / 1000;
  const maneuvers = route.legs.flatMap((leg) => leg.maneuvers).slice(0, 10);
  setIntegrationContent(`
    <div class="integration-heading"><div><h2>${escapeHTML(origin.name)} → ${escapeHTML(destination.name)}</h2><p>${escapeHTML(route.mode)} route</p></div><button class="integration-button" type="button" data-directions-again>Change</button></div>
    <div class="route-summary"><strong>${minutes} min · ${kilometers.toFixed(kilometers < 10 ? 1 : 0)} km</strong><p>Provider-derived estimate; live navigation and traffic acceptance are not implemented.</p></div>
    ${maneuvers.length ? `<ol class="maneuver-list">${maneuvers.map((maneuver) => `<li><span>${escapeHTML(maneuver.instruction)}</span><small>${Math.round(maneuver.distanceMeters)} m</small></li>`).join('')}</ol>` : '<p class="integration-note">No maneuver list was returned for this route.</p>'}
  `);
  document.querySelectorAll<HTMLButtonElement>('[data-directions-again]').forEach((button) => button.addEventListener('click', renderDirectionsForm));
  try {
    drawRoute(route);
  } catch {
    showToast('The route was calculated, but its geometry could not be rendered.');
  }
};

const calculateDirections = async (form: HTMLFormElement): Promise<void> => {
  const data = new FormData(form);
  const originQuery = String(data.get('origin') ?? '').trim();
  const destinationQuery = String(data.get('destination') ?? '').trim();
  const mode = String(data.get('mode') ?? 'drive');
  if (!originQuery || !destinationQuery) {
    showToast('Enter both an origin and destination.');
    return;
  }
  setIntegrationContent('<div class="integration-heading"><div><h2>Planning route…</h2><p>Resolving the requested places and route.</p></div></div>');
  try {
    const [originResults, destinationResults] = await Promise.all([api.search(originQuery, 1), api.search(destinationQuery, 1)]);
    const origin = originResults[0];
    const destination = destinationResults[0];
    if (!origin || !destination) {
      showToast('Maps could not resolve both route endpoints.');
      renderDirectionsForm();
      return;
    }
    const route = await api.route(mode, [
      { latitude: origin.latitude, longitude: origin.longitude },
      { latitude: destination.latitude, longitude: destination.longitude },
    ]);
    renderRoute(origin, destination, route);
  } catch (error) {
    handleApiError(error);
    renderDirectionsForm();
  }
};

function renderDirectionsForm(): void {
  if (!capabilities.geocoding || !capabilities.routing) {
    setIntegrationContent(`<div class="integration-heading"><div><h2>Directions unavailable</h2><p>${!apiReachable ? 'The Maps API is unavailable.' : 'Directions require both an approved geocoder and routing provider configuration.'}</p></div></div>`);
    return;
  }
  if (!requireAuthentication('Sign in to resolve route endpoints and request a route.')) return;
  setIntegrationContent(`
    <div class="integration-heading"><div><h2>Directions</h2><p>Resolve two places and calculate a provider-backed route.</p></div></div>
    <form class="integration-form" data-direction-form>
      <label>From<input name="origin" autocomplete="off" maxlength="240" required placeholder="Origin"></label>
      <label>To<input name="destination" autocomplete="off" maxlength="240" required placeholder="Destination"></label>
      <label>Travel mode<select name="mode"><option value="drive">Drive</option><option value="walk">Walk</option><option value="bicycle">Bicycle</option><option value="transit">Transit / multimodal</option></select></label>
      <div class="integration-actions"><button class="integration-button primary" type="submit">Get route</button></div>
    </form>
  `);
  document.querySelectorAll<HTMLFormElement>('[data-direction-form]').forEach((form) => {
    form.addEventListener('submit', (event) => {
      event.preventDefault();
      void calculateDirections(form);
    });
  });
}

const renderCollectionItems = async (collection: Collection): Promise<void> => {
  try {
    const items = await api.listCollectionItems(collection.id);
    setIntegrationContent(`
      <div class="integration-heading"><div><h2>${escapeHTML(collection.name)}</h2><p>${escapeHTML(collection.role)} · ${items.length} item${items.length === 1 ? '' : 's'}</p></div><button class="integration-button" type="button" data-back-collections>Back</button></div>
      ${collection.description ? `<p class="integration-note">${escapeHTML(collection.description)}</p>` : ''}
      <div class="integration-list">${items.length ? items.map((item) => `<button type="button" class="collection-item-card" data-item-lat="${item.latitude}" data-item-lng="${item.longitude}"><strong>${escapeHTML(item.name)}</strong><p>${escapeHTML(item.address || item.note || 'Saved map item')}</p></button>`).join('') : '<p class="integration-note">This collection has no items yet.</p>'}</div>
      <p class="integration-note">Member invitation discovery is intentionally not exposed here until GoreeCloud Identity defines an approved recipient-directory contract.</p>
    `);
    document.querySelectorAll<HTMLButtonElement>('[data-back-collections]').forEach((button) => button.addEventListener('click', () => void renderCollections()));
    document.querySelectorAll<HTMLButtonElement>('[data-item-lat]').forEach((button) => {
      button.addEventListener('click', () => {
        const latitude = Number(button.dataset.itemLat);
        const longitude = Number(button.dataset.itemLng);
        if (!Number.isFinite(latitude) || !Number.isFinite(longitude)) return;
        map.flyTo({ center: [longitude, latitude], zoom: Math.max(map.getZoom(), 14) });
      });
    });
  } catch (error) {
    handleApiError(error);
  }
};

const renderCollections = async (): Promise<void> => {
  if (!requireAuthentication('Sign in to load and create Maps collections.')) return;
  setIntegrationContent('<div class="integration-heading"><div><h2>Shared collections</h2><p>Loading your authorized collections…</p></div></div>');
  try {
    const collections = await api.listCollections();
    setIntegrationContent(`
      <div class="integration-heading"><div><h2>Shared collections</h2><p>${collections.length} visible collection${collections.length === 1 ? '' : 's'}.</p></div></div>
      <form class="integration-form" data-new-collection-form>
        <label>New collection<input name="name" maxlength="160" required placeholder="Collection name"></label>
        <label>Description<textarea name="description" maxlength="1000" placeholder="Optional description"></textarea></label>
        <div class="integration-actions"><button class="integration-button primary" type="submit">Create collection</button></div>
      </form>
      <div class="integration-list" style="margin-top: .8rem">${collections.map((collection, index) => `<button type="button" class="collection-card" data-collection-index="${index}"><strong>${escapeHTML(collection.name)}</strong><p>${escapeHTML(collection.description || 'No description')}</p><div class="collection-meta"><span class="status-pill">${escapeHTML(collection.role)}</span><span class="status-pill">${escapeHTML(collection.sharingMode)}</span></div></button>`).join('')}</div>
      <p class="integration-note">Human-friendly invitations remain blocked on an approved GoreeCloud Identity recipient-discovery contract; Maps does not substitute an administrative identity directory.</p>
    `);
    document.querySelectorAll<HTMLFormElement>('[data-new-collection-form]').forEach((form) => {
      form.addEventListener('submit', (event) => {
        event.preventDefault();
        const data = new FormData(form);
        const name = String(data.get('name') ?? '').trim();
        const description = String(data.get('description') ?? '').trim();
        if (!name) return;
        void api
          .createCollection(name, description)
          .then(() => renderCollections())
          .catch(handleApiError);
      });
    });
    document.querySelectorAll<HTMLButtonElement>('[data-collection-index]').forEach((button) => {
      button.addEventListener('click', () => {
        const collection = collections[Number(button.dataset.collectionIndex)];
        if (collection) void renderCollectionItems(collection);
      });
    });
  } catch (error) {
    handleApiError(error);
  }
};

const validCoordinatePair = (latitude: number, longitude: number): boolean =>
  Number.isFinite(latitude) && Number.isFinite(longitude) && latitude >= -90 && latitude <= 90 && longitude >= -180 && longitude <= 180;

const focusSavedPlace = (place: SavedPlace): void => {
  map.flyTo({ center: [place.longitude, place.latitude], zoom: Math.max(map.getZoom(), 14) });
  selectedMarker?.remove();
  selectedMarker = new maplibregl.Marker().setLngLat([place.longitude, place.latitude]).addTo(map);
};

const renderSavedPlaceDeleteConfirmation = (place: SavedPlace): void => {
  setIntegrationContent(`
    <div class="integration-heading">
      <div><h2>Delete ${escapeHTML(place.name)}?</h2><p>This removes the saved place from your private Maps data.</p></div>
    </div>
    <p class="integration-note">The delete is revision-checked. If this place changed in another session, Maps will reject the stale delete instead of silently discarding newer data.</p>
    <div class="integration-actions">
      <button class="integration-button" type="button" data-cancel-delete-saved>Cancel</button>
      <button class="integration-button primary" type="button" data-confirm-delete-saved>Delete place</button>
    </div>
  `);
  document.querySelectorAll<HTMLButtonElement>('[data-cancel-delete-saved]').forEach((button) => {
    button.addEventListener('click', () => renderSavedPlaceEditor(place));
  });
  document.querySelectorAll<HTMLButtonElement>('[data-confirm-delete-saved]').forEach((button) => {
    button.addEventListener('click', () => {
      button.disabled = true;
      void api
        .deleteSavedPlace(place.id, place.revision)
        .then(() => {
          showToast('Saved place deleted.');
          return renderSavedPlaces();
        })
        .catch((error) => {
          handleApiError(error);
          void renderSavedPlaces();
        });
    });
  });
};

const renderSavedPlaceEditor = (place: SavedPlace): void => {
  setIntegrationContent(`
    <div class="integration-heading">
      <div><h2>Edit saved place</h2><p>Revision ${place.revision} · private to your Maps account.</p></div>
      <button class="integration-button" type="button" data-back-saved>Back</button>
    </div>
    <form class="integration-form" data-edit-saved-form>
      <label>Name<input name="name" maxlength="240" required value="${escapeHTML(place.name)}"></label>
      <label>Address<input name="address" maxlength="1000" value="${escapeHTML(place.address ?? '')}" placeholder="Optional address"></label>
      <label>Latitude<input name="latitude" type="number" min="-90" max="90" step="any" required value="${place.latitude}"></label>
      <label>Longitude<input name="longitude" type="number" min="-180" max="180" step="any" required value="${place.longitude}"></label>
      <label>Note<textarea name="note" maxlength="4000" placeholder="Optional private note">${escapeHTML(place.note ?? '')}</textarea></label>
      <div class="integration-actions">
        <button class="integration-button primary" type="submit">Save changes</button>
        <button class="integration-button" type="button" data-show-saved-on-map>Show on map</button>
        <button class="integration-button" type="button" data-delete-saved>Delete</button>
      </div>
    </form>
  `);
  document.querySelectorAll<HTMLButtonElement>('[data-back-saved]').forEach((button) => button.addEventListener('click', () => void renderSavedPlaces()));
  document.querySelectorAll<HTMLButtonElement>('[data-show-saved-on-map]').forEach((button) => button.addEventListener('click', () => focusSavedPlace(place)));
  document.querySelectorAll<HTMLButtonElement>('[data-delete-saved]').forEach((button) => button.addEventListener('click', () => renderSavedPlaceDeleteConfirmation(place)));
  document.querySelectorAll<HTMLFormElement>('[data-edit-saved-form]').forEach((form) => {
    form.addEventListener('submit', (event) => {
      event.preventDefault();
      const data = new FormData(form);
      const name = String(data.get('name') ?? '').trim();
      const address = String(data.get('address') ?? '').trim();
      const note = String(data.get('note') ?? '').trim();
      const latitude = Number(data.get('latitude'));
      const longitude = Number(data.get('longitude'));
      if (!name || !validCoordinatePair(latitude, longitude)) {
        showToast('Enter a name and valid latitude/longitude coordinates.');
        return;
      }
      const submit = form.querySelector<HTMLButtonElement>('button[type="submit"]');
      if (submit) submit.disabled = true;
      void api
        .updateSavedPlace(place.id, {
          name,
          address,
          latitude,
          longitude,
          note,
          expectedRevision: place.revision,
        })
        .then((updated) => {
          focusSavedPlace(updated);
          showToast('Saved place updated.');
          return renderSavedPlaces();
        })
        .catch((error) => {
          handleApiError(error);
          void renderSavedPlaces();
        });
    });
  });
};

const renderSavedPlaces = async (): Promise<void> => {
  if (!apiReachable) {
    setIntegrationContent('<div class="integration-heading"><div><h2>Saved places unavailable</h2><p>The same-origin Maps API is not reachable.</p></div></div>');
    return;
  }
  if (!requireAuthentication('Sign in to manage your private saved places.')) return;
  setIntegrationContent('<div class="integration-heading"><div><h2>Saved places</h2><p>Loading your private saved places…</p></div></div>');
  try {
    const places = await api.listSavedPlaces();
    const center = map.getCenter();
    setIntegrationContent(`
      <div class="integration-heading"><div><h2>Saved places</h2><p>${places.length} private place${places.length === 1 ? '' : 's'}.</p></div></div>
      <form class="integration-form" data-new-saved-form>
        <label>Name<input name="name" maxlength="240" required placeholder="Place name"></label>
        <label>Address<input name="address" maxlength="1000" placeholder="Optional address"></label>
        <label>Latitude<input name="latitude" type="number" min="-90" max="90" step="any" required value="${center.lat.toFixed(6)}"></label>
        <label>Longitude<input name="longitude" type="number" min="-180" max="180" step="any" required value="${center.lng.toFixed(6)}"></label>
        <label>Note<textarea name="note" maxlength="4000" placeholder="Optional private note"></textarea></label>
        <p class="integration-note">Coordinates default to the current map center. This workflow does not request device location or read GoreeCloud Location history.</p>
        <div class="integration-actions"><button class="integration-button primary" type="submit">Save place</button></div>
      </form>
      <div class="integration-list" style="margin-top: .8rem">
        ${
          places.length
            ? places
                .map(
                  (place, index) => `<button type="button" class="collection-item-card" data-saved-index="${index}">
                    <strong>${escapeHTML(place.name)}</strong>
                    <p>${escapeHTML(place.address || `${place.latitude.toFixed(5)}, ${place.longitude.toFixed(5)}`)}</p>
                    <div class="result-meta"><span class="status-pill">private</span><span class="status-pill">rev ${place.revision}</span></div>
                  </button>`,
                )
                .join('')
            : '<p class="integration-note">No saved places yet. Add one from the current map center.</p>'
        }
      </div>
    `);
    document.querySelectorAll<HTMLFormElement>('[data-new-saved-form]').forEach((form) => {
      form.addEventListener('submit', (event) => {
        event.preventDefault();
        const data = new FormData(form);
        const name = String(data.get('name') ?? '').trim();
        const address = String(data.get('address') ?? '').trim();
        const note = String(data.get('note') ?? '').trim();
        const latitude = Number(data.get('latitude'));
        const longitude = Number(data.get('longitude'));
        if (!name || !validCoordinatePair(latitude, longitude)) {
          showToast('Enter a name and valid latitude/longitude coordinates.');
          return;
        }
        const submit = form.querySelector<HTMLButtonElement>('button[type="submit"]');
        if (submit) submit.disabled = true;
        void api
          .createSavedPlace({ provider: 'manual', name, address, latitude, longitude, note })
          .then((created) => {
            focusSavedPlace(created);
            showToast('Place saved.');
            return renderSavedPlaces();
          })
          .catch((error) => {
            handleApiError(error);
            if (submit) submit.disabled = false;
          });
      });
    });
    document.querySelectorAll<HTMLButtonElement>('[data-saved-index]').forEach((button) => {
      button.addEventListener('click', () => {
        const place = places[Number(button.dataset.savedIndex)];
        if (!place) return;
        focusSavedPlace(place);
        renderSavedPlaceEditor(place);
      });
    });
  } catch (error) {
    handleApiError(error);
  }
};

map.on('error', () => {
  console.warn('Map renderer reported an error.');
  showToast('The current map style could not be fully loaded.');
});

document.querySelectorAll<HTMLButtonElement>('[data-map-action]').forEach((button) => {
  button.addEventListener('click', () => {
    switch (button.dataset.mapAction) {
      case 'zoom-in':
        map.zoomIn();
        break;
      case 'zoom-out':
        map.zoomOut();
        break;
      case 'reset-bearing':
        map.easeTo({ bearing: 0, pitch: 0, duration: 280 });
        break;
    }
  });
});

const search = document.querySelector<HTMLInputElement>('#map-search');
const clearSearch = document.querySelector<HTMLButtonElement>('.search-clear');
search?.addEventListener('input', () => {
  if (clearSearch) clearSearch.hidden = search.value.length === 0;
});
search?.form?.addEventListener('submit', (event) => {
  event.preventDefault();
  const query = search.value.trim();
  if (query) void performSearch(query);
});
clearSearch?.addEventListener('click', () => {
  if (!search) return;
  search.value = '';
  clearSearch.hidden = true;
  search.focus();
  renderExplore();
});

document.querySelectorAll<HTMLButtonElement>('[data-category]').forEach((button) => {
  button.addEventListener('click', () => {
    if (!search) return;
    search.value = button.dataset.category ?? '';
    clearSearch?.removeAttribute('hidden');
    if (search.value) void performSearch(search.value);
  });
});

document.querySelector<HTMLButtonElement>('.account-button')?.addEventListener('click', () => {
  if (identity.authenticated) {
    identity.disconnect();
    updateAccountButton();
    renderExplore();
    showToast('Maps disconnected its in-memory access token. The GoreeCloud Identity SSO session may remain active.');
    return;
  }
  if (!identity.configured) {
    showToast('GoreeCloud Identity is not registered for this Maps build.');
    return;
  }
  void startSignIn();
});

document.querySelector<HTMLButtonElement>('.locate-control')?.addEventListener('click', () => {
  showToast('Current location remains delegated to GoreeCloud Location and is not implemented in Maps yet.');
});

document.querySelectorAll<HTMLButtonElement>('[data-action]').forEach((button) => {
  button.addEventListener('click', () => {
    if (button.dataset.action === 'directions') renderDirectionsForm();
    else if (button.dataset.action === 'shared') void renderCollections();
    else void renderSavedPlaces();
  });
});

document.querySelectorAll<HTMLButtonElement>('.nav-item').forEach((item) => {
  item.addEventListener('click', () => {
    document.querySelectorAll<HTMLButtonElement>('.nav-item').forEach((candidate) => {
      const selected = candidate === item;
      candidate.classList.toggle('is-selected', selected);
      if (selected) candidate.setAttribute('aria-current', 'page');
      else candidate.removeAttribute('aria-current');
    });
    if (item.dataset.panel === 'explore') renderExplore();
    else if (item.dataset.panel === 'shared') void renderCollections();
    else void renderSavedPlaces();
  });
});

const reduceMotion = window.matchMedia('(prefers-reduced-motion: reduce)');
const applyMotionPreference = (): void => {
  document.documentElement.dataset.reducedMotion = String(reduceMotion.matches);
};
applyMotionPreference();
reduceMotion.addEventListener('change', applyMotionPreference);

const bootstrap = async (): Promise<void> => {
  try {
    if ((await identity.initialize()) === 'authenticated') showToast('Signed in to Maps with GoreeCloud Identity.');
  } catch (error) {
    showToast(error instanceof Error ? error.message : 'GoreeCloud Identity sign-in failed.');
  }
  updateAccountButton();
  try {
    capabilities = await api.capabilities();
    apiReachable = true;
  } catch {
    apiReachable = false;
  }
  await applyConfiguredMapData();
  updateProviderStatus();
  renderExplore();
};

void bootstrap();