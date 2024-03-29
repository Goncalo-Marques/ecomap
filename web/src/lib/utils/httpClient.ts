import createClient from "openapi-fetch";
import type { paths } from "../../../api/ecomap/http";

const httpClient = createClient<paths>({
	baseUrl: "/api",
});

export default httpClient;
