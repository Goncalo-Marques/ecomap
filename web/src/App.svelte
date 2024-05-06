<script lang="ts">
	import { onMount } from "svelte";
	import { Router, Route, navigate } from "svelte-routing";
	import SignIn from "./routes/signIn/SignIn.svelte";
	import Dashboard from "./routes/backOffice/dashboard/Dashboard.svelte";
	import Map from "./routes/backOffice/map/Map.svelte";
	import Routes from "./routes/backOffice/routes/Routes.svelte";
	import Trucks from "./routes/backOffice/trucks/Trucks.svelte";
	import Reports from "./routes/backOffice/reports/Reports.svelte";
	import EmployeesRouter from "./routes/backOffice/employees/EmployeesRouter.svelte";
	import Forbidden from "./routes/clientErrors/Forbidden.svelte";
	import NotFound from "./routes/clientErrors/NotFound.svelte";
	import BackOfficeLayout from "./routes/backOffice/components/BackOfficeLayout.svelte";
	import {
		BackOfficeRouterPaths,
		backOfficeBasename,
		CommonRoutes,
	} from "./routes/constants/routes";
	import url from "./lib/stores/url";
	import ContainersRouter from "./routes/backOffice/containers/ContainersRouter.svelte";
	import Toast from "./lib/components/Toast.svelte";
	import WarehousesRouter from "./routes/backOffice/warehouses/WarehousesRouter.svelte";

	onMount(() => {
		if ($url.pathname === "/") {
			navigate(CommonRoutes.SIGN_IN, { replace: true });
		}
	});
</script>

<Toast>
	<Router>
		<Route path={`/${backOfficeBasename}/*`}>
			<Router>
				<BackOfficeLayout>
					<Route
						path={BackOfficeRouterPaths.EMPLOYEES}
						component={EmployeesRouter}
					/>
					<Route path={BackOfficeRouterPaths.REPORTS} component={Reports} />
					<Route path={BackOfficeRouterPaths.TRUCKS} component={Trucks} />
					<Route path={`${BackOfficeRouterPaths.WAREHOUSES}/*`}>
						<WarehousesRouter />
					</Route>
					<Route path={`${BackOfficeRouterPaths.CONTAINERS}/*`}>
						<ContainersRouter />
					</Route>
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
</Toast>
