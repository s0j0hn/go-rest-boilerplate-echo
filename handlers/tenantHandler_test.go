package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	userModel "github.com/tkc/go-echo-server-sandbox/db/models/tenant"
)

var (
	userJSON = `{"Name":"post Name","Uuid":"27"}`
)

func TestGetUser(t *testing.T) {
	e := echo.New()
	req := new(http.Request)
	rec := httptest.NewRecorder()
	c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))
	c.SetPath("/users/:ID")
	c.SetParamNames("ID")
	c.SetParamValues("10")

	u := userModel.TenantModel{}
	h := CreateHandler(u)

	if assert.NoError(t, h.GetUser(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		// t.Log(rec.Body.String())
	}
}

func TestCreateUser(t *testing.T) {

	e := echo.New()
	req, err := http.NewRequest(echo.POST, "/tenant", strings.NewReader(userJSON))
	if assert.NoError(t, err) {
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))

		u := userModel.TenantModel{}
		h := CreateHandler(u)

		if assert.NoError(t, h.CreateUser(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			//t.Log(rec.Body.String())
		}
	}
}
