<script lang="ts">
	/**
	 * A space-separated list of the classes of the element.
	 * @default ""
	 */
	let className: string = "";
	export { className as class };

	/**
	 * The label of the section.
	 */
	export let label: string;
</script>

<section class={`flex flex-col gap-2 ${className}`}>
	<header>
		<h2 class="font-semibold">{label}</h2>
	</header>
	<slot />
</section>
