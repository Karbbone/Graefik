import { Outlet } from 'react-router-dom'
import { useAuth } from '@/features/auth'
import { ThemeToggle } from '@/shared/ui/ThemeToggle'

// AppLayout : coquille des pages protégées (barre de navigation + contenu).
export function AppLayout() {
  const { user, logout } = useAuth()

  return (
    <div className="min-h-full bg-base-200">
      <div className="navbar bg-base-100 border-b border-base-300 px-4">
        <div className="flex-1">
          <span className="text-xl font-bold text-primary">Graefik</span>
        </div>
        <div className="flex items-center gap-2">
          <ThemeToggle />
          {user && <span className="text-sm opacity-70">{user.username}</span>}
          <button
            type="button"
            className="btn btn-sm btn-ghost"
            onClick={() => void logout()}
          >
            Déconnexion
          </button>
        </div>
      </div>

      <main className="mx-auto max-w-3xl p-6">
        <Outlet />
      </main>
    </div>
  )
}
