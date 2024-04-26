package service

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/goncalo-marques/ecomap/server/internal/authn"
	"github.com/goncalo-marques/ecomap/server/internal/domain"
	"github.com/goncalo-marques/ecomap/server/internal/logging"
)

const (
	descriptionFailedCreateEmployee        = "service: failed to create employee"
	descriptionFailedGetEmployeeByID       = "service: failed to get employee by id"
	descriptionFailedGetEmployeeByUsername = "service: failed to get employee by username"
	descriptionFailedGetEmployeeSignIn     = "service: failed to get employee sign-in"
)

// CreateEmployee creates a new employee with the specified data.
func (s *service) CreateEmployee(ctx context.Context, editableEmployee domain.EditableEmployeeWithPassword) (domain.Employee, error) {
	logAttrs := []any{
		slog.String(logging.ServiceMethod, "CreateEmployee"),
		slog.String(logging.EmployeeUsername, string(editableEmployee.Username)),
		slog.String(logging.EmployeeFirstName, string(editableEmployee.FirstName)),
		slog.String(logging.EmployeeLastName, string(editableEmployee.LastName)),
		slog.String(logging.EmployeeRole, string(editableEmployee.Role)),
		slog.Time(logging.EmployeeDateOfBirth, editableEmployee.DateOfBirth),
		slog.String(logging.EmployeePhoneNumber, string(editableEmployee.PhoneNumber)),
		slog.Any(logging.EmployeeGeoJSON, editableEmployee.GeoJSON),
		slog.Time(logging.EmployeeScheduleStart, editableEmployee.ScheduleStart),
		slog.Time(logging.EmployeeScheduleEnd, editableEmployee.ScheduleEnd),
	}

	editableEmployee.Username = domain.Username(replaceSpacesWithHyphen(string(editableEmployee.Username)))
	editableEmployee.Username = domain.Username(strings.ToLower(string(editableEmployee.Username)))
	editableEmployee.FirstName = domain.Name(removeExtraSpaces(string(editableEmployee.FirstName)))
	editableEmployee.LastName = domain.Name(removeExtraSpaces(string(editableEmployee.LastName)))

	if !editableEmployee.Username.Valid() {
		return domain.Employee{}, logInfoAndWrapError(ctx, &domain.ErrFieldValueInvalid{FieldName: domain.FieldUsername}, descriptionInvalidFieldValue, logAttrs...)
	}
	if !s.authnService.ValidPassword([]byte(editableEmployee.Password)) {
		return domain.Employee{}, logInfoAndWrapError(ctx, &domain.ErrFieldValueInvalid{FieldName: domain.FieldPassword}, descriptionInvalidFieldValue, logAttrs...)
	}
	if !editableEmployee.FirstName.Valid() {
		return domain.Employee{}, logInfoAndWrapError(ctx, &domain.ErrFieldValueInvalid{FieldName: domain.FieldFirstName}, descriptionInvalidFieldValue, logAttrs...)
	}
	if !editableEmployee.LastName.Valid() {
		return domain.Employee{}, logInfoAndWrapError(ctx, &domain.ErrFieldValueInvalid{FieldName: domain.FieldLastName}, descriptionInvalidFieldValue, logAttrs...)
	}
	if !editableEmployee.Role.Valid() {
		return domain.Employee{}, logInfoAndWrapError(ctx, &domain.ErrFieldValueInvalid{FieldName: domain.FieldRole}, descriptionInvalidFieldValue, logAttrs...)
	}
	if !editableEmployee.PhoneNumber.Valid() {
		return domain.Employee{}, logInfoAndWrapError(ctx, &domain.ErrFieldValueInvalid{FieldName: domain.FieldPhoneNumber}, descriptionInvalidFieldValue, logAttrs...)
	}
	if editableEmployee.ScheduleStart.After(editableEmployee.ScheduleEnd) {
		return domain.Employee{}, logInfoAndWrapError(ctx, &domain.ErrFieldValueInvalid{FieldName: domain.FieldScheduleStart}, descriptionInvalidFieldValue, logAttrs...)
	}

	hashedPassword, err := s.authnService.HashPassword([]byte(editableEmployee.Password))
	if err != nil {
		return domain.Employee{}, logAndWrapError(ctx, err, descriptionFailedHashPassword, logAttrs...)
	}

	editableEmployee.Password = domain.Password(hashedPassword)

	var geometry domain.GeoJSONGeometryPoint
	if feature, ok := editableEmployee.GeoJSON.(domain.GeoJSONFeature); ok {
		if g, ok := feature.Geometry.(domain.GeoJSONGeometryPoint); ok {
			geometry = g
		}
	}

	var employee domain.Employee

	err = s.readWriteTx(ctx, func(tx pgx.Tx) error {
		var roadID *int
		road, err := s.store.GetRoadByGeometry(ctx, tx, geometry)
		if err != nil {
			if !errors.Is(err, domain.ErrRoadNotFound) {
				return err
			}
		} else {
			roadID = &road.ID
		}

		var municipalityID *int
		municipality, err := s.store.GetMunicipalityByGeometry(ctx, tx, geometry)
		if err != nil {
			if !errors.Is(err, domain.ErrMunicipalityNotFound) {
				return err
			}
		} else {
			municipalityID = &municipality.ID
		}

		id, err := s.store.CreateEmployee(ctx, tx, editableEmployee, roadID, municipalityID)
		if err != nil {
			return err
		}

		employee, err = s.store.GetEmployeeByID(ctx, tx, id)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrEmployeeAlreadyExists):
			return domain.Employee{}, logInfoAndWrapError(ctx, err, descriptionFailedCreateEmployee, logAttrs...)
		default:
			return domain.Employee{}, logAndWrapError(ctx, err, descriptionFailedCreateEmployee, logAttrs...)
		}
	}

	return employee, nil
}

