import { useTasks } from '../hooks/useTasks'
import { TaskForm } from './TaskForm'

export function TasksPanel() {
  const { tasks, loading, error, add } = useTasks()

  return (
    <section className="card">
      <h2>Tâches</h2>

      <TaskForm onAdd={add} />

      {error && <p className="status status--error">{error}</p>}

      {loading ? (
        <p className="status">Chargement&hellip;</p>
      ) : tasks.length === 0 ? (
        <p className="status">Aucune tâche pour le moment.</p>
      ) : (
        <ul className="task-list">
          {tasks.map((task) => (
            <li key={task.id}>{task.title}</li>
          ))}
        </ul>
      )}
    </section>
  )
}
