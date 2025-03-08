import { includeIgnoreFile } from "@eslint/compat";
import eslint from "@eslint/js";
import simpleImportSort from "eslint-plugin-simple-import-sort";
import pluginVue from "eslint-plugin-vue";
import globals from "globals";
import path from "node:path";
import { fileURLToPath } from "node:url";
import tseslint from "typescript-eslint";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const gitignorePath = path.resolve(__dirname, ".gitignore");

export default tseslint.config([
  includeIgnoreFile(gitignorePath),
  { ignores: ["*.d.ts", "**/coverage", "**/dist"] },
  {
    extends: [
      eslint.configs.recommended,
      ...tseslint.configs.recommended,
      ...pluginVue.configs["flat/recommended"],
    ],
    files: ["**/*.{ts,vue}"],
    languageOptions: {
      ecmaVersion: "latest",
      sourceType: "module",
      parserOptions: {
        parser: tseslint.parser,
      },
    },
    rules: {
      "vue/multi-word-component-names": "off",
      "vue/attributes-order": [
        "error",
        {
          order: [
            "DEFINITION", // is, v-is
            "LIST_RENDERING", // v-for
            "CONDITIONALS", // v-if, v-else-if, v-else, v-show, v-cloak
            "RENDER_MODIFIERS", // v-pre, v-once
            "GLOBAL", // id
            "UNIQUE", // ref, key, slot
            "TWO_WAY_BINDING", // v-model
            "OTHER_DIRECTIVES", // v-custom-directives
            "OTHER_ATTR", // class, style, data-*
            "EVENTS", // @click, @change
            "CONTENT", // v-html, v-text
          ],
          alphabetical: false,
        },
      ],
    },
  },
  {
    plugins: {
      "simple-import-sort": simpleImportSort,
    },
    rules: {
      "simple-import-sort/exports": "error",
      "simple-import-sort/imports": [
        "error",
        {
          groups: [
            ["^vue(.*)$"],
            ["^@/stores/(.*)$"],
            ["^@/composables/(.*)$"],
            ["^@/utils/(.*)$"],
            ["^@/services/(.*)$"],
            ["^@/components/(.*)$"],
            ["^@/components/icons/(.*)$"],
            ["^[./]"],
          ],
        },
      ],
    },
  },
]);
