export type ProviderCapabilities = {
  geocoding: boolean;
  routing: boolean;
};

export type SearchResult = {
  id: string;
  name: string;
  label: string;
  latitude: number;
  longitude: number;
  category?: string;
  type?: string;
};

export type RoutePoint = {
  latitude: number;
  longitude: number;
};

export type RouteManeuver = {
  type: number;
  instruction: string;
  verbalInstruction?: string;
  distanceMeters: number;
  durationSeconds: number;
  beginShapeIndex: number;
  endShapeIndex: number;
  streetNames?: string[];
};

export type RouteLeg = {
  distanceMeters: number;
  durationSeconds: number;
  shape: string;
  maneuvers: RouteManeuver[];
};

export type Route = {
  mode: string;
  distanceMeters: number;
  durationSeconds: number;
  legs: RouteLeg[];
};

export type Collection = {
  id: string;
  name: string;
  description?: string;
  sharingMode: string;
  role: 'owner' | 'editor' | 'viewer';
  revision: number;
  updatedAt: string;
};

export type CollectionItem = {
  id: string;
  collectionId: string;
  provider: string;
  providerPlaceId?: string;
  name: string;
  address?: string;
  latitude: number;
  longitude: number;
  note?: string;
  sortKey: number;
  revision: number;
  createdAt: string;
  updatedAt: string;
};

export class MapsApiError extends Error {
  constructor(
    readonly code: string,
    readonly status: number,
    message: string,
  ) {
    super(message);
    this.name = 'MapsApiError';
  }
}

const configuredBasePath = import.meta.env.VITE_MAPS_API_BASE_PATH?.trim() || '/api/v1';
if (!configuredBasePath.startsWith('/') || configuredBasePath.startsWith('//')) {
  throw new Error('VITE_MAPS_API_BASE_PATH must be a same-origin absolute path.');
}
const apiBasePath = configuredBasePath.replace(/\/+$/, '');

export class MapsAPI {
  constructor(private readonly token: () => string | null) {}

  async capabilities(): Promise<ProviderCapabilities> {
    const response = await this.request<{ providers: ProviderCapabilities }>('/capabilities', {}, false);
    return response.providers;
  }

  async me(): Promise<{ id: string; displayName?: string; createdAt: string }> {
    return this.request('/me');
  }

  async search(query: string, limit = 10): Promise<SearchResult[]> {
    const params = new URLSearchParams({ q: query, limit: String(limit) });
    const response = await this.request<{ results: SearchResult[] }>(`/search?${params.toString()}`);
    return response.results;
  }

  async route(mode: string, locations: RoutePoint[]): Promise<Route> {
    return this.request('/routes', {
      method: 'POST',
      body: JSON.stringify({ mode, locations }),
    });
  }

  async listCollections(): Promise<Collection[]> {
    const response = await this.request<{ collections: Collection[] }>('/collections');
    return response.collections;
  }

  async createCollection(name: string, description = ''): Promise<Collection> {
    return this.request('/collections', {
      method: 'POST',
      body: JSON.stringify({ name, description }),
    });
  }

  async listCollectionItems(collectionID: string): Promise<CollectionItem[]> {
    const response = await this.request<{ items: CollectionItem[] }>(
      `/collections/${encodeURIComponent(collectionID)}/items`,
    );
    return response.items;
  }

  private async request<T>(path: string, init: RequestInit = {}, authenticated = true): Promise<T> {
    const headers = new Headers(init.headers);
    headers.set('Accept', 'application/json');
    if (init.body !== undefined) headers.set('Content-Type', 'application/json');
    if (authenticated) {
      const accessToken = this.token();
      if (!accessToken) throw new MapsApiError('auth_required', 401, 'Sign in to GoreeCloud Identity to continue.');
      headers.set('Authorization', `Bearer ${accessToken}`);
    }

    const response = await fetch(`${apiBasePath}${path}`, {
      ...init,
      headers,
      credentials: 'same-origin',
      cache: 'no-store',
    });
    if (!response.ok) {
      let code = 'request_failed';
      let message = `Maps API request failed with status ${response.status}.`;
      try {
        const payload = (await response.json()) as { error?: { code?: string; message?: string } };
        if (payload.error?.code) code = payload.error.code;
        if (payload.error?.message) message = payload.error.message;
      } catch {
        // Preserve the bounded generic error when the response is not JSON.
      }
      throw new MapsApiError(code, response.status, message);
    }
    return (await response.json()) as T;
  }
}
