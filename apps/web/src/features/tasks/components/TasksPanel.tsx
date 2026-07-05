import { useTasks } from "../hooks/useTasks";
import { TaskForm } from "./TaskForm";

export function TasksPanel() {
  const { tasks, loading, error, add } = useTasks();

  return (
    <div className="card bg-base-100 shadow">
      <div className="card-body">
        <h2 className="card-title">Tâches</h2>

        <TaskForm onAdd={add} />

        {error && (
          <div role="alert" className="alert alert-error text-sm">
            {error}
          </div>
        )}

        {loading ? (
          <span className="loading loading-dots loading-md" />
        ) : tasks.length === 0 ? (
          <p className="opacity-60">Aucune tâche pour le moment.</p>
        ) : (
          <ul className="menu bg-base-200 rounded-box w-full">
            {tasks.map((task) => (
              <li key={task.id}>
                <span>{task.title}</span>
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  );
}
