package tenant

import (
	libUuid "github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gitlab.com/s0j0hn/go-rest-boilerplate-echo/database"
	"gorm.io/gorm"
	"log"
	"os"
	"testing"
)

var DbClient *gorm.DB
var validTenantID = "6fcec554-9861-4965-bf7d-036be545a92e"

func TestMain(m *testing.M) {
	DbClient = database.ConnectForTests()
	err := refreshTenantTable()
	if err != nil {
		log.Fatal(err)
		return
	}

	os.Exit(m.Run())
}

func refreshTenantTable() (err error) {
	err = DbClient.Exec("DROP TABLE IF EXISTS tenant").Error
	if err != nil {
		log.Fatalf("Cannot refresh tenant table: %v", err)
		return err
	}

	err = DbClient.AutoMigrate(&Model{})
	if err != nil {
		return err
	}

	return nil
}

func seedTenants() error {
	tenants := []Model{
		{
			Name: "Bob",
		},
		{
			Name: "Alice",
		},
	}

	for i := range tenants {
		tenant, err := tenants[i].Save()
		if err != nil {
			return err
		}
		log.Printf("Seed with Tenant ID: %s", tenant.UUID)
	}
	return nil
}

func seedOneTenant() {
	err := refreshTenantTable()
	if err != nil {
		log.Fatalf("Cannot seed tenant table: %v", err)
		return
	}

	tenant := Model{
		Name: "Greg",
		UUID: libUuid.MustParse(validTenantID),
	}

	tenantSaved, err := tenant.Save()
	if err != nil {
		log.Fatalf("Cannot seed tenant table: %v", err)
		return
	}
	log.Printf("Seed with Tenant ID: %s", tenantSaved.UUID)
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

	tenantInstance := Model{}

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

	newUser := Model{
		UUID: libUuid.New(),
		Name: "Test",
	}

	savedUser, err := newUser.Save()
	if err != nil {
		t.Errorf("this is the error getting the users: %v\n", err)
		return
	}

	assert.Equal(t, newUser.ID, savedUser.ID)
	assert.Equal(t, newUser.Name, savedUser.Name)
	assert.Equal(t, newUser.UUID, savedUser.UUID)
	t.Log("End TestSaveTenant")
}

func TestWrongSaveTenant(t *testing.T) {
	err := refreshTenantTable()
	if err != nil {
		t.Fatal(err)
		return
	}

	newUser := Model{
		UUID: libUuid.New(),
		Name: "",
	}

	savedUser, err := newUser.Save()
	if assert.Error(t, err) {
		assert.Nil(t, savedUser)
		assert.Equal(t, "NOT NULL constraint failed: tenant.name", err.Error())
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
	tenantInstance := Model{UUID: libUuid.MustParse(validTenantID)}

	foundTenant, err := tenantInstance.GetOne()
	assert.NoError(t, err)
	assert.Equal(t, foundTenant.UUID.String(), validTenantID)
	t.Log("End TestTenantGetById")
}

func TestGetWrongTenantByID(t *testing.T) {
	err := refreshTenantTable()
	if err != nil {
		t.Fatal(err)
		return
	}

	seedOneTenant()
	tenantInstance := Model{UUID: libUuid.New()}

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

	existingTenant := Model{
		UUID: libUuid.MustParse(validTenantID),
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
	tenantInstance := Model{
		UUID: libUuid.MustParse(validTenantID),
	}

	isDeleted, err := tenantInstance.Delete()
	if err != nil {
		t.Errorf("Error test deleting the tenant: %v\n", err)
		return
	}

	if assert.NoError(t, err) {
		assert.Equal(t, true, isDeleted)
		t.Log("End TestDeleteTenant")
	}

	//err = refreshTenantTable()
	//if err != nil {
	//	t.Fatal(err)
	//	return
	//}
	//
	//tenantInstance = Model{
	//	UUID: libUuid.New(),
	//}
	//
	//transaction := database.Connect().Begin()
	//
	//isDeleted, err = tenantInstance.Delete()
	//
	//if assert.Error(t, err) {
	//	assert.NoError(t, transaction.Error)
	//	assert.Equal(t, "no uuid specified", err.Error())
	//	assert.Equal(t, isDeleted, false)
	//	t.Log("End TestDeleteTenant")
	//}
}
