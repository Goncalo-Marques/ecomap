<script lang="ts">
	type Size = "x-small" | "small" | "medium" | "large" | "x-large" | "xx-large";

	/**
	 * The icon name.
	 */
	export let name: string;

	/**
	 * The icon size.
	 * @default "medium"
	 */
	export let size: Size = "medium";
</script>

<span class={`material-symbols-rounded ${size}`}>{name}</span>

<style>
	span {
		-webkit-user-select: none;
		user-select: none;
	}
	.x-small {
		font-size: 0.875rem;
	}
	.small {
		font-size: 1rem;
	}
	.medium {
		font-size: 1.5rem;
	}
	.large {
		font-size: 2rem;
	}
	.x-large {
		font-size: 2.5rem;
	}
	.xx-large {
		font-size: 3rem;
	}
</style>
