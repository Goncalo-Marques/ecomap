import js from "@eslint/js";
import prettier from "eslint-config-prettier";
import simpleImportSort from "eslint-plugin-simple-import-sort";
import svelte from "eslint-plugin-svelte";
import globals from "globals";
import ts from "typescript-eslint";

import svelteConfig from "./svelte.config.js";

/** @type {import('eslint').Linter.Config[]} */
export default [
	js.configs.recommended,
	...ts.configs.recommended,
	...svelte.configs["flat/recommended"],
	prettier,
	...svelte.configs["flat/prettier"],
	{
		languageOptions: {
			globals: {
				...globals.browser,
				...globals.node,
			},
		},
	},
	{
		files: ["**/*.svelte"],
		languageOptions: {
			parserOptions: {
				svelteConfig,
				parser: ts.parser,
				svelteFeatures: {
					experimentalGenerics: true,
				},
			},
		},
	},
	{
		plugins: {
			"simple-import-sort": simpleImportSort,
		},
		rules: {
			"simple-import-sort/imports": [
				"error",
				{
					groups: [
						// Side effect imports.
						["^\\u0000"],
						// Svelte imports.
						["^svelte"],
						// Node.js builtins prefixed with `node:`.
						["^node:"],
						// Packages.
						// Things that start with a letter (or digit or underscore), or `@` followed by a letter.
						["^@?\\w"],
						// Absolute imports and other imports such as Svelte-style `$lib`.
						// Anything not matched in another group.
						["^"],
						// Relative imports.
						// Anything that starts with a dot.
						["^\\."],
					],
				},
			],
			"simple-import-sort/exports": "error",
		},
	},
	{
		ignores: [".svelte-kit/", "dist/"],
	},
];
