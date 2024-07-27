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
	<dt class="text-gray-500">{label}</dt>
	<dd>{value ?? MISSING_VALUE_SYMBOL}</dd>
</div>
