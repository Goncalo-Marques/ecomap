/**
 * Basename for back office routes.
 */
export const backOfficeBasename = "back-office";

/**
 * Paths of back office router.
 */
export enum BackOfficeRouterPaths {
	DASHBOARD = "dashboard",
	MAP = "map",
	ROUTES = "routes",
	CONTAINERS = "containers",
	WAREHOUSES = "warehouses",
	TRUCKS = "trucks",
	LANDFILLS = "landfills",
	EMPLOYEES = "employees",
}

/**
 * Back office routes.
 */
export enum BackOfficeRoutes {
	DASHBOARD = `/${backOfficeBasename}/${BackOfficeRouterPaths.DASHBOARD}`,
	MAP = `/${backOfficeBasename}/${BackOfficeRouterPaths.MAP}`,
	ROUTES = `/${backOfficeBasename}/${BackOfficeRouterPaths.ROUTES}`,
	CONTAINERS = `/${backOfficeBasename}/${BackOfficeRouterPaths.CONTAINERS}`,
	WAREHOUSES = `/${backOfficeBasename}/${BackOfficeRouterPaths.WAREHOUSES}`,
	TRUCKS = `/${backOfficeBasename}/${BackOfficeRouterPaths.TRUCKS}`,
	LANDFILLS = `/${backOfficeBasename}/${BackOfficeRouterPaths.LANDFILLS}`,
	EMPLOYEES = `/${backOfficeBasename}/${BackOfficeRouterPaths.EMPLOYEES}`,
}

/**
 * Common routes.
 */
export enum CommonRoutes {
	SIGN_IN = "/signin",
	FORBIDDEN = "/forbidden",
}
