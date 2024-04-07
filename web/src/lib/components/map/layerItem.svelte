<script lang="ts">
	import { Layer } from "ol/layer";
	import { map } from "./mapStore";
	import Dot from "./dot.svelte";

	export let layer: Layer;

	let nodeRef: any;

	let visible: boolean = true;

	function refreshMap() {
		$map?.getAllLayers().forEach(obj => {
			obj.getSource()?.changed();
		});
	}

	function remove() {
		nodeRef.parentNode.removeChild(nodeRef);
		console.log("Remove: ", layer);
		$map?.removeLayer(layer);
		refreshMap();
	}

	function toggleVisibility() {
		if (visible) {
			visible = false;
			layer.setVisible(false);
		} else {
			visible = true;
			layer.setVisible(true);
		}
	}
</script>

<div class="layer-item" bind:this={nodeRef}>
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

	<label class="switch">
		<input type="checkbox" bind:checked={visible} on:click={toggleVisibility}/>
		<span class="slider round"></span>
	</label>
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
	
	.switch { 
		box-sizing: border-box;
		margin-left: auto;
		
		position: relative;
		display: inline-block;
		width: 40px;
		height: 20px;
	}

	.switch input {
		opacity: 0;
		width: 0;
		height: 0;
	}

	.slider {
		position: absolute;
		cursor: pointer;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background-color: var(--white);
		-webkit-transition: 0.2s;
		transition: 0.2s;
	}

	.slider:before {
		position: absolute;
		content: "";
		height: 12px;
		width: 12px;
		bottom: 4px;

		-webkit-transition: 0.2s;
		transition: 0.2s;
	}

	input:checked + .slider {
		outline: 0;
		background-color: var(--green-100);
	}

	input:focus + .slider {
		box-shadow: 0 0 1px var(--green-100);
	}

	input:checked + .slider:before {
		-webkit-transform: translateX(22px);
		-ms-transform: translateX(22px);
		transform: translateX(22px);
		background-color: var(--green-800);

	}

	input + .slider::before{
		transform: translateX(4px);
		background-color: var(--gray-400);
	}

	input + .slider {
		outline: 1px solid var(--gray-400);
	}

	/* Rounded sliders */
	.slider.round {
		border-radius: 12px;
	}

	.slider.round:before {
		border-radius: 50%;
	}
</style>
