import { useState } from "react";

type Props = {
  size?: number;
  className?: string;
};

// BrandLogo affiche le logo (/logo.png). Si le fichier n'est pas présent,
// il retombe sur un emblème SVG aux couleurs de la marque (clin d'œil mascotte).
export function BrandLogo({ size = 48, className = "" }: Props) {
  const [failed, setFailed] = useState(false);

  if (!failed) {
    return (
      <img
        src="/logo.png"
        alt="Logo Graefik"
        width={size}
        height={size}
        onError={() => setFailed(true)}
        style={{ width: size, height: size }}
        className={className}
      />
    );
  }

  return (
    <svg
      viewBox="0 0 100 100"
      width={size}
      height={size}
      role="img"
      aria-label="Logo Graefik"
      className={className}
      style={{ width: size, height: size }}
    >
      <rect
        x="6"
        y="6"
        width="88"
        height="88"
        rx="26"
        fill="var(--color-gfteal)"
      />
      <rect
        x="6"
        y="6"
        width="88"
        height="88"
        rx="26"
        fill="none"
        stroke="var(--color-gfink)"
        strokeWidth="5"
      />
      {/* lunettes */}
      <circle
        cx="37"
        cy="45"
        r="15"
        fill="var(--color-gfcream)"
        stroke="var(--color-gfink)"
        strokeWidth="5"
      />
      <circle
        cx="70"
        cy="45"
        r="13"
        fill="var(--color-gfcream)"
        stroke="var(--color-gfink)"
        strokeWidth="5"
      />
      <line
        x1="51"
        y1="44"
        x2="58"
        y2="44"
        stroke="var(--color-gfink)"
        strokeWidth="5"
      />
      {/* yeux */}
      <circle cx="37" cy="46" r="4.5" fill="var(--color-gfink)" />
      <circle cx="70" cy="46" r="4" fill="var(--color-gfink)" />
      {/* cravate */}
      <path d="M50 62 l6 6 -6 20 -6 -20 z" fill="var(--color-gfink)" />
    </svg>
  );
}
