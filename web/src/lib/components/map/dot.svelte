<script lang="ts">
	type Size = "x-small" | "small" | "medium" | "large";

	/**
	 * Dot height
	 * @default "12px"
	 */
	export let height: string = "12px";

	/**
	 * Dot width
	 * @default "12px"
	 */
	export let width: string = "12px";

	/**
	 * Dot color.
	 */
	export let color: string;
</script>

<span class="dot" style="--dot-color: {color}; --dot-height: {height}; --dot-width: {width};" />

<style>
	span {
		-webkit-user-select: none;
		user-select: none;
	}
	.dot {
		height: var(--dot-height);
		width: var(--dot-width);
		background-color: var(--dot-color);
		border-radius: 50%;
		display: inline-block;
	}
</style>
