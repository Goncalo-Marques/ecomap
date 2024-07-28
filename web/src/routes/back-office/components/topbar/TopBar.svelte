<script lang="ts">
	import { BackOfficeRoutes } from "../../../constants/routes";
	import EmployeeMenu from "./EmployeeMenu.svelte";
	import Logo from "./Logo.svelte";
</script>

<nav
	class="root flex h-[var(--top-bar-height)] justify-between justify-items-center border border-b-gray-300 bg-white px-7 py-4"
>
	<a href={BackOfficeRoutes.DASHBOARD}>
		<div class="flex justify-center gap-1">
			<Logo />
			<h1 class="text-xl font-semibold">EcoMap</h1>
		</div>
	</a>
	<EmployeeMenu />
</nav>

<style lang="postcss">
	:root {
		--top-bar-height: 4rem;
	}
</style>
