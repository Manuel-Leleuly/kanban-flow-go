package main

import (
	"os"

	dbhelper "github.com/Manuel-Leleuly/kanban-flow-go/helpers/db"
	"github.com/Manuel-Leleuly/kanban-flow-go/initializer"
	"github.com/Manuel-Leleuly/kanban-flow-go/routes"
	"github.com/sirupsen/logrus"
)

func init() {
	initializer.LoadEnvVariables()
}

func main() {
	// check if client secret is set
	if os.Getenv("CLIENT_SECRET") == "" {
		logrus.Fatal("[Error] client secret is not set")
	}

	db := dbhelper.NewDBClient()

	if err := db.ConnectToDB(os.Getenv("DB_NAME")); err != nil {
		logrus.Fatal("[Error] failed to connect to db due to: " + err.Error())
	}

	if err := db.SyncDatabase(); err != nil {
		logrus.Fatal("[Error] failed to sync database due to: " + err.Error())
	}

	server := routes.GetRoutes(db)
	if err := server.Run(":3005"); err != nil {
		logrus.Fatal("[Error] failed to start Gin server due to: " + err.Error())
	}
}
