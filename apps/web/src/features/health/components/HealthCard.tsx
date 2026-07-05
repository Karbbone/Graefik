import { useHealth } from '../hooks/useHealth'

export function HealthCard() {
  const { health, error } = useHealth()

  return (
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
        <ul className="kv">
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
        </ul>
      )}
    </section>
  )
}
