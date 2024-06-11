import Map from "ol/Map";
import View from "ol/View";

import { Vector as VectorSource, Cluster, OSM } from "ol/source";
import { WebGLTile as TileLayer, Layer } from "ol/layer";
import { fromLonLat } from "ol/proj";

import WebGLVectorLayerRenderer from "ol/renderer/webgl/VectorLayer";
import WebGLPointsLayer from "ol/layer/WebGLPoints";

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
	DEFAULT_PIN_ICON_SRC,
	DEFAULT_MIN_ZOOM,
	OL_PROJECTION,
} from "../../constants/map";
import type { Geometry } from "ol/geom";
import type Feature from "ol/Feature";
import type { Options } from "ol/source/Vector";
import Rotate from "ol/control/Rotate";
import { getCssVariable } from "../../utils/cssVars";
import type {
	CreateMapOptions,
	MapHelperClusterLayerOptions,
} from "../../../domain/components/map";

/**
 * Default style for vector layer.
 */
const defaultVectorStyle: VectorStyle = {
	"stroke-color": getCssVariable("--white"),
	"fill-color": getCssVariable("--cyan-500"),
};

/**
 * Default style for WebGl point layer.
 */
const defaultIconStyle: WebGLStyle = {
	"icon-src": DEFAULT_PIN_ICON_SRC,
};

/**
 * Style for cluster symbol.
 */
const clusterStyle = new Style({
	text: new Text({
		font: getCssVariable("--text-sm-semibold"),
		textAlign: "center",
		fill: new Fill({
			color: getCssVariable("--black"),
		}),
	}),
});

/**
 * Style for cluster circle.
 */
const clusterCircle = new Circle({
	radius: 20,
	fill: new Fill({
		color: getCssVariable("--white"),
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
	 * Add a cluster layer into map.
	 *
	 * @param features Layer features.
	 * @param options Layer options.
	 */
	public addClusterLayer(
		features: Feature<Geometry>[],
		options?: MapHelperClusterLayerOptions,
	) {
		const iconStyle = new Style({
			image: new Icon({
				src: options?.iconSrc ?? DEFAULT_PIN_ICON_SRC,
			}),
		});
		const selectedIconStyle = new Style({
			image: new Icon({
				src: options?.selectedIconSrc ?? DEFAULT_PIN_ICON_SRC,
			}),
		});

		const cluster = new VectorLayer({
			source: new Cluster({
				distance: 50,
				minDistance: 10,
				source: new VectorSource({
					features,
				}),
			}),
			style(feature: FeatureLike) {
				const features: Feature[] = feature.get("features");
				const size = features.length;

				if (options?.clusterBorderColor) {
					clusterCircle.setStroke(
						new Stroke({
							color: options.clusterBorderColor,
							width: 3,
						}),
					);
					clusterStyle.setImage(clusterCircle);
				} else {
					clusterCircle.setStroke(
						new Stroke({
							color: getCssVariable("--gray-400"),
							width: 3,
						}),
					);
					clusterStyle.setImage(clusterCircle);
				}

				if (size >= 2) {
					clusterStyle.getText()?.setText(size.toString());
					return clusterStyle;
				}

				if (features[0].get("selected")) {
					return selectedIconStyle;
				}

				return iconStyle;
			},
		});

		if (options?.layerName) {
			cluster.set(nameLayerKey, options.layerName);
		}

		if (options?.layerColor) {
			cluster.set(colorLayerKey, options.layerColor);
		}

		this.map.addLayer(cluster);

		let hoverFeature: FeatureLike;
		this.map.on("pointermove", e => {
			cluster.getFeatures(e.pixel).then(features => {
				if (cluster.isVisible() && features[0] !== hoverFeature) {
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
export function createMap(options: CreateMapOptions): Map {
	const baseLayer = new TileLayer({
		source: new OSM(),
		visible: true,
		zIndex: 0,
	});

	baseLayer.set(nameLayerKey, mapLayerName);

	const {
		lon,
		lat,
		zoom,
		maxZoom = DEFAULT_MAX_ZOOM,
		minZoom = DEFAULT_MIN_ZOOM,
		projection = OL_PROJECTION,
	} = options;

	return new Map({
		controls: [new Rotate()],
		layers: [baseLayer],
		view: new View({
			center: fromLonLat([lon, lat]),
			zoom,
			maxZoom,
			minZoom,
			projection,
			// Locks the map on the iberian peninsula
			extent: [
				-2159435.3010021457, 3990778.5878774817, 863857.4518866497,
				5984975.69547515,
			],
		}),
	});
}
