package tenantModel

import (
	"errors"
	libUuid "github.com/google/uuid"
	databaseManager "gitlab.com/s0j0hn/go-rest-boilerplate-echo/database"
	"gorm.io/gorm"
)

type TenantModel struct {
	gorm.Model
	Uuid libUuid.UUID `gorm:"unique_index;not null"`
	Name string       `gorm:"unique;not null;type:varchar(100);default:null"`
}

func (TenantModel) TableName() string {
	return "tenant"
}

func (tenantModel *TenantModel) BeforeCreate(ctx *gorm.DB) (err error) {
	if tenantModel.Uuid.String() == "00000000-0000-0000-0000-000000000000" {
		err = ctx.Statement.Set("Uuid", libUuid.New()).Error
		if err != nil {
			return err
		}
		return
	}
	return
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
	transaction := databaseManager.Connect().Begin()

	if transaction.Error != nil {
		return nil, transaction.Error
	}

	err := transaction.Create(&tenantModel).Error
	if err != nil {
		transaction.Rollback()
		return nil, err
	}

	transaction.Commit()
	return tenantModel, nil
}

func (tenantModel *TenantModel) Update() (*TenantModel, error) {
	transaction := databaseManager.Connect().Begin()

	if transaction.Error != nil {
		return nil, transaction.Error
	}

	err := transaction.Model(&tenantModel).Where(TenantModel{Uuid: tenantModel.Uuid}).Updates(&tenantModel).Error
	if err != nil {
		transaction.Rollback()
		return nil, err
	}

	transaction.Commit()
	return tenantModel, nil
}

func (tenantModel *TenantModel) GetOne() (*TenantModel, error) {
	err := databaseManager.Connect().Where(&TenantModel{Uuid: tenantModel.Uuid}).First(&tenantModel).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("tenant not found in database")
	}

	return tenantModel, nil
}

func (tenantModel *TenantModel) Delete() (bool, error) {
	libUuid.MustParse(tenantModel.Uuid.String())

	if tenantModel.Uuid.String() == "00000000-0000-0000-0000-000000000000" {
		return false, errors.New("no uuid specified")
	}

	err := databaseManager.Connect().First(&tenantModel).Where(&TenantModel{Uuid: tenantModel.Uuid}).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}

	transaction := databaseManager.Connect().Begin()

	if transaction.Error != nil {
		return false, transaction.Error
	}

	err = transaction.Unscoped().Model(&tenantModel).Delete(&tenantModel).Error
	if err != nil {
		transaction.Rollback()
		return false, err
	}

	transaction.Commit()

	return true, nil
}
