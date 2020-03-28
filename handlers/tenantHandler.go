package handlers

import (
	uuid "github.com/satori/go.uuid"
	"net/http"

	"github.com/labstack/echo"
	tenantModel "github.com/tkc/go-echo-server-sandbox/db/models/tenant"
)

type (
	handler struct {
		tenantModel tenantModel.TenantModel
	}
	resultJson struct {
		id   uuid.UUID
		name string
	}
	postTenantData struct {
		ID   uuid.UUID `json:"id" form:"id" validate:"required"`
		Name string    `json:"name" form:"name" validate:"required"`
	}
)

func CreateHandler(u tenantModel.TenantModel) *handler {
	return &handler{u}
}

func (h handler) GetOneById(c echo.Context) error {
	tenantId, err := uuid.FromString(c.Param("ID"))
	if err != nil {
		return c.JSON(http.StatusOK, err.Error())
	}
	
	user := h.tenantModel.GetOne(tenantId)
	
	return c.JSON(http.StatusOK, resultJson{
		id:   tenantId,
		name: user.Name,
	})
}

func (h handler) Create(c echo.Context) error {
	post := new(postTenantData)
	if err := c.Bind(post); err != nil {
		return c.JSON(http.StatusOK, err.Error())
	}
	if err := c.Validate(post); err != nil {
		return c.JSON(http.StatusOK, err.Error())
	}
	user := h.tenantModel.Create(post.Name)
	return c.JSON(http.StatusOK, user.ID)
}

func (h handler) Update(c echo.Context) error {
	post := new(postTenantData)
	if err := c.Bind(post); err != nil {
		return c.JSON(http.StatusOK, err.Error())
	}
	if err := c.Validate(post); err != nil {
		return c.JSON(http.StatusOK, err.Error())
	}
	user := h.tenantModel.Create(post.Name)
	return c.JSON(http.StatusOK, user.ID)
}

func (h handler) DeleteById(c echo.Context) error {
	tenantId, err := uuid.FromString(c.Param("ID"))
	if err != nil {
		return c.JSON(http.StatusOK, err.Error())
	}

	h.tenantModel.Delete(tenantId)
	return c.JSON(http.StatusOK, true)
}
