<script lang="ts">
	/**
	 * The value of the option.
	 * @default ""
	 */
	export let value: string = "";
</script>

<option class="text-gray-900" {value}>
	<slot />
</option>
