<script lang="ts">
	import { MISSING_VALUE_SYMBOL } from "../constants/symbols";

	/**
	 * The label of the field.
	 */
	export let label: string;

	/**
	 * The value of the field.
	 * @default null
	 */
	export let value: number | string | null = null;
</script>

<div>
	<dt>{label}</dt>
	<dd>{value ?? MISSING_VALUE_SYMBOL}</dd>
</div>

<style>
	dt {
		color: var(--gray-500);
	}
</style>
