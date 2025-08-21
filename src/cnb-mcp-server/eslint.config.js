import js from '@eslint/js';
import tseslint from 'typescript-eslint';
import eslintPluginPrettierRecommended from 'eslint-plugin-prettier/recommended';

/**
 * https://typescript-eslint.io/packages/typescript-eslint/#usage-with-other-plugins
 */
export default tseslint.config(
  {
    // config with just ignores is the replacement for `.eslintignore`
    ignores: ['dist/**', 'node_modules/**', '**/*.d.ts']
  },
  js.configs.recommended,
  tseslint.configs.recommended,
  {
    plugins: {
      '@typescript-eslint': tseslint.plugin
    },
    languageOptions: {
      parser: tseslint.parser,
      parserOptions: {
        projectService: true,
        tsconfigRootDir: import.meta.dirname
      },
      sourceType: 'module'
    },
    rules: {
      '@typescript-eslint/no-floating-promises': 'error',
      '@typescript-eslint/no-misused-promises': 'warn',
      '@typescript-eslint/no-unused-vars': 'warn',
      '@typescript-eslint/member-ordering': 'warn',
      '@typescript-eslint/prefer-optional-chain': 'warn',
      '@typescript-eslint/indent': ['error', 2]
    }
  },
  eslintPluginPrettierRecommended,
  {
    // disable type-aware linting on JS files
    files: ['**/*.js'],
    extends: [tseslint.configs.disableTypeChecked]
  }
);
