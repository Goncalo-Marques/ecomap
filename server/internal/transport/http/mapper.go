package http

import (
	"time"

	oapitypes "github.com/oapi-codegen/runtime/types"

	spec "github.com/goncalo-marques/ecomap/server/api/ecomap"
	"github.com/goncalo-marques/ecomap/server/internal/domain"
)

// dateFromTime returns a standardized date based on the time model.
func dateFromTime(time time.Time) oapitypes.Date {
	return oapitypes.Date{
		Time: time,
	}
}

// jwtFromJWTToken returns a standardized JWT based on the JWT token.
func jwtFromJWTToken(token string) spec.JWT {
	return spec.JWT{
		Token: token,
	}
}

// userPostToDomainEditableUserWithPassword returns a domain editable user with password based on the standardized user
// post.
func userPostToDomainEditableUserWithPassword(userPost spec.UserPost) domain.EditableUserWithPassword {
	return domain.EditableUserWithPassword{
		EditableUser: domain.EditableUser{
			Username:  domain.Username(userPost.Username),
			FirstName: domain.Name(userPost.FirstName),
			LastName:  domain.Name(userPost.LastName),
		},
		Password: domain.Password(userPost.Password),
	}
}

// userFromDomainUser returns a standardized user based on the domain model.
func userFromDomainUser(user domain.User) spec.User {
	return spec.User{
		Id:           user.ID,
		Username:     string(user.Username),
		FirstName:    string(user.FirstName),
		LastName:     string(user.LastName),
		CreatedTime:  user.CreatedTime,
		ModifiedTime: user.ModifiedTime,
	}
}

// employeeFromDomain returns a standardized employee based on the domain model.
func employeeFromDomain(employee domain.Employee) spec.Employee {
	return spec.Employee{
		Id:          employee.ID,
		Name:        string(employee.FirstName),
		DateOfBirth: dateFromTime(employee.DateOfBirth),
	}
}
