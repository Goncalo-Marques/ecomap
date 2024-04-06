import Map from "ol/Map";
import { writable, get, derived } from "svelte/store";
import View from "ol/View";
import GeoJSON from "ol/format/GeoJSON";

import { Vector as VectorSource, XYZ, Cluster, OSM } from "ol/source";
import { WebGLTile as TileLayer, Layer } from "ol/layer";
import { fromLonLat } from "ol/proj";

import WebGLVectorLayerRenderer from "ol/renderer/webgl/VectorLayer.js";
import WebGLPointsLayer from "ol/layer/WebGLPoints.js";

import { boundingExtent, type Extent } from "ol/extent";

import { Circle, Fill, Icon, Stroke, Style, Text } from "ol/style.js";
import { Vector as VectorLayer } from "ol/layer.js";

import type { FeatureLike } from "ol/Feature";
import type { Options as OptionsLayer } from "ol/layer/Layer";
import type { VectorStyle } from "ol/render/webgl/VectorStyleRenderer";
import type { WebGLStyle } from "ol/style/webgl";

export const map = writable<Map | null>(null);

const docElement = document.documentElement;
const style = getComputedStyle(docElement);

let cssVars = {
	text_sm_semibold: style.getPropertyValue('--text-sm-semibold'),
	indigo_400 : style.getPropertyValue('--indigo-400')
}

const defaultVectorStyle: VectorStyle = {
	"stroke-color": "#fff",
	"fill-color": "#3980a885",
};

const defaultIconStyle: WebGLStyle = {
	"icon-src": "/images/logo.svg",
};

const defaultClusterIcon = new Style({
	image: new Icon({
		src: "/images/logo.svg",
	}),
});

const defaultClusterSymbol: Style = new Style({
	image: new Circle({
		radius: 20,
		stroke: new Stroke({
			color: cssVars.indigo_400,
			width: 2
		}),
		fill: new Fill({
			color: "#fff",
		}),
	}),
	text: new Text({
		text: "",
		font: cssVars.text_sm_semibold,
		textAlign: "center",
		fill: new Fill({
			color: "#000",
		}),
	}),
});

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
		source: new OSM(),
		visible: true,
		zIndex: 0
	});

	baseLayer.set("layer-name", layerName )

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
export function addVectorLayer(
	url: string,
	layerName: string,
	style: VectorStyle = defaultVectorStyle,
) {
	const mapValue = get(map);

	const vectorLayer = new WebGLLayer(
		{
			source: new VectorSource({
				url: url,
				format: new GeoJSON(),
			}),
			zIndex: mapValue?.getAllLayers().length
		},
		style,
	);

	vectorLayer.set("layer-name", layerName )
	mapValue?.addLayer(vectorLayer);
}

/**
 * Add's point's vector layer into map
 *
 * @param url receives geojson data
 */
export function addPointLayer(
	url: string,
	layerName: string,
	style: WebGLStyle = defaultIconStyle,
) {
	const mapValue = get(map);

	const pointsLayer = new WebGLPointsLayer({
		source: new VectorSource({
			url: url,
			format: new GeoJSON(),
		}),
		style: style,
		zIndex: mapValue?.getAllLayers().length
	});

	pointsLayer.set("layer-name", layerName )
	mapValue?.addLayer(pointsLayer);
}

/**
 * Add's clusterLayer into map
 *
 * @param url receives geojson data
 */
export function addClusterLayer(
	url: string,
	layerName: string,
	clusterStyle: Style = defaultClusterSymbol,
	iconStyle: Style = defaultClusterIcon
) {
	const mapValue = get(map);

	const cluster = new VectorLayer({
		zIndex: mapValue?.getAllLayers().length,
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

			clusterStyle.getText()?.setText(size.toString());


			return size >= 2 ? clusterStyle : iconStyle;
		},
	});

	cluster.set("layer-name", layerName )

	mapValue?.addLayer(cluster);

	mapValue?.on("click", e => {
		cluster.getFeatures(e.pixel).then(clickedFeatures => {
			if (clickedFeatures.length) {
				const features: FeatureLike[] = clickedFeatures[0].get("features");
				if (features.length > 1) {
					console.log("fetatures: ", features);
					
					const extent: Extent = boundingExtent(features.map((r: any)=> r.getGeometry().getCoordinates()),
					);
					mapValue.getView().fit(extent, { duration: 800, padding: [50, 50, 50, 50] });
				} else {
					console.log("APENAS 1: ", features);
				}
			}
		});
	});

	let hoverFeature: FeatureLike;
	mapValue?.on("pointermove", e => {
		cluster.getFeatures(e.pixel).then(features => {
			if (features[0] !== hoverFeature) {
				hoverFeature = features[0];

				mapValue.getTargetElement().style.cursor = hoverFeature ? "pointer": "";

				cluster.getSource()?.changed();
			}
		});
	});
}
