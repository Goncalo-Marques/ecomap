<script lang="ts">
	import { Layer } from "ol/layer";
	import Dot from "./dot.svelte";
	import Switch from "../switch.svelte";

	/**
	 * Object layer reference to this layerItem 
	 */
	export let layer: Layer;

	/**
	 * Switch value
	 */
	let visible: boolean = true;

	/**
	 * Changes layer visibility based on @visible value
	 */
	function toggleVisibility() {
		visible = !visible;
		layer.setVisible(visible);
	}
</script>

<div class="layer-item">
	{#if layer.get("layer-color")}
		<Dot color={layer.get("layer-color")} />
	{:else}
		<Dot color={"white"} />
	{/if}

	{#if layer.get("layer-name")}
		<h2>{layer.get("layer-name")}</h2>
	{:else}
		<h2>#UNDEFINED</h2>
	{/if}

	<Switch checked={visible} onClick={toggleVisibility} />
</div>

<style>
	* {
		box-sizing: border-box;
	}

	h2 {
		font: var(--text-sm-regular);
	}

	.layer-item {
		display: flex;
		gap: 8px;
		height: 24px;
		align-items: center;
	}
</style>
