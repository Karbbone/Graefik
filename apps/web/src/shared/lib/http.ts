// Client HTTP minimal partagé par toutes les features.
// Les appels passent par le proxy Vite (/api -> backend Go) configuré dans vite.config.ts.
export async function apiFetch<T>(
  path: string,
  init?: RequestInit,
): Promise<T> {
  const res = await fetch(path, {
    headers: {
      "Content-Type": "application/json",
      ...init?.headers,
    },
    ...init,
  });

  if (!res.ok) {
    throw new Error(`Requête ${path} échouée (${res.status})`);
  }

  return res.json() as Promise<T>;
}
