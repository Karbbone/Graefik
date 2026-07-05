import { useEffect, useState } from 'react'
import './App.css'

type Health = {
  status: string
  service: string
  go: string
  time: string
}

function App() {
  const [health, setHealth] = useState<Health | null>(null)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    fetch('/api/health')
      .then((res) => {
        if (!res.ok) throw new Error(`HTTP ${res.status}`)
        return res.json() as Promise<Health>
      })
      .then(setHealth)
      .catch((e: unknown) =>
        setError(e instanceof Error ? e.message : 'Erreur inconnue'),
      )
  }, [])

  return (
    <main className="app">
      <header className="hero">
        <h1>Graefik</h1>
        <p className="subtitle">
          Monorepo <strong>Go (Echo)</strong> + <strong>React (Vite)</strong>{' '}
          &mdash; Turborepo &amp; Docker
        </p>
      </header>

      <section className="card">
        <h2>Statut de l&apos;API</h2>

        {error && (
          <p className="status status--error">
            Impossible de joindre l&apos;API : {error}
          </p>
        )}

        {!error && !health && (
          <p className="status">Connexion &agrave; l&apos;API&hellip;</p>
        )}

        {health && (
          <ul className="health">
            <li>
              <span>Statut</span>
              <strong className="badge badge--ok">{health.status}</strong>
            </li>
            <li>
              <span>Service</span>
              <strong>{health.service}</strong>
            </li>
            <li>
              <span>Version Go</span>
              <strong>{health.go}</strong>
            </li>
            <li>
              <span>Horodatage</span>
              <strong>{health.time}</strong>
            </li>
          </ul>
        )}
      </section>

      <footer className="footer">
        Modifie <code>apps/web/src/App.tsx</code> ou{' '}
        <code>apps/api/internal/handler/health.go</code> &mdash; le hot reload
        s&apos;occupe du reste.
      </footer>
    </main>
  )
}

export default App
