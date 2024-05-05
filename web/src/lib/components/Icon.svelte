<script lang="ts">
	import type { IconSize } from "../../domain/components/icon";

	/**
	 * The icon name.
	 */
	export let name: string;

	/**
	 * The icon size.
	 * @default "medium"
	 */
	export let size: IconSize = "medium";
</script>

<span class={`material-symbols-rounded ${size}`}>{name}</span>

<style>
	span {
		-webkit-user-select: none;
		user-select: none;
	}
	.xx-small {
		line-height: 0.875rem;
		font-size: 0.875rem;
	}
	.x-small {
		line-height: 1rem;
		font-size: 1rem;
	}
	.small {
		line-height: 1.25rem;
		font-size: 1.25rem;
	}
	.medium {
		line-height: 1.5rem;
		font-size: 1.5rem;
	}
	.large {
		line-height: 2rem;
		font-size: 2rem;
	}
	.x-large {
		line-height: 2.5rem;
		font-size: 2.5rem;
	}
	.xx-large {
		line-height: 3rem;
		font-size: 3rem;
	}
</style>
