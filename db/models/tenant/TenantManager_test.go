package userModel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	_user := TenantModel{}
	user := _user.Create("test", 30)
	assert.Equal(t, user.Name, "test")
	assert.Equal(t, user.ID, 30)
}

func TestFetch(t *testing.T) {
	_user := TenantModel{}
	user := _user.GetOne(10)
	assert.Equal(t, int(user.ID), 10)
}

func TestAll(t *testing.T) {
	_user := TenantModel{}
	users := _user.GetAll(10, 0)
	for i, user := range users {
		t.Log(i)
		t.Log(user.Name)
		t.Log(user)
	}
}

func TestGetName(t *testing.T) {
	_user := TenantModel{}
	_user.Name = "test"
	name := _user.GetName()
	assert.Equal(t, name, "test")
}

func TestMapByName(t *testing.T) {
	_user := TenantModel{}
	users := _user.MapByName("test")
	for _, v := range users {
		assert.Equal(t, v.Name, "test")
	}
}
