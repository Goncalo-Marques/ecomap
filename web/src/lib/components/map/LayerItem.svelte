<script lang="ts">
	import { Layer } from "ol/layer";
	import Dot from "./Dot.svelte";
	import Switch from "../Switch.svelte";
	import { colorLayerReference, nameLayerReference } from "./mapUtils";

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
	 * Changes layer visibility based on {@link visible} value.
	 */
	function toggleVisibility() {
		visible = !visible;
		layer.setVisible(visible);
	}
</script>

<div class="layer-item">
	{#if layer.get(colorLayerReference)}
		<Dot color={layer.get(colorLayerReference)} size="x-small" />
	{:else}
		<Dot color="white" size="x-small" />
	{/if}

	{#if layer.get(nameLayerReference)}
		<h2>{layer.get(nameLayerReference)}</h2>
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
