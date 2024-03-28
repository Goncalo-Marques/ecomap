import Map from "ol/Map";
import { writable, get } from "svelte/store";

import View from "ol/View";
import GeoJSON from "ol/format/GeoJSON";

import { Vector as VectorSource, XYZ } from "ol/source";
import { WebGLTile as TileLayer, Layer } from "ol/layer";
import { fromLonLat } from "ol/proj";

import WebGLVectorLayerRenderer from "ol/renderer/webgl/VectorLayer.js";
import WebGLPointsLayer from "ol/layer/WebGLPoints.js";

export const map = writable<Map|null>(null);

const styles: any = {
    'route-1': {
        'stroke-color': ['*', ['get', 'COLOR'], [220, 220, 220]],
        'stroke-width': 2,
        'stroke-offset': -1,
        'fill-color': ['*', ['get', 'COLOR'], [255, 255, 255, 0.6]],
    },
    'route-2': {
        'stroke-color': '#e609d7',
        'fill-color': '#f0b10585',
    }
};

class WebGLLayer extends Layer {
	private id: string;

    constructor(options: any) {
        super(options);
        this.id = options.id;
    }

    createRenderer(): any {
        return new WebGLVectorLayerRenderer(this, {
            style: styles[this.id]
        });
    }
}

export function createMap(
	lon: number,
	lat: number,
	zoom: number,
	projection: string = "EPSG:3857",
) {
	map.set(
		new Map({
			layers: [
				new TileLayer({
					source: new XYZ({
						url: "https://tile.openstreetmap.org/{z}/{x}/{y}.png",
						tileSize: 256,
						crossOrigin: "anonymous",
					}),
					visible: true,
				})
			],
			view: new View({
				center: fromLonLat([lon, lat]),
				zoom: zoom,
				projection: projection,
			}),
		}),
	);
}

/**
 * Add's vector layer to map with geojson
 *
 * @param url
 */
export function addVectorLayer(url: string) {
	const mapValue = get(map);

	const vectorLayer = new WebGLLayer({
		id: 'route-2',
		source: new VectorSource({
			url: url,
			format: new GeoJSON(),
		}),
	});

	mapValue?.addLayer(vectorLayer);
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

	mapValue?.addLayer(pointsLayer);
}
