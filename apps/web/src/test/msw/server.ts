import { setupServer } from 'msw/node'
import { handlers } from './handlers'

// Serveur MSW utilisé dans les tests Node (Vitest + jsdom).
export const server = setupServer(...handlers)
