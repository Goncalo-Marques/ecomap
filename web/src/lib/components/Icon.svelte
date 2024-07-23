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

	/**
	 * A space-separated list of the classes of the element.
	 * @default ""
	 */
	let className: string = "";
	export { className as class };
</script>

<span
	class={[
		"material-symbols-rounded select-none",
		"data-[size=xx-small]:text-sm data-[size=xx-small]:leading-[0.875rem]",
		"data-[size=x-small]:text-base data-[size=x-small]:leading-4",
		"data-[size=small]:text-xl data-[size=small]:leading-5",
		"data-[size=medium]:text-2xl data-[size=medium]:leading-6",
		"data-[size=large]:text-[2rem] data-[size=large]:leading-8",
		"data-[size=x-large]:text-[2.5rem] data-[size=x-large]:leading-10",
		"data-[size=xx-large]:text-5xl data-[size=xx-large]:leading-[3rem]",
		className,
	].join(" ")}
	data-size={size}>{name}</span
>
