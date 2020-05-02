package tenantModel

import (
	"errors"
	libUuid "github.com/google/uuid"
	"github.com/jinzhu/gorm"
	databaseManager "gitlab.com/s0j0hn/go-rest-boilerplate-echo/database"
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

func (tenantModel *TenantModel) GetAll() (*[]TenantModel, error) {
	var tenants []TenantModel
	err := databaseManager.Connect().Find(&tenants).Error
	if err != nil {
		return nil, err
	}
	return &tenants, nil
}

func (tenantModel *TenantModel) Save() (*TenantModel, error) {
	transation := databaseManager.Connect().Begin()

	if transation.Error != nil {
		return nil, transation.Error
	}

	err := transation.Create(&tenantModel).Error
	if err != nil {
		transation.Rollback()
		return nil, err
	}

	transation.Commit()
	return tenantModel, nil
}

func (tenantModel *TenantModel) Update() (*TenantModel, error) {
	transaction := databaseManager.Connect().Begin()

	if transaction.Error != nil {
		return &TenantModel{}, transaction.Error
	}

	err := transaction.Model(&tenantModel).Update(TenantModel{ Name: tenantModel.Name, Uuid: tenantModel.Uuid}).Error
	if err != nil {
		transaction.Rollback()
		return nil, err
	}

	transaction.Commit()
	return tenantModel, nil
}

func (tenantModel *TenantModel) GetOne() (*TenantModel, error) {
	err := databaseManager.Connect().First(&tenantModel).Error
	if err != nil {
		return nil, err
	}

	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("tenant not found in database")
	}

	return tenantModel, nil
}

func (tenantModel *TenantModel) Delete() (bool, error) {
	transaction := databaseManager.Connect().Begin()

	if transaction.Error != nil {
		return false, transaction.Error
	}

	err := transaction.First(&tenantModel).Unscoped().Delete(&tenantModel).Error
	if err != nil {
		transaction.Rollback()
		return false, err
	}

	transaction.Commit()
	return true, nil
}
