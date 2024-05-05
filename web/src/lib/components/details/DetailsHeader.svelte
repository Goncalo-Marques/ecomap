<script lang="ts">
	import { Link } from "svelte-routing";
	import Icon from "../Icon.svelte";

	/**
	 * The title displayed in the header.
	 */
	export let title: string;

	/**
	 * The route to navigate to when the back button is clicked.
	 */
	export let to: string;
</script>

<header>
	<Link {to} class="back">
		<Icon name="arrow_back" size="medium" />
	</Link>
	<h1>{title}</h1>
	<div class="actions">
		<slot />
	</div>
</header>

<style>
	header {
		display: flex;
		align-items: center;
		flex-wrap: wrap;
		gap: 1rem;
	}

	header > :global(.back) {
		display: contents;
		color: var(--green-700);
	}

	header > h1 {
		flex: 1;
		font: var(--text-2xl-semibold);
	}

	.actions {
		display: flex;
		gap: 1rem;
		align-items: center;
	}
</style>
