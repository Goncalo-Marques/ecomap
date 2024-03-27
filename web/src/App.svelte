<script lang="ts">
	import { Router, Route, navigate } from "svelte-routing";
	import SignIn from "./routes/signIn/SignIn.svelte";
	import Dashboard from "./routes/dashboard/Dashboard.svelte";
	import Map from "./routes/map/Map.svelte";
	import Routes from "./routes/routes/Routes.svelte";
	import Containers from "./routes/containers/Containers.svelte";
	import Warehouses from "./routes/warehouses/Warehouses.svelte";
	import Trucks from "./routes/trucks/Trucks.svelte";
	import Reports from "./routes/reports/Reports.svelte";
	import Employees from "./routes/employees/Employees.svelte";
	import NotFound from "./routes/notFound/NotFound.svelte";
	import Layout from "./routes/components/Layout.svelte";
	import {
		AppRouterPaths,
		AppRoutes,
		appBasename,
	} from "./routes/constants/appRoutes";
	import url from "./lib/utils/url";

	// Redirect to dashboard page if pathname is at root level.
	if (
		$url.pathname === `/${appBasename}` ||
		$url.pathname === `/${appBasename}/`
	) {
		navigate(AppRoutes.DASHBOARD);
	}
</script>

<Router>
	<Route path={AppRoutes.SIGN_IN} component={SignIn} />
	<Route path={`/${appBasename}/*`}>
		<Router>
			<Layout>
				<Route path={AppRouterPaths.EMPLOYEES} component={Employees} />
				<Route path={AppRouterPaths.REPORTS} component={Reports} />
				<Route path={AppRouterPaths.TRUCKS} component={Trucks} />
				<Route path={AppRouterPaths.WAREHOUSES} component={Warehouses} />
				<Route path={AppRouterPaths.CONTAINERS} component={Containers} />
				<Route path={AppRouterPaths.ROUTES} component={Routes} />
				<Route path={AppRouterPaths.MAP} component={Map} />
				<Route path={AppRouterPaths.DASHBOARD} component={Dashboard} />
				<Route component={NotFound} />
			</Layout>
		</Router>
	</Route>
	<Route component={NotFound} />
</Router>
