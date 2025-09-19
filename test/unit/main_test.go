package unit

import (
	"testing"

	dbhelper "github.com/Manuel-Leleuly/kanban-flow-go/helpers/db"
	testhelper "github.com/Manuel-Leleuly/kanban-flow-go/helpers/test"
	"github.com/Manuel-Leleuly/kanban-flow-go/models"
)

var D *models.DBInstance = dbhelper.NewDBClient()

func TestMain(m *testing.M) {
	if !D.IsDBConnected() {
		if err := testhelper.ConnectToTestDB(D); err != nil {
			panic("[Error] failed to connect to test db due to: " + err.Error())
		}
	}

	// delete all data
	if err := testhelper.DeleteAllTestTickets(D); err != nil {
		panic("[Error] failed to delete all test tickets before running test due to: " + err.Error())
	}

	if err := testhelper.DeleteAllTestUsers(D); err != nil {
		panic("[Error] failed to delete all test users before running test due to: " + err.Error())
	}

	// create all data
	if err := testhelper.CreateTestUser(D); err != nil {
		panic("[Error] failed to create test user due to: " + err.Error())
	}

	if err := testhelper.CreateTestTicket(D); err != nil {
		panic("[Error] failed to create test ticket due to: " + err.Error())
	}

	m.Run()

	if err := testhelper.DeleteAllTestTickets(D); err != nil {
		panic("[Error] failed to delete all test tickets after running test due to: " + err.Error())
	}

	if err := testhelper.DeleteAllTestUsers(D); err != nil {
		panic("[Error] failed to delete all test users after running test due to: " + err.Error())
	}
}
