package service

import (
	"context"

	"github.com/google/uuid"
	"gitlab.com/s0j0hn/go-rest-boilerplate-echo/internal/domain"
)

// TenantService defines the business logic interface for tenant operations
type TenantService interface {
	GetAllTenants(ctx context.Context) ([]domain.TenantResponse, error)
	GetTenantByID(ctx context.Context, id uuid.UUID) (*domain.TenantResponse, error)
	CreateTenant(ctx context.Context, req domain.TenantCreateRequest) (*domain.TenantResponse, error)
	UpdateTenant(ctx context.Context, req domain.TenantUpdateRequest) (*domain.TenantResponse, error)
	DeleteTenant(ctx context.Context, id uuid.UUID) error
}

// Services holds all service interfaces
type Services struct {
	Tenant TenantService
}