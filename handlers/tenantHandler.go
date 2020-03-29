package handlers

import (
	libUuid "github.com/google/uuid"
	"log"
	"net/http"

	"github.com/labstack/echo"
	tenantModel "gitlab.com/s0j0hn/go-rest-boilerplate-echo/db/models/tenant"
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

func (h handler) GetOneById(c echo.Context) error {
	tenantId, err := libUuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	h.tenantModel.Uuid = tenantId

	tenant, err := h.tenantModel.GetOne()
	if err != nil {
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
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if err := c.Validate(post); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	h.tenantModel.Uuid = post.ID
	h.tenantModel.Name = post.Name

	log.Print(h.tenantModel)

	tenant, err := h.tenantModel.Save()
	if err != nil {
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
		return c.JSON(http.StatusOK, err.Error())
	}

	if err := c.Validate(post); err != nil {
		return c.JSON(http.StatusOK, err.Error())
	}

	h.tenantModel.Uuid = post.ID
	h.tenantModel.Name = post.Name

	tenant, err := h.tenantModel.Update()
	if err != nil {
		return c.JSON(http.StatusOK, err.Error())
	}

	return c.JSON(http.StatusOK, resultJson{
		ID:   tenant.Uuid,
		Name: tenant.Name,
	})
}

func (h handler) DeleteById(c echo.Context) error {
	tenantId, err := libUuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusOK, err.Error())
	}

	h.tenantModel.Uuid = tenantId

	isDeleted, err := h.tenantModel.Delete()
	if err != nil {
		return c.JSON(http.StatusOK, err.Error())
	}
	return c.JSON(http.StatusOK, isDeleted)
}
