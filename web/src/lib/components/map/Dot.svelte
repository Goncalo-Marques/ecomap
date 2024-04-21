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

<span class={`dot ${size}`} style="--dot-color: {color}" />

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

	.x-small {
		height: 0.75rem;
		width: 0.75rem;
	}
	.small {
		height: 1rem;
		width: 1rem;
	}
	.medium {
		height: 1.5rem;
		width: 1.5rem;
	}
	.large {
		height: 2rem;
		width: 2rem;
	}
	.x-large {
		height: 2.5rem;
		width: 2.5rem;
	}
	.xx-large {
		height: 3rem;
		width: 3rem;
	}
</style>
