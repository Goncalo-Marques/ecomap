<script lang="ts">
	import Map from "ol/Map";
	import View from "ol/View";
	import TileLayer from "ol/layer/Tile";
	import XYZ from "ol/source/XYZ";
	import { transform } from "ol/proj";

    /**
     * 
     */
    export let map: Map

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

	function createMap(node:any) {
		map = new Map({
			target: map_id,
			layers: [
				new TileLayer({
					source: new XYZ({
						url: "https://tile.openstreetmap.org/{z}/{x}/{y}.png",
					}),
				}),
			],
			view: new View({
				center: transform([lon, lat], 'EPSG:4326', 'EPSG:3857'),
				zoom: zoom,
			}),
		});
	};
</script>

<div id={map_id} use:createMap style="--map_width: {map_width}; --map_height: {map_height}"/>

<style>
    div {
        width: var(--map_width);
        height: var(--map_height);
        border: 2px solid var(--blue-500);
    }
</style>
