<script lang="ts">
	type Size = "x-small" | "small" | "medium" | "large" | "x-large" | "xx-large";

	/**
	 * Dot's size.
	 * @default "medium"
	 */
	export let size: Size = "medium";

	/**
	 * Dot's color.
	 */
	export let color: string;
</script>

<span
	data-size={size}
	class="inline-block select-none rounded-full data-[size=large]:size-8 data-[size=medium]:size-6 data-[size=small]:size-4 data-[size=x-large]:size-10 data-[size=x-small]:size-3 data-[size=xx-large]:size-12"
	style:background-color={color}
/>
