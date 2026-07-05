import { useHealth } from "../hooks/useHealth";

export function HealthCard() {
  const { health, error } = useHealth();
  const online = !!health && !error;
  const label = online ? "En ligne" : error ? "Hors ligne" : "Connexion…";

  return (
    <div className="rounded-box border border-base-300 bg-base-100 p-5 shadow-sm">
      <div className="flex items-center justify-between">
        <span className="text-xs font-semibold uppercase tracking-widest opacity-50">
          API
        </span>
        <span className="relative flex size-2.5">
          {online && (
            <span className="absolute inline-flex size-full animate-ping rounded-full bg-success opacity-70" />
          )}
          <span
            className={`relative inline-flex size-2.5 rounded-full ${
              online ? "bg-success" : error ? "bg-error" : "bg-warning"
            }`}
          />
        </span>
      </div>
      <div className="mt-3 font-display text-2xl font-extrabold text-base-content">
        {label}
      </div>
      <div className="mt-1 font-mono text-xs opacity-60">
        {health ? health.go : error ? error : "—"}
      </div>
    </div>
  );
}
