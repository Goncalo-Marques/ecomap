<script lang="ts">
	import type { IconSize } from "$domain/components/icon";

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

<div role="progressbar" class={`text-green-700 ${className}`}>
	<Icon class="animate-spin" name="progress_activity" {size} />
</div>
