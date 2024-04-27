<script lang="ts">
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
</script>

<label class={className}>
	{#if label}
		{label}
	{/if}
	<slot />
	{#if helperText}
		<span class={`helper-text ${error ? "error" : ""}`}>{helperText}</span>
	{/if}
</label>

<style>
	label {
		display: flex;
		flex-direction: column;
		color: var(--gray-500);
	}

	.helper-text {
		font: var(--text-xs-regular);
		color: var(--gray-400);

		&.error {
			color: var(--red-500);
		}
	}
</style>
