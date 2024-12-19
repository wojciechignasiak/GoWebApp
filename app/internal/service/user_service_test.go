package service

import (
	apperror "app/internal/app_error"
	"testing"
)

func stringPtr(s string) *string {
	return &s
}

func getStringPointerValue(ptr *string) string {
	if ptr == nil {
		return "<nil>"
	}
	return *ptr
}

var ValidateUsernameTestcases = []struct {
	name     string
	username string
	expected *apperror.AppError
}{
	{"success", "wojciechignasiak96", nil},
	{
		"failure - username to short",
		"wk96",
		&apperror.AppError{
			StatusCode:      400,
			Message:         "Username must contain between 5 and 20 characters",
			StructAndMethod: "UserService.validateUsername()",
			Argument:        stringPtr("username: wk96"),
			ChildAppError:   nil,
			ChildError:      nil},
	},
	{
		"failure - username to long",
		"tooLongUsername1996!!!",
		&apperror.AppError{
			StatusCode:      400,
			Message:         "Username must contain between 5 and 20 characters",
			StructAndMethod: "UserService.validateUsername()",
			Argument:        stringPtr("username: tooLongUsername1996!!!"),
			ChildAppError:   nil,
			ChildError:      nil},
	},
}

func TestValidateUsername(t *testing.T) {
	us := &userService{}
	for _, tc := range ValidateUsernameTestcases {
		t.Run(tc.name, func(t *testing.T) {
			expected := tc.expected
			got := us.validateUsername(tc.username)

			if got == nil && expected == nil {
				return
			}

			if (got == nil && expected != nil) || (got != nil && expected == nil) {
				t.Errorf("scenario: %s, expected: %v, got: %v", tc.name, expected, got)
				return
			}

			if got.StatusCode != expected.StatusCode ||
				got.Message != expected.Message ||
				got.StructAndMethod != expected.StructAndMethod ||
				(got.Argument != nil && expected.Argument != nil && *got.Argument != *expected.Argument) ||
				(got.ChildAppError != nil && expected.ChildAppError != nil && *got.ChildAppError != *expected.ChildAppError) ||
				(got.ChildError != nil && expected.ChildError != nil && *got.ChildError != *expected.ChildError) {
				t.Errorf(
					`scenario: %s, mismatch found:
					Expected:
						StatusCode:      %d
						Message:         %s
						StructAndMethod: %s
						Argument:        %v
						ChildAppError:   %v
						ChildError:      %v
					Got:
						StatusCode:      %d
						Message:         %s
						StructAndMethod: %s
						Argument:        %v
						ChildAppError:   %v
						ChildError:      %v`,
					tc.name,
					expected.StatusCode,
					expected.Message,
					expected.StructAndMethod,
					getStringPointerValue(expected.Argument),
					expected.ChildAppError,
					expected.ChildError,
					got.StatusCode,
					got.Message,
					got.StructAndMethod,
					getStringPointerValue(got.Argument),
					got.ChildAppError,
					got.ChildError,
				)
			}
		})
	}
}

var ValidateEmailsTestcases = []struct {
	name         string
	email        string
	confirmEmail string
	expected     *apperror.AppError
}{
	{"success", "wojciech_ignasiak@icloud.com", "wojciech_ignasiak@icloud.com", nil},
	{
		"failure - emails do not match",
		"wojciech_ignasiak@icloud.com",
		"ignasiak_wojciech@icloud.com",
		&apperror.AppError{
			StatusCode:      400,
			Message:         "Provided emails do not match",
			StructAndMethod: "UserService.validateEmails()",
			Argument:        stringPtr("email: wojciech_ignasiak@icloud.com, confirmEmail: ignasiak_wojciech@icloud.com"),
			ChildAppError:   nil,
			ChildError:      nil},
	},
	{
		"failure - invalid email format",
		"wojciech_ignasiakicloud.com",
		"wojciech_ignasiakicloud.com",
		&apperror.AppError{
			StatusCode:      400,
			Message:         "Invalid email format",
			StructAndMethod: "UserService.validateEmails()",
			Argument:        stringPtr("email: wojciech_ignasiakicloud.com, confirmEmail: wojciech_ignasiakicloud.com"),
			ChildAppError:   nil,
			ChildError:      nil},
	},
}

