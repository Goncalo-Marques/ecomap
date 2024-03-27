import Map from "ol/Map";
import { writable, get } from "svelte/store";

import View from "ol/View";
import GeoJSON from "ol/format/GeoJSON";

import {Vector as VectorSource, XYZ} from "ol/source";
import {Vector as VectorLayer, Tile as TileLayer} from "ol/layer";
import { fromLonLat } from "ol/proj";
import LayerGroup from "ol/layer/Group";

export const map = writable<Map>()

const osmStandard = new TileLayer({
    source: new XYZ({
        url: "https://tile.openstreetmap.org/{z}/{x}/{y}.png",
        tileSize: 256,
        crossOrigin: 'anonymous'
    }),
    visible: true
})

const osmHumanitarian = new TileLayer({
    source: new XYZ({
        url: "https://tile.openstreetmap.fr/hot/{z}/{x}/{y}.png",
        tileSize: 256,
        crossOrigin: 'anonymous',
    }),
    visible: false,
})

const layerGroup = new LayerGroup({
    layers: [
        osmStandard, osmHumanitarian
    ]
})

export function createMap(lon: number, lat: number, zoom: number, projection :string = 'EPSG:3857') {
    map.set(
        new Map({
            view: new View({
                center: fromLonLat([lon, lat]),
                zoom: zoom,   
                projection: projection,
                // [min x, min y, max x, max y]
                extent: [-1354248.9461922427, 4274625.428689052, 523429.8051994869, 5593519.232428095],
            }),
        })
    )

    get(map).addLayer(layerGroup)
}

export function mapEvents() {
    const mapValue = get(map)
    
    // Example value
    mapValue.on('click',(e) => {
        mapValue.forEachFeatureAtPixel(e.pixel, (feature) => {
            console.log(feature.get('concelho'));
        })
    })
}

export function setStandard() {
    layerGroup.getLayersArray()[0].setVisible(true)
    layerGroup.getLayersArray()[1].setVisible(false)
}

export function setHumanitarian() {
    layerGroup.getLayersArray()[0].setVisible(false)
    layerGroup.getLayersArray()[1].setVisible(true)
}

/**
 * Add's vector layer to map with geojson
 * 
 * @param url 
 */
export function addVectorLayer(url :string) {
    const mapValue = get(map)

    const vectorLayer = new VectorLayer({
        source: new VectorSource({
            url: url,
            format: new GeoJSON(),
        })
    });

    mapValue.addLayer(vectorLayer);
}

/**
 * Add's Raster XYZ Tile layer to map
 * 
 * @param url 
 */
export function addRasterLayerXYZ(url :string) {
    const mapValue = get(map)

    const layer = new TileLayer({
        source: new XYZ({
            url: url,
            tileSize: 256,
            crossOrigin: 'anonymous'
        }),
        visible: true,
    })

    mapValue.addLayer(layer);
}