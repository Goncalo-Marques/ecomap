<script lang="ts">
	/**
	 * The value of the option.
	 * @default ""
	 */
	export let value: string = "";
</script>

<option {value}>
	<slot />
</option>
