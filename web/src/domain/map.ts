import type { Container } from "./container";
import type { Truck } from "./truck";
import type { Warehouse } from "./warehouse";

/**
 * Represents a group of resources in the same location.
 */
export interface ResourceGroupLocation {
	/**
	 * The containers in the location.
	 */
	containers: Container[];

	/**
	 * The trucks in the location.
	 */
	trucks: Truck[];

	/**
	 * The warehouses in the location.
	 */
	warehouses: Warehouse[];

	/**
	 * The way name of the location.
	 */
	wayName?: string;

	/**
	 * The municipality name of the location.
	 */
	municipalityName?: string;
}
