import Map from "ol/Map";
import { writable } from "svelte/store";

export const map = writable<Map>()