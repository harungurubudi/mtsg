package user

import (
	"time"

	"github.com/google/uuid"
	"github.com/harungurubudi/mtsg/internal/domain/tenant"
)

// UserID represents the unique identity of a user.
type UserID uuid.UUID

// CipherText represents a hashed password.
type CipherText string

// UserRole represents the role of a user.
type UserRole string

const (
	// UserRoleAdmin indicates the user is an admin.
	UserRoleAdmin UserRole = "admin"
	// UserRoleMember indicates the user is a member.
	UserRoleMember UserRole = "member"
	// UserRoleDefault is the default role for a user.
	UserRoleDefault UserRole = UserRoleMember
)

// UserStatus represents the status of a user.
type UserStatus string

const (
	// UserStatusActive indicates the user is active.
	UserStatusActive UserStatus = "active"
	// UserStatusInactive indicates the user is inactive.
	UserStatusInactive UserStatus = "inactive"
	// UserStatusDefault is the default status for a user.
	UserStatusDefault UserStatus = UserStatusActive
)

// User represents a user entity in the system.
type User struct {
	ID         UserID          `json:"id"`         // Unique identifier for the user
	TenantID   tenant.TenantID `json:"tenant_id"`  // Tenant association (referenced from tenant package)
	Email      Email           `json:"email"`      // User's email address
	Name       string          `json:"name"`       // User's name
	Role       UserRole        `json:"role"`       // User's role (admin or member)
	CipherText CipherText      `json:"-"`          // Hashed password, not marshaled
	Status     UserStatus      `json:"status"`     // Status of the user (active or inactive)
	CreatedAt  time.Time       `json:"created_at"` // Timestamp of user creation
	UpdatedAt  time.Time       `json:"updated_at"` // Timestamp of last update
	DeletedAt  *time.Time      `json:"-"`          // Timestamp of deletion (if soft deleted), not marshaled
	IsDeleted  bool            `json:"-"`          // Soft delete flag, not marshaled
}

// NewUser creates a new User with the given constructor input fields and sensible defaults.
func NewUser(tenantID tenant.TenantID, email Email, name string, role UserRole) *User {
	id := uuid.New()
	now := time.Now().UTC()
	if role == "" {
		role = UserRoleDefault
	}
	return &User{
		ID:         UserID(id),
		TenantID:   tenantID,
		Email:      email,
		Name:       name,
		Role:       role,
		CipherText: "",
		Status:     UserStatusDefault,
		CreatedAt:  now,
		UpdatedAt:  now,
		DeletedAt:  nil,
		IsDeleted:  false,
	}
}
