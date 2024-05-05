import Map from "ol/Map";
import View from "ol/View";
import SimpleGeometry from "ol/geom/SimpleGeometry";
import { type Coordinate } from "ol/coordinate";

import { Vector as VectorSource, Cluster, OSM } from "ol/source";
import { WebGLTile as TileLayer, Layer } from "ol/layer";
import { fromLonLat } from "ol/proj";

import WebGLVectorLayerRenderer from "ol/renderer/webgl/VectorLayer";
import WebGLPointsLayer from "ol/layer/WebGLPoints";

import { boundingExtent, type Extent } from "ol/extent";

import { Circle, Fill, Icon, Stroke, Style, Text } from "ol/style";
import { Vector as VectorLayer } from "ol/layer";

import type { FeatureLike } from "ol/Feature";
import type { Options as OptionsLayer } from "ol/layer/Layer";
import type { VectorStyle } from "ol/render/webgl/VectorStyleRenderer";
import type { WebGLStyle } from "ol/style/webgl";
import {
	mapLayerName,
	colorLayerKey,
	nameLayerKey,
	DEFAULT_MAX_ZOOM,
} from "../../constants/map";
import type { Geometry } from "ol/geom";
import type Feature from "ol/Feature";
import type { Options } from "ol/source/Vector";
import Rotate from "ol/control/Rotate";

const docElement = document.documentElement;
const style = getComputedStyle(docElement);

/**
 * Variables retrieved from css vars.
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
	"fill-color": "#3980a895",
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
		font: cssVars.text_sm_semibold,
		textAlign: "center",
		fill: new Fill({
			color: "#000",
		}),
	}),
});

/**
 * WebGl Vector layer for Open Layers.
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

export class MapHelper {
	public constructor(private readonly map: Map) {}

	/**
	 * Add a vector layer to the map.
	 *
	 * @param sourceOptions Options of the vector layer source.
	 * @param layerName Name for layer.
	 * @param [layerColor="#15803D"] Color that identifies the layer.
	 * @param [style=defaultIconStyle] Style for new layer.
	 */
	public addVectorLayer(
		sourceOptions: Options<Feature<Geometry>>,
		layerName: string,
		layerColor: string = "#15803D",
		style: VectorStyle = defaultVectorStyle,
	) {
		const vectorLayer = new WebGLLayer(
			{
				source: new VectorSource(sourceOptions),
				zIndex: this.map.getAllLayers().length,
			},
			style,
		);

		vectorLayer.set(nameLayerKey, layerName);
		vectorLayer.set(colorLayerKey, layerColor);
		this.map.addLayer(vectorLayer);
	}

	/**
	 * Add a point vector layer into map.
	 *
	 * @param sourceOptions Options of the point layer source.
	 * @param layerName Name for layer.
	 * @param [layerColor="#15803D"] Color that identifies the layer.
	 * @param [style=defaultIconStyle] Style for new layer.
	 */
	public addPointLayer(
		sourceOptions: Options<Feature<Geometry>>,
		layerName: string,
		layerColor: string = "#15803D",
		style: WebGLStyle = defaultIconStyle,
	) {
		const pointsLayer = new WebGLPointsLayer({
			source: new VectorSource(sourceOptions),
			style: style,
			zIndex: this.map.getAllLayers().length,
		});

		pointsLayer.set(nameLayerKey, layerName);
		pointsLayer.set(colorLayerKey, layerColor);
		this.map.addLayer(pointsLayer);
	}

	/**
	 * Add a clusterLayer into map.
	 *
	 * @param sourceOptions Options of the cluster layer source.
	 * @param layerName Name for layer.
	 * @param [layerColor="#15803D"] Color that identifies the layer.
	 * @param [clusterStyle=defaultClusterSymbol] Style for the cluster nodes.
	 * @param [iconStyle=defaultClusterIcon] Style for each independent node.
	 */
	public addClusterLayer(
		sourceOptions: Options<Feature<Geometry>>,
		layerName: string,
		layerColor: string = "#15803D",
		clusterStyle: Style = defaultClusterSymbol,
		iconStyle: Style = defaultClusterIcon,
	) {
		const cluster = new VectorLayer({
			zIndex: this.map.getAllLayers().length,
			source: new Cluster({
				distance: 50,
				minDistance: 10,
				source: new VectorSource(sourceOptions),
			}),
			style: (feature: FeatureLike) => {
				const size = feature.get("features").length;

				clusterStyle.getText()?.setText(size.toString());

				return size >= 2 ? clusterStyle : iconStyle;
			},
		});

		cluster.set(nameLayerKey, layerName);
		cluster.set(colorLayerKey, layerColor);

		this.map.addLayer(cluster);

		this.map.on("click", e => {
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

						this.map
							.getView()
							.fit(extent, { duration: 800, padding: [50, 50, 50, 50] });
					}
				}
			});
		});

		let hoverFeature: FeatureLike;
		this.map.on("pointermove", e => {
			cluster.getFeatures(e.pixel).then(features => {
				if (features[0] !== hoverFeature) {
					hoverFeature = features[0];

					this.map.getTargetElement().style.cursor = hoverFeature
						? "pointer"
						: "";

					cluster.getSource()?.changed();
				}
			});
		});
	}
}

/**
 * Creates a new Open Layers map.
 *
 * @param lon Center longitude.
 * @param lat Center latitude.
 * @param zoom Default zoom.
 * @param projection Projection used, ex: EPSG:3857.
 * @returns Map.
 */
export function createMap(
	lon: number,
	lat: number,
	zoom: number,
	projection: string = "EPSG:3857",
): Map {
	const baseLayer = new TileLayer({
		source: new OSM(),
		visible: true,
		zIndex: 0,
	});

	baseLayer.set(nameLayerKey, mapLayerName);

	return new Map({
		controls: [new Rotate()],
		layers: [baseLayer],
		view: new View({
			center: fromLonLat([lon, lat]),
			zoom: zoom,
			maxZoom: DEFAULT_MAX_ZOOM,
			projection: projection,
			// Locks the map on the iberian peninsula
			extent: [
				-2159435.3010021457, 3990778.5878774817, 863857.4518866497,
				5984975.69547515,
			],
		}),
	});
}
