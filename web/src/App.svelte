<script lang="ts">
	import { onMount } from "svelte";
	import { Router, Route, navigate } from "svelte-routing";
	import SignIn from "./routes/signIn/SignIn.svelte";
	import Dashboard from "./routes/backOffice/dashboard/Dashboard.svelte";
	import Map from "./routes/backOffice/map/Map.svelte";
	import RoutesRouter from "./routes/backOffice/routes/RoutesRouter.svelte";
	import LandfillsRouter from "./routes/backOffice/landfills/LandfillsRouter.svelte";
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
	import TrucksRouter from "./routes/backOffice/trucks/TrucksRouter.svelte";

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
					<Route path={`${BackOfficeRouterPaths.LANDFILLS}/*`}>
						<LandfillsRouter />
					</Route>
					<Route path={`${BackOfficeRouterPaths.EMPLOYEES}/*`}>
						<EmployeesRouter />
					</Route>
					<Route path={`${BackOfficeRouterPaths.TRUCKS}/*`}>
						<TrucksRouter />
					</Route>
					<Route path={`${BackOfficeRouterPaths.WAREHOUSES}/*`}>
						<WarehousesRouter />
					</Route>
					<Route path={`${BackOfficeRouterPaths.CONTAINERS}/*`}>
						<ContainersRouter />
					</Route>
					<Route path={`${BackOfficeRouterPaths.ROUTES}/*`}>
						<RoutesRouter />
					</Route>
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