// GetEmployeeByID returns the employee with the specified identifier.
func (s *service) GetEmployeeByID(ctx context.Context, id uuid.UUID) (domain.Employee, error) {
	logAttrs := []any{
		slog.String(logging.ServiceMethod, "GetEmployeeByID"),
		slog.String(logging.EmployeeID, id.String()),
	}

	var employee domain.Employee
	var err error

	err = s.readOnlyTx(ctx, func(tx pgx.Tx) error {
		employee, err = s.store.GetEmployeeByID(ctx, tx, id)
		return err
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrEmployeeNotFound):
			return domain.Employee{}, logInfoAndWrapError(ctx, err, descriptionFailedGetEmployeeByID, logAttrs...)
		default:
			return domain.Employee{}, logAndWrapError(ctx, err, descriptionFailedGetEmployeeByID, logAttrs...)
		}
	}

	return employee, nil
}

// SignInEmployee returns a JSON Web Token for the specified username and password.
func (s *service) SignInEmployee(ctx context.Context, username domain.Username, password domain.Password) (string, error) {
	logAttrs := []any{
		slog.String(logging.ServiceMethod, "SignInEmployee"),
		slog.String(logging.EmployeeUsername, string(username)),
	}

	username = domain.Username(strings.ToLower(string(username)))

	var signIn domain.SignIn
	var err error

	err = s.readOnlyTx(ctx, func(tx pgx.Tx) error {
		signIn, err = s.store.GetEmployeeSignIn(ctx, tx, username)
		return err
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrEmployeeNotFound):
			return "", logInfoAndWrapError(ctx, domain.ErrCredentialsIncorrect, descriptionFailedGetEmployeeSignIn, logAttrs...)
		default:
			return "", logAndWrapError(ctx, err, descriptionFailedGetEmployeeSignIn, logAttrs...)
		}
	}

	valid, err := s.authnService.CheckPasswordHash([]byte(password), []byte(signIn.Password))
	if err != nil {
		return "", logAndWrapError(ctx, err, descriptionFailedCheckPasswordHash, logAttrs...)
	}

	if !valid {
		return "", logInfoAndWrapError(ctx, domain.ErrCredentialsIncorrect, descriptionFailedCheckPasswordHash, logAttrs...)
	}

	var employee domain.Employee

	err = s.readOnlyTx(ctx, func(tx pgx.Tx) error {
		employee, err = s.store.GetEmployeeByUsername(ctx, tx, username)
		return err
	})
	if err != nil {
		return "", logAndWrapError(ctx, err, descriptionFailedGetEmployeeByUsername, logAttrs...)
	}

	var role authn.SubjectRole
	switch employee.Role {
	case domain.EmployeeRoleWasteOperator:
		role = authn.SubjectRoleWasteOperator
	case domain.EmployeeRoleManager:
		role = authn.SubjectRoleManager
	}

	token, err := s.authnService.NewJWT(employee.ID.String(), []authn.SubjectRole{role})
	if err != nil {
		return "", logAndWrapError(ctx, err, descriptionFailedCreateJWT, logAttrs...)
	}

	return token, nil
}
