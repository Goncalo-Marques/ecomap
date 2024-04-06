<script lang="ts">
	import type { Layer } from "ol/layer";
	import Icon from "../Icon.svelte";
	import { map } from "./mapStore";
	import { onMount } from "svelte";

	export let layer: Layer;

	let nodeRef: any;

    let visible : boolean = true

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

	function zIndexUp() {
		let index: number | undefined = layer.getZIndex();

		if (index) {
			layer.setZIndex(index + 1);
			refreshMap();
		}
	}
	
    function setVisibility() {
        if (visible) {
            visible = false
            layer.setVisible(false)
        }else{
            visible = true
            layer.setVisible(true)
        }
    }
</script>

<div class="layer-item" bind:this={nodeRef}>
	<h4>{layer.get("layer-name")}</h4>

	<div class="all-buttons">
        <button class="buttons material-symbols-rounded" on:click={setVisibility}>
            {#if visible}
                <Icon name="visibility" />
            {:else}
                <Icon name="visibility_off" />
            {/if}
        </button>
	</div>
</div>

<style>
	* {
		box-sizing: border-box;
	}

	.layer-item {
		display: flex;
		gap: 1em;
		align-items: center;
		padding: 0.5rem 0.75rem;
	}

	button {
		height: 24px;
	}

	.all-buttons {
		margin-left: auto;
		display: flex;
		gap: 5px;
		flex-direction: row;
	}

	.buttons {
		background-color: var(--gray-100);
		border-radius: 4px;
	}

    .buttons:hover{
        cursor: pointer;
    }
</style>
