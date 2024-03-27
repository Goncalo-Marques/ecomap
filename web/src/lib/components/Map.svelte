<script lang="ts">    
    import {map, createMap} from '../utils/mapStore'
	import { onMount } from "svelte";

    /**
     * Zoom value for map view  
     * @default 7
     */
    export let zoom: number = 7

    /**
     * Center latitude of map
     * @default 40
     */
    export let lat: number = 40

    /**
     * Center longitude of map
     * @default -4
     */
    export let lon: number = -4

    const map_id: string = "map_id"

    // Create Map
    createMap(lon, lat, zoom)

    // Mount map into div
    onMount(() => {
        $map.setTarget(map_id)
	});
</script>

<div id={map_id}>
    <a href="https://www.openstreetmap.org/copyright" target="_blank" rel="noopener noreferrer">© OpenStreetMap contributors</a>
</div>

<style>
    div {
        width:  100%;
        height: 100%;
        position: relative;
    }
    a {
        position:absolute;
        color: var(--gray-950);
        background-color: var(--white);
        font-size: 12px;
        padding: 1px .5em;
        border-radius: 0px 5px;
        opacity: 0.8;
        left:0;
        bottom:0;
        z-index:999;
    }
</style>
