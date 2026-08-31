const API_BASE = import.meta.env.VITE_API_BASE || '/api';

interface ApiEnvelope<T> {
  success: boolean;
  code: number;
  message: string;
  data: T;
  errors?: unknown;
}

export async function apiPost<T>(path: string, body?: unknown): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: body ? JSON.stringify(body) : undefined,
  });
  const env: ApiEnvelope<T> = await res.json();
  if (!res.ok || !env.success) {
    throw new Error(env.message || `API error ${res.status}`);
  }
  return env.data;
}

export async function apiGet<T>(path: string): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, { headers: { 'Content-Type': 'application/json' } });
  const env: ApiEnvelope<T> = await res.json();
  if (!res.ok || !env.success) {
    throw new Error(env.message || `API error ${res.status}`);
  }
  return env.data;
}
