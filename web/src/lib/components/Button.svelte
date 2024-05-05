<script lang="ts">
	import type {
		HTMLButtonAttributes,
		MouseEventHandler,
	} from "svelte/elements";
	import Icon from "./Icon.svelte";
	import type { IconSize } from "../../domain/components/icon";

	type Variant = "primary" | "secondary" | "tertiary";
	type Size = "small" | "medium" | "large";
	type ActionType = "default" | "danger";

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
	 * The type of action that will be performed by the button.
	 * @default "default"
	 */
	export let actionType: ActionType = "default";

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

	/**
	 * Map of icon sizes for each button size.
	 */
	const iconSizes: Record<Size, IconSize> = {
		small: "x-small",
		medium: "small",
		large: "medium",
	};
</script>

<button
	class={`${size} ${variant} ${actionType} ${className} ${onlyIcon ? "only-icon" : ""}`}
	{disabled}
	{type}
	on:click={onClick}
>
	{#if startIcon}
		<Icon name={startIcon} size={iconSizes[size]} />
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
		border-radius: 0.25rem;
		transition-property: var(--transition-colors);
		transition-duration: var(--transition-duration);
		transition-timing-function: var(--transition-timing-function);

		&:disabled {
			opacity: 0.6;
		}
	}

	.primary {
		&.default {
			color: var(--white);
			background-color: var(--green-700);

			&:hover:enabled {
				background-color: var(--green-800);
			}

			&:active:enabled {
				background-color: var(--green-900);
			}
		}

		&.danger {
			color: var(--white);
			background-color: var(--red-700);

			&:hover:enabled {
				background-color: var(--red-800);
			}

			&:active:enabled {
				background-color: var(--red-900);
			}
		}
	}
	.secondary {
		&.default {
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

		&.danger {
			color: var(--red-700);
			background-color: var(--white);
			border: 1px solid var(--red-700);

			&:hover:enabled {
				color: var(--red-800);
				border-color: var(--red-800);
				background-color: var(--red-50);
			}

			&:active:enabled {
				color: var(--red-800);
				border-color: var(--red-800);
				background-color: var(--red-100);
			}
		}
	}
	.tertiary {
		&.default {
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

		&.danger {
			color: var(--red-700);
			background-color: var(--white);

			&:hover:enabled {
				color: var(--red-800);
				background-color: var(--red-50);
			}

			&:active:enabled {
				color: var(--red-800);
				background-color: var(--red-100);
			}
		}
	}

	.small {
		font: var(--text-xs-semibold);
		padding: 0.375rem 0.5rem;

		&.only-icon {
			padding: 0.375rem;
		}
	}
	.medium {
		font: var(--text-sm-semibold);
		padding: 0.5rem 0.75rem;

		&.only-icon {
			padding: 0.5rem;
		}
	}
	.large {
		font: var(--text-base-semibold);
		padding: 0.625rem 1rem;

		&.only-icon {
			padding: 0.625rem;
		}
	}
</style>
