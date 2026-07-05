import js from '@eslint/js'
import globals from 'globals'
import reactHooks from 'eslint-plugin-react-hooks'
import reactRefresh from 'eslint-plugin-react-refresh'
import boundaries from 'eslint-plugin-boundaries'
import tseslint from 'typescript-eslint'

export default tseslint.config(
  { ignores: ['dist', 'coverage'] },

  {
    files: ['**/*.{ts,tsx}'],
    extends: [js.configs.recommended, tseslint.configs.recommended],
    plugins: {
      'react-hooks': reactHooks,
      'react-refresh': reactRefresh,
    },
    languageOptions: {
      ecmaVersion: 2022,
      globals: globals.browser,
    },
    rules: {
      ...reactHooks.configs['recommended-latest'].rules,
      'react-refresh/only-export-components': [
        'warn',
        { allowConstantExport: true },
      ],
    },
  },

  // Fichiers de config exécutés sous Node.
  {
    files: ['vite.config.ts', 'eslint.config.js'],
    languageOptions: {
      globals: globals.node,
    },
  },

  // --- Frontières feature-based (eslint-plugin-boundaries) ---
  {
    files: ['src/**/*.{ts,tsx}'],
    plugins: { boundaries },
    settings: {
      'import/resolver': {
        typescript: { alwaysTryTypes: true, project: './tsconfig.app.json' },
      },
      'boundaries/include': ['src/**/*'],
      'boundaries/ignore': [
        'src/main.tsx',
        'src/test/**/*',
        'src/**/*.test.{ts,tsx}',
        'src/vite-env.d.ts',
      ],
      'boundaries/elements': [
        { type: 'app', pattern: 'src/app/**/*' },
        { type: 'pages', pattern: 'src/pages/**/*' },
        {
          type: 'feature',
          pattern: 'src/features/*/**/*',
          capture: ['featureName'],
        },
        { type: 'shared', pattern: 'src/shared/**/*' },
      ],
    },
    rules: {
      // Une feature n'importe QUE du shared ou sa propre feature.
      // app/pages composent les features. shared reste autonome.
      'boundaries/dependencies': [
        'error',
        {
          default: 'disallow',
          rules: [
            {
              from: [{ type: 'app' }],
              allow: [
                { to: { type: 'app' } },
                { to: { type: 'pages' } },
                { to: { type: 'feature' } },
                { to: { type: 'shared' } },
              ],
            },
            {
              from: [{ type: 'pages' }],
              allow: [
                { to: { type: 'pages' } },
                { to: { type: 'feature' } },
                { to: { type: 'shared' } },
              ],
            },
            {
              from: [{ type: 'feature' }],
              allow: [
                { to: { type: 'shared' } },
                {
                  to: {
                    type: 'feature',
                    captured: { featureName: '{{ from.captured.featureName }}' },
                  },
                },
              ],
            },
            {
              from: [{ type: 'shared' }],
              allow: [{ to: { type: 'shared' } }],
            },
          ],
        },
      ],
    },
  },
)
