<script lang="ts" generics="TValue">
	/**
	 * The label of the radio button.
	 * @default null
	 */
	export let label: string | null = null;

	/**
	 * Indicates that the radio button is the default selected one in the group.
	 * @default false
	 */
	export let checked: boolean = false;

	/**
	 * The ID of the radio button.
	 * @default null
	 */
	export let id: string | null = null;

	/**
	 * The name of the radio button.
	 * @default null
	 */
	export let name: string | null = null;

	/**
	 * The value of the radio button.
	 * @default null
	 */
	export let value: TValue | null = null;

	/**
	 * The value of the selected radio button in a group.
	 * @default null
	 */
	export let group: TValue | null = null;
</script>

<label>
	<input type="radio" {id} {name} {value} {checked} bind:group />
	{label}
</label>

<style>
	label {
		display: flex;
		gap: 0.5rem;
		font: var(--text-base-regular);
	}
</style>
