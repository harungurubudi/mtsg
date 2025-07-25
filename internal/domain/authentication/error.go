package authentication

import "errors"

// Error types for specific authentication scenarios
var (
	ErrInvalidCredential     = errors.New("invalid credentials")
	ErrUserNotFound          = errors.New("user not found")
	ErrAccountLocked         = errors.New("account is locked")
	ErrPasswordExpired       = errors.New("password has expired")
	ErrUserInactive          = errors.New("user account is inactive")
	ErrInvalidAuthentication = errors.New("invalid authentication")
	ErrForbidden             = errors.New("forbidden")
	ErrTenantMismatch        = errors.New("user does not belong to required tenant")
)
