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
)

func CreateHandler(tenant tenantModel.TenantModel) *handler {
	return &handler{tenant}
}

func (h handler) GetAll(c echo.Context) error {
	tenants, err := h.tenantModel.GetAll()
	if err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	var results []resultJson
	for _, tenant := range *tenants {
		results = append(results, resultJson{ID: tenant.Uuid, Name: tenant.Name})
	}
	return c.JSON(http.StatusOK, results)
}

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
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, &resultJson{
		ID:   tenant.Uuid,
		Name: tenant.Name,
	})
}

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
