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

<TopBar />
<div class="flex h-[calc(100vh-var(--top-bar-height))]">
	<SideBar />
	<div class="grid flex-1 overflow-auto">
		<slot />
	</div>
</div>
