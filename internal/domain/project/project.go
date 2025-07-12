package project

import (
    "time"
    "github.com/google/uuid"
    "github.com/harungurubudi/mtsg/internal/domain/tenant"
    "github.com/harungurubudi/mtsg/internal/domain/user"
)

// ProjectID represents the unique identity of a project.
type ProjectID uuid.UUID

// ProjectStatus represents the status of a project.
type ProjectStatus string

const (
    // ProjectStatusActive indicates the project is active.
    ProjectStatusActive ProjectStatus = "active"
    // ProjectStatusInactive indicates the project is inactive.
    ProjectStatusInactive ProjectStatus = "inactive"
    // ProjectStatusArchived indicates the project is archived.
    ProjectStatusArchived ProjectStatus = "archived"
    // ProjectStatusDraft indicates the project is in draft state.
    ProjectStatusDraft ProjectStatus = "draft"
    // ProjectStatusDefault is the default status for a project.
    ProjectStatusDefault ProjectStatus = ProjectStatusActive
)

// Project represents a project entity in the system.
type Project struct {
    ID          ProjectID       `json:"id"`            // Unique identifier for the project
    TenantID    tenant.TenantID `json:"tenant_id"`     // Tenant association (referenced from tenant package)
    Name        string          `json:"name"`          // Project name
    Description string          `json:"description"`   // Project description
    CreatedBy   user.UserID     `json:"created_by"`    // User who created the project (referenced from user package)
    Status      ProjectStatus   `json:"status"`        // Status of the project (active, inactive, archived, draft)
    CreatedAt   time.Time       `json:"created_at"`    // Timestamp of project creation
    UpdatedAt   time.Time       `json:"updated_at"`    // Timestamp of last update
    DeletedAt   *time.Time      `json:"-"`             // Timestamp of deletion (if soft deleted), not marshaled
    IsDeleted   bool            `json:"-"`             // Soft delete flag, not marshaled
}

// NewProject creates a new Project with the given constructor input fields and sensible defaults.
func NewProject(tenantID tenant.TenantID, name string, createdBy user.UserID, status ProjectStatus) *Project {
    id := uuid.New()
    now := time.Now().UTC()
    if status == "" {
        status = ProjectStatusDefault
    }
    return &Project{
        ID:          ProjectID(id),
        TenantID:    tenantID,
        Name:        name,
        Description: "",
        CreatedBy:   createdBy,
        Status:      status,
        CreatedAt:   now,
        UpdatedAt:   now,
        DeletedAt:   nil,
        IsDeleted:   false,
    }
} 