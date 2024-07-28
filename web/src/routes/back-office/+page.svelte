<script lang="ts">
	import { goto } from "$app/navigation";

	import { BackOfficeRoutes } from "../constants/routes";

	// Redirect to dashboard page if pathname is at back office root level.
	goto(BackOfficeRoutes.DASHBOARD);
</script>
