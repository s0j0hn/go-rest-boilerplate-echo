package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gitlab.com/s0j0hn/go-rest-boilerplate-echo/internal/domain"
	"gitlab.com/s0j0hn/go-rest-boilerplate-echo/internal/service"
)

// TenantHandler handles tenant-related HTTP requests
type TenantHandler struct {
	tenantService service.TenantService
}

// NewTenantHandler creates a new tenant handler
func NewTenantHandler(tenantService service.TenantService) *TenantHandler {
	return &TenantHandler{
		tenantService: tenantService,
	}
}

// GetAllTenants godoc
// @Summary List all tenants
// @Description Get a list of all tenants
// @Tags tenants
// @Accept json
// @Produce json
// @Success 200 {array} domain.TenantResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /tenants [get]
func (h *TenantHandler) GetAllTenants(c echo.Context) error {
	tenants, err := h.tenantService.GetAllTenants(c.Request().Context())
	if err != nil {
		return err // Let middleware handle it
	}

	return c.JSON(http.StatusOK, tenants)
}

// GetTenantByID godoc
// @Summary Get tenant by ID
// @Description Get a specific tenant by UUID
// @Tags tenants
// @Accept json
// @Produce json
// @Param id path string true "Tenant UUID"
// @Success 200 {object} domain.TenantResponse
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Router /tenants/{id} [get]
func (h *TenantHandler) GetTenantByID(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return service.NewValidationError("Invalid UUID format")
	}

	tenant, err := h.tenantService.GetTenantByID(c.Request().Context(), id)
	if err != nil {
		return err // Let middleware handle it
	}

	return c.JSON(http.StatusOK, tenant)
}

// CreateTenant godoc
// @Summary Create a new tenant
// @Description Create a new tenant with the provided data
// @Tags tenants
// @Accept json
// @Produce json
// @Param tenant body domain.TenantCreateRequest true "Tenant creation data"
// @Success 201 {object} domain.TenantResponse
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 409 {object} middleware.ErrorResponse
// @Router /tenants [post]
func (h *TenantHandler) CreateTenant(c echo.Context) error {
	var req domain.TenantCreateRequest

	if err := c.Bind(&req); err != nil {
		return service.NewBadRequestError("Invalid request body")
	}

	if err := c.Validate(req); err != nil {
		return err // Validation errors will be handled by middleware
	}

	tenant, err := h.tenantService.CreateTenant(c.Request().Context(), req)
	if err != nil {
		return err // Let middleware handle it
	}

	return c.JSON(http.StatusCreated, tenant)
}

// UpdateTenant godoc
// @Summary Update an existing tenant
// @Description Update a tenant with the provided data
// @Tags tenants
// @Accept json
// @Produce json
// @Param tenant body domain.TenantUpdateRequest true "Tenant update data"
// @Success 200 {object} domain.TenantResponse
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Failure 409 {object} middleware.ErrorResponse
// @Router /tenants [put]
func (h *TenantHandler) UpdateTenant(c echo.Context) error {
	var req domain.TenantUpdateRequest

	if err := c.Bind(&req); err != nil {
		return service.NewBadRequestError("Invalid request body")
	}

	if err := c.Validate(req); err != nil {
		return err // Validation errors will be handled by middleware
	}

	tenant, err := h.tenantService.UpdateTenant(c.Request().Context(), req)
	if err != nil {
		return err // Let middleware handle it
	}

	return c.JSON(http.StatusOK, tenant)
}

// DeleteTenant godoc
// @Summary Delete a tenant
// @Description Delete a tenant by UUID
// @Tags tenants
// @Accept json
// @Produce json
// @Param id path string true "Tenant UUID"
// @Success 204
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Router /tenants/{id} [delete]
func (h *TenantHandler) DeleteTenant(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return service.NewValidationError("Invalid UUID format")
	}

	if err := h.tenantService.DeleteTenant(c.Request().Context(), id); err != nil {
		return err // Let middleware handle it
	}

	return c.NoContent(http.StatusNoContent)
}