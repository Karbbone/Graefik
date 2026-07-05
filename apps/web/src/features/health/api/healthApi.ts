import { apiFetch } from '@/shared/lib/http'
import type { Health } from '../model/health'

export function getHealth(): Promise<Health> {
  return apiFetch<Health>('/api/health')
}
