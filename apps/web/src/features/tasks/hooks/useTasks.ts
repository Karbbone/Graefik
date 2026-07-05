import { useCallback, useEffect, useState } from 'react'
import { createTask, listTasks } from '../api/tasksApi'
import type { Task } from '../model/task'

export function useTasks() {
  const [tasks, setTasks] = useState<Task[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    listTasks()
      .then(setTasks)
      .catch((e: unknown) =>
        setError(e instanceof Error ? e.message : 'Erreur inconnue'),
      )
      .finally(() => setLoading(false))
  }, [])

  const add = useCallback(async (title: string) => {
    const task = await createTask(title)
    setTasks((prev) => [...prev, task])
  }, [])

  return { tasks, loading, error, add }
}
