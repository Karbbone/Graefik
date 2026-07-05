import { HomePage } from '@/pages/HomePage'
import './App.css'

export function App() {
  return (
    <main className="app">
      <header className="hero">
        <h1>Graefik</h1>
        <p className="subtitle">
          Monorepo <strong>Go</strong> (Echo, hexagonal) +{' '}
          <strong>React</strong> (Vite, feature-based)
        </p>
      </header>

      <HomePage />

      <footer className="footer">
        Modifie une <code>feature</code> dans <code>src/features</code> ou un{' '}
        <code>service</code> dans <code>apps/api/internal/core</code> &mdash; le
        hot reload s&apos;occupe du reste.
      </footer>
    </main>
  )
}
