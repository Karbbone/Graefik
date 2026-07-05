import { Outlet } from "react-router-dom";
import { useAuth } from "@/features/auth";
import { BrandLogo } from "@/shared/ui/BrandLogo";
import { ThemeToggle } from "@/shared/ui/ThemeToggle";

// AppLayout : coquille des pages protégées (barre de navigation + contenu).
export function AppLayout() {
  const { user, logout } = useAuth();

  return (
    <div className="min-h-full bg-base-200">
      <header className="sticky top-0 z-20 border-b border-base-300 bg-base-100/80 backdrop-blur">
        <div className="mx-auto flex max-w-5xl items-center justify-between px-6 py-3">
          <div className="flex items-center gap-2.5">
            <BrandLogo size={32} />
            <span className="font-display text-xl font-extrabold tracking-tight">
              graefik
            </span>
          </div>

          <div className="flex items-center gap-2">
            <ThemeToggle />
            {user && (
              <div className="hidden items-center gap-2 rounded-full border border-base-300 bg-base-200 py-1 pl-1 pr-3 sm:flex">
                <span className="grid size-6 place-items-center rounded-full bg-primary text-xs font-bold text-primary-content">
                  {user.username.charAt(0).toUpperCase()}
                </span>
                <span className="text-sm font-medium">{user.username}</span>
              </div>
            )}
            <button
              type="button"
              className="btn btn-sm btn-ghost"
              onClick={() => void logout()}
            >
              Déconnexion
            </button>
          </div>
        </div>
      </header>

      <main className="mx-auto max-w-5xl px-6 py-10">
        <Outlet />
      </main>
    </div>
  );
}
