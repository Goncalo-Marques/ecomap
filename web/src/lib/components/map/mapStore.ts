import Map from "ol/Map";
import { writable, get } from "svelte/store";
import View from "ol/View";
import SimpleGeometry from "ol/geom/SimpleGeometry";
import { type Coordinate } from "ol/coordinate";
import GeoJSON from "ol/format/GeoJSON";

import { Vector as VectorSource, Cluster, OSM } from "ol/source";
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

/**
 * Variables gotten from css vars.
 */
const cssVars = {
	text_sm_semibold: style.getPropertyValue("--text-sm-semibold"),
	indigo_400: style.getPropertyValue("--indigo-400"),
};

/**
 * Default style for vector layer.
 */
const defaultVectorStyle: VectorStyle = {
	"stroke-color": "#fff",
	"fill-color": "#3980a885",
};

/**
 * Default style for WebGl point layer.
 */
const defaultIconStyle: WebGLStyle = {
	"icon-src": "/images/logo.svg",
};

/**
 * Default style for cluster point layer.
 */
const defaultClusterIcon = new Style({
	image: new Icon({
		src: "/images/logo.svg",
	}),
});

/**
 * Default style for cluster symbol.
 */
const defaultClusterSymbol: Style = new Style({
	image: new Circle({
		radius: 20,
		stroke: new Stroke({
			color: cssVars.indigo_400,
			width: 2,
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

	createRenderer(): WebGLVectorLayerRenderer {
		return new WebGLVectorLayerRenderer(this, {
			style: this.style,
		});
	}
}

/**
 * Set map store as a new Map
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
) {
	const baseLayer = new TileLayer({
		source: new OSM(),
		visible: true,
		zIndex: 0,
	});

	baseLayer.set("layer-name", "baseLayer");

	map.set(
		new Map({
			layers: [baseLayer],
			view: new View({
				center: fromLonLat([lon, lat]),
				zoom: zoom,
				projection: projection,
				// Locks the map on the iberian peninsula
				extent: [
					-2159435.3010021457, 3990778.5878774817, 863857.4518866497,
					5984975.69547515,
				],
			}),
		}),
	);
}

/**
 * Add a vector layer to the map
 *
 * @param url receives geojson data
 * @param layerName name for layer
 * @param layerColor color that identifies the layer
 * @param style style for new layer
 */
export function addVectorLayer(
	url: string,
	layerName: string,
	layerColor: string,
	style: VectorStyle = defaultVectorStyle,
) {
	const mapValue = get(map);

	const vectorLayer = new WebGLLayer(
		{
			source: new VectorSource({
				url: url,
				format: new GeoJSON(),
			}),
			zIndex: mapValue?.getAllLayers().length,
		},
		style,
	);

	vectorLayer.set("layer-name", layerName);
	vectorLayer.set("layer-color", layerColor);
	mapValue?.addLayer(vectorLayer);
}

/**
 * Add a point vector layer into map
 *
 * @param url receives geojson data
 * @param layerName name for layer
 * @param layerColor color that identifies the layer
 * @param style style for new layer
 */
export function addPointLayer(
	url: string,
	layerName: string,
	layerColor: string,
	style: WebGLStyle = defaultIconStyle,
) {
	const mapValue = get(map);

	const pointsLayer = new WebGLPointsLayer({
		source: new VectorSource({
			url: url,
			format: new GeoJSON(),
		}),
		style: style,
		zIndex: mapValue?.getAllLayers().length,
	});

	pointsLayer.set("layer-name", layerName);
	pointsLayer.set("layer-color", layerColor);
	mapValue?.addLayer(pointsLayer);
}

/**
 * Add a clusterLayer into map
 *
 * @param url receives geojson data
 * @param layerName name for layer
 * @param layerColor color that identifies the layer
 * @param clusterStyle style for the cluster nodes
 * @param iconStyle style for each independent node
 */
export function addClusterLayer(
	url: string,
	layerName: string,
	layerColor: string,
	clusterStyle: Style = defaultClusterSymbol,
	iconStyle: Style = defaultClusterIcon,
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

	cluster.set("layer-name", layerName);
	cluster.set("layer-color", layerColor);

	mapValue?.addLayer(cluster);

	mapValue?.on("click", e => {
		cluster.getFeatures(e.pixel).then(clickedFeatures => {
			if (clickedFeatures.length) {
				const features: FeatureLike[] = clickedFeatures[0].get("features");
				if (features.length > 1) {
					const coordinates: Coordinate[] = [];

					for (const feature of features) {
						const geom = feature.getGeometry();
						if (!(geom instanceof SimpleGeometry)) {
							continue;
						}

						const coord = geom.getCoordinates();
						if (!coord) {
							continue;
						}

						coordinates.push(coord);
					}

					const extent: Extent = boundingExtent(coordinates);

					mapValue
						.getView()
						.fit(extent, { duration: 800, padding: [50, 50, 50, 50] });
				} else {
					//Code for single icon clicks
				}
			}
		});
	});

	let hoverFeature: FeatureLike;
	mapValue?.on("pointermove", e => {
		cluster.getFeatures(e.pixel).then(features => {
			if (features[0] !== hoverFeature) {
				hoverFeature = features[0];

				mapValue.getTargetElement().style.cursor = hoverFeature
					? "pointer"
					: "";

				cluster.getSource()?.changed();
			}
		});
	});
}
