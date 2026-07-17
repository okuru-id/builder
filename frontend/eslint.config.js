import js from '@eslint/js'
import globals from 'globals'
import tseslint from 'typescript-eslint'
import pluginVue from 'eslint-plugin-vue'
import { includeIgnoreFile } from '@eslint/compat'
import { resolve } from 'node:path'

const gitignorePath = resolve(import.meta.dirname, '../.gitignore')

export default tseslint.config(
  includeIgnoreFile(gitignorePath),
  js.configs.recommended,
  ...tseslint.configs.recommended,
  ...pluginVue.configs['flat/recommended'],
  {
    files: ['*.vue', '**/*.vue'],
    languageOptions: {
      parserOptions: { parser: tseslint.parser },
    },
  },
  {
    languageOptions: {
      globals: { ...globals.browser },
    },
    rules: {
      'vue/multi-word-component-names': 'off',
      'vue/max-attributes-per-line': 'off',
      'vue/singleline-html-element-content-newline': 'off',
      'vue/html-self-closing': 'off',
      'no-undef': 'off',
      'no-empty': 'off',
      'no-sparse-arrays': 'off',
      'no-cond-assign': 'off',
      'no-prototype-builtins': 'off',
      'no-redeclare': 'off',
      'no-setter-return': 'off',
      'no-self-assign': 'off',
      'no-func-assign': 'off',
      'no-unassigned-vars': 'off',
      'getter-return': 'off',
      'preserve-caught-error': 'off',
      'no-unsafe-optional-chaining': 'off',
      'no-useless-assignment': 'off',
      'no-useless-escape': 'off',
      '@typescript-eslint/no-explicit-any': 'warn',
      '@typescript-eslint/no-unused-vars': ['warn', { argsIgnorePattern: '^_' }],
      '@typescript-eslint/no-unused-expressions': 'off',
      '@typescript-eslint/no-this-alias': 'off',
    },
  },
)
