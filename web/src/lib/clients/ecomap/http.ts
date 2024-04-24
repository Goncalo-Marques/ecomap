import { navigate } from "svelte-routing";
import createClient, { type Middleware } from "openapi-fetch";
import type { components, paths } from "../../../../api/ecomap/http";
import { getToken } from "../../utils/auth";
import { CommonRoutes } from "../../../routes/constants/routes";

const middleware: Middleware = {
	onRequest(request) {
		const token = getToken();

		// Set Authorization header.
		request.headers.set("Authorization", `Bearer ${token}`);

		return request;
	},
	async onResponse(response) {
		let body = await response.text();

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

		if (!response.ok && !body) {
			const error: components["schemas"]["Error"] = {
				code: "internal_server_error",
			};
			body = JSON.stringify(error);
		}

		return new Response(body, response);
	},
};

const ecomapHttpClient = createClient<paths>({
	baseUrl: "/api",
});
ecomapHttpClient.use(middleware);

export default ecomapHttpClient;
