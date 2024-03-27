<script lang="ts">
	import type { HTMLButtonAttributes } from "svelte/elements";

	type Variant = "primary" | "secondary" | "tertiary";
	type Size = "medium" | "large";

	/**
	 * A space-separated list of the classes of the element.
	 * @default ""
	 */
	let className: string = "";
	export { className as class };

	/**
	 * Prevents the user from interacting with the button: it cannot be pressed or focused.
	 * @default false
	 */
	export let disabled: boolean = false;

	/**
	 * The size of the button.
	 * @default "medium"
	 */
	export let size: Size = "medium";

	/**
	 * The type of the button.
	 * @default "button"
	 */
	export let type: HTMLButtonAttributes["type"] = "button";

	/**
	 * The variant of the button.
	 * @default "primary"
	 */
	export let variant: Variant = "primary";
</script>

<button class={`${size} ${variant} ${className}`} {disabled} {type}>
	<slot />
</button>

<style>
	button {
		display: flex;
		justify-content: center;
		align-items: center;
		gap: 0.25rem;
		font-weight: 600;
		border-radius: 0.25rem;
		transition-property: var(--transition-colors);
		transition-duration: var(--transition-duration);
		transition-timing-function: var(--transition-timing-function);

		&:disabled {
			opacity: 0.6;
		}
	}

	.primary {
		color: var(--white);
		background-color: var(--green-700);

		&:hover:enabled {
			background-color: var(--green-800);
		}

		&:active:enabled {
			background-color: var(--green-900);
		}
	}
	.secondary {
		color: var(--green-700);
		background-color: transparent;
		border: 1px solid var(--green-700);

		&:hover:enabled {
			color: var(--green-800);
			border-color: var(--green-800);
			background-color: var(--green-50);
		}

		&:active:enabled {
			color: var(--green-800);
			border-color: var(--green-800);
			background-color: var(--green-100);
		}
	}
	.tertiary {
		color: var(--green-700);
		background-color: transparent;

		&:hover:enabled {
			color: var(--green-800);
			background-color: var(--green-50);
		}

		&:active:enabled {
			color: var(--green-800);
			background-color: var(--green-100);
		}
	}

	.medium {
		padding: 0.5rem 0.75rem;
	}
	.large {
		padding: 0.625rem 1rem;
	}
</style>
