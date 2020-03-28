package base

import (
	"github.com/jinzhu/gorm"
	libUuid "github.com/satori/go.uuid"
	"time"
)

type BaseModel struct {
	ID         string     `sql:"type:uuid;primary_key;default:uuid_generate_v4()"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updateAt"`
	DeletedAt *time.Time `sql:"index" json:"deletedAt"`
}

func (base *BaseModel) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("ID", libUuid.NewV4().String())
}