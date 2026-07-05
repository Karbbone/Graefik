// ESLint est réduit à sa seule responsabilité que Biome ne couvre pas :
// l'enforcement des frontières de l'architecture feature-based
// (eslint-plugin-boundaries). Le formatage et le reste du lint sont gérés par Biome.
import boundaries from "eslint-plugin-boundaries";
import tseslint from "typescript-eslint";

export default tseslint.config(
  { ignores: ["dist", "coverage"] },
  {
    files: ["src/**/*.{ts,tsx}"],
    languageOptions: {
      parser: tseslint.parser,
    },
    plugins: { boundaries },
    settings: {
      "import/resolver": {
        typescript: { alwaysTryTypes: true, project: "./tsconfig.app.json" },
      },
      "boundaries/include": ["src/**/*"],
      "boundaries/ignore": [
        "src/main.tsx",
        "src/test/**/*",
        "src/**/*.test.{ts,tsx}",
        "src/vite-env.d.ts",
      ],
      "boundaries/elements": [
        { type: "app", pattern: "src/app/**/*" },
        { type: "pages", pattern: "src/pages/**/*" },
        {
          type: "feature",
          pattern: "src/features/*/**/*",
          capture: ["featureName"],
        },
        { type: "shared", pattern: "src/shared/**/*" },
      ],
    },
    rules: {
      "boundaries/dependencies": [
        "error",
        {
          default: "disallow",
          rules: [
            {
              from: [{ type: "app" }],
              allow: [
                { to: { type: "app" } },
                { to: { type: "pages" } },
                { to: { type: "feature" } },
                { to: { type: "shared" } },
              ],
            },
            {
              from: [{ type: "pages" }],
              allow: [
                { to: { type: "pages" } },
                { to: { type: "feature" } },
                { to: { type: "shared" } },
              ],
            },
            {
              from: [{ type: "feature" }],
              allow: [
                { to: { type: "shared" } },
                {
                  to: {
                    type: "feature",
                    captured: {
                      featureName: "{{ from.captured.featureName }}",
                    },
                  },
                },
              ],
            },
            {
              from: [{ type: "shared" }],
              allow: [{ to: { type: "shared" } }],
            },
          ],
        },
      ],
    },
  },
);
