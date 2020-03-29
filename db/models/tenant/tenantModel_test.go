package tenantModel

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	libUuid "github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gitlab.com/s0j0hn/go-rest-boilerplate-echo/db"
	"log"
	"os"
	"testing"
)

var DbClient *gorm.DB

func TestMain(m *testing.M) {
	DbClient = db.DatabaseConnect()

	os.Exit(m.Run())
}

func refreshTenantTable() error {
	err := DbClient.DropTableIfExists(&TenantModel{}).Error
	if err != nil {
		return err
	}

	err = DbClient.AutoMigrate(&TenantModel{}).Error
	if err != nil {
		return err
	}

	return nil
}

func seedTenants() error {
	tenants := []TenantModel{
		{
			Name: "Bob",
		},
		{
			Name: "Alice",
		},
	}

	for i, _ := range tenants {
		tenant, err := tenants[i].Save()
		if err != nil {
			return err
		}
		log.Printf("Seed with Tenant ID: %s", tenant.Uuid)
	}
	return nil
}

func seedOneTenant() *TenantModel {
	err := refreshTenantTable()
	if err != nil {
		log.Fatal(err)
	}

	tenant := TenantModel{
		Name: "Greg",
		Uuid: libUuid.MustParse("39b0b2fc-749f-46f3-8960-453418e72b2e"),
	}

	tenantSaved, err := tenant.Save()
	if err != nil {
		log.Fatalf("Cannot seed tenant table: %v", err)
	}
	log.Printf("Seed with Tenant ID: %s", tenantSaved.Uuid)
	return nil
}

func TestGetAllTenants(t *testing.T) {
	err := refreshTenantTable()
	if err != nil {
		log.Fatal(err)
	}

	err = seedTenants()
	if err != nil {
		log.Fatal(err)
	}

	tenantInstance := TenantModel{}

	tenants := tenantInstance.GetAll(2)
	assert.Equal(t, len(tenants), 2)
	log.Printf("End TestgetAllTenants")
}

func TestSaveTenant(t *testing.T) {
	err := refreshTenantTable()
	if err != nil {
		log.Fatal(err)
	}

	newUser := TenantModel{
		Uuid:     libUuid.New(),
		Name:     "Test",
	}

	savedUser, err := newUser.Save()
	if err != nil {
		t.Errorf("this is the error getting the users: %v\n", err)
		return
	}

	assert.Equal(t, newUser.ID, savedUser.ID)
	assert.Equal(t, newUser.Name, savedUser.Name)
	assert.Equal(t, newUser.Uuid, savedUser.Uuid)
	log.Printf("End TestSaveTenant")
}

func TestGetTenantByID(t *testing.T) {
	err := refreshTenantTable()
	if err != nil {
		log.Fatal(err)
	}

	seedOneTenant()
	tenantInstance := TenantModel{Uuid: libUuid.MustParse("39b0b2fc-749f-46f3-8960-453418e72b2e")}

	foundTenant, err := tenantInstance.GetOne()
	if err != nil {
		t.Errorf("Error getting one tenant: %v\n", err)
		return
	}
	assert.Equal(t, foundTenant.Uuid.String(), "39b0b2fc-749f-46f3-8960-453418e72b2e")
	log.Printf("End TestTenantGetById")
}

func TestUpdateTenant(t *testing.T) {
	err := refreshTenantTable()
	if err != nil {
		log.Fatal(err)
	}

	seedOneTenant()

	newTenant := TenantModel{
		Uuid:     libUuid.MustParse("39b0b2fc-749f-46f3-8960-453418e72b2e"),
		Name:	  "Gregory",
	}

	updatedTenant, err := newTenant.Update()
	if err != nil {
		t.Errorf("Error test updating the tenant: %v\n", err)
		return
	}


	assert.Equal(t, updatedTenant.ID, newTenant.ID)
	assert.Equal(t, updatedTenant.Name, newTenant.Name)
	log.Printf("End TestUpdateTenant")

}

func TestDeleteTenant(t *testing.T) {
	err := refreshTenantTable()
	if err != nil {
		log.Fatal(err)
	}

	seedOneTenant()
	tenantInstance := TenantModel{
		Uuid:     libUuid.MustParse("39b0b2fc-749f-46f3-8960-453418e72b2e"),
	}

	isDeleted, err := tenantInstance.Delete()
	if err != nil {
		t.Errorf("Error test deleting the tenant: %v\n", err)
		return
	}

	assert.Equal(t, isDeleted, true)
	log.Printf("End TestDeleteTenant")

}