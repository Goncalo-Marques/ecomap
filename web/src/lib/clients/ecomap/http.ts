import { navigate } from "svelte-routing";
import createClient from "openapi-fetch";
import type { components, paths } from "../../../../api/ecomap/http";
import { getToken } from "../../utils/auth";
import { CommonRoutes } from "../../../routes/constants/routes";

const ecomapHttpClient = createClient<paths>({
	baseUrl: "/api",
	// Custom fetch implementation to handle native fetch exceptions.
	async fetch(input, init) {
		try {
			// Append Authorization header.
			const headers = new Headers({
				...init?.headers,
				Authorization: `Bearer ${getToken()}`,
			});

			// Build request init with the new headers.
			const requestInit: RequestInit = {
				...init,
				headers,
			};

			const response = await fetch(input, requestInit);

			switch (response.status) {
				case 401:
					// Only redirect if the page is not the sign in page.
					if (location.pathname !== CommonRoutes.SIGN_IN) {
						navigate(CommonRoutes.SIGN_IN);
					}
					break;

				case 403:
					navigate(CommonRoutes.FORBIDDEN, { replace: true });
					break;
			}

			return response;
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

export default ecomapHttpClient;
