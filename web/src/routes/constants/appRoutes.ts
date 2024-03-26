/**
 * Basename for application routes.
 */
export const appBasename = "app";

/**
 * Paths of /app router.
 */
export enum AppRouterPaths {
	DASHBOARD = "dashboard",
	MAP = "map",
	ROUTES = "routes",
	CONTAINERS = "containers",
	WAREHOUSES = "warehouses",
	TRUCKS = "trucks",
	REPORTS = "reports",
	EMPLOYEES = "employees",
}

/**
 * Application routes.
 */
export enum AppRoutes {
	DASHBOARD = `/${appBasename}/${AppRouterPaths.DASHBOARD}`,
	MAP = `/${appBasename}/${AppRouterPaths.MAP}`,
	ROUTES = `/${appBasename}/${AppRouterPaths.ROUTES}`,
	CONTAINERS = `/${appBasename}/${AppRouterPaths.CONTAINERS}`,
	WAREHOUSES = `/${appBasename}/${AppRouterPaths.WAREHOUSES}`,
	TRUCKS = `/${appBasename}/${AppRouterPaths.TRUCKS}`,
	REPORTS = `/${appBasename}/${AppRouterPaths.REPORTS}`,
	EMPLOYEES = `/${appBasename}/${AppRouterPaths.EMPLOYEES}`,
	SIGN_IN = "/signin",
}
