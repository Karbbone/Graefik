import { LoginForm } from "@/features/auth";
import { BrandLogo } from "@/shared/ui/BrandLogo";
import { ThemeToggle } from "@/shared/ui/ThemeToggle";

export function LoginPage() {
  return (
    <div className="gf-grain gf-mesh relative flex min-h-full items-center justify-center p-6">
      <div className="gf-grid absolute inset-0" aria-hidden />

      <div className="absolute right-4 top-4 z-10">
        <ThemeToggle />
      </div>

      <div className="gf-reveal relative w-full max-w-sm">
        <div className="mb-6 flex flex-col items-center text-center">
          <div className="relative mb-3 grid size-20 place-items-center">
            <div
              className="gf-glow absolute inset-0 rounded-full blur-xl"
              aria-hidden
            />
            <BrandLogo size={72} className="relative gf-float" />
          </div>
          <h1 className="font-display text-3xl font-extrabold tracking-tight">
            Graefik
          </h1>
          <p className="mt-1 text-sm text-base-content/60">
            Observabilité self-hosted pour vos reverse proxies.
          </p>
        </div>

        <div className="rounded-box border border-base-300 bg-base-100/90 p-6 shadow-xl backdrop-blur">
          <LoginForm />
        </div>

        <p className="mt-4 text-center text-xs text-base-content/50">
          Premier lancement&nbsp;? Le mot de passe de{" "}
          <code className="font-mono">graefik</code> s&apos;affiche dans les
          logs du conteneur.
        </p>
      </div>
    </div>
  );
}
