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
	descriptionFailedCreateEmployee         = "service: failed to create employee"
	descriptionFailedGetEmployeeByID        = "service: failed to get employee by id"
	descriptionFailedGetEmployeeByUsername  = "service: failed to get employee by username"
	descriptionFailedGetEmployeeSignIn      = "service: failed to get employee sign-in"
	descriptionFailedPatchEmployee          = "service: failed to patch employee"
	descriptionFailedUpdateEmployeePassword = "service: failed to update employee password"
	descriptionFailedDeleteEmployeeByID     = "service: failed to delete employee by id"
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

// PatchEmployee modifies the employee with the specified identifier.
func (s *service) PatchEmployee(ctx context.Context, id uuid.UUID, editableEmployee domain.EditableEmployeePatch) (domain.Employee, error) {
	logAttrs := []any{
		slog.String(logging.ServiceMethod, "PatchEmployee"),
		slog.String(logging.EmployeeID, id.String()),
	}

	if editableEmployee.Username != nil {
		username := domain.Username(replaceSpacesWithHyphen(string(*editableEmployee.Username)))
		username = domain.Username(strings.ToLower(string(username)))
		editableEmployee.Username = &username
	}
	if editableEmployee.FirstName != nil {
		firstName := domain.Name(removeExtraSpaces(string(*editableEmployee.FirstName)))
		editableEmployee.FirstName = &firstName
	}
	if editableEmployee.LastName != nil {
		lastName := domain.Name(removeExtraSpaces(string(*editableEmployee.LastName)))
		editableEmployee.LastName = &lastName
	}

	if editableEmployee.Username != nil && !editableEmployee.Username.Valid() {
		return domain.Employee{}, logInfoAndWrapError(ctx, &domain.ErrFieldValueInvalid{FieldName: domain.FieldUsername}, descriptionInvalidFieldValue, logAttrs...)
	}
	if editableEmployee.FirstName != nil && !editableEmployee.FirstName.Valid() {
		return domain.Employee{}, logInfoAndWrapError(ctx, &domain.ErrFieldValueInvalid{FieldName: domain.FieldFirstName}, descriptionInvalidFieldValue, logAttrs...)
	}
	if editableEmployee.LastName != nil && !editableEmployee.LastName.Valid() {
		return domain.Employee{}, logInfoAndWrapError(ctx, &domain.ErrFieldValueInvalid{FieldName: domain.FieldLastName}, descriptionInvalidFieldValue, logAttrs...)
	}
	if editableEmployee.PhoneNumber != nil && !editableEmployee.PhoneNumber.Valid() {
		return domain.Employee{}, logInfoAndWrapError(ctx, &domain.ErrFieldValueInvalid{FieldName: domain.FieldPhoneNumber}, descriptionInvalidFieldValue, logAttrs...)
	}

	var geometry domain.GeoJSONGeometryPoint
	if editableEmployee.GeoJSON != nil {
		if feature, ok := editableEmployee.GeoJSON.(domain.GeoJSONFeature); ok {
			if g, ok := feature.Geometry.(domain.GeoJSONGeometryPoint); ok {
				geometry = g
			}
		}
	}

	var employee domain.Employee
	var err error

	err = s.readWriteTx(ctx, func(tx pgx.Tx) error {
		var roadID *int
		var municipalityID *int

		if editableEmployee.GeoJSON != nil {
			road, err := s.store.GetRoadByGeometry(ctx, tx, geometry)
			if err != nil {
				if !errors.Is(err, domain.ErrRoadNotFound) {
					return err
				}
			} else {
				roadID = &road.ID
			}

			municipality, err := s.store.GetMunicipalityByGeometry(ctx, tx, geometry)
			if err != nil {
				if !errors.Is(err, domain.ErrMunicipalityNotFound) {
					return err
				}
			} else {
				municipalityID = &municipality.ID
			}
		}

		err = s.store.PatchEmployee(ctx, tx, id, editableEmployee, roadID, municipalityID)
		if err != nil {
			return err
		}

		employee, err = s.store.GetEmployeeByID(ctx, tx, id)
		if err != nil {
			return err
		}

		if employee.ScheduleStart.After(employee.ScheduleEnd) {
			if editableEmployee.ScheduleStart != nil {
				return &domain.ErrFieldValueInvalid{FieldName: domain.FieldScheduleStart}
			} else {
				return &domain.ErrFieldValueInvalid{FieldName: domain.FieldScheduleEnd}
			}
		}

		return nil
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrEmployeeNotFound),
			errors.Is(err, domain.ErrEmployeeAlreadyExists):
			return domain.Employee{}, logInfoAndWrapError(ctx, err, descriptionFailedPatchEmployee, logAttrs...)
		default:
			return domain.Employee{}, logAndWrapError(ctx, err, descriptionFailedPatchEmployee, logAttrs...)
		}
	}

	return employee, nil
}

