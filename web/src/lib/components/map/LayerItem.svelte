<script lang="ts">
	import { Layer } from "ol/layer";
	import Dot from "./Dot.svelte";
	import Switch from "../Switch.svelte";

	/**
	 * Layer reference.
	 */
	export let layer: Layer;

	/**
	 * Indicates if layers is visible.
	 *
	 * @default true
	 */
	let visible: boolean = true;

	/**
	 * Property for layer color.
	 *
	 */
	const colorReference = "layer-color";

	/**
	 * Property for layer name.
	 */
	const nameReference = "layer-name";

	/**
	 * Changes layer visibility based on {@link visible} value.
	 */
	function toggleVisibility() {
		visible = !visible;
		layer.setVisible(visible);
	}
</script>

<div class="layer-item">
	{#if layer.get(colorReference)}
		<Dot color={layer.get(colorReference)} size="x-small" />
	{:else}
		<Dot color={"white"} size="x-small" />
	{/if}

	{#if layer.get(nameReference)}
		<h2>{layer.get(nameReference)}</h2>
	{/if}

	<Switch checked={visible} onClick={toggleVisibility} />
</div>

<style>
	h2 {
		font: var(--text-sm-regular);
	}

	.layer-item {
		display: flex;
		gap: 0.5rem;
		height: 1.5rem;
		align-items: center;
	}
</style>
