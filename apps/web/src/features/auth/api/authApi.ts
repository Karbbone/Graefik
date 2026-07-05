import { apiFetch } from '@/shared/lib/http'
import type { User } from '../model/user'

export function login(username: string, password: string): Promise<User> {
  return apiFetch<User>('/api/auth/login', {
    method: 'POST',
    body: JSON.stringify({ username, password }),
  })
}

// logout renvoie 204 (pas de corps) : on n'utilise pas apiFetch (qui parse du JSON).
export async function logout(): Promise<void> {
  await fetch('/api/auth/logout', { method: 'POST' })
}

export function getCurrentUser(): Promise<User> {
  return apiFetch<User>('/api/auth/me')
}
