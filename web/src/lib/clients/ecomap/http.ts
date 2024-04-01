import { navigate } from "svelte-routing";
import createClient from "openapi-fetch";
import type { components, paths } from "../../../../api/ecomap/http";
import { getToken } from "../../utils/auth";
import { CommonRoutes } from "../../../routes/constants/routes";

const ecomapHttpClient = createClient<paths>({
	baseUrl: "/api",
	headers: {
		Authorization: `Bearer ${getToken()}`,
	},
	// Custom fetch implementation to handle native fetch exceptions.
	async fetch(input, init) {
		try {
			const response = await fetch(input, init);

			switch (response.status) {
				case 401:
					navigate(CommonRoutes.SIGN_IN);
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