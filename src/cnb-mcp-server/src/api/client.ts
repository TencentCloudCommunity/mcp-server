export default class CnbApiClient {
  private static instance: CnbApiClient | null = null;

  private _baseUrl: string;
  private _token: string;

  private constructor(options: CnbApiOptions) {
    this._baseUrl = options.baseUrl;
    this._token = options.token;
  }

  get baseUrl() {
    return this._baseUrl;
  }

  static initialize(options: CnbApiOptions): void {
    if (!CnbApiClient.instance) {
      CnbApiClient.instance = new CnbApiClient(options);
    }
  }

  static getInstance(): CnbApiClient {
    if (!CnbApiClient.instance) {
      throw new Error('CnbApiClient not initialized. Call CnbApiClient.initialize(baseUrl, token) first.');
    }
    return CnbApiClient.instance;
  }

  async request<T>(
    method: string,
    path: string,
    body?: JsonValue,
    config?: {
      header?: Record<string, string>;
    },
    responseType = 'json' // 'json' | 'text' | 'raw'
  ): Promise<T> {
    const url = `${this._baseUrl}${path}`;
    const headers = {
      Authorization: `Bearer ${this._token}`,
      Accept: 'application/vnd.cnb.api+json',
      ...(config?.header || {})
    };

    const options: RequestInit = {
      method,
      headers,
      body: body ? JSON.stringify(body) : undefined
    };

    const response = await fetch(url, options);

    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(
        `API request failed: ${response.status} \n${response.headers.get('traceparent')} \n${errorText} )`
      );
    }

    if (responseType === 'raw') {
      return response as T;
    }

    if (responseType === 'text') {
      return response.text() as Promise<T>;
    }

    return response.json() as Promise<T>;
  }
}

export interface CnbApiOptions {
  baseUrl: string;
  token: string;
}

type JsonPrimitive = string | number | boolean | null;
type JsonArray = JsonValue[];
interface JsonObject {
  [key: string]: JsonValue;
}
// Any JSON value
type JsonValue = JsonPrimitive | JsonObject | JsonArray;
