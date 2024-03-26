<script lang="ts">
	import VectorSource from "ol/source/Vector";
	import { map } from "../utils/store";
	import VectorLayer from "ol/layer/Vector";
	import GeoJSON from "ol/format/GeoJSON";
	import Style from "ol/style/Style";
	import Stroke from "ol/style/Stroke";
	import Fill from "ol/style/Fill";
    
	function addLayer() {
		const vectorLayer = new VectorLayer({
			source: new VectorSource({
				url: "http://localhost:8000/concelhos.geojson",
				format: new GeoJSON(),
			}),
			style: new Style({
				stroke: new Stroke({
					color: "#000000",
				}),
				fill: new Fill({
					color: "#f7000050",
				}),
			}),
		});

		$map.addLayer(vectorLayer);
	}
</script>
