import { useEffect, useState, type ReactNode } from 'react'
import {
  getInitialTheme,
  THEME_STORAGE_KEY,
  ThemeContext,
  type Theme,
} from './theme'

export function ThemeProvider({ children }: { children: ReactNode }) {
  const [theme, setTheme] = useState<Theme>(getInitialTheme)

  useEffect(() => {
    document.documentElement.setAttribute('data-theme', theme)
    localStorage.setItem(THEME_STORAGE_KEY, theme)
  }, [theme])

  const toggle = () =>
    setTheme((t) => (t === 'graefik-dark' ? 'graefik-light' : 'graefik-dark'))

  return (
    <ThemeContext.Provider
      value={{ theme, isDark: theme === 'graefik-dark', toggle }}
    >
      {children}
    </ThemeContext.Provider>
  )
}
