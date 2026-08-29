type OIDCDiscovery = {
  issuer: string;
  authorization_endpoint: string;
  token_endpoint: string;
};

type TokenResponse = {
  access_token?: string;
  token_type?: string;
  expires_in?: number;
};

const pkceVerifierKey = 'goreecloud.maps.oidc.pkce_verifier';
const stateKey = 'goreecloud.maps.oidc.state';

const trimTrailingSlash = (value: string): string => value.replace(/\/+$/, '');

const validateEndpoint = (value: string, label: string): URL => {
  let parsed: URL;
  try {
    parsed = new URL(value);
  } catch {
    throw new Error(`${label} must be an absolute URL.`);
  }
  const localDevelopment =
    parsed.protocol === 'http:' && (parsed.hostname === 'localhost' || parsed.hostname === '127.0.0.1');
  if (parsed.protocol !== 'https:' && !localDevelopment) {
    throw new Error(`${label} must use HTTPS outside local development.`);
  }
  if (parsed.username || parsed.password || parsed.search || parsed.hash) {
    throw new Error(`${label} must not contain credentials, query parameters, or fragments.`);
  }
  return parsed;
};

const base64URL = (bytes: Uint8Array): string => {
  let binary = '';
  for (const byte of bytes) binary += String.fromCharCode(byte);
  return btoa(binary).replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/g, '');
};

const randomValue = (bytes = 32): string => {
  const value = new Uint8Array(bytes);
  crypto.getRandomValues(value);
  return base64URL(value);
};

const sha256Challenge = async (verifier: string): Promise<string> => {
  const digest = await crypto.subtle.digest('SHA-256', new TextEncoder().encode(verifier));
  return base64URL(new Uint8Array(digest));
};

export class IdentityClient {
  private readonly issuer = import.meta.env.VITE_IDENTITY_ISSUER_URL?.trim() ?? '';
  private readonly clientID = import.meta.env.VITE_IDENTITY_CLIENT_ID?.trim() ?? '';
  private readonly scope = import.meta.env.VITE_IDENTITY_SCOPES?.trim() || 'openid profile';
  private readonly redirectURI =
    import.meta.env.VITE_IDENTITY_REDIRECT_URI?.trim() || `${window.location.origin}${window.location.pathname}`;
  private discovery: OIDCDiscovery | null = null;
  private accessToken: string | null = null;
  private accessTokenExpiresAt = 0;

  get configured(): boolean {
    return this.issuer !== '' && this.clientID !== '';
  }

  get authenticated(): boolean {
    return this.getAccessToken() !== null;
  }

  getAccessToken(): string | null {
    if (!this.accessToken) return null;
    if (this.accessTokenExpiresAt > 0 && Date.now() >= this.accessTokenExpiresAt - 30_000) {
      this.accessToken = null;
      this.accessTokenExpiresAt = 0;
      return null;
    }
    return this.accessToken;
  }

  async initialize(): Promise<'none' | 'authenticated'> {
    const current = new URL(window.location.href);
    const code = current.searchParams.get('code');
    const error = current.searchParams.get('error');
    if (!code && !error) return 'none';
    if (!this.configured) {
      this.clearCallbackParameters(current);
      throw new Error('Identity callback received, but GoreeCloud Identity is not configured for Maps.');
    }
    if (error) {
      const description = current.searchParams.get('error_description')?.trim();
      this.clearTransientState();
      this.clearCallbackParameters(current);
      throw new Error(description || 'GoreeCloud Identity did not complete sign-in.');
    }
    if (!code) {
      this.clearTransientState();
      this.clearCallbackParameters(current);
      throw new Error('GoreeCloud Identity callback did not include an authorization code.');
    }

    const expectedState = sessionStorage.getItem(stateKey);
    const verifier = sessionStorage.getItem(pkceVerifierKey);
    const receivedState = current.searchParams.get('state');
    if (!expectedState || !verifier || !receivedState || expectedState !== receivedState) {
      this.clearTransientState();
      this.clearCallbackParameters(current);
      throw new Error('Identity callback state validation failed.');
    }

    const discovery = await this.getDiscovery();
    const body = new URLSearchParams({
      grant_type: 'authorization_code',
      client_id: this.clientID,
      code,
      redirect_uri: this.validatedRedirectURI(),
      code_verifier: verifier,
    });
    const response = await fetch(discovery.token_endpoint, {
      method: 'POST',
      headers: { 'Content-Type': 'application/x-www-form-urlencoded', Accept: 'application/json' },
      body,
      credentials: 'omit',
      cache: 'no-store',
    });
    if (!response.ok) {
      this.clearTransientState();
      this.clearCallbackParameters(current);
      throw new Error('GoreeCloud Identity token exchange failed.');
    }

    const token = (await response.json()) as TokenResponse;
    if (!token.access_token || !token.token_type || token.token_type.toLowerCase() !== 'bearer') {
      this.clearTransientState();
      this.clearCallbackParameters(current);
      throw new Error('GoreeCloud Identity did not return a usable bearer access token.');
    }

    this.accessToken = token.access_token;
    this.accessTokenExpiresAt =
      typeof token.expires_in === 'number' && token.expires_in > 0 ? Date.now() + token.expires_in * 1000 : 0;
    this.clearTransientState();
    this.clearCallbackParameters(current);
    return 'authenticated';
  }

