import { HealthCard } from '@/features/health'
import { TasksPanel } from '@/features/tasks'

// Une page compose des features. Elle n'implémente pas de logique métier.
export function HomePage() {
  return (
    <div className="flex flex-col gap-6">
      <div>
        <h1 className="text-3xl font-bold">Tableau de bord</h1>
        <p className="opacity-60">
          Monitoring self-hosted pour reverse proxy (démo d&apos;architecture).
        </p>
      </div>
      <HealthCard />
      <TasksPanel />
    </div>
  )
}
