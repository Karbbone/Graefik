import { Navigate, Outlet } from 'react-router-dom'
import { useAuth } from '@/features/auth'

// ProtectedRoute : bloque l'accès aux routes enfants tant que l'utilisateur
// n'est pas authentifié (redirige vers /login).
export function ProtectedRoute() {
  const { user, loading } = useAuth()

  if (loading) {
    return (
      <div className="flex h-full items-center justify-center">
        <span className="loading loading-spinner loading-lg text-primary" />
      </div>
    )
  }

  if (!user) {
    return <Navigate to="/login" replace />
  }

  return <Outlet />
}
