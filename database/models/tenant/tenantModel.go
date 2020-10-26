package tenantmodel

import (
	"errors"
	libUuid "github.com/google/uuid"
	databaseManager "gitlab.com/s0j0hn/go-rest-boilerplate-echo/database"
	"gorm.io/gorm"
)

// TenantModel is a tenant model description.
type TenantModel struct {
	gorm.Model
	UUID libUuid.UUID `gorm:"unique_index;not null"`
	Name string       `gorm:"unique;not null;type:varchar(100);default:null"`
}

// TableName used to set the table name.
func (TenantModel) TableName() string {
	return "tenant"
}

// BeforeCreate used to transform some params before saving to database.
func (tenantModel *TenantModel) BeforeCreate(ctx *gorm.DB) (err error) {
	if tenantModel.UUID.String() == "00000000-0000-0000-0000-000000000000" {
		err = ctx.Statement.Set("UUID", libUuid.New()).Error
		if err != nil {
			return err
		}
		return
	}
	return
}

// GetAll is used to get all elements for database.
func (tenantModel *TenantModel) GetAll() (*[]TenantModel, error) {
	var tenants []TenantModel
	err := databaseManager.Connect().Find(&tenants).Error
	if err != nil {
		return nil, err
	}
	return &tenants, nil
}

// Save is used to write data into database.
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
// Update is used to write data into database.
func (tenantModel *TenantModel) Update() (*TenantModel, error) {
	transaction := databaseManager.Connect().Begin()

	if transaction.Error != nil {
		return nil, transaction.Error
	}

	err := transaction.Model(&tenantModel).Where(TenantModel{UUID: tenantModel.UUID}).Updates(&tenantModel).Error
	if err != nil {
		transaction.Rollback()
		return nil, err
	}

	transaction.Commit()
	return tenantModel, nil
}

// GetOne is used to retrieve element from database.
func (tenantModel *TenantModel) GetOne() (*TenantModel, error) {
	err := databaseManager.Connect().Where(&TenantModel{UUID: tenantModel.UUID}).First(&tenantModel).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("tenant not found in database")
	}

	return tenantModel, nil
}

// Delete is used to drop data from database.
func (tenantModel *TenantModel) Delete() (bool, error) {
	libUuid.MustParse(tenantModel.UUID.String())

	if tenantModel.UUID.String() == "00000000-0000-0000-0000-000000000000" {
		return false, errors.New("no uuid specified")
	}

	err := databaseManager.Connect().First(&tenantModel).Where(&TenantModel{UUID: tenantModel.UUID}).Error

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
