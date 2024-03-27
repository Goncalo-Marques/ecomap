import Map from "ol/Map";
import { writable, get } from "svelte/store";

import View from "ol/View";
import GeoJSON from "ol/format/GeoJSON";

import { Vector as VectorSource, XYZ } from "ol/source";
import { WebGLTile as TileLayer, Layer } from "ol/layer";
import { fromLonLat } from "ol/proj";
import LayerGroup from "ol/layer/Group";

import WebGLVectorLayerRenderer from "ol/renderer/webgl/VectorLayer.js";
import WebGLPointsLayer from "ol/layer/WebGLPoints.js";

export const map = writable<Map>();

class WebGLLayer extends Layer {
	// eslint-disable-next-line @typescript-eslint/no-explicit-any
	createRenderer(): any {
		return new WebGLVectorLayerRenderer(this, {
			style: {
				"stroke-color": "#000000",
				"fill-color": "#f7000050",
			},
		});
	}
}

const osmStandard = new TileLayer({
	source: new XYZ({
		url: "https://tile.openstreetmap.org/{z}/{x}/{y}.png",
		tileSize: 256,
		crossOrigin: "anonymous",
	}),
	visible: true,
});

const osmHumanitarian = new TileLayer({
	source: new XYZ({
		url: "https://tile.openstreetmap.fr/hot/{z}/{x}/{y}.png",
		tileSize: 256,
		crossOrigin: "anonymous",
	}),
	visible: false,
});

const layerGroup = new LayerGroup({
	layers: [osmStandard, osmHumanitarian],
});

export function createMap(
	lon: number,
	lat: number,
	zoom: number,
	projection: string = "EPSG:3857",
) {
	map.set(
		new Map({
			view: new View({
				center: fromLonLat([lon, lat]),
				zoom: zoom,
				projection: projection,
				extent: [
					-1354248.9461922427, 4274625.428689052, 523429.8051994869,
					5593519.232428095,
				],
			}),
		}),
	);

	get(map).addLayer(layerGroup);
}

/**
 * Add's vector layer to map with geojson
 *
 * @param url
 */
export function addVectorLayer(url: string) {
	const mapValue = get(map);

	const vectorLayer = new WebGLLayer({
		source: new VectorSource({
			url: url,
			format: new GeoJSON(),
		}),
	});

	mapValue.addLayer(vectorLayer);
}

/**
 * Add's point's vector layer to map with geojson
 *
 * @param url
 */
export function addPointLayer(url: string) {
	const mapValue = get(map);

	const pointsLayer = new WebGLPointsLayer({
		source: new VectorSource({
			url: url,
			format: new GeoJSON(),
		}),
		style: {
			"icon-src": "/images/logo.svg",
		},
	});

	mapValue.addLayer(pointsLayer);
}
