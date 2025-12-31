package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gitlab.com/s0j0hn/go-rest-boilerplate-echo/internal/domain"
	"gitlab.com/s0j0hn/go-rest-boilerplate-echo/internal/repository"
	"gitlab.com/s0j0hn/go-rest-boilerplate-echo/rabbitmq"
)

// tenantService implements TenantService interface
type tenantService struct {
	tenantRepo  repository.TenantRepository
	taskManager *rabbitmq.TaskClient
}

// NewTenantService creates a new tenant service
func NewTenantService(tenantRepo repository.TenantRepository, taskManager *rabbitmq.TaskClient) TenantService {
	return &tenantService{
		tenantRepo:  tenantRepo,
		taskManager: taskManager,
	}
}

// GetAllTenants retrieves all tenants
func (s *tenantService) GetAllTenants(ctx context.Context) ([]domain.TenantResponse, error) {
	tenants, err := s.tenantRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all tenants: %w", err)
	}

	responses := make([]domain.TenantResponse, 0, len(tenants))
	for _, tenant := range tenants {
		responses = append(responses, tenant.ToResponse())
	}

	return responses, nil
}

// GetTenantByID retrieves a tenant by ID
func (s *tenantService) GetTenantByID(ctx context.Context, id uuid.UUID) (*domain.TenantResponse, error) {
	tenant, err := s.tenantRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant by ID: %w", err)
	}

	if tenant == nil {
		return nil, NewNotFoundError("tenant", id.String())
	}

	response := tenant.ToResponse()
	return &response, nil
}

// CreateTenant creates a new tenant
func (s *tenantService) CreateTenant(ctx context.Context, req domain.TenantCreateRequest) (*domain.TenantResponse, error) {
	// Check if tenant with same name already exists
	exists, err := s.tenantRepo.ExistsByName(ctx, req.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to check tenant existence: %w", err)
	}

	if exists {
		return nil, NewConflictError("tenant", "name", req.Name)
	}

	// Create new tenant
	tenant := &domain.Tenant{
		Name: req.Name,
	}

	// Start async task for tenant creation
	go func() {
		if s.taskManager != nil {
			task := rabbitmq.CreateNewTask([]string{"tenant", "create"}, fmt.Sprintf("Creating tenant: %s", req.Name))
			// Simulate async processing
			err := s.tenantRepo.Create(context.Background(), tenant)
			if err != nil {
				s.taskManager.FailTask(task)
				return
			}
			s.taskManager.CompleteTask(task)
		}
	}()

	if err := s.tenantRepo.Create(ctx, tenant); err != nil {
		return nil, fmt.Errorf("failed to create tenant: %w", err)
	}

	response := tenant.ToResponse()
	return &response, nil
}

// UpdateTenant updates an existing tenant
func (s *tenantService) UpdateTenant(ctx context.Context, req domain.TenantUpdateRequest) (*domain.TenantResponse, error) {
	// Check if tenant exists
	existingTenant, err := s.tenantRepo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant for update: %w", err)
	}

	if existingTenant == nil {
		return nil, NewNotFoundError("tenant", req.ID.String())
	}

	// Check if another tenant with the same name exists
	if existingTenant.Name != req.Name {
		exists, err := s.tenantRepo.ExistsByName(ctx, req.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to check tenant name existence: %w", err)
		}

		if exists {
			return nil, NewConflictError("tenant", "name", req.Name)
		}
	}

	// Update tenant
	existingTenant.Name = req.Name

	if err := s.tenantRepo.Update(ctx, existingTenant); err != nil {
		return nil, fmt.Errorf("failed to update tenant: %w", err)
	}

	response := existingTenant.ToResponse()
	return &response, nil
}

// DeleteTenant deletes a tenant by ID
func (s *tenantService) DeleteTenant(ctx context.Context, id uuid.UUID) error {
	// Check if tenant exists
	tenant, err := s.tenantRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get tenant for deletion: %w", err)
	}

	if tenant == nil {
		return NewNotFoundError("tenant", id.String())
	}

	if err := s.tenantRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete tenant: %w", err)
	}

	return nil
}