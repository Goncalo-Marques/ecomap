<script lang="ts">
	import Map from "../../lib/components/map/Map.svelte";
	import { onDestroy, onMount } from "svelte";
	import {
		map,
		addClusterLayer,
		addVectorLayer,
	} from "../../lib/components/map/mapStore";
	import type { Layer } from "ol/layer";

	let show = false
	let layers : Layer[]

	onMount(() => {
		addClusterLayer("http://localhost:8000/contentores.geojson");
		setTimeout(() => {
			addVectorLayer("http://localhost:8000/rede_viaria.geojson");
		}, 2000);
	});

	function handleMessage(event: any) {
		show = true
		layers = event.detail.layers
	}
	
</script>

<main>
	<Map on:message={handleMessage}/>
	{#if show}
		{#each layers as layer}
			<p>{layer}</p>
		{/each}
	{/if}
</main>

<style>
	main {
		width: 100%;
	}
</style>
