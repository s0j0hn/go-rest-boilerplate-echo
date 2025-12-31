package repository

import (
	"context"

	"github.com/google/uuid"
	"gitlab.com/s0j0hn/go-rest-boilerplate-echo/internal/domain"
)

// TenantRepository defines the interface for tenant data access
type TenantRepository interface {
	GetAll(ctx context.Context) ([]domain.Tenant, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Tenant, error)
	Create(ctx context.Context, tenant *domain.Tenant) error
	Update(ctx context.Context, tenant *domain.Tenant) error
	Delete(ctx context.Context, id uuid.UUID) error
	ExistsByName(ctx context.Context, name string) (bool, error)
}

// Repository holds all repository interfaces
type Repository struct {
	Tenant TenantRepository
}