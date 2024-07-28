import { sveltekit } from "@sveltejs/kit/vite";
import { defineConfig } from "vitest/config";

export default defineConfig({
	plugins: [sveltekit()],
	server: {
		proxy: {
			"/api": {
				target: "http://127.0.0.1:8080",
				changeOrigin: true,
			},
		},
	},
	test: {
		environment: "jsdom",
		include: ["**/__tests__/*.test.ts"],
	},
});
