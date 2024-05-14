<script lang="ts">
	import { Link } from "svelte-routing";
	import Icon from "../../../lib/components/Icon.svelte";
	import Button from "../../../lib/components/Button.svelte";

	/**
	 * The icon to be displayed.
	 */
	export let icon: string;

	/**
	 * The title of the card.
	 */
	export let title: string;

	/**
	 * The resource link.
	 */
	export let resourceLink: string;
</script>

<article>
	<div class="top-bar">
		<div class="title">
			<Icon name={icon} />
			<h3>{title}</h3>
		</div>
		<Link to={resourceLink} style="display:contents">
			<Button startIcon="open_in_new" variant="tertiary" />
		</Link>
	</div>
	<slot />
</article>

<style>
	article {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		border-radius: 0.25rem;
		padding: 1rem;
		border: 1px solid var(--gray-300);
	}

	.top-bar {
		display: flex;
		align-items: center;
		gap: 0.5rem;
	}

	.title {
		flex: 1;
		display: flex;
		align-items: center;
		gap: 0.25rem;

		& h3 {
			flex: 1;
			font: var(--text-base-regular);
		}
	}
</style>