// UpdateEmployeePassword updates the password of the employee with the specified username.
func (s *service) UpdateEmployeePassword(ctx context.Context, username domain.Username, oldPassword, newPassword domain.Password) error {
	logAttrs := []any{
		slog.String(logging.ServiceMethod, "UpdateEmployeePassword"),
		slog.String(logging.EmployeeUsername, string(username)),
	}

	username = domain.Username(strings.ToLower(string(username)))

	if !s.authnService.ValidPassword([]byte(newPassword)) {
		return logInfoAndWrapError(ctx, &domain.ErrFieldValueInvalid{FieldName: domain.FieldNewPassword}, descriptionInvalidFieldValue, logAttrs...)
	}

	var signIn domain.SignIn
	var err error

	err = s.readOnlyTx(ctx, func(tx pgx.Tx) error {
		signIn, err = s.store.GetEmployeeSignIn(ctx, tx, username)
		return err
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrEmployeeNotFound):
			return logInfoAndWrapError(ctx, domain.ErrCredentialsIncorrect, descriptionFailedGetEmployeeSignIn, logAttrs...)
		default:
			return logAndWrapError(ctx, err, descriptionFailedGetEmployeeSignIn, logAttrs...)
		}
	}

	valid, err := s.authnService.CheckPasswordHash([]byte(oldPassword), []byte(signIn.Password))
	if err != nil {
		return logAndWrapError(ctx, err, descriptionFailedCheckPasswordHash, logAttrs...)
	}

	if !valid {
		return logInfoAndWrapError(ctx, domain.ErrCredentialsIncorrect, descriptionFailedCheckPasswordHash, logAttrs...)
	}

	hashedPassword, err := s.authnService.HashPassword([]byte(newPassword))
	if err != nil {
		return logAndWrapError(ctx, err, descriptionFailedHashPassword, logAttrs...)
	}

	newPassword = domain.Password(hashedPassword)

	err = s.readWriteTx(ctx, func(tx pgx.Tx) error {
		err = s.store.UpdateEmployeePassword(ctx, tx, username, newPassword)
		return err
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrEmployeeNotFound):
			return logInfoAndWrapError(ctx, domain.ErrCredentialsIncorrect, descriptionFailedUpdateEmployeePassword, logAttrs...)
		default:
			return logAndWrapError(ctx, err, descriptionFailedUpdateEmployeePassword, logAttrs...)
		}
	}

	return nil
}

// ResetEmployeePassword resets the password of the employee with the specified username.
func (s *service) ResetEmployeePassword(ctx context.Context, username domain.Username, newPassword domain.Password) error {
	logAttrs := []any{
		slog.String(logging.ServiceMethod, "ResetEmployeePassword"),
		slog.String(logging.EmployeeUsername, string(username)),
	}

	username = domain.Username(strings.ToLower(string(username)))

	if !s.authnService.ValidPassword([]byte(newPassword)) {
		return logInfoAndWrapError(ctx, &domain.ErrFieldValueInvalid{FieldName: domain.FieldNewPassword}, descriptionInvalidFieldValue, logAttrs...)
	}

	hashedPassword, err := s.authnService.HashPassword([]byte(newPassword))
	if err != nil {
		return logAndWrapError(ctx, err, descriptionFailedHashPassword, logAttrs...)
	}

	newPassword = domain.Password(hashedPassword)

	err = s.readWriteTx(ctx, func(tx pgx.Tx) error {
		err = s.store.UpdateEmployeePassword(ctx, tx, username, newPassword)
		return err
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrEmployeeNotFound):
			return logInfoAndWrapError(ctx, err, descriptionFailedUpdateEmployeePassword, logAttrs...)
		default:
			return logAndWrapError(ctx, err, descriptionFailedUpdateEmployeePassword, logAttrs...)
		}
	}

	return nil
}

// DeleteEmployeeByID deletes the employee with the specified identifier.
func (s *service) DeleteEmployeeByID(ctx context.Context, id uuid.UUID) (domain.Employee, error) {
	logAttrs := []any{
		slog.String(logging.ServiceMethod, "DeleteEmployeeByID"),
		slog.String(logging.EmployeeID, id.String()),
	}

	var employee domain.Employee
	var err error

	err = s.readWriteTx(ctx, func(tx pgx.Tx) error {
		employee, err = s.store.GetEmployeeByID(ctx, tx, id)
		if err != nil {
			return err
		}

		err = s.store.DeleteEmployeeByID(ctx, tx, id)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrEmployeeNotFound):
			return domain.Employee{}, logInfoAndWrapError(ctx, err, descriptionFailedDeleteEmployeeByID, logAttrs...)
		default:
			return domain.Employee{}, logAndWrapError(ctx, err, descriptionFailedDeleteEmployeeByID, logAttrs...)
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
