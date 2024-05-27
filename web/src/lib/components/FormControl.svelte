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

<label class={className}>
	{#if label}
		<span {title} class="header">
			<span class="input-label">{label}</span>
			{#if title}
				<Icon name="info" size="xx-small" />
			{/if}
		</span>
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
	}

	.header {
		display: flex;
		align-items: center;
		gap: 0.25rem;

		& .input-label {
			color: var(--gray-500);
		}
	}

	.helper-text {
		font: var(--text-xs-regular);
		color: var(--gray-400);

		&.error {
			color: var(--red-500);
		}
	}
</style>
