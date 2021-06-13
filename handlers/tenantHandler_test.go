package handlers

import (
	"encoding/json"
	"github.com/NeowayLabs/wabbit/amqptest/server"
	"github.com/go-playground/validator/v10"
	libUuid "github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"gitlab.com/s0j0hn/go-rest-boilerplate-echo/database"
	tenantModel "gitlab.com/s0j0hn/go-rest-boilerplate-echo/database/models/tenant"
	"gitlab.com/s0j0hn/go-rest-boilerplate-echo/rabbitmq"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var (
	mockDBTenant              = tenantModel.Model{Name: ""}
	createTenantString        = `{"id":"39b0b2fc-749f-46f3-8960-453418e72b2e","name":"NAME"}`
	allTenantsString          = `[{"id":"39b0b2fc-749f-46f3-8960-453418e72b2e","name":"NAME"}]`
	updatedTenantString       = `{"id":"39b0b2fc-749f-46f3-8960-453418e72b2e","name":"NAME2"}`
	updatedWrongTenantString  = `{"id":"yolo","name":"NAME2"}`
	updatedWrongTenantString2 = `{"id":"39b0b2fc-749f-46f3-8960-453418e72b2e","name":""}`
	updatedWrongTenantString3 = `{"name":"TEST"}`
	createWrongTenantString   = `{"id":"yolo","name":111}`
	validTenantID             = "39b0b2fc-749f-46f3-8960-453418e72b2e"
)

type (
	CustomValidator struct {
		validator *validator.Validate
	}
)

var DbClient *gorm.DB
var TaskManager *rabbitmq.TaskClient = nil
var ZeroLogger zerolog.Logger


func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return nil
}


func TestMain(m *testing.M) {
	fakeServer := server.NewServer("amqp://127.0.0.1:5672/%2f")
	err := fakeServer.Start()
	if err != nil {
		panic(err)
	}

	ZeroLogger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	doneChannel := make(chan bool)
	messagesChannel := make(chan []byte)
	rbbtMQClient := rabbitmq.NewAMQPClient("testQueue", "testQueue", "amqp://127.0.0.1:5672/%2f", ZeroLogger, doneChannel, messagesChannel, false)
	doneChannel <- true

	TaskManager = rabbitmq.NewTaskManagerClient(rbbtMQClient)

	DbClient = database.ConnectForTests()

	returnCode := m.Run()
	os.Exit(returnCode)
}

func refreshTenantTable(t *testing.T) {
	ZeroLogger.Printf("Reset table")
	err := DbClient.Exec("DROP TABLE IF EXISTS tenant").Error
	if err != nil {
		t.Errorf("Error drop tenants models: %v\n", err)
		return
	}


	err = DbClient.AutoMigrate(&tenantModel.Model{})
	if err != nil {
		t.Errorf("Error migrate tenants models: %v\n", err)
		return
	}
	ZeroLogger.Printf("Reset done")

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
	h := &handlerTenant{mockDBTenant, TaskManager}

	// Assertions
	if assert.NoError(t, h.Create(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		var response ResultTask
		assert.NoError(t, json.Unmarshal([]byte(rec.Body.String()), &response))
		assert.NotNil(t, response.TaskID)
	}

	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(createWrongTenantString))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetPath("/tenants")
	h = &handlerTenant{mockDBTenant, TaskManager}

	// Assertions
	if assert.NoError(t, h.Create(c)) {
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, "\"code=400, message=Unmarshal type error: expected=string, got=number, field=name, offset=23, internal=json: cannot unmarshal number into Go struct field tenantData.name of type string\"\n", rec.Body.String())
	}

	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(createTenantString))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetPath("/tenants")
	h = &handlerTenant{mockDBTenant, TaskManager}

	// Assertions
	if assert.NoError(t, h.Create(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		var response ResultTask
		assert.NoError(t, json.Unmarshal([]byte(rec.Body.String()), &response))
		assert.NotNil(t, response.TaskID)
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
	c.SetParamValues(validTenantID)
	h := &handlerTenant{mockDBTenant, TaskManager}

	// Assertions
	if assert.NoError(t, h.GetOneByID(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetPath("/tenants/:id")
	c.SetParamNames("id")
	c.SetParamValues("yolo")
	h = &handlerTenant{mockDBTenant, TaskManager}

	// Assertions
	if assert.NoError(t, h.GetOneByID(c)) {
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
	h = &handlerTenant{mockDBTenant, TaskManager}

	// Assertions
	if assert.NoError(t, h.GetOneByID(c)) {
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
	h := &handlerTenant{mockDBTenant, TaskManager}

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
	h := &handlerTenant{mockDBTenant, TaskManager}

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
	h = &handlerTenant{mockDBTenant, TaskManager}

	// Assertions
	if assert.NoError(t, h.Update(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "\"code=500, message=Key: 'tenantData.ID' Error:Field validation for 'ID' failed on the 'uuid4' tag\"\n", rec.Body.String())
	}

	req = httptest.NewRequest(http.MethodPut, "/", strings.NewReader(updatedWrongTenantString))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetPath("/tenants")
	h = &handlerTenant{mockDBTenant, TaskManager}

	// Assertions
	if assert.NoError(t, h.Update(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "\"code=500, message=Key: 'tenantData.ID' Error:Field validation for 'ID' failed on the 'uuid4' tag\"\n", rec.Body.String())
	}

	req = httptest.NewRequest(http.MethodPut, "/", strings.NewReader(updatedWrongTenantString2))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetPath("/tenants")
	h = &handlerTenant{mockDBTenant, TaskManager}

	// Assertions
	if assert.NoError(t, h.Update(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "\"code=500, message=Key: 'tenantData.Name' Error:Field validation for 'Name' failed on the 'required' tag\"\n", rec.Body.String())
	}

	req = httptest.NewRequest(http.MethodPut, "/", strings.NewReader(updatedWrongTenantString3))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetPath("/tenants")
	h = &handlerTenant{mockDBTenant, TaskManager}

	// Assertions
	if assert.NoError(t, h.Update(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "\"code=500, message=Key: 'tenantData.ID' Error:Field validation for 'ID' failed on the 'required' tag\"\n", rec.Body.String())
	}
}

func TestCreateHandler(t *testing.T) {
	h := CreateHandlerTenant(mockDBTenant, TaskManager)
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
	c.SetParamValues(validTenantID)
	h := &handlerTenant{mockDBTenant, TaskManager}

	// Assertions
	if assert.NoError(t, h.DeleteByID(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "true\n", rec.Body.String())
	}

	req = httptest.NewRequest(http.MethodDelete, "/", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetPath("/tenants/:id")
	c.SetParamNames("id")
	c.SetParamValues("yolo")
	h = &handlerTenant{mockDBTenant, TaskManager}

	// Assertions
	if assert.NoError(t, h.DeleteByID(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, "\"invalid UUID length: 4\"\n", rec.Body.String())
	}

	req = httptest.NewRequest(http.MethodDelete, "/", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetPath("/tenants/:id")
	c.SetParamNames("id")
	c.SetParamValues(libUuid.New().String())
	h = &handlerTenant{mockDBTenant, TaskManager}

	// Assertions
	if assert.NoError(t, h.DeleteByID(c)) {
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
	h := &handlerTenant{mockDBTenant, TaskManager}

	// Assertions
	if assert.NoError(t, h.GetAll(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "[]\n", rec.Body.String())
	}
}
