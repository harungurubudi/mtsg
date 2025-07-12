package tenant

import (
	"time"

	"github.com/google/uuid"
)

// TenantID represents the unique identity of a tenant.
type TenantID uuid.UUID

// TenantStatus represents the status of a tenant.
type TenantStatus string

const (
	// TenantStatusActive indicates the tenant is active.
	TenantStatusActive TenantStatus = "active"
	// TenantStatusInactive indicates the tenant is inactive.
	TenantStatusInactive TenantStatus = "inactive"
	// TenantStatusDefault is the default status for a tenant.
	TenantStatusDefault TenantStatus = TenantStatusActive
)

// Tenant represents a logical customer or organization using the application.
type Tenant struct {
	ID        TenantID     `json:"id"`         // Unique identifier for the tenant
	Name      string       `json:"name"`       // Tenant name
	Status    TenantStatus `json:"status"`     // Status of the tenant (active or inactive)
	CreatedAt time.Time    `json:"created_at"` // Timestamp of tenant creation
	UpdatedAt time.Time    `json:"updated_at"` // Timestamp of last update
	DeletedAt *time.Time   `json:"-"`          // Timestamp of deletion (if soft deleted), not marshaled
	IsDeleted bool         `json:"-"`          // Soft delete flag, not marshaled
}

// NewTenant creates a new Tenant with the given name and status.
// If status is empty, it defaults to TenantStatusActive.
func NewTenant(name string, status TenantStatus) *Tenant {
	id := uuid.New()
	now := time.Now().UTC()
	if status == "" {
		status = TenantStatusDefault
	}
	return &Tenant{
		ID:        TenantID(id),
		Name:      name,
		Status:    status,
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: nil,
		IsDeleted: false,
	}
}
