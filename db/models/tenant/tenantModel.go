package tenantModel

import (
	"errors"
	"github.com/jinzhu/gorm"
	libUuid "github.com/satori/go.uuid"
	databaseManager "gitlab.com/s0j0hn/go-rest-boilerplate-echo/db"
	"time"
)

type TenantModel struct {
	ID        libUuid.UUID  `gorm:"primary_key; unique; 
			    type:uuid; column:id; 
			    default:uuid_generate_v4()`
	Name      string       `gorm:"unique;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
func (TenantModel) TableName() string {
	return "tenant"
}

func (tenantModel *TenantModel) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("ID", libUuid.NewV4().String())
}

func (tenantModel *TenantModel) GetAll(limit int) []TenantModel {
	var users []TenantModel
	databaseManager.Client.Limit(limit).Order("id desc").Find(&users)
	return users
}

func (tenantModel *TenantModel) Save() (*TenantModel, error) {
	transation := databaseManager.Client.Begin()

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
	transation := databaseManager.Client.Begin()

	if transation.Error != nil {
		return &TenantModel{}, transation.Error
	}

	err := transation.Update(tenantModel.ID, tenantModel.Name).Error
	if err != nil {
		transation.Rollback()
		return &TenantModel{}, err
	}

	transation.Commit()
	return tenantModel, nil
}

func (tenantModel *TenantModel) GetOne() (*TenantModel, error) {
	err := databaseManager.Client.First(&tenantModel).Error
	if err != nil {
		return &TenantModel{}, err
	}

	if gorm.IsRecordNotFoundError(err) {
		return &TenantModel{}, errors.New("Tenant not found in database")
	}

	return tenantModel, nil
}

func (tenantModel *TenantModel) Delete() (bool, error) {
	transation := databaseManager.Client.Begin()

	if transation.Error != nil {
		return false, transation.Error
	}

	err := transation.Delete(&tenantModel).Error
	if err != nil {
		transation.Rollback()
		return false, err
	}

	transation.Commit()
	return true, nil
}
