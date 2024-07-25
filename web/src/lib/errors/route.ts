/**
 * Error that indicates that a route does not have any containers associated
 * with it.
 */
export class RouteWithNoContainersError extends Error {
	constructor() {
		super("Route does not have any containers associated");
	}
}
