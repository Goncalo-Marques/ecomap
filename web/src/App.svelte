<script lang="ts">
	import { Router, Route, navigate } from "svelte-routing";
	import SignIn from "./routes/signIn/SignIn.svelte";
	import Dashboard from "./routes/backOffice/dashboard/Dashboard.svelte";
	import Map from "./routes/backOffice/map/Map.svelte";
	import Routes from "./routes/backOffice/routes/Routes.svelte";
	import Containers from "./routes/backOffice/containers/Containers.svelte";
	import Warehouses from "./routes/backOffice/warehouses/Warehouses.svelte";
	import Trucks from "./routes/backOffice/trucks/Trucks.svelte";
	import Reports from "./routes/backOffice/reports/Reports.svelte";
	import Employees from "./routes/backOffice/employees/Employees.svelte";
	import NotFound from "./routes/notFound/NotFound.svelte";
	import BackOfficeLayout from "./routes/backOffice/components/BackOfficeLayout.svelte";
	import {
		BackOfficeRouterPaths,
		backOfficeBasename,
		CommonRoutes,
		BackOfficeRoutes,
	} from "./routes/constants/routes";
	import PrivateRouter from "./routes/components/PrivateRouter.svelte";
	import { isAuthenticated } from "./lib/utils/auth";
	import url from "./lib/utils/url";

	// Redirect to back office dashboard page if user is authenticated and URL pathname is at the root level.
	if (isAuthenticated() && $url.pathname === "/") {
		navigate(BackOfficeRoutes.DASHBOARD);
	}
</script>

<Router>
	<Route path={`/${backOfficeBasename}/*`}>
		<PrivateRouter>
			<BackOfficeLayout>
				<Route path={BackOfficeRouterPaths.EMPLOYEES} component={Employees} />
				<Route path={BackOfficeRouterPaths.REPORTS} component={Reports} />
				<Route path={BackOfficeRouterPaths.TRUCKS} component={Trucks} />
				<Route path={BackOfficeRouterPaths.WAREHOUSES} component={Warehouses} />
				<Route path={BackOfficeRouterPaths.CONTAINERS} component={Containers} />
				<Route path={BackOfficeRouterPaths.ROUTES} component={Routes} />
				<Route path={BackOfficeRouterPaths.MAP} component={Map} />
				<Route path={BackOfficeRouterPaths.DASHBOARD} component={Dashboard} />
				<Route component={NotFound} />
			</BackOfficeLayout>
		</PrivateRouter>
	</Route>
	<Route path={CommonRoutes.SIGN_IN} component={SignIn} />
	<Route component={NotFound} />
</Router>
