<script lang="ts">
	import Map from "ol/Map";
	import View from "ol/View";
	import TileLayer from "ol/layer/Tile";
	import XYZ from "ol/source/XYZ";
	import { transform } from "ol/proj";
    
    import {map} from '../utils/store'
	import { onMount } from "svelte";

    /**
     * 
     */
    export let zoom: number

    /**
     * 
     */
    export let lat: number

    /**
     * 
     */
    export let lon: number

    /**
     * 
     */
    export let map_width: string

    /**
     * 
     */
    export let map_height: string

    const map_id: string = "map_id"

    onMount(() => {
        $map.setTarget(map_id)
	});

    $map = new Map({
			layers: [
				new TileLayer({
					source: new XYZ({
						url: "https://tile.openstreetmap.org/{z}/{x}/{y}.png",
                        tileSize: 256,
                        crossOrigin: 'anonymous'
					}),
				}),
			],
			view: new View({
				center: transform([lon, lat], 'EPSG:4326', 'EPSG:3857'),
				zoom: zoom,
			}),
	});
</script>

<div id={map_id} style="--map_width: {map_width}; --map_height: {map_height}">
    <a style="position:absolute;left:0;bottom:0;z-index:999;" href="https://www.openstreetmap.org/copyright" target="_blank">© OpenStreetMap contributors</a>
</div>

<style>
    div {
        width: var(--map_width);
        height: var(--map_height);
        position: relative;
    }
    a {
        color: var(--gray-950);
        background-color: var(--white);
        margin: 0;
        padding: 1px .5em;
        font-size: 12px;
        border-radius: 0px 5px;
        opacity: 0.8;
    }
</style>
