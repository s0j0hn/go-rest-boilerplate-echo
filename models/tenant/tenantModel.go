package userModel

import (
	uuid "github.com/satori/go.uuid"
	databaseManager "github.com/tkc/go-echo-server-sandbox/db"
	. "github.com/tkc/go-echo-server-sandbox/models/base"
	"time"
)

type TenantModel struct {
	BaseModel
	Name      string `gorm:"unique;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (u *TenantModel) GetName() string {
	return u.Name
}

func (u *TenantModel) GetAll(limit int, offset int) []TenantModel {
	var users []TenantModel
	_databaseManager := databaseManager.GetClient()
	_databaseManager.Limit(limit).Offset(offset).Order("id desc").Find(&users)
	return users
}

func (u *TenantModel) Create(name string) TenantModel {
	var tenant TenantModel
	_databaseManager := databaseManager.GetClient()
	tenant = TenantModel{Name: name}
	_databaseManager.Create(&tenant)
	return tenant
}

func (u *TenantModel) Update(id uuid.UUID, name string) TenantModel {
	var tenant TenantModel
	_databaseManager := databaseManager.GetClient()
	tenant = TenantModel{}
	_databaseManager.Model(&tenant).Update(id, name)
	return tenant
}

func (u *TenantModel) GetOne(id uuid.UUID) TenantModel {
	var tenant TenantModel
	_databaseManager := databaseManager.GetClient()
	tenant = TenantModel{}
	_databaseManager.First(&tenant, id)
	return tenant
}

func (u *TenantModel) MapByName(name string) []TenantModel {
	var users []TenantModel
	_databaseManager := databaseManager.GetClient()
	_databaseManager.Where(map[string]interface{}{"name": name}).Find(&users)
	return users
}

func (u *TenantModel) Delete(id uuid.UUID) {
	var tenant TenantModel
	_databaseManager := databaseManager.GetClient()
	tenant = TenantModel{}
	_databaseManager.Delete(&tenant, id)
}
