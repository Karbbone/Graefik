import { HealthCard } from '@/features/health'
import { TasksPanel } from '@/features/tasks'

// Une page compose des features. Elle n'implémente pas de logique métier.
export function HomePage() {
  return (
    <>
      <HealthCard />
      <TasksPanel />
    </>
  )
}
