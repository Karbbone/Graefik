import { http, HttpResponse } from 'msw'
import type { Task } from '@/features/tasks'

// Handlers MSW par défaut, partagés par les tests.
// Un test peut les surcharger via server.use(...).
export const handlers = [
  http.get('/api/health', () =>
    HttpResponse.json({
      status: 'ok',
      service: 'graefik-api',
      go: 'go1.26',
      time: '2026-01-01T00:00:00Z',
    }),
  ),

  http.get('/api/tasks', () => HttpResponse.json<Task[]>([])),

  http.post('/api/tasks', async ({ request }) => {
    const body = (await request.json()) as { title: string }
    return HttpResponse.json<Task>(
      {
        id: 'generated-id',
        title: body.title,
        done: false,
        createdAt: '2026-01-01T00:00:00Z',
      },
      { status: 201 },
    )
  }),
]
