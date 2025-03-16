package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"gitlab.com/s0j0hn/go-rest-boilerplate-echo/internal/domain"
)

// tenantRepository implements TenantRepository interface
type tenantRepository struct {
	db *gorm.DB
}

// NewTenantRepository creates a new tenant repository
func NewTenantRepository(db *gorm.DB) TenantRepository {
	return &tenantRepository{
		db: db,
	}
}

// GetAll retrieves all tenants from database
func (r *tenantRepository) GetAll(ctx context.Context) ([]domain.Tenant, error) {
	var tenants []domain.Tenant

	if err := r.db.WithContext(ctx).Find(&tenants).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch tenants: %w", err)
	}

	return tenants, nil
}

// GetByID retrieves a tenant by UUID
func (r *tenantRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Tenant, error) {
	var tenant domain.Tenant

	if err := r.db.WithContext(ctx).Where("uuid = ?", id).First(&tenant).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Return nil instead of error for not found
		}
		return nil, fmt.Errorf("failed to fetch tenant by ID %s: %w", id, err)
	}

	return &tenant, nil
}

// Create creates a new tenant
func (r *tenantRepository) Create(ctx context.Context, tenant *domain.Tenant) error {
	if err := r.db.WithContext(ctx).Create(tenant).Error; err != nil {
		return fmt.Errorf("failed to create tenant: %w", err)
	}

	return nil
}

// Update updates an existing tenant
func (r *tenantRepository) Update(ctx context.Context, tenant *domain.Tenant) error {
	result := r.db.WithContext(ctx).Model(tenant).Where("uuid = ?", tenant.UUID).Updates(tenant)

	if result.Error != nil {
		return fmt.Errorf("failed to update tenant: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("tenant with ID %s not found", tenant.UUID)
	}

	return nil
}

// Delete deletes a tenant by UUID (soft delete)
func (r *tenantRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Where("uuid = ?", id).Delete(&domain.Tenant{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete tenant: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("tenant with ID %s not found", id)
	}

	return nil
}

// ExistsByName checks if a tenant with the given name exists
func (r *tenantRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	var count int64

	if err := r.db.WithContext(ctx).Model(&domain.Tenant{}).Where("name = ?", name).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check tenant existence by name: %w", err)
	}

	return count > 0, nil
}