package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/goncalo-marques/ecomap/server/internal/authn"
	"github.com/goncalo-marques/ecomap/server/internal/domain"
	"github.com/goncalo-marques/ecomap/server/internal/logging"
)

const (
	descriptionInvalidFieldValue       = "service: invalid field value"
	descriptionInvalidFilterValue      = "service: invalid filter value"
	descriptionFailedHashPassword      = "service: failed to hash password"
	descriptionFailedCheckPasswordHash = "service: failed to check password hash"
	descriptionFailedCreateJWT         = "service: failed to create jwt"
)

// AuthenticationService defines the authentication service interface.
type AuthenticationService interface {
	ValidPassword(password []byte) bool
	HashPassword(password []byte) ([]byte, error)
	CheckPasswordHash(password, hash []byte) (bool, error)

	NewJWT(subject string, subjectRoles []authn.SubjectRole) (string, error)
}

// Store defines the store interface.
type Store interface {
	CreateUser(ctx context.Context, tx pgx.Tx, editableUser domain.EditableUserWithPassword) (uuid.UUID, error)
	ListUsers(ctx context.Context, tx pgx.Tx, filter domain.UsersPaginatedFilter) (domain.PaginatedResponse[domain.User], error)
	GetUserByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (domain.User, error)
	GetUserByUsername(ctx context.Context, tx pgx.Tx, username domain.Username) (domain.User, error)
	GetUserSignIn(ctx context.Context, tx pgx.Tx, username domain.Username) (domain.SignIn, error)
	PatchUser(ctx context.Context, tx pgx.Tx, id uuid.UUID, editableUser domain.EditableUserPatch) error
	UpdateUserPassword(ctx context.Context, tx pgx.Tx, username domain.Username, password domain.Password) error
	DeleteUserByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) error

	CreateUserContainerBookmark(ctx context.Context, tx pgx.Tx, userID, containerID uuid.UUID) error
	ListUserContainerBookmarks(ctx context.Context, tx pgx.Tx, userID uuid.UUID, filter domain.UserContainerBookmarksPaginatedFilter) (domain.PaginatedResponse[domain.Container], error)
	DeleteUserContainerBookmark(ctx context.Context, tx pgx.Tx, userID, containerID uuid.UUID) error

	CreateEmployee(ctx context.Context, tx pgx.Tx, editableEmployee domain.EditableEmployeeWithPassword, roadID, municipalityID *int) (uuid.UUID, error)
	ListEmployees(ctx context.Context, tx pgx.Tx, filter domain.EmployeesPaginatedFilter) (domain.PaginatedResponse[domain.Employee], error)
	GetEmployeeByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (domain.Employee, error)
	GetEmployeeByUsername(ctx context.Context, tx pgx.Tx, username domain.Username) (domain.Employee, error)
	GetEmployeeSignIn(ctx context.Context, tx pgx.Tx, username domain.Username) (domain.SignIn, error)
	PatchEmployee(ctx context.Context, tx pgx.Tx, id uuid.UUID, editableEmployee domain.EditableEmployeePatch, roadID, municipalityID *int) error
	UpdateEmployeePassword(ctx context.Context, tx pgx.Tx, username domain.Username, password domain.Password) error
	DeleteEmployeeByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) error

	CreateContainer(ctx context.Context, tx pgx.Tx, editableContainer domain.EditableContainer, roadID, municipalityID *int) (uuid.UUID, error)
	ListContainers(ctx context.Context, tx pgx.Tx, filter domain.ContainersPaginatedFilter) (domain.PaginatedResponse[domain.Container], error)
	GetContainerByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (domain.Container, error)
	PatchContainer(ctx context.Context, tx pgx.Tx, id uuid.UUID, editableContainer domain.EditableContainerPatch, roadID, municipalityID *int) error
	DeleteContainerByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) error

	CreateTruck(ctx context.Context, tx pgx.Tx, editableTruck domain.EditableTruck, roadID, municipalityID *int) (uuid.UUID, error)
	ListTrucks(ctx context.Context, tx pgx.Tx, filter domain.TrucksPaginatedFilter) (domain.PaginatedResponse[domain.Truck], error)
	GetTruckByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (domain.Truck, error)
	PatchTruck(ctx context.Context, tx pgx.Tx, id uuid.UUID, editableTruck domain.EditableTruckPatch, roadID, municipalityID *int) error
	DeleteTruckByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) error

	CreateWarehouse(ctx context.Context, tx pgx.Tx, editableWarehouse domain.EditableWarehouse, roadID, municipalityID *int) (uuid.UUID, error)
	ListWarehouses(ctx context.Context, tx pgx.Tx, filter domain.WarehousesPaginatedFilter) (domain.PaginatedResponse[domain.Warehouse], error)
	GetWarehouseByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (domain.Warehouse, error)
	PatchWarehouse(ctx context.Context, tx pgx.Tx, id uuid.UUID, editableWarehouse domain.EditableWarehousePatch, roadID, municipalityID *int) error
	DeleteWarehouseByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) error

	CreateWarehouseTruck(ctx context.Context, tx pgx.Tx, warehouseID, truckID uuid.UUID) error
	ListWarehouseTrucks(ctx context.Context, tx pgx.Tx, warehouseID uuid.UUID, filter domain.WarehouseTrucksPaginatedFilter) (domain.PaginatedResponse[domain.Truck], error)
	DeleteWarehouseTruck(ctx context.Context, tx pgx.Tx, warehouseID, truckID uuid.UUID) error

	CreateLandfill(ctx context.Context, tx pgx.Tx, editableLandfill domain.EditableLandfill, roadID, municipalityID *int) (uuid.UUID, error)
	ListLandfills(ctx context.Context, tx pgx.Tx, filter domain.LandfillsPaginatedFilter) (domain.PaginatedResponse[domain.Landfill], error)
	GetLandfillByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (domain.Landfill, error)
	GetLandfillClosestGeometry(ctx context.Context, tx pgx.Tx, geometry domain.GeoJSONGeometryPoint) (domain.Landfill, error)
	PatchLandfill(ctx context.Context, tx pgx.Tx, id uuid.UUID, editableLandfill domain.EditableLandfillPatch, roadID, municipalityID *int) error
	DeleteLandfillByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) error

	CreateRoute(ctx context.Context, tx pgx.Tx, editableRoute domain.EditableRoute) (uuid.UUID, error)
	ListRoutes(ctx context.Context, tx pgx.Tx, filter domain.RoutesPaginatedFilter) (domain.PaginatedResponse[domain.Route], error)
	GetRouteByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (domain.Route, error)
	PatchRoute(ctx context.Context, tx pgx.Tx, id uuid.UUID, editableRoute domain.EditableRoutePatch) error
	DeleteRouteByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) error

	CreateRouteContainer(ctx context.Context, tx pgx.Tx, routeID, containerID uuid.UUID) error
	ListRouteContainers(ctx context.Context, tx pgx.Tx, routeID uuid.UUID, filter domain.RouteContainersPaginatedFilter) (domain.PaginatedResponse[domain.Container], error)
	DeleteRouteContainer(ctx context.Context, tx pgx.Tx, routeID, containerID uuid.UUID) error

	CreateRouteEmployee(ctx context.Context, tx pgx.Tx, routeID, employeeID uuid.UUID, editableRouteEmployee domain.EditableRouteEmployee) error
	ListRouteEmployees(ctx context.Context, tx pgx.Tx, routeID uuid.UUID, filter domain.RouteEmployeesPaginatedFilter) (domain.PaginatedResponse[domain.RouteEmployee], error)
	DeleteRouteEmployee(ctx context.Context, tx pgx.Tx, routeID, employeeID uuid.UUID) error

	GetRoadByGeometry(ctx context.Context, tx pgx.Tx, geometry domain.GeoJSONGeometryPoint) (domain.Road, error)
	CreateTemporaryTableRoadNetworkWithBuffer(ctx context.Context, tx pgx.Tx, tableName string, verticesGeometry []domain.GeoJSONGeometryPoint) error
	CreateVerticesCloseToRoadNetwork(ctx context.Context, tx pgx.Tx, roadNetworkTableName string, verticesGeometry []domain.GeoJSONGeometryPoint) ([]int, error)
	GetRoadVerticesTSP(ctx context.Context, tx pgx.Tx, roadNetworkTableName string, vertexIDs []int, startVertexID, endVertexID int, directed bool) ([]int, error)
	GetRoadsGeometryAStar(ctx context.Context, tx pgx.Tx, roadNetworkTableName string, seqVertexIDs []int, directed bool) ([]domain.GeoJSONGeometryLineString, error)

	GetMunicipalityByGeometry(ctx context.Context, tx pgx.Tx, geometry domain.GeoJSONGeometryPoint) (domain.Municipality, error)

	NewTx(ctx context.Context, isoLevel pgx.TxIsoLevel, accessMode pgx.TxAccessMode) (pgx.Tx, error)
}

// service defines the service structure.
type service struct {
	authnService AuthenticationService
	store        Store
}

// New returns a new http handler.
func New(authnService AuthenticationService, store Store) *service {
	return &service{
		authnService: authnService,
		store:        store,
	}
}

// rollbackFunc returns a function to rollback a transaction.
func rollbackFunc(ctx context.Context, tx pgx.Tx) func() {
	return func() {
		err := tx.Rollback(ctx)
		if err != nil {
			logging.Logger.ErrorContext(ctx, "service: failed to rollback transaction", logging.Error(err))
		}
	}
}

// readOnlyTx returns a read only transaction wrapper.
func (s *service) readOnlyTx(ctx context.Context, f func(pgx.Tx) error) error {
	tx, err := s.store.NewTx(ctx, pgx.ReadCommitted, pgx.ReadOnly)
	if err != nil {
		return err
	}
	defer rollbackFunc(ctx, tx)()

	if err := f(tx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// readWriteTx returns a read and write transaction wrapper.
func (s *service) readWriteTx(ctx context.Context, f func(pgx.Tx) error) error {
	tx, err := s.store.NewTx(ctx, pgx.RepeatableRead, pgx.ReadWrite)
	if err != nil {
		return err
	}
	defer rollbackFunc(ctx, tx)()

	if err := f(tx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// logInfoAndWrapError logs the error at the info level and returns the error wrapped with the provided description.
func logInfoAndWrapError(ctx context.Context, err error, description string, logAttrs ...any) error {
	logAttrs = append(logAttrs, logging.Error(err))
	logging.Logger.InfoContext(ctx, description, logAttrs...)
	return fmt.Errorf("%s: %w", description, err)
}

// logAndWrapError logs the error and returns the error wrapped with the provided description.
func logAndWrapError(ctx context.Context, err error, description string, logAttrs ...any) error {
	logAttrs = append(logAttrs, logging.Error(err))
	logging.Logger.ErrorContext(ctx, description, logAttrs...)
	return fmt.Errorf("%s: %w", description, err)
}

// replaceSpacesWithHyphen returns s with no extra spaces and separates it with a hyphen.
func replaceSpacesWithHyphen(s string) string {
	return strings.Join(strings.Fields(s), "-")
}

// removeExtraSpaces returns s with no extra spaces.
func removeExtraSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

// geometryPointFromGeoJSON returns the geometry point from the provided GeoJSON.
func geometryPointFromGeoJSON(geoJSON domain.GeoJSON) domain.GeoJSONGeometryPoint {
	if geoJSON == nil {
		return domain.GeoJSONGeometryPoint{}
	}

	feature, ok := geoJSON.(domain.GeoJSONFeature)
	if !ok {
		return domain.GeoJSONGeometryPoint{}
	}

	geometry, ok := feature.Geometry.(domain.GeoJSONGeometryPoint)
	if !ok {
		return domain.GeoJSONGeometryPoint{}
	}

	return geometry
}
