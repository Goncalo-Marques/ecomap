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
		font-size: 0.875rem;
	}
	.x-small {
		font-size: 1rem;
	}
	.small {
		font-size: 1.25rem;
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