func TestValidateEmails(t *testing.T) {
	us := &userService{}
	for _, tc := range ValidateEmailsTestcases {
		t.Run(tc.name, func(t *testing.T) {
			expected := tc.expected
			got := us.validateEmails(tc.email, tc.confirmEmail)

			if got == nil && expected == nil {
				return
			}

			if (got == nil && expected != nil) || (got != nil && expected == nil) {
				t.Errorf("scenario: %s, expected: %v, got: %v", tc.name, expected, got)
				return
			}

			if got.StatusCode != expected.StatusCode ||
				got.Message != expected.Message ||
				got.StructAndMethod != expected.StructAndMethod ||
				(got.Argument != nil && expected.Argument != nil && *got.Argument != *expected.Argument) ||
				(got.ChildAppError != nil && expected.ChildAppError != nil && *got.ChildAppError != *expected.ChildAppError) ||
				(got.ChildError != nil && expected.ChildError != nil && *got.ChildError != *expected.ChildError) {
				t.Errorf(
					`scenario: %s, mismatch found:
				Expected:
					StatusCode:      %d
					Message:         %s
					StructAndMethod: %s
					Argument:        %v
					ChildAppError:   %v
					ChildError:      %v
				Got:
					StatusCode:      %d
					Message:         %s
					StructAndMethod: %s
					Argument:        %v
					ChildAppError:   %v
					ChildError:      %v`,
					tc.name,
					expected.StatusCode,
					expected.Message,
					expected.StructAndMethod,
					getStringPointerValue(expected.Argument),
					expected.ChildAppError,
					expected.ChildError,
					got.StatusCode,
					got.Message,
					got.StructAndMethod,
					getStringPointerValue(got.Argument),
					got.ChildAppError,
					got.ChildError,
				)
			}
		})
	}
}

var ValidatePasswordsTestcases = []struct {
	name            string
	password        string
	confirmPassword string
	expected        *apperror.AppError
}{
	{"success", "!hardPassw0rd.", "!hardPassw0rd.", nil},
	{
		"failure - passwords are not the same",
		"a!hardPassw0rd.",
		"!hardPassw0rd.",
		&apperror.AppError{
			StatusCode:      400,
			Message:         "Provided passwords are not the same",
			StructAndMethod: "UserService.validatePasswords()",
			Argument:        nil,
			ChildAppError:   nil,
			ChildError:      nil},
	},
	{
		"failure - password too short",
		"!har1",
		"!har1",
		&apperror.AppError{
			StatusCode:      400,
			Message:         "Password must contain at least 8 characters",
			StructAndMethod: "UserService.validatePasswords()",
			Argument:        nil,
			ChildAppError:   nil,
			ChildError:      nil},
	},
	{
		"failure - missing special characters",
		"hardpassword",
		"hardpassword",
		&apperror.AppError{
			StatusCode:      403,
			Message:         "Password must contain at least one digit and one special character",
			StructAndMethod: "UserService.validatePasswords()",
			Argument:        nil,
			ChildAppError:   nil,
			ChildError:      nil},
	},
}

func TestValidatePasswords(t *testing.T) {
	us := &userService{}
	for _, tc := range ValidatePasswordsTestcases {
		t.Run(tc.name, func(t *testing.T) {
			expected := tc.expected
			got := us.validatePasswords(tc.password, tc.confirmPassword)

			if got == nil && expected == nil {
				return
			}

			if (got == nil && expected != nil) || (got != nil && expected == nil) {
				t.Errorf("scenario: %s, expected: %v, got: %v", tc.name, expected, got)
				return
			}

			if got.StatusCode != expected.StatusCode ||
				got.Message != expected.Message ||
				got.StructAndMethod != expected.StructAndMethod ||
				(got.Argument != nil && expected.Argument != nil && *got.Argument != *expected.Argument) ||
				(got.ChildAppError != nil && expected.ChildAppError != nil && *got.ChildAppError != *expected.ChildAppError) ||
				(got.ChildError != nil && expected.ChildError != nil && *got.ChildError != *expected.ChildError) {
				t.Errorf(
					`scenario: %s, mismatch found:
				Expected:
					StatusCode:      %d
					Message:         %s
					StructAndMethod: %s
					Argument:        %v
					ChildAppError:   %v
					ChildError:      %v
				Got:
					StatusCode:      %d
					Message:         %s
					StructAndMethod: %s
					Argument:        %v
					ChildAppError:   %v
					ChildError:      %v`,
					tc.name,
					expected.StatusCode,
					expected.Message,
					expected.StructAndMethod,
					getStringPointerValue(expected.Argument),
					expected.ChildAppError,
					expected.ChildError,
					got.StatusCode,
					got.Message,
					got.StructAndMethod,
					getStringPointerValue(got.Argument),
					got.ChildAppError,
					got.ChildError,
				)
			}
		})
	}
}
