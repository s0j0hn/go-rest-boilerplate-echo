package tenantModel

import (
	libUuid "github.com/google/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/stretchr/testify/assert"
	"gitlab.com/s0j0hn/go-rest-boilerplate-echo/database"
	"log"
	"os"
	"testing"
)

var DbClient *gorm.DB

func TestMain(m *testing.M) {
	DbClient = database.ConnectForTests()
	err := refreshTenantTable()
	if err != nil {
		return
	}

	os.Exit(m.Run())
}

func refreshTenantTable() error {
	err := DbClient.DropTableIfExists(&TenantModel{}).Error
	if err != nil {
		log.Fatalf("Cannot refresh tenant table: %v", err)
		return err
	}

	err = DbClient.AutoMigrate(&TenantModel{}).Error
	if err != nil {
		log.Fatalf("Cannot automigrate tenant table: %v", err)
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

func seedOneTenant() {
	err := refreshTenantTable()
	if err != nil {
		log.Fatalf("Cannot seed tenant table: %v", err)
		return
	}

	tenant := TenantModel{
		Name: "Greg",
		Uuid: libUuid.MustParse("6fcec554-9861-4965-bf7d-036be545a92e"),
	}

	tenantSaved, err := tenant.Save()
	if err != nil {
		log.Fatalf("Cannot seed tenant table: %v", err)
		return
	}
	log.Printf("Seed with Tenant ID: %s", tenantSaved.Uuid)
}

func TestGetAllTenants(t *testing.T) {
	err := refreshTenantTable()
	if err != nil {
		t.Fatal(err)
	}

	err = seedTenants()
	if err != nil {
		t.Fatal(err)
	}

	tenantInstance := TenantModel{}

	tenants, err := tenantInstance.GetAll()
	if err != nil {
		t.Errorf("Error getting all tenants: %v\n", err)
		return
	}

	assert.Equal(t, len(*tenants), 2)
	t.Log("End TestgetAllTenants")
}

func TestSaveTenant(t *testing.T) {
	err := refreshTenantTable()
	if err != nil {
		t.Fatal(err)
		return
	}

	newUser := TenantModel{
		Uuid: libUuid.New(),
		Name: "Test",
	}

	savedUser, err := newUser.Save()
	if err != nil {
		t.Errorf("this is the error getting the users: %v\n", err)
		return
	}

	assert.Equal(t, newUser.ID, savedUser.ID)
	assert.Equal(t, newUser.Name, savedUser.Name)
	assert.Equal(t, newUser.Uuid, savedUser.Uuid)
	t.Log("End TestSaveTenant")
}

func TestWrongSaveTenant(t *testing.T) {
	err := refreshTenantTable()
	if err != nil {
		t.Fatal(err)
		return
	}

	newUser := TenantModel{
		Uuid: libUuid.New(),
		Name: "",
	}

	savedUser, err := newUser.Save()
	if assert.Error(t, err) {
		assert.Nil(t, savedUser)
		assert.Equal(t, "name can't be empty", err.Error())
		t.Log("End TestWrongSaveTenant")
	}
}

func TestGetTenantByID(t *testing.T) {
	err := refreshTenantTable()
	if err != nil {
		t.Fatal(err)
		return
	}

	seedOneTenant()
	tenantInstance := TenantModel{Uuid: libUuid.MustParse("6fcec554-9861-4965-bf7d-036be545a92e")}

	foundTenant, err := tenantInstance.GetOne()
	if err != nil {
		t.Errorf("Error getting one tenant: %v\n", err)
		return
	}
	assert.Equal(t, foundTenant.Uuid.String(), "6fcec554-9861-4965-bf7d-036be545a92e")
	t.Log("End TestTenantGetById")
}

func TestGetWrongTenantByID(t *testing.T) {
	err := refreshTenantTable()
	if err != nil {
		t.Fatal(err)
		return
	}

	seedOneTenant()
	tenantInstance := TenantModel{Uuid: libUuid.MustParse("6fcec554-9861-4965-bf7d-036be545a93e")}

	foundTenant, err := tenantInstance.GetOne()

	if assert.Nil(t, foundTenant) {
		assert.Equal(t, "tenant not found in database", err.Error())
		t.Log("End TestTenantGetById")
	}
}

func TestUpdateTenant(t *testing.T) {
	err := refreshTenantTable()
	if err != nil {
		t.Fatal(err)
		return
	}

	seedOneTenant()

	existingTenant := TenantModel{
		Uuid: libUuid.MustParse("6fcec554-9861-4965-bf7d-036be545a92e"),
		Name: "Gregory",
	}

	updatedTenant, err := existingTenant.Update()
	if err != nil {
		t.Errorf("Error test updating the tenant: %v\n", err)
		return
	}

	assert.Equal(t, updatedTenant.ID, existingTenant.ID)
	assert.Equal(t, updatedTenant.Name, existingTenant.Name)
	t.Log("End TestUpdateTenant")
}

func TestDeleteTenant(t *testing.T) {
	err := refreshTenantTable()
	if err != nil {
		t.Fatal(err)
		return
	}

	seedOneTenant()
	tenantInstance := TenantModel{
		Uuid: libUuid.MustParse("39b0b2fc-749f-46f3-8960-453418e72b2e"),
	}

	isDeleted, err := tenantInstance.Delete()
	if err != nil {
		t.Errorf("Error test deleting the tenant: %v\n", err)
		return
	}

	assert.Equal(t, isDeleted, true)
	t.Log("End TestDeleteTenant")
}
