<script lang="ts">
	import type { IconSize } from "../../domain/components/icon";
	import Icon from "./Icon.svelte";

	/**
	 * A space-separated list of the classes of the element.
	 * @default ""
	 */
	let className: string = "";
	export { className as class };

	/**
	 * The spinner size.
	 * @default "xx-large"
	 */
	export let size: IconSize = "xx-large";
</script>

<div role="progressbar" class={`loading ${className}`}>
	<Icon name="progress_activity" {size} />
</div>

<style>
	.loading {
		color: var(--green-700);
	}
	.loading > :global(.material-symbols-rounded) {
		animation: var(--animation-spin);
	}
</style>
