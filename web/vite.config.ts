/// <reference types="vitest" />
import { defineConfig } from "vite";
import { svelte } from "@sveltejs/vite-plugin-svelte";

// https://vitejs.dev/config/
export default defineConfig({
	plugins: [svelte()],
	server: {
		proxy: {
			"/api": {
				target: "https://server-7fzc7ivuwa-ue.a.run.app",
				changeOrigin: true,
			},
		},
	},
	test: {
		environment: "jsdom",
		include: ["**/__tests__/*.test.ts"],
	},
});
