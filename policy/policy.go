package policy

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/rbac/default-role-manager"
	"github.com/casbin/casbin/v2/util"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
	"log"
)

// InitPolicy is used to initialise all policy manager needs to function.
func InitPolicy(gormClient *gorm.DB) (*casbin.Enforcer, error) {
	// Initialize a Gorm adapter and use it in a Casbin enforcer:
	// The adapter will use the MySQL database named "casbin".
	// If it doesn't exist, the adapter will create it automatically.
	casbinGormAdapter, err := gormadapter.NewAdapterByDB(gormClient) // Your driver and data source.
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// We need a new roleManager to force the route params verification ex: (/tenants/:id).
	var roleManager = defaultrolemanager.NewRoleManager(2)
	roleManager.AddMatchingFunc("KeyMatch2", util.KeyMatch2)

	// Create Policy enforcer with our customized model.
	policyEnforcer, err := casbin.NewEnforcer("config/keymatch_model", casbinGormAdapter)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	policyEnforcer.SetRoleManager(roleManager)

	// Logs for casbin.
	policyEnforcer.EnableLog(true)

	// Load the policy from DB.
	err = policyEnforcer.LoadPolicy()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// Save the policy back to DB.
	err = policyEnforcer.SavePolicy()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return policyEnforcer, nil
}

// AddCreatePolicy is used to add policy for specified user.
func AddCreatePolicy(policyEnforcer *casbin.Enforcer, user string, url string) {
	isAdded, err := policyEnforcer.AddPolicy(user, url, "POST")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Added policy create with result: %t", isAdded)
}

// AddUpdatePolicy is used to add policy for specified user.
func AddUpdatePolicy(policyEnforcer *casbin.Enforcer, user string, url string) {
	isAdded, err := policyEnforcer.AddPolicy(user, url, "PUT")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Added policy update with result: %t", isAdded)
}

// AddDeletePolicy is used to add policy for specified user.
func AddDeletePolicy(policyEnforcer *casbin.Enforcer, user string, url string) {
	isAdded, err := policyEnforcer.AddPolicy(user, url+"/:id", "DELETE")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Added policy delete with result: %t", isAdded)
}

// AddGetPolicy is used to add policy for specified user.
func AddGetPolicy(policyEnforcer *casbin.Enforcer, user string, url string) {
	isAdded, err := policyEnforcer.AddPolicy(user, url, "GET")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Added policy get with result: %t", isAdded)

	AddGetByIDPolicy(policyEnforcer, user, url+"/:id")
}

// AddGetByIDPolicy is used to add policy for specified user.
func AddGetByIDPolicy(policyEnforcer *casbin.Enforcer, user string, url string) {
	isAdded, err := policyEnforcer.AddPolicy(user, url, "GET")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Added policy by id with result: %t", isAdded)
}
