<script lang="ts">
	import { onMount } from "svelte";
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
	import Forbidden from "./routes/clientErrors/Forbidden.svelte";
	import NotFound from "./routes/clientErrors/NotFound.svelte";
	import BackOfficeLayout from "./routes/backOffice/components/BackOfficeLayout.svelte";
	import {
		BackOfficeRouterPaths,
		backOfficeBasename,
		CommonRoutes,
	} from "./routes/constants/routes";

	onMount(() => {
		navigate(CommonRoutes.SIGN_IN, { replace: true });
	});
</script>

<Router>
	<Route path={`/${backOfficeBasename}/*`}>
		<Router>
			<BackOfficeLayout>
				<Route path={BackOfficeRouterPaths.EMPLOYEES} component={Employees} />
				<Route path={BackOfficeRouterPaths.REPORTS} component={Reports} />
				<Route path={BackOfficeRouterPaths.TRUCKS} component={Trucks} />
				<Route path={BackOfficeRouterPaths.WAREHOUSES} component={Warehouses} />
				<Route path={BackOfficeRouterPaths.CONTAINERS} component={Containers} />
				<Route path={BackOfficeRouterPaths.ROUTES} component={Routes} />
				<Route path={BackOfficeRouterPaths.MAP} component={Map} />
				<Route path={BackOfficeRouterPaths.DASHBOARD} component={Dashboard} />
				<Route path={BackOfficeRouterPaths.DASHBOARD} component={NotFound} />
				<Route component={NotFound} />
			</BackOfficeLayout>
		</Router>
	</Route>
	<Route path={CommonRoutes.SIGN_IN} component={SignIn} />
	<Route path={CommonRoutes.FORBIDDEN} component={Forbidden} />
	<Route component={NotFound} />
</Router>
