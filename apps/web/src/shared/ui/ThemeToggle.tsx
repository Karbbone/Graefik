import { useTheme } from "../theme/theme";

export function ThemeToggle() {
  const { isDark, toggle } = useTheme();

  return (
    <button
      type="button"
      className="btn btn-ghost btn-circle"
      onClick={toggle}
      aria-label={isDark ? "Passer en thème clair" : "Passer en thème sombre"}
    >
      {isDark ? "☀️" : "🌙"}
    </button>
  );
}
