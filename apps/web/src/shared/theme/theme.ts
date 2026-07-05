import { createContext, useContext } from "react";

export type Theme = "graefik-light" | "graefik-dark";

export const THEME_STORAGE_KEY = "graefik-theme";

export type ThemeContextValue = {
  theme: Theme;
  isDark: boolean;
  toggle: () => void;
};

export const ThemeContext = createContext<ThemeContextValue | null>(null);

export function useTheme(): ThemeContextValue {
  const ctx = useContext(ThemeContext);
  if (!ctx) {
    throw new Error("useTheme doit être utilisé dans un ThemeProvider");
  }
  return ctx;
}

// getInitialTheme : préférence sauvegardée, sinon préférence système.
export function getInitialTheme(): Theme {
  const saved = localStorage.getItem(THEME_STORAGE_KEY);
  if (saved === "graefik-light" || saved === "graefik-dark") {
    return saved;
  }
  return window.matchMedia("(prefers-color-scheme: dark)").matches
    ? "graefik-dark"
    : "graefik-light";
}
