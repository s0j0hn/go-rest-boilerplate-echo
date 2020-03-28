package handlers

import (
	uuid "github.com/satori/go.uuid"
	"net/http"

	"github.com/labstack/echo"
	tenantModel "github.com/tkc/go-echo-server-sandbox/models/tenant"
	"github.com/google/uuid"
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
	userId, err := uuid.FromString(c.Param("ID"))
	if err != nil {
		return c.JSON(http.StatusOK, err.Error())
	}
	
	user := h.tenantModel.GetOne(userId)
	
	return c.JSON(http.StatusOK, resultJson{
		id:   userId,
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
	id := uuid.Must(c.Param("ID"))
	h.tenantModel.Delete(id)
	return c.JSON(http.StatusOK, true)
}
