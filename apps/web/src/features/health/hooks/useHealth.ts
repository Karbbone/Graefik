import { useEffect, useState } from 'react'
import { getHealth } from '../api/healthApi'
import type { Health } from '../model/health'

type State = {
  health: Health | null
  error: string | null
}

export function useHealth(): State {
  const [state, setState] = useState<State>({ health: null, error: null })

  useEffect(() => {
    getHealth()
      .then((health) => setState({ health, error: null }))
      .catch((e: unknown) =>
        setState({
          health: null,
          error: e instanceof Error ? e.message : 'Erreur inconnue',
        }),
      )
  }, [])

  return state
}
