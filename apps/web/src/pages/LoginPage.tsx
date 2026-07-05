import { LoginForm } from '@/features/auth'
import { ThemeToggle } from '@/shared/ui/ThemeToggle'

export function LoginPage() {
  return (
    <div className="min-h-full bg-base-200 flex items-center justify-center p-6">
      <div className="card w-full max-w-sm bg-base-100 shadow-xl">
        <div className="card-body">
          <div className="flex items-center justify-between">
            <h1 className="text-2xl font-bold text-primary">Graefik</h1>
            <ThemeToggle />
          </div>
          <p className="text-sm opacity-60 mb-2">
            Connectez-vous pour accéder au tableau de bord.
          </p>
          <LoginForm />
        </div>
      </div>
    </div>
  )
}
