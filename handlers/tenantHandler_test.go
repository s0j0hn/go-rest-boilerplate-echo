package handlers

import (
	libUuid "github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
	"gitlab.com/s0j0hn/go-rest-boilerplate-echo/database"
	tenantModel "gitlab.com/s0j0hn/go-rest-boilerplate-echo/database/models/tenant"
	"gopkg.in/go-playground/validator.v9"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var (
	mockDBTenant             = tenantModel.TenantModel{}
	createTenantString       = `{"id":"39b0b2fc-749f-46f3-8960-453418e72b2e","name":"NAME"}`
	createTenantStringName       = `{"id":"39b0b2fc-749f-46f3-8960-453418e72b3e","name":"NAME"}`
	allTenantsString         = `[{"id":"39b0b2fc-749f-46f3-8960-453418e72b2e","name":"NAME"}]`
	updatedTenantString      = `{"id":"39b0b2fc-749f-46f3-8960-453418e72b2e","name":"NAME2"}`
	updatedWrongTenantString = `{"id":"yolo","name":"NAME2"}`
	updatedWrongTenantString2 = `{"id":"39b0b2fc-749f-46f3-8960-453418e72b2e","name":""}`
	updatedWrongTenantString3 = `{"name":"TEST"}`
	createWrongTenantString  = `{"id":"yolo","name":111}`
	validTenantId		     = "39b0b2fc-749f-46f3-8960-453418e72b2e"
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
	DbClient = database.ConnectForTests()
	os.Exit(m.Run())
}

func refreshTenantTable(t *testing.T) {
	err := DbClient.Exec("DROP TABLE IF EXISTS tenants").Error
	if err != nil {
		t.Errorf("Error drop tenants handler: %v\n", err)
		return
	}

	DbClient.AutoMigrate(&tenantModel.TenantModel{})
}

func TestCreateTenant(t *testing.T) {
	refreshTenantTable(t)
	e := echo.New()
	e.Use(middleware.Logger())
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(createTenantString))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/tenants")
	h := &handler{mockDBTenant}

	// Assertions
	if assert.NoError(t, h.Create(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, createTenantString+"\n", rec.Body.String())
	}

	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(createWrongTenantString))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetPath("/tenants")
	h = &handler{mockDBTenant}

	// Assertions
	if assert.NoError(t, h.Create(c)) {
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, "\"code=400, message=invalid UUID length: 4, internal=invalid UUID length: 4\"\n", rec.Body.String())
	}

	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(createTenantString))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetPath("/tenants")
	h = &handler{mockDBTenant}

	// Assertions
	if assert.NoError(t, h.Create(c)) {
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, "\"UNIQUE constraint failed: tenant.uuid\"\n", rec.Body.String())
	}

	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(createTenantStringName))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetPath("/tenants")
	h = &handler{mockDBTenant}

	// Assertions
	if assert.NoError(t, h.Create(c)) {
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, "\"UNIQUE constraint failed: tenant.name\"\n", rec.Body.String())
	}
}

func TestGetTenant(t *testing.T) {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/tenants/:id")
	c.SetParamNames("id")
	c.SetParamValues(validTenantId)
	h := &handler{mockDBTenant}

	// Assertions
	if assert.NoError(t, h.GetOneById(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetPath("/tenants/:id")
	c.SetParamNames("id")
	c.SetParamValues("yolo")
	h = &handler{mockDBTenant}

	// Assertions
	if assert.NoError(t, h.GetOneById(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "\"invalid UUID length: 4\"\n", rec.Body.String())
	}

	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetPath("/tenants/:id")
	c.SetParamNames("id")
	c.SetParamValues(libUuid.New().String())
	h = &handler{mockDBTenant}

	// Assertions
	if assert.NoError(t, h.GetOneById(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, "null\n", rec.Body.String())
	}
}

func TestGetAllTenants(t *testing.T) {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/tenants")
	h := &handler{mockDBTenant}

	// Assertions
	if assert.NoError(t, h.GetAll(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, allTenantsString+"\n", rec.Body.String())
	}
}

func TestUpdateTenant(t *testing.T) {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(updatedTenantString))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/tenants")
	h := &handler{mockDBTenant}

	// Assertions
	if assert.NoError(t, h.Update(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, updatedTenantString+"\n", rec.Body.String())
	}

	req = httptest.NewRequest(http.MethodPut, "/", strings.NewReader(updatedWrongTenantString))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetPath("/tenants")
	h = &handler{mockDBTenant}

	// Assertions
	if assert.NoError(t, h.Update(c)) {
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, "\"code=400, message=invalid UUID length: 4, internal=invalid UUID length: 4\"\n", rec.Body.String())
	}

	req = httptest.NewRequest(http.MethodPut, "/", strings.NewReader(updatedWrongTenantString))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetPath("/tenants")
	h = &handler{mockDBTenant}

	// Assertions
	if assert.NoError(t, h.Update(c)) {
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, "\"code=400, message=invalid UUID length: 4, internal=invalid UUID length: 4\"\n", rec.Body.String())
	}

	req = httptest.NewRequest(http.MethodPut, "/", strings.NewReader(updatedWrongTenantString2))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetPath("/tenants")
	h = &handler{mockDBTenant}

	// Assertions
	if assert.NoError(t, h.Update(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "\"Key: 'postTenantData.Name' Error:Field validation for 'Name' failed on the 'required' tag\"\n", rec.Body.String())
	}

	req = httptest.NewRequest(http.MethodPut, "/", strings.NewReader(updatedWrongTenantString3))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetPath("/tenants")
	h = &handler{mockDBTenant}

	// Assertions
	if assert.NoError(t, h.Update(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "\"Key: 'postTenantData.ID' Error:Field validation for 'ID' failed on the 'required' tag\"\n", rec.Body.String())
	}
}

func TestCreateHandler(t *testing.T) {
	h := CreateHandler(mockDBTenant)
	assert.NotNil(t, h)
	assert.NotNil(t, h.tenantModel)
}

func TestDeleteTenant(t *testing.T) {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/tenants/:id")
	c.SetParamNames("id")
	c.SetParamValues(validTenantId)
	h := &handler{mockDBTenant}

	// Assertions
	if assert.NoError(t, h.DeleteById(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "true\n", rec.Body.String())
	}

	req = httptest.NewRequest(http.MethodDelete, "/", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetPath("/tenants/:id")
	c.SetParamNames("id")
	c.SetParamValues("yolo")
	h = &handler{mockDBTenant}

	// Assertions
	if assert.NoError(t, h.DeleteById(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, "\"invalid UUID length: 4\"\n", rec.Body.String())
	}

	req = httptest.NewRequest(http.MethodDelete, "/", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetPath("/tenants/:id")
	c.SetParamNames("id")
	c.SetParamValues(libUuid.New().String())
	h = &handler{mockDBTenant}

	// Assertions
	if assert.NoError(t, h.DeleteById(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "false\n", rec.Body.String())
	}
}

func TestGetAllNoTenants(t *testing.T) {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/tenants")
	h := &handler{mockDBTenant}

	// Assertions
	if assert.NoError(t, h.GetAll(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "[]\n", rec.Body.String())
	}
}