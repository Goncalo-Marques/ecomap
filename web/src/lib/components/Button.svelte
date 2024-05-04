<script lang="ts">
	import type {
		HTMLButtonAttributes,
		MouseEventHandler,
	} from "svelte/elements";
	import Icon from "./Icon.svelte";

	type Variant = "primary" | "secondary" | "tertiary";
	type Size = "small" | "medium" | "large";

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
	 * Callback fired when button is clicked.
	 * @default null
	 */
	export let onClick: MouseEventHandler<HTMLButtonElement> | null = null;

	/**
	 * The size of the button.
	 * @default "medium"
	 */
	export let size: Size = "medium";

	/**
	 * The name of the icon positioned at the beginning of the button.
	 * @default null
	 */
	export let startIcon: string | null = null;

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

	/**
	 * Indicates if the button only contains an icon.
	 * When the button doesn't contain a slot, the styles for the button differ from the standard ones.
	 */
	let onlyIcon = !$$props.$$slots;
</script>

<button
	class={`${size} ${variant} ${className} ${onlyIcon ? "only-icon" : ""}`}
	{disabled}
	{type}
	on:click={onClick}
>
	{#if startIcon}
		<Icon name={startIcon} size={onlyIcon ? "medium" : "small"} />
	{/if}
	<slot />
</button>

<style>
	button {
		position: relative;
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
		background-color: var(--white);
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
		background-color: var(--white);

		&:hover:enabled {
			color: var(--green-800);
			background-color: var(--green-50);
		}

		&:active:enabled {
			color: var(--green-800);
			background-color: var(--green-100);
		}
	}

	.small {
		font-size: 0.75rem;
		padding: 0.375rem 0.5rem;
	}
	.medium {
		padding: 0.5rem 0.75rem;
	}
	.large {
		padding: 0.625rem 1rem;
	}

	.only-icon {
		padding: 0.625rem;
	}
</style>