  async signIn(): Promise<void> {
    if (!this.configured) throw new Error('GoreeCloud Identity is not configured for this Maps build.');
    const discovery = await this.getDiscovery();
    const verifier = randomValue(64);
    const state = randomValue(32);
    const challenge = await sha256Challenge(verifier);
    sessionStorage.setItem(pkceVerifierKey, verifier);
    sessionStorage.setItem(stateKey, state);

    const authorization = new URL(discovery.authorization_endpoint);
    authorization.searchParams.set('response_type', 'code');
    authorization.searchParams.set('client_id', this.clientID);
    authorization.searchParams.set('redirect_uri', this.validatedRedirectURI());
    authorization.searchParams.set('scope', this.scope);
    authorization.searchParams.set('state', state);
    authorization.searchParams.set('code_challenge', challenge);
    authorization.searchParams.set('code_challenge_method', 'S256');
    window.location.assign(authorization.toString());
  }

  disconnect(): void {
    this.accessToken = null;
    this.accessTokenExpiresAt = 0;
    this.clearTransientState();
  }

  private async getDiscovery(): Promise<OIDCDiscovery> {
    if (this.discovery) return this.discovery;
    const issuerURL = validateEndpoint(this.issuer, 'Identity issuer');
    const issuer = `${trimTrailingSlash(issuerURL.toString())}/`;
    const discoveryURL = new URL('.well-known/openid-configuration', issuer);
    const response = await fetch(discoveryURL, { headers: { Accept: 'application/json' }, cache: 'no-store' });
    if (!response.ok) throw new Error('GoreeCloud Identity discovery is unavailable.');
    const discovery = (await response.json()) as Partial<OIDCDiscovery>;
    if (!discovery.issuer || !discovery.authorization_endpoint || !discovery.token_endpoint) {
      throw new Error('GoreeCloud Identity discovery is incomplete.');
    }
    if (trimTrailingSlash(discovery.issuer) !== trimTrailingSlash(this.issuer)) {
      throw new Error('GoreeCloud Identity discovery returned an unexpected issuer.');
    }
    validateEndpoint(discovery.authorization_endpoint, 'Identity authorization endpoint');
    validateEndpoint(discovery.token_endpoint, 'Identity token endpoint');
    this.discovery = discovery as OIDCDiscovery;
    return this.discovery;
  }

  private validatedRedirectURI(): string {
    const redirect = new URL(this.redirectURI, window.location.origin);
    if (redirect.origin !== window.location.origin) {
      throw new Error('Maps OIDC redirect URI must remain on the current application origin.');
    }
    if (redirect.username || redirect.password || redirect.hash) {
      throw new Error('Maps OIDC redirect URI is invalid.');
    }
    return redirect.toString();
  }

  private clearTransientState(): void {
    sessionStorage.removeItem(pkceVerifierKey);
    sessionStorage.removeItem(stateKey);
  }

  private clearCallbackParameters(current: URL): void {
    for (const key of ['code', 'state', 'session_state', 'iss', 'error', 'error_description']) {
      current.searchParams.delete(key);
    }
    window.history.replaceState({}, document.title, `${current.pathname}${current.search}${current.hash}`);
  }
}
