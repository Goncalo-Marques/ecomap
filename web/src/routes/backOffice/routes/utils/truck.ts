/**
 * Retrieves a truck name given its make, model and license plate.
 * @param make Truck make.
 * @param model Truck model.
 * @param licensePlate Truck license plate.
 * @returns Truck name.
 */
export function getTruckName(
	make: string,
	model: string,
	licensePlate: string,
) {
	return `${make} ${model} (${licensePlate})`;
}
