import maplibregl from 'maplibre-gl';
import 'maplibre-gl/dist/maplibre-gl.css';
import './styles.css';

const app = document.querySelector<HTMLDivElement>('#app');

if (!app) {
  throw new Error('Maps application root was not found.');
}

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

      <button class="account-button" type="button" aria-label="Account and profile">
        <span aria-hidden="true">GC</span>
      </button>
    </header>

    <aside class="desktop-panel surface" aria-label="Explore Maps">
      <div class="panel-heading">
        <p class="eyebrow">Explore</p>
        <h1>Find somewhere worth going.</h1>
        <p class="panel-copy">Search, save places, build routes, and share collections. Live place data will appear when an approved provider is configured.</p>
      </div>

      <div class="category-row" aria-label="Place categories">
        <button type="button" class="chip">Food</button>
        <button type="button" class="chip">Coffee</button>
        <button type="button" class="chip">Shopping</button>
        <button type="button" class="chip">Outdoors</button>
      </div>

      <section class="status-card surface-raised" aria-labelledby="provider-title">
        <div class="status-dot" aria-hidden="true"></div>
        <div>
          <h2 id="provider-title">Map data provider not configured</h2>
          <p>The renderer is running with a local empty style. Configure <code>VITE_MAP_STYLE_URL</code> with an approved GoreeCloud map style endpoint to load geographic data.</p>
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
          <h2>Maps is ready for a provider.</h2>
        </div>
        <button class="directions-button" type="button" data-action="directions">Directions</button>
      </div>
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

const configuredStyle = import.meta.env.VITE_MAP_STYLE_URL?.trim();
const style = configuredStyle || '/map-style.json';

const map = new maplibregl.Map({
  container: 'map',
  style,
  center: [0, 20],
  zoom: 1.7,
  attributionControl: false,
  pitchWithRotate: true,
  dragRotate: true,
  touchPitch: true,
});

map.addControl(new maplibregl.AttributionControl({ compact: true }), 'bottom-right');

const showToast = (message: string): void => {
  const toast = document.querySelector<HTMLDivElement>('.toast');
  if (!toast) return;

  toast.textContent = message;
  toast.hidden = false;
  window.setTimeout(() => {
    toast.hidden = true;
  }, 3200);
};

map.on('error', (event) => {
  console.error('Map renderer error', event.error);
  showToast('The current map provider could not be loaded.');
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
  if (!query) return;
  showToast('Place search will activate with the GoreeCloud geocoding provider.');
});

clearSearch?.addEventListener('click', () => {
  if (!search) return;
  search.value = '';
  clearSearch.hidden = true;
  search.focus();
});

document.querySelector<HTMLButtonElement>('.locate-control')?.addEventListener('click', () => {
  showToast('Current location will be supplied by GoreeCloud Location after permission integration.');
});

document.querySelectorAll<HTMLButtonElement>('[data-action]').forEach((button) => {
  button.addEventListener('click', () => {
    const action = button.dataset.action;
    if (action === 'directions') {
      showToast('Route planning is the next provider-backed milestone.');
      return;
    }
    showToast(`${action === 'saved' ? 'Saved places' : 'Shared maps'} will become available with the multi-user API.`);
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

    if (item.dataset.panel !== 'explore') {
      showToast(`${item.dataset.panel === 'saved' ? 'Saved places' : 'Shared maps'} are not connected yet.`);
    }
  });
});

const reduceMotion = window.matchMedia('(prefers-reduced-motion: reduce)');
const applyMotionPreference = (): void => {
  document.documentElement.dataset.reducedMotion = String(reduceMotion.matches);
};
applyMotionPreference();
reduceMotion.addEventListener('change', applyMotionPreference);
