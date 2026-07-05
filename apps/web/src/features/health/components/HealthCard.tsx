import { useHealth } from '../hooks/useHealth'

export function HealthCard() {
  const { health, error } = useHealth()

  return (
    <div className="card bg-base-100 shadow">
      <div className="card-body">
        <h2 className="card-title">Statut de l&apos;API</h2>

        {error && (
          <div role="alert" className="alert alert-error text-sm">
            Impossible de joindre l&apos;API : {error}
          </div>
        )}

        {!error && !health && (
          <span className="loading loading-dots loading-md" />
        )}

        {health && (
          <div className="flex flex-wrap items-center gap-3">
            <span className="badge badge-success">{health.status}</span>
            <span className="text-sm opacity-70">{health.service}</span>
            <span className="text-sm opacity-70">{health.go}</span>
          </div>
        )}
      </div>
    </div>
  )
}
