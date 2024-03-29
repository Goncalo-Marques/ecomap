import Map from "ol/Map";
import { writable, get, derived } from "svelte/store";
import View from "ol/View";
import GeoJSON from "ol/format/GeoJSON";

import { Vector as VectorSource, XYZ, Cluster } from "ol/source";
import { WebGLTile as TileLayer, Layer } from "ol/layer";
import { fromLonLat } from "ol/proj";

import WebGLVectorLayerRenderer, {
	type Options,
} from "ol/renderer/webgl/VectorLayer.js";
import WebGLPointsLayer from "ol/layer/WebGLPoints.js";

import { boundingExtent } from "ol/extent";

import { Circle, Fill, Icon, Stroke, Style, Text } from "ol/style.js";
import { Vector as VectorLayer } from "ol/layer.js";

import type { FeatureLike } from "ol/Feature";
import type { Options as OptionsLayer } from "ol/layer/Layer";
import type { VectorStyle } from "ol/render/webgl/VectorStyleRenderer";
import type { WebGLStyle } from "ol/style/webgl";
import { createEventDispatcher } from "svelte";

export const map = writable<Map | null>(null);


const defaultVectorStyle: VectorStyle = {
	"stroke-color": "#e609d7",
	"fill-color": "#f0b10585",
}

const defaultIconStyle: WebGLStyle = {
	'icon-src': '/images/logo.svg' 
}

/**
 * Custom vector layer
 */
class WebGLLayer extends Layer {
	private style: VectorStyle;

	constructor(options: OptionsLayer, style: VectorStyle) {
		super(options);
		this.style = style;
	}

	createRenderer(): any {
		return new WebGLVectorLayerRenderer(this, {
			style: this.style,
		});
	}
}

/**
 *
 * @param size Number os icons inside each Cluster
 * @returns Icon or Circle Style according to size number
 */
function createClusterIcon(size: number) {
	const icon = new Style({
		image: new Icon({
			src: "/images/logo.svg",
		}),
	});

	const circle: Style = new Style({
		image: new Circle({
			radius: 20,
			stroke: new Stroke({
				color: "#fff",
			}),
			fill: new Fill({
				color: "#68b083",
			}),
		}),
		text: new Text({
			text: size.toString(),
			fill: new Fill({
				color: "#fff",
			}),
		}),
	});

	if (size >= 2) {
		return circle;
	}

	return icon;
}

/**
 * Set's map store as a new Map
 *
 * @param lon Center longitude
 * @param lat Center latitude
 * @param zoom Default zoom
 * @param projection Projection used, ex: EPSG:3857
 */
export function createMap(
	lon: number,
	lat: number,
	zoom: number,
	projection: string = "EPSG:3857",
	layerName: string = "baseLayer",
) {
	const baseLayer = new TileLayer({
		source: new XYZ({
			url: "https://tile.openstreetmap.org/{z}/{x}/{y}.png",
			tileSize: 256,
		}),
		visible: true,
	});

	baseLayer.setProperties({ "layer-name": layerName });

	map.set(
		new Map({
			layers: [baseLayer],
			view: new View({
				center: fromLonLat([lon, lat]),
				zoom: zoom,
				projection: projection,
				extent: [
					-2159435.3010021457, 3990778.5878774817, 863857.4518866497,
					5984975.69547515,
				],
			}),
		}),
	);
}

/**
 * Add's vector layer into map
 *
 * @param url receives geojson data
 */
export function addVectorLayer(url: string, style: VectorStyle = defaultVectorStyle) {
	const mapValue = get(map);

	const vectorLayer = new WebGLLayer(
		{
			source: new VectorSource({
				url: url,
				format: new GeoJSON(),
			}),
		},
		style,
	);

	mapValue?.addLayer(vectorLayer);
}

/**
 * Add's point's vector layer into map
 *
 * @param url receives geojson data
 */
export function addPointLayer(url: string, style: WebGLStyle = defaultIconStyle) {
	const mapValue = get(map);

	const pointsLayer = new WebGLPointsLayer({
		source: new VectorSource({
			url: url,
			format: new GeoJSON(),
		}),
		style: style
	});

	mapValue?.addLayer(pointsLayer);
}

/**
 * Add's clusterLayer into map
 *
 * @param url receives geojson data
 */
export function addClusterLayer(url: string) {
	const mapValue = get(map);

	const cluster = new VectorLayer({
		source: new Cluster({
			distance: 50,
			minDistance: 10,
			source: new VectorSource({
				url: url,
				format: new GeoJSON(),
			}),
		}),
		style: (feature: FeatureLike) => {
			const size = feature.get("features").length;
			let style = createClusterIcon(size);
			return style;
		},
	});

	mapValue?.addLayer(cluster);

	mapValue?.on("click", e => {
		cluster.getFeatures(e.pixel).then(clickedFeatures => {
			if (clickedFeatures.length) {
				const features = clickedFeatures[0].get("features");
				if (features.length > 1) {
					const extent = boundingExtent(
						features.map((r: any) => r.getGeometry().getCoordinates()),
					);
					mapValue
						.getView()
						.fit(extent, { duration: 800, padding: [50, 50, 50, 50] });
				} else {
					console.log("APENAS 1: ", features);
				}
			}
		});
	});
}
