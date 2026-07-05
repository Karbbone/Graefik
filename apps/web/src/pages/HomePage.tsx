import type { CSSProperties, ReactNode } from "react";
import { useAuth } from "@/features/auth";
import { HealthCard } from "@/features/health";
import { BrandLogo } from "@/shared/ui/BrandLogo";

// Délai de révélation (variable CSS --gf-d) typé pour le style inline.
const delay = (ms: number): CSSProperties =>
  ({ "--gf-d": `${ms}ms` }) as CSSProperties;

export function HomePage() {
  const { user } = useAuth();

  return (
    <div className="flex flex-col gap-10">
      {/* Hero */}
      <section className="gf-grain relative overflow-hidden rounded-box border border-base-300 gf-mesh">
        <div className="gf-grid absolute inset-0" aria-hidden />
        <div className="relative grid gap-8 p-8 sm:p-12 md:grid-cols-[1.4fr_1fr] md:items-center">
          <div>
            <span
              className="gf-reveal inline-flex items-center gap-2 rounded-full border border-base-content/10 bg-base-100/70 px-3 py-1 text-xs font-semibold uppercase tracking-widest text-gfteal-deep backdrop-blur"
              style={delay(0)}
            >
              ● Observabilité self-hosted
            </span>

            <h1
              className="gf-reveal mt-5 font-display text-5xl font-extrabold leading-[1.05] tracking-tight text-base-content sm:text-6xl"
              style={delay(80)}
            >
              Bonjour,{" "}
              <span className="text-gfteal-deep">
                {user?.username ?? "graefik"}
              </span>
              .
            </h1>

            <p
              className="gf-reveal mt-4 max-w-md text-base leading-relaxed text-base-content/70"
              style={delay(160)}
            >
              Graefik transforme les métriques brutes de vos reverse proxies en
              tableaux de bord clairs&nbsp;: trafic, latence et santé de chaque
              routeur, sans assembler toute une stack.
            </p>

            <div
              className="gf-reveal mt-7 flex flex-wrap items-center gap-3"
              style={delay(240)}
            >
              <button type="button" className="btn btn-primary" disabled>
                Connecter une source
                <span className="badge badge-sm border-0 bg-base-100/25 text-primary-content">
                  bientôt
                </span>
              </button>
              <a
                href="https://github.com/Karbbone/Graefik"
                target="_blank"
                rel="noreferrer noopener"
                className="btn btn-ghost"
              >
                Documentation
              </a>
            </div>
          </div>

          {/* Mascotte */}
          <div className="relative mx-auto flex aspect-square w-52 items-center justify-center sm:w-64">
            <div
              className="gf-glow absolute inset-0 rounded-full blur-2xl"
              aria-hidden
            />
            <div
              className="gf-ring absolute inset-4 rounded-full border border-gfteal/50"
              aria-hidden
            />
            <div
              className="gf-ring absolute inset-4 rounded-full border border-gfteal/40"
              style={{ animationDelay: "1.6s" }}
              aria-hidden
            />
            <BrandLogo
              size={210}
              className="gf-float relative drop-shadow-xl"
            />
          </div>
        </div>
      </section>

      {/* Statistiques */}
      <section className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <HealthCard />
        <StatCard label="Sources" value="0" hint="reverse proxies connectés" />
        <StatCard label="Routeurs" value="0" hint="routeurs surveillés" />
        <StatCard label="Alertes" value="0" hint="règles actives" />
      </section>

      {/* Feuille de route */}
      <section>
        <div className="mb-4 flex items-end justify-between">
          <h2 className="font-display text-2xl font-bold">
            Bientôt sur Graefik
          </h2>
          <span className="text-sm opacity-50">feuille de route</span>
        </div>
        <div className="grid gap-4 md:grid-cols-3">
          <RoadmapCard
            icon="📈"
            title="Métriques temps réel"
            text="Collecte et historisation des métriques de vos proxies dans une base time-series."
          />
          <RoadmapCard
            icon="🌐"
            title="Vue par routeur"
            text="Trafic, latence, codes HTTP et taux d'erreur, routeur par routeur, service par service."
          />
          <RoadmapCard
            icon="🔔"
            title="Alertes"
            text="Des seuils simples et des notifications quand une métrique dérape."
          />
        </div>
      </section>
    </div>
  );
}

function StatCard({
  label,
  value,
  hint,
}: {
  label: string;
  value: string;
  hint: string;
}) {
  return (
    <div className="rounded-box border border-base-300 bg-base-100 p-5 shadow-sm">
      <span className="text-xs font-semibold uppercase tracking-widest opacity-50">
        {label}
      </span>
      <div className="mt-3 font-mono text-3xl font-semibold text-base-content">
        {value}
      </div>
      <div className="mt-1 text-xs opacity-60">{hint}</div>
    </div>
  );
}

function RoadmapCard({
  icon,
  title,
  text,
}: {
  icon: ReactNode;
  title: string;
  text: string;
}) {
  return (
    <div className="group rounded-box border border-base-300 bg-base-100 p-6 shadow-sm transition-all hover:-translate-y-1 hover:border-gfteal hover:shadow-md">
      <div className="flex items-center justify-between">
        <span className="grid size-11 place-items-center rounded-field bg-gfteal/15 text-xl">
          {icon}
        </span>
        <span className="badge badge-ghost badge-sm">bientôt</span>
      </div>
      <h3 className="mt-4 font-display text-lg font-bold">{title}</h3>
      <p className="mt-1 text-sm leading-relaxed text-base-content/70">
        {text}
      </p>
    </div>
  );
}
