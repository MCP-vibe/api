package entities

import (
	"api/internal/errors"
	"context"
	"net/http"
	"time"
)

var (
	ErrOrganizationNotFound = errors.NewError(http.StatusNotFound, "organization_not_found")
)

type OrganizationRepository interface {
	Create(ctx context.Context, organization Organization) (Organization, error)
	// Update(ctx context.Context, organization Organization) error
	FindByID(ctx context.Context, id uint32) (Organization, error)
	FindAll(ctx context.Context) ([]Organization, int64, error)
}

type Organization struct {
	id        uint32
	name      string
	createdAt time.Time
}

func NewOrganization(
	id uint32,
	name string,
	createdAt time.Time,
) Organization {
	return Organization{
		id:        id,
		name:      name,
		createdAt: createdAt,
	}
}

func NewOrganizationCreate(
	name string,
	createdAt time.Time,
) Organization {
	return Organization{
		name:      name,
		createdAt: createdAt,
	}
}

func (u Organization) ID() uint32 {
	return u.id
}

func (u Organization) Name() string {
	return u.name
}

func (u Organization) CreatedAt() time.Time {
	return u.createdAt
}

// -----------------------------------------------------------------------------

func (u *Organization) SetID(id uint32) {
	u.id = id
}

func (u *Organization) SetName(name string) {
	u.name = name
}
