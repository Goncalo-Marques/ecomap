<script lang="ts">
	import { Link } from "svelte-routing";
	import Button from "./Button.svelte";

	/**
	 * The title of the bottom sheet.
	 */
	export let title: string;

	/**
	 * The page to with the resource details.
	 */
	export let resourceLink: string | null = null;
</script>

<div class="layout">
	<div class="sheet">
		<div class="top-bar">
			<h2>{title}</h2>
			{#if resourceLink}
				<Link to={resourceLink} style="display: contents">
					<Button variant="tertiary" startIcon="open_in_new" />
				</Link>
			{/if}
		</div>
		<div class="details">
			<slot />
		</div>
	</div>
</div>

<style>
	.layout {
		position: absolute;
		left: 0;
		bottom: 0;
		width: 100%;
		display: flex;
		padding-inline: 2.5rem;
	}

	.sheet {
		flex: 1;
		display: flex;
		flex-direction: column;
		gap: 1rem;
		max-height: 20rem;
		border-top-left-radius: 0.25rem;
		border-top-right-radius: 0.25rem;
		background-color: var(--white);
		box-shadow: var(--shadow-lg);

		& h2 {
			font: var(--text-2xl-semibold);
		}
	}

	.top-bar {
		display: flex;
		padding: 2rem 2rem 0 2rem;

		& h2 {
			flex: 1;
		}
	}

	.details {
		--column-count: 4;
		--column-gap: 0.5rem;
		--column-min-width: 13rem;

		overflow: auto;
		padding: 0 2rem 2rem 2rem;
		display: grid;
		grid-template-columns: repeat(
			auto-fill,
			minmax(
				max(
					var(--column-min-width),
					(100% - var(--column-gap)) / var(--column-count)
				),
				1fr
			)
		);
		gap: var(--column-gap);
	}
</style>
