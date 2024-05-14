package logging

// Attribute key const.
const (
	BuildGitHash   = "build.gitHash"
	BuildTimestamp = "build.timestamp"
	HostName       = "host.name"

	ServiceMethod = "service.method"

	UserID        = "user.id"
	UserUsername  = "user.username"
	UserFirstName = "user.firstName"
	UserLastName  = "user.lastName"

	EmployeeID            = "employee.id"
	EmployeeUsername      = "employee.username"
	EmployeeFirstName     = "employee.firstName"
	EmployeeLastName      = "employee.lastName"
	EmployeeRole          = "employee.role"
	EmployeeDateOfBirth   = "employee.dateOfBirth"
	EmployeePhoneNumber   = "employee.phoneNumber"
	EmployeeGeoJSON       = "employee.geoJson"
	EmployeeScheduleStart = "employee.scheduleStart"
	EmployeeScheduleEnd   = "employee.scheduleEnd"

	ContainerID       = "container.id"
	ContainerCategory = "container.category"

	TruckID             = "truck.id"
	TruckMake           = "truck.make"
	TruckModel          = "truck.model"
	TruckLicensePlate   = "truck.licensePlate"
	TruckPersonCapacity = "truck.personCapacity"

	WarehouseID            = "warehouse.id"
	WarehouseTruckCapacity = "warehouse.truckCapacity"

	LandfillID = "landfill.id"

	RouteID                   = "route.id"
	RouteName                 = "route.name"
	RouteTruckID              = "route.truckID"
	RouteDepartureWarehouseID = "route.departureWarehouseID"
	RouteArrivalWarehouseID   = "route.arrivalWarehouseID"

	RouteEmployeeEmployeeRole = "routeEmployee.employeeRole"
)
