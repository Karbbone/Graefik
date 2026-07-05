import { useState } from 'react'

type Props = {
  onAdd: (title: string) => void | Promise<void>
}

export function TaskForm({ onAdd }: Props) {
  const [title, setTitle] = useState('')

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    const value = title.trim()
    if (!value) return
    await onAdd(value)
    setTitle('')
  }

  return (
    <form onSubmit={handleSubmit} className="task-form">
      <input
        aria-label="Titre de la tâche"
        placeholder="Nouvelle tâche…"
        value={title}
        onChange={(e) => setTitle(e.target.value)}
      />
      <button type="submit">Ajouter</button>
    </form>
  )
}
