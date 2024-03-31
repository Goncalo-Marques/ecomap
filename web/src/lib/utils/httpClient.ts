import createClient from "openapi-fetch";
import type { components, paths } from "../../../api/ecomap/http";

const httpClient = createClient<paths>({
	baseUrl: "/api",
	// Custom fetch implementation to handle native fetch exceptions.
	async fetch(input, init) {
		try {
			return await fetch(input, init);
		} catch {
			const body: components["schemas"]["Error"] = {
				code: "internal_server_error",
			};
			return new Response(JSON.stringify(body), {
				status: 500,
				statusText: "Internal Server Error",
			});
		}
	},
});

export default httpClient;
