import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAuth } from '../context/auth'

export function LoginForm() {
  const { login } = useAuth()
  const navigate = useNavigate()
  const [username, setUsername] = useState('graefik')
  const [password, setPassword] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [submitting, setSubmitting] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)
    setSubmitting(true)
    try {
      await login(username, password)
      navigate('/', { replace: true })
    } catch {
      setError('Identifiants invalides')
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <form onSubmit={handleSubmit} className="flex flex-col gap-4">
      <label className="form-control w-full">
        <span className="label-text mb-1">Identifiant</span>
        <input
          className="input input-bordered w-full"
          aria-label="Identifiant"
          autoComplete="username"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
        />
      </label>

      <label className="form-control w-full">
        <span className="label-text mb-1">Mot de passe</span>
        <input
          type="password"
          className="input input-bordered w-full"
          aria-label="Mot de passe"
          autoComplete="current-password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
        />
      </label>

      {error && (
        <div role="alert" className="alert alert-error text-sm">
          {error}
        </div>
      )}

      <button type="submit" className="btn btn-primary" disabled={submitting}>
        {submitting ? (
          <span className="loading loading-spinner loading-sm" />
        ) : (
          'Se connecter'
        )}
      </button>
    </form>
  )
}
