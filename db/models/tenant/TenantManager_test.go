package tenantModel

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	libUuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"gitlab.com/s0j0hn/go-rest-boilerplate-echo/db"
	"log"
	"os"
	"testing"
)

var DbClient *gorm.DB
var instanceTenant = TenantModel{}

func TestMain(m *testing.M) {
	//err = godotenv.Load(os.ExpandEnv("../../.env"))
	//if err != nil {
	//	log.Fatalf("Error getting env %v\n", err)
	//}
	DbClient = db.InitClient()

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
		TenantModel{
			Name: "Bob",
		},
		TenantModel{
			Name: "Alice",
		},
	}

	for i, _ := range tenants {
		err := DbClient.Model(&TenantModel{}).Create(&tenants[i]).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func seedOneTenant() *TenantModel {
	refreshTenantTable()

	tenant := TenantModel{
		Name: "Greg",
	}

	tenantSaved, err := tenant.Save()
	if err != nil {
		log.Fatalf("Cannot seed tenant table: %v", err)
	}
	return tenantSaved
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
		ID:       libUuid.NewV4(),
		Name:     "Test",
	}

	savedUser, err := newUser.Save()
	if err != nil {
		t.Errorf("this is the error getting the users: %v\n", err)
		return
	}

	assert.Equal(t, newUser.ID, savedUser.ID)
	assert.Equal(t, newUser.Name, savedUser.Name)
	log.Printf("End TestSaveTenant")
}

func TestGetTenantByID(t *testing.T) {

	err := refreshTenantTable()
	if err != nil {
		log.Fatal(err)
	}

	tenant := seedOneTenant()

	foundTenant, err := tenant.GetOne()
	if err != nil {
		t.Errorf("this is the error getting one tenant: %v\n", err)
		return
	}
	assert.Equal(t, foundTenant.ID, tenant.ID)
	assert.Equal(t, foundTenant.Name, tenant.Name)
	log.Printf("End TestTenantGetById")
}

func TestUpdateTenant(t *testing.T) {
	err := refreshTenantTable()
	if err != nil {
		log.Fatal(err)
	}

	tenant := seedOneTenant()

	newTenant := TenantModel{
		ID:       tenant.ID,
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

	tenant := seedOneTenant()

	isDeleted, err := tenant.Delete()
	if err != nil {
		t.Errorf("Error test deleting the tenant: %v\n", err)
		return
	}

	assert.Equal(t, isDeleted, true)
	log.Printf("End TestDeleteTenant")

}