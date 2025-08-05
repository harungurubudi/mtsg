package repository

import (
	"context"

	"github.com/google/uuid"
	tenantdomain "github.com/harungurubudi/mtsg/internal/domain/tenant"
	userdomain "github.com/harungurubudi/mtsg/internal/domain/user"
)

// UserRepository defines the interface for user data access.
// See internal/repository/.cursor/rules/user.mdc for requirements and best practices.
type UserRepository interface {
	// GetOneByEmail retrieves a user by email. Returns the user if found, or an error if not found.
	GetOneByEmail(ctx context.Context, email userdomain.Email) (*userdomain.User, error)
	// GetOneByID retrieves a user by ID. Returns the user if found, or an error if not found.
	GetOneByID(ctx context.Context, userID userdomain.UserID) (*userdomain.User, error)
}

// UserPersistence is the current in-memory implementation of UserRepository for testing and simulation.
type UserPersistence struct {
	users []userdomain.User
}

// NewUserPersistence creates a new UserPersistence with two hardcoded users.
func NewUserPersistence() *UserPersistence {
	uid1 := userdomain.UserID(uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"))
	uid2 := userdomain.UserID(uuid.MustParse("123e4567-e89b-12d3-a456-426614174001"))
	tid := tenantdomain.TenantID(uuid.MustParse("123e4567-e89b-12d3-a456-426614174002"))

	password1 := userdomain.Password("password1")
	hashedPassword1, _ := userdomain.NewCipherText(password1)
	password2 := userdomain.Password("password2")
	hashedPassword2, _ := userdomain.NewCipherText(password2)

	return &UserPersistence{
		users: []userdomain.User{
			{
				ID:         uid1,
				TenantID:   tid,
				Email:      "user1@example.com",
				Name:       "User One",
				Role:       "admin",
				CipherText: hashedPassword1,
				Status:     "active",
			},
			{
				ID:         uid2,
				TenantID:   tid,
				Email:      "user2@example.com",
				Name:       "User Two",
				Role:       "member",
				CipherText: hashedPassword2,
				Status:     "active",
			},
		},
	}
}

// GetOneByEmail returns the user with the given email, or an error if not found.
func (r *UserPersistence) GetOneByEmail(ctx context.Context, email userdomain.Email) (*userdomain.User, error) {
	for _, u := range r.users {
		if u.Email == email {
			return &u, nil
		}
	}
	return nil, userdomain.ErrUserNotFound
}

// GetOneByID returns the user with the given ID, or an error if not found.
func (r *UserPersistence) GetOneByID(ctx context.Context, userID userdomain.UserID) (*userdomain.User, error) {
	for _, u := range r.users {
		if u.ID == userID {
			return &u, nil
		}
	}
	return nil, userdomain.ErrUserNotFound
}
