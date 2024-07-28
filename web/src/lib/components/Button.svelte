<script lang="ts">
	import type {
		HTMLButtonAttributes,
		MouseEventHandler,
	} from "svelte/elements";

	import type { IconSize } from "$domain/components/icon";

	import Icon from "./Icon.svelte";

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
	{disabled}
	{type}
	class={[
		"relative flex items-center justify-center gap-1 rounded font-semibold transition-colors disabled:opacity-60",

		// Sizes.
		"data-[size=small]:data-[icon=true]:p-[0.375rem] data-[size=small]:px-2 data-[size=small]:py-[0.375rem] data-[size=small]:text-xs",
		"data-[size=medium]:data-[icon=true]:p-2 data-[size=medium]:px-3 data-[size=medium]:py-2 data-[size=medium]:text-sm",
		"data-[size=large]:data-[icon=true]:p-[0.625rem] data-[size=large]:px-4 data-[size=large]:py-[0.625rem] data-[size=large]:text-base",

		// Primary default variant.
		"data-[variant=primary]:data-[actiontype=default]:bg-green-700 data-[variant=primary]:data-[actiontype=default]:text-white",
		"enabled:hover:data-[variant=primary]:data-[actiontype=default]:bg-green-800",
		"enabled:active:data-[variant=primary]:data-[actiontype=default]:bg-green-900",

		// Primary danger variant.
		"data-[variant=primary]:data-[actiontype=danger]:bg-red-700 data-[variant=primary]:data-[actiontype=danger]:text-white",
		"enabled:hover:data-[variant=primary]:data-[actiontype=danger]:bg-red-800",
		"enabled:active:data-[variant=primary]:data-[actiontype=danger]:bg-red-900",

		// Secondary default variant.
		"data-[variant=secondary]:data-[actiontype=default]:border data-[variant=secondary]:data-[actiontype=default]:border-green-700 data-[variant=secondary]:data-[actiontype=default]:bg-white data-[variant=secondary]:data-[actiontype=default]:text-green-700",
		"enabled:hover:data-[variant=secondary]:data-[actiontype=default]:border-green-800 enabled:hover:data-[variant=secondary]:data-[actiontype=default]:bg-green-50 enabled:hover:data-[variant=secondary]:data-[actiontype=default]:text-green-800",
		"enabled:active:data-[variant=secondary]:data-[actiontype=default]:border-green-800 enabled:active:data-[variant=secondary]:data-[actiontype=default]:bg-green-100 enabled:active:data-[variant=secondary]:data-[actiontype=default]:text-green-800",

		// Secondary danger variant.
		"data-[variant=secondary]:data-[actiontype=danger]:border data-[variant=secondary]:data-[actiontype=danger]:border-red-700 data-[variant=secondary]:data-[actiontype=danger]:bg-white data-[variant=secondary]:data-[actiontype=danger]:text-red-700",
		"enabled:hover:data-[variant=secondary]:data-[actiontype=danger]:border-red-800 enabled:hover:data-[variant=secondary]:data-[actiontype=danger]:bg-red-50 enabled:hover:data-[variant=secondary]:data-[actiontype=danger]:text-red-800",
		"enabled:active:data-[variant=secondary]:data-[actiontype=danger]:border-red-800 enabled:active:data-[variant=secondary]:data-[actiontype=danger]:bg-red-100 enabled:active:data-[variant=secondary]:data-[actiontype=danger]:text-red-800",

		// Tertiary default variant.
		"data-[variant=tertiary]:data-[actiontype=default]:bg-white data-[variant=tertiary]:data-[actiontype=default]:text-green-700",
		"enabled:hover:data-[variant=tertiary]:data-[actiontype=default]:bg-green-50 enabled:hover:data-[variant=tertiary]:data-[actiontype=default]:text-green-800",
		"enabled:active:data-[variant=tertiary]:data-[actiontype=default]:bg-green-100 enabled:active:data-[variant=tertiary]:data-[actiontype=default]:text-green-800",

		// Tertiary danger variant.
		"data-[variant=tertiary]:data-[actiontype=danger]:bg-white data-[variant=tertiary]:data-[actiontype=danger]:text-red-700",
		"enabled:hover:data-[variant=tertiary]:data-[actiontype=danger]:bg-red-50 enabled:hover:data-[variant=tertiary]:data-[actiontype=danger]:text-red-800",
		"enabled:active:data-[variant=tertiary]:data-[actiontype=danger]:bg-red-100 enabled:active:data-[variant=tertiary]:data-[actiontype=danger]:text-red-800",

		className,
	].join(" ")}
	data-actiontype={actionType}
	data-variant={variant}
	data-size={size}
	data-icon={onlyIcon}
	on:click={onClick}
>
	{#if startIcon}
		<Icon name={startIcon} size={iconSizes[size]} />
	{/if}
	<slot />
</button>
