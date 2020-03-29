package tenantModel

import (
	"errors"
	libUuid "github.com/google/uuid"
	"github.com/jinzhu/gorm"
	databaseManager "gitlab.com/s0j0hn/go-rest-boilerplate-echo/db"
)

type TenantModel struct {
	gorm.Model
	Uuid      libUuid.UUID `gorm:"unique_index;not null"`
	Name      string       `gorm:"unique;not null"`
}

func (TenantModel) TableName() string {
	return "tenant"
}

func (tenantModel *TenantModel) BeforeCreate(scope *gorm.Scope) error {
	if tenantModel.Uuid.String() == "00000000-0000-0000-0000-000000000000" {
		return scope.SetColumn("Uuid", libUuid.New())
	}
	return scope.SetColumn("Uuid", tenantModel.Uuid)
}

func (tenantModel *TenantModel) GetAll(limit int) []TenantModel {
	var tenant []TenantModel
	databaseManager.DatabaseConnect().Limit(limit).Order("id desc").Find(&tenant)
	return tenant
}

func (tenantModel *TenantModel) Save() (*TenantModel, error) {
	transation := databaseManager.DatabaseConnect().Begin()

	if transation.Error != nil {
		return &TenantModel{}, transation.Error
	}

	err := transation.Create(&tenantModel).Error
	if err != nil {
		transation.Rollback()
		return &TenantModel{}, err
	}

	transation.Commit()
	return tenantModel, nil
}

func (tenantModel *TenantModel) Update() (*TenantModel, error) {
	transation := databaseManager.DatabaseConnect().Begin()

	if transation.Error != nil {
		return &TenantModel{}, transation.Error
	}

	err := transation.Update(tenantModel.Uuid, tenantModel.Name).Error
	if err != nil {
		transation.Rollback()
		return &TenantModel{}, err
	}

	transation.Commit()
	return tenantModel, nil
}

func (tenantModel *TenantModel) GetOne() (*TenantModel, error) {
	err := databaseManager.DatabaseConnect().First(&tenantModel).Error
	if err != nil {
		return &TenantModel{}, err
	}

	if gorm.IsRecordNotFoundError(err) {
		return &TenantModel{}, errors.New("Tenant not found in database")
	}

	return tenantModel, nil
}

func (tenantModel *TenantModel) Delete() (bool, error) {
	transation := databaseManager.DatabaseConnect().Begin()

	if transation.Error != nil {
		return false, transation.Error
	}

	err := transation.Unscoped().Delete(&tenantModel).Error
	if err != nil {
		transation.Rollback()
		return false, err
	}

	transation.Commit()
	return true, nil
}
