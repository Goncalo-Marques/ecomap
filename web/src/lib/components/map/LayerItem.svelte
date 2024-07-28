<script lang="ts">
	import { Layer } from "ol/layer";

	import { colorLayerKey, nameLayerKey } from "../../constants/map";
	import Switch from "../Switch.svelte";
	import Dot from "./Dot.svelte";

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

<div class="flex h-6 items-center gap-2">
	{#if layer.get(colorLayerKey)}
		<Dot color={layer.get(colorLayerKey)} size="x-small" />
	{:else}
		<Dot color="white" size="x-small" />
	{/if}

	{#if layer.get(nameLayerKey)}
		<h2 class="text-sm">{layer.get(nameLayerKey)}</h2>
	{/if}

	<Switch checked={layer.getVisible()} onClick={toggleVisibility} />
</div>
