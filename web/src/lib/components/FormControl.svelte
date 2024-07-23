<script lang="ts">
	import Icon from "./Icon.svelte";

	/**
	 * A space-separated list of the classes of the element.
	 * @default ""
	 */
	let className: string = "";
	export { className as class };

	/**
	 * Indicates if the input contains an error.
	 * @default false
	 */
	export let error: boolean = false;

	/**
	 * The text below the input that contains additional information about how the input will be used.
	 * @default null
	 */
	export let helperText: string | null = null;

	/**
	 * The label for the input.
	 * @default null
	 */
	export let label: string | null = null;

	/**
	 * The text representing advisory information related to the element it belongs to.
	 * @default null
	 */
	export let title: string | null = null;
</script>

<label class={`flex flex-col ${className}`}>
	{#if label}
		<span {title} class="flex items-center gap-1">
			<span class="text-gray-500">{label}</span>
			{#if title}
				<Icon name="info" size="xx-small" />
			{/if}
		</span>
	{/if}
	<slot />
	{#if helperText}
		<span data-error={error} class="text-xs data-[error=true]:text-red-500">
			{helperText}
		</span>
	{/if}
</label>
