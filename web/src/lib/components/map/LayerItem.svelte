<script lang="ts">
	import { Layer } from "ol/layer";
	import Dot from "./Dot.svelte";
	import Switch from "../Switch.svelte";
	import { colorLayerKey, nameLayerKey } from "../../constants/map";

	/**
	 * Layer reference.
	 */
	export let layer: Layer;

	/**
	 * Toggles the layer visibility.
	 */
	function toggleVisibility() {
		layer.setVisible(!layer.getVisible());
	}
</script>

<div class="layer-item">
	{#if layer.get(colorLayerKey)}
		<Dot color={layer.get(colorLayerKey)} size="x-small" />
	{:else}
		<Dot color="white" size="x-small" />
	{/if}

	{#if layer.get(nameLayerKey)}
		<h2>{layer.get(nameLayerKey)}</h2>
	{/if}

	<Switch checked={layer.getVisible()} onClick={toggleVisibility} />
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
