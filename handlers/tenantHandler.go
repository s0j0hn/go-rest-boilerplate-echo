package handlers

import (
	libUuid "github.com/google/uuid"
	"github.com/labstack/echo/v4"
	tenantModel "gitlab.com/s0j0hn/go-rest-boilerplate-echo/database/models/tenant"
	"net/http"
)

type (
	handler struct {
		tenantModel tenantModel.TenantModel
	}
	resultJson struct {
		ID   libUuid.UUID `json:"id" form:"id" validate:"required"`
		Name string       `json:"name" form:"name" validate:"required"`
	}
	postTenantData struct {
		ID   libUuid.UUID `json:"id" form:"id" validate:"required"`
		Name string       `json:"name" form:"name" validate:"required"`
	}

	errorResult struct {
		Message string `json:"message" comment:"Something went wrong"`
	}
)

func CreateHandler(tenant tenantModel.TenantModel) *handler {
	return &handler{tenant}
}

// GetAll godoc
// @Summary List tenants
// @Description get tenants
// @Tags tenants
// @Accept  json
// @Produce  json
// @Success 200 {array} handlers.resultJson
// @Failure 500 {object} handlers.errorResult
// @Router /tenants [get]
func (h handler) GetAll(c echo.Context) error {
	tenants, err := h.tenantModel.GetAll()
	if err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusInternalServerError, errorResult{Message: err.Error()})
	}

	var results []resultJson
	for _, tenant := range *tenants {
		results = append(results, resultJson{ID: tenant.Uuid, Name: tenant.Name})
	}

	if len(results) == 0 {
		return c.JSON(http.StatusOK, []resultJson{})

	}

	return c.JSON(http.StatusOK, results)
}

// GetOneById godoc
// @Summary Show a tenant info
// @Description get tenant by id
// @Tags tenants
// @ID get-tenant-by-id
// @Accept  json
// @Produce  json
// @Param id path string true "Tenant Id"
// @Success 200 {object} handlers.resultJson
// @Failure 400 {object} handlers.errorResult
// @Failure 404 {object} handlers.errorResult
// @Router /tenants/{id} [get]
func (h handler) GetOneById(c echo.Context) error {
	tenantId, err := libUuid.Parse(c.Param("id"))
	if err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	h.tenantModel.Uuid = tenantId

	tenant, err := h.tenantModel.GetOne()
	if err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusNotFound, nil)
	}

	return c.JSON(http.StatusOK, &resultJson{
		ID:   tenant.Uuid,
		Name: tenant.Name,
	})
}

// Create godoc
// @Summary Create a tenant
// @Description create by json tenant
// @Tags tenants
// @Accept  json
// @Produce  json
// @Param tenant body handlers.postTenantData true "Add tenant"
// @Success 201 {object} handlers.resultJson
// @Failure 400 {object} handlers.errorResult
// @Failure 500 {object} handlers.errorResult
// @Router /tenants [post]
func (h handler) Create(c echo.Context) error {
	post := new(postTenantData)
	if err := c.Bind(post); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if err := c.Validate(post); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	h.tenantModel.Uuid = post.ID
	h.tenantModel.Name = post.Name

	tenant, err := h.tenantModel.Save()
	if err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, resultJson{
		ID:   tenant.Uuid,
		Name: tenant.Name,
	})
}

// Update godoc
// @Summary Update a tenant
// @Description update by json tenant
// @Tags tenants
// @Accept  json
// @Produce  json
// @Param tenant body handlers.postTenantData true "Update tenant"
// @Success 200 {object} handlers.resultJson
// @Failure 400 {object} handlers.errorResult
// @Failure 500 {object} handlers.errorResult
// @Router /tenants [put]
func (h handler) Update(c echo.Context) error {
	post := new(postTenantData)

	if err := c.Bind(post); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if err := c.Validate(post); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	h.tenantModel.Uuid = post.ID
	h.tenantModel.Name = post.Name

	tenant, err := h.tenantModel.Update()
	if err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, resultJson{
		ID:   tenant.Uuid,
		Name: tenant.Name,
	})
}

// DeleteById godoc
// @Summary Delete tenant
// @Description delete tenant by id
// @Tags tenants
// @ID delete-tenant-by-id
// @Produce  json
// @Param id path string true "Tenant Id"
// @Success 200 {object} handlers.resultJson
// @Failure 404 {object} handlers.errorResult
// @Router /tenants/{id} [delete]
func (h handler) DeleteById(c echo.Context) error {
	tenantId, err := libUuid.Parse(c.Param("id"))
	if err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusNotFound, err.Error())
	}

	h.tenantModel.Uuid = tenantId

	isDeleted, err := h.tenantModel.Delete()
	if err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, isDeleted)
}
