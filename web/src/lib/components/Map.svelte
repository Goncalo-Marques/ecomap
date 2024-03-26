<script lang="ts">
	import Map from "ol/Map";
	import View from "ol/View";
	import TileLayer from "ol/layer/Tile";
	import XYZ from "ol/source/XYZ";
	import { transform } from "ol/proj";
    
    import {map} from '../utils/store'
	import { onMount } from "svelte";

    /**
     * Zoom value for map view  
     * @default 5
     */
    export let zoom: number = 5

    /**
     * Center latitude of map
     * @default 50
     */
    export let lat: number = 50

    /**
     * Center longitude of map
     * @default 20
     */
    export let lon: number = 20

    /**
     * Map Viewport width size
     * @default 100vw
     */
    export let map_width: string = '100vw'

    /**
     * Map Viewport height size
     * @default 100vh
     */
    export let map_height: string = '100vh'

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
    <a href="https://www.openstreetmap.org/copyright" target="_blank" rel="noopener noreferrer">© OpenStreetMap contributors</a>
</div>

<style>
    div {
        width: var(--map_width);
        height: var(--map_height);
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
