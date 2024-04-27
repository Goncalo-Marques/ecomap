<script lang="ts">
	import type { MouseEventHandler } from "svelte/elements";
	import type { IconSize } from "../../domain/components/icon";
	import Icon from "./Icon.svelte";

	/**
	 * TODO.
	 */
	export let icon: string;

	/**
	 * TODO.
	 */
	export let size: IconSize = "medium";

	/**
	 * TODO.
	 */
	export let onClick: MouseEventHandler<HTMLButtonElement> | null = null;
</script>

<button on:click={onClick}>
	<Icon name={icon} {size} />
</button>

<style>
	button {
		display: flex;
		justify-content: center;
		align-items: center;
		border-radius: 50%;
		padding: 0.25rem;

		&:hover {
			background-color: var(--gray-100);
		}
		&:active {
			background-color: var(--gray-200);
		}
	}
</style>
