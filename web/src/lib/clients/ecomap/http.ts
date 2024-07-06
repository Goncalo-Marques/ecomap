import { navigate } from "svelte-routing";
import createClient, { type Middleware } from "openapi-fetch";
import type { components, paths } from "../../../../api/ecomap/http";
import { clearToken, getToken } from "../../utils/auth";
import { CommonRoutes } from "../../../routes/constants/routes";

/**
 * Endpoints that should be ignored if the status code from the server response is 401.
 * Used to prevent the employee token from being deleted and the employee from being redirected to the sign in page.
 */
const UNAUTHORIZED_IGNORED_ENDPOINTS: (keyof paths)[] = ["/employees/password"];

const middleware: Middleware = {
	onRequest(request) {
		const token = getToken();

		// Set Authorization header.
		request.headers.set("Authorization", `Bearer ${token}`);

		return request;
	},
	async onResponse(response, options) {
		let body = (await response.text()) || null;

		switch (response.status) {
			case 401: {
				const responseUrl = new URL(response.url);

				// Check if it's an endpoint that should be ignored.
				const isIgnoredEndpoint = UNAUTHORIZED_IGNORED_ENDPOINTS.some(
					endpoint => {
						const pathname = `${options.baseUrl}${endpoint}`;
						return pathname === responseUrl.pathname;
					},
				);
				if (isIgnoredEndpoint) {
					break;
				}

				clearToken();
				// Only redirect if the page is not the sign in page.
				if (location.pathname !== CommonRoutes.SIGN_IN) {
					navigate(CommonRoutes.SIGN_IN);
				}
				break;
			}

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
	signal: AbortSignal.timeout(1000 * 60), // 1 minute timeout.
});

ecomapHttpClient.use(middleware);

export default ecomapHttpClient;
