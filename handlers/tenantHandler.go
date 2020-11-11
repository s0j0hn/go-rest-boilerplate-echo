package handlers

import (
	libUUID "github.com/google/uuid"
	"github.com/labstack/echo/v4"
	tenantModel "gitlab.com/s0j0hn/go-rest-boilerplate-echo/database/models/tenant"
	"gitlab.com/s0j0hn/go-rest-boilerplate-echo/rabbitmq"
	"net/http"
)

type (
	// Handler is a default handler as there is no generics.
	Handler struct {
		tenantModel tenantModel.Model
		taskManager *rabbitmq.TaskClient
	}

	resultJSON struct {
		ID   libUUID.UUID `json:"id" form:"id" validate:"required"`
		Name string       `json:"name" form:"name" validate:"required"`
	}

	postTenantData struct {
		ID   libUUID.UUID `json:"id" form:"id" validate:"required"`
		Name string       `json:"name" form:"name" validate:"required"`
	}

	updateTenantData struct {
		ID   libUUID.UUID `json:"id" form:"id" validate:"required"`
		Name string       `json:"name" form:"name" validate:"required"`
	}

	errorResult struct {
		Message string `json:"message" comment:"Something went wrong"`
	}
)

// CreateHandler is always in each Handler
func CreateHandler(tenant tenantModel.Model, taskClient *rabbitmq.TaskClient) *Handler {
	return &Handler{tenant,taskClient}
}

// GetAll godoc
// @Summary List tenants
// @Description get tenants
// @Tags tenants
// @Accept  json
// @Produce  json
// @Success 200 {array} handlers.resultJSON
// @Failure 500 {object} handlers.errorResult
// @Router /tenants [get]
func (h Handler) GetAll(c echo.Context) error {
	tenants, err := h.tenantModel.GetAll()
	if err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusInternalServerError, errorResult{Message: err.Error()})
	}

	var results []resultJSON
	for _, tenant := range *tenants {
		results = append(results, resultJSON{ID: tenant.UUID, Name: tenant.Name})
	}

	if len(results) == 0 {
		return c.JSON(http.StatusOK, []resultJSON{})

	}

	return c.JSON(http.StatusOK, results)
}

// GetOneByID godoc
// @Summary Show a tenant info
// @Description get tenant by id
// @Tags tenants
// @ID get-tenant-by-id
// @Accept  json
// @Produce  json
// @Param id path string true "Tenant ID"
// @Success 200 {object} handlers.resultJSON
// @Failure 400 {object} handlers.errorResult
// @Failure 404 {object} handlers.errorResult
// @Router /tenants/{id} [get]
func (h Handler) GetOneByID(c echo.Context) error {
	tenantID, err := libUUID.Parse(c.Param("id"))
	if err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	h.tenantModel.UUID = tenantID

	tenant, err := h.tenantModel.GetOne()
	if err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusNotFound, nil)
	}

	return c.JSON(http.StatusOK, &resultJSON{
		ID:   tenant.UUID,
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
// @Success 201 {object} handlers.resultJSON
// @Failure 400 {object} handlers.errorResult
// @Failure 500 {object} handlers.errorResult
// @Router /tenants [post]
func (h Handler) Create(c echo.Context) error {
	post := new(postTenantData)
	if err := c.Bind(post); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if err := c.Validate(post); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	taskBytes := rabbitmq.CreateNewTask([]string{"test", "test2"}, "Status is OK")
	err := h.taskManager.PushNewTask(taskBytes)
	if err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	h.tenantModel.UUID = post.ID
	h.tenantModel.Name = post.Name

	tenant, err := h.tenantModel.Save()
	if err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, resultJSON{
		ID:   tenant.UUID,
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
// @Success 200 {object} handlers.resultJSON
// @Failure 400 {object} handlers.errorResult
// @Failure 500 {object} handlers.errorResult
// @Router /tenants [put]
func (h Handler) Update(c echo.Context) error {
	post := new(updateTenantData)

	if err := c.Bind(post); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if err := c.Validate(post); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	h.tenantModel.UUID = post.ID
	h.tenantModel.Name = post.Name

	tenant, err := h.tenantModel.Update()
	if err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, resultJSON{
		ID:   tenant.UUID,
		Name: tenant.Name,
	})
}

// DeleteByID godoc
// @Summary Delete tenant
// @Description delete tenant by id
// @Tags tenants
// @ID delete-tenant-by-id
// @Produce  json
// @Param id path string true "Tenant ID"
// @Success 200 {object} handlers.resultJSON
// @Failure 404 {object} handlers.errorResult
// @Router /tenants/{id} [delete]
func (h Handler) DeleteByID(c echo.Context) error {
	tenantID, err := libUUID.Parse(c.Param("id"))
	if err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusNotFound, err.Error())
	}

	h.tenantModel.UUID = tenantID

	isDeleted, err := h.tenantModel.Delete()
	if err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, isDeleted)
}
