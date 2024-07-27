/// <reference types="vitest" />
import { defineConfig } from "vite";
import { svelte } from "@sveltejs/vite-plugin-svelte";

// https://vitejs.dev/config/
export default defineConfig({
	plugins: [svelte()],
	build: {
		rollupOptions: {
			output: {
				manualChunks(id) {
					const isVendorChunk = id.includes("node_modules");
					if (isVendorChunk) {
						const vendor = id.split("node_modules/")[1].split("/")[0];
						switch (vendor) {
							case "svelte":
							case "svelte-routing":
								return "vendor-svelte";

							case "ol":
								return "vendor-ol";

							case "chart.js":
								return "vendor-chart.js";

							default:
								return "vendor-misc";
						}
					}
				},
			},
		},
	},
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
