<script lang="ts">
	import { Link } from "svelte-routing";
	import Logo from "./Logo.svelte";
	import { BackOfficeRoutes } from "../../../constants/routes";
	import EmployeeMenu from "./EmployeeMenu.svelte";
</script>

<nav
	class="root flex h-[var(--top-bar-height)] justify-between justify-items-center border border-b-gray-300 bg-white px-7 py-4"
>
	<Link to={BackOfficeRoutes.DASHBOARD}>
		<div class="flex justify-center gap-1">
			<Logo />
			<h1 class="text-xl font-semibold">EcoMap</h1>
		</div>
	</Link>
	<EmployeeMenu />
</nav>

<style lang="postcss">
	:root {
		--top-bar-height: 4rem;
	}
</style>
