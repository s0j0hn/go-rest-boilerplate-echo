package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Tenant represents the business domain model for tenant
type Tenant struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
	UUID      uuid.UUID      `json:"uuid" gorm:"unique;not null;type:uuid;default:gen_random_uuid()"`
	Name      string         `json:"name" gorm:"unique;not null;size:100" validate:"required,min=2,max=100"`
}

// TableName specifies the database table name
func (Tenant) TableName() string {
	return "tenant"
}

// BeforeCreate generates UUID if not provided
func (t *Tenant) BeforeCreate(tx *gorm.DB) error {
	if t.UUID == uuid.Nil {
		t.UUID = uuid.New()
	}
	return nil
}

// TenantCreateRequest represents the request to create a tenant
type TenantCreateRequest struct {
	Name string `json:"name" validate:"required,min=2,max=100"`
}

// TenantUpdateRequest represents the request to update a tenant
type TenantUpdateRequest struct {
	ID   uuid.UUID `json:"id" validate:"required"`
	Name string    `json:"name" validate:"required,min=2,max=100"`
}

// TenantResponse represents the response when returning tenant data
type TenantResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// ToResponse converts domain model to response model
func (t *Tenant) ToResponse() TenantResponse {
	return TenantResponse{
		ID:   t.UUID,
		Name: t.Name,
	}
}