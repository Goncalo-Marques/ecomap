<script lang="ts">
	import { Link } from "svelte-routing";
	import Logo from "./Logo.svelte";
	import { BackOfficeRoutes } from "../../../constants/routes";
	import EmployeeMenu from "./EmployeeMenu.svelte";
</script>

<nav>
	<Link to={BackOfficeRoutes.DASHBOARD}>
		<div class="logo-container">
			<Logo />
			<h1>EcoMap</h1>
		</div>
	</Link>
	<EmployeeMenu />
</nav>

<style>
	:root {
		--top-bar-height: 4rem;
	}

	nav {
		height: var(--top-bar-height);
		display: flex;
		justify-content: space-between;
		justify-items: center;
		padding: 1rem 1.75rem;
		background-color: var(--white);
		border-bottom: 1px solid var(--gray-300);
	}

	.logo-container {
		display: flex;
		justify-items: center;
		gap: 0.25rem;
	}

	.logo-container h1 {
		font: var(--text-xl-semibold);
	}
</style>
