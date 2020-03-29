package handlers

import (
	libUuid "github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
	"gitlab.com/s0j0hn/go-rest-boilerplate-echo/db"
	tenantModel "gitlab.com/s0j0hn/go-rest-boilerplate-echo/db/models/tenant"
	"gopkg.in/go-playground/validator.v9"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var (
	mockDBTenant = tenantModel.TenantModel{ Name: "NAME", Uuid: libUuid.MustParse("39b0b2fc-749f-46f3-8960-453418e72b2e")}
	tenantJSON   = `{"id":"39b0b2fc-749f-46f3-8960-453418e72b2e","name":"NAME"}`
)

type (
	CustomValidator struct {
		validator *validator.Validate
	}
)

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

var DbClient *gorm.DB

func TestMain(m *testing.M) {
	DbClient = db.DatabaseConnect()
	err := refreshTenantTable()
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}

func refreshTenantTable() error {
	err := DbClient.DropTableIfExists(&tenantModel.TenantModel{}).Error
	if err != nil {
		return err
	}

	err = DbClient.AutoMigrate(&tenantModel.TenantModel{}).Error
	if err != nil {
		return err
	}

	return nil
}


func TestCreateTenant(t *testing.T) {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tenantJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/tenants")
	h := &handler{mockDBTenant}

	// Assertions
	if assert.NoError(t, h.Create(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}

func TestGetTenant(t *testing.T) {
	// Setup
	e := echo.New()
	e.Use(middleware.Logger())
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/tenants/:id")
	c.SetParamNames("id")
	c.SetParamValues("39b0b2fc-749f-46f3-8960-453418e72b2e")
	h := &handler{mockDBTenant}

	// Assertions
	if assert.NoError(t, h.GetOneById(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestDeleteTenant(t *testing.T) {
	// Setup
	e := echo.New()
	e.Use(middleware.Logger())
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/tenants/:id")
	c.SetParamNames("id")
	c.SetParamValues("39b0b2fc-749f-46f3-8960-453418e72b2e")
	h := &handler{mockDBTenant}

	// Assertions
	if assert.NoError(t, h.DeleteById(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}
