package handlers

import (
	libUUID "github.com/google/uuid"
	"github.com/labstack/echo/v4"
	tenantModel "gitlab.com/s0j0hn/go-rest-boilerplate-echo/database/models/tenant"
	"gitlab.com/s0j0hn/go-rest-boilerplate-echo/rabbitmq"
	"net/http"
)

type (
	// HandlerTenant is a default handler as there is no generics.
	HandlerTenant struct {
		tenantModel tenantModel.ModelTenant
		taskManager *rabbitmq.TaskClient
	}

	resultJSON struct {
		ID   libUUID.UUID `json:"id" validate:"required"`
		Name string       `json:"name" validate:"required"`
	}

	// ResultTask Response given by task processing endpoint
	ResultTask struct {
		TaskID libUUID.UUID `json:"taskId" validate:"required"`
	}

	tenantData struct {
		ID   string `json:"id" validate:"required,uuid4"`
		Name string `json:"name" validate:"required"`
	}

	errorResult struct {
		Message string `json:"message" comment:"Something went wrong"`
	}
)

// CreateHandlerTenant is always in each HandlerTenant
func CreateHandlerTenant(tenant tenantModel.ModelTenant, taskClient *rabbitmq.TaskClient) *HandlerTenant {
	return &HandlerTenant{tenant, taskClient}
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
func (h HandlerTenant) GetAll(c echo.Context) error {
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
func (h HandlerTenant) GetOneByID(c echo.Context) error {
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

	if tenant == nil {
		return c.JSON(http.StatusNotFound, nil)
	}

	if tenant.Name == "" {
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
// @Param tenant body handlers.tenantData true "Add tenant"
// @Success 201 {object} handlers.ResultTask
// @Failure 400 {object} handlers.errorResult
// @Failure 500 {object} handlers.errorResult
// @Router /tenants [post]
func (h HandlerTenant) Create(c echo.Context) error {
	newTenantData := new(tenantData)

	if err := c.Bind(newTenantData); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if err := c.Validate(newTenantData); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	id, err := libUUID.Parse(newTenantData.ID)
	if err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	h.tenantModel.UUID = id
	h.tenantModel.Name = newTenantData.Name

	task := rabbitmq.CreateNewTask([]string{"create", "tenant"}, "Creating tenant "+newTenantData.Name)
	err = h.taskManager.PushTask(task)
	if err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	go func() {
		_, err = h.tenantModel.Save()
		if err != nil {
			c.Logger().Error(err.Error())
			err = h.taskManager.FailTask(task)
			if err != nil {
				c.Logger().Error(err.Error())
			}
		} else {
			err = h.taskManager.CompleteTask(task)
			if err != nil {
				c.Logger().Error(err.Error())
			}
		}
	}()

	return c.JSON(http.StatusCreated, ResultTask{
		TaskID: task.ID,
	})
}

// Update godoc
// @Summary Update a tenant
// @Description update by json tenant
// @Tags tenants
// @Accept  json
// @Produce  json
// @Param tenant body handlers.tenantData true "Update tenant"
// @Success 200 {object} handlers.resultJSON
// @Failure 400 {object} handlers.errorResult
// @Failure 500 {object} handlers.errorResult
// @Router /tenants [put]
func (h HandlerTenant) Update(c echo.Context) error {
	post := new(tenantData)

	if err := c.Bind(post); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if err := c.Validate(post); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	id, err := libUUID.Parse(post.ID)
	if err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	h.tenantModel.UUID = id
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
func (h HandlerTenant) DeleteByID(c echo.Context) error {
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
