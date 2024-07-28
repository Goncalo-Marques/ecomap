<script lang="ts">
	import { goto } from "$app/navigation";
	import { CommonRoutes } from "$lib/constants/routes";

	goto(CommonRoutes.SIGN_IN, { replaceState: true });
</script>
