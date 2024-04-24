<script lang="ts">
	import { navigate } from "svelte-routing";
	import url from "../../../lib/stores/url";
	import { BackOfficeRoutes, backOfficeBasename } from "../../constants/routes";
	import SideBar from "./sidebar/SideBar.svelte";
	import TopBar from "./topbar/TopBar.svelte";

	// Redirect to dashboard page if pathname is at back office root level.
	if (
		$url.pathname === `/${backOfficeBasename}` ||
		$url.pathname === `/${backOfficeBasename}/`
	) {
		navigate(BackOfficeRoutes.DASHBOARD);
	}
</script>

<div class="main">
	<TopBar />
	<div class="container">
		<SideBar />
		<slot />
	</div>
</div>

<style>
	.main {
		min-height: 100vh;
		display: flex;
		flex-direction: column;
	}

	.container {
		flex: 1;
		display: flex;
	}
</style>
