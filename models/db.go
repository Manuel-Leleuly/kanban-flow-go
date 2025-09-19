package models

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBInstance struct {
	*gorm.DB
}

func (d *DBInstance) ConnectToDB(dbName string) error {
	if dbName == "" {
		return errors.New("DB name is empty")
	}

	dialect := postgres.Open(getGORMDatabaseUrl(dbName, map[string]string{
		"ssl_mode":        "required",
		"channel_binding": "required",
	}))

	gormConfig := &gorm.Config{}

	if strings.ToLower(os.Getenv("ENABLE_DB_LOGGER")) == "true" {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	}

	db, err := gorm.Open(dialect, gormConfig)
	if err != nil {
		return err
	}

	postgresDB, err := db.DB()
	if err != nil {
		return err
	}

	postgresDB.SetMaxOpenConns(100)
	postgresDB.SetMaxIdleConns(10)
	postgresDB.SetConnMaxLifetime(30 * time.Minute)
	postgresDB.SetConnMaxIdleTime(5 * time.Minute)

	d.DB = db

	return nil
}

func (d *DBInstance) IsDBConnected() bool {
	return d.DB != nil
}

func (d *DBInstance) SyncDatabase() error {
	if !d.IsDBConnected() {
		return errors.New("DB is not initialized")
	}

	d.DB.AutoMigrate(&User{}, &Ticket{})

	return nil
}

func (d *DBInstance) CheckDBConnection(c *gin.Context) {
	if !d.IsDBConnected() {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"Message": "DB is not initialized",
		})
	} else {
		c.Next()
	}
}

func (d *DBInstance) MakeHTTPHandleFunc(f func(db *DBInstance, c *gin.Context)) gin.HandlerFunc {
	return func(c *gin.Context) {
		f(d, c)
	}
}

// helpers
func getGORMDatabaseUrl(dbName string, params map[string]string) string {
	dbUrl := fmt.Sprintf("postgresql://%s:%s@%s/%s", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_DOMAIN"), dbName)

	queryParams := url.Values{}
	for k, v := range params {
		queryParams.Add(k, v)
	}

	return dbUrl + "?" + queryParams.Encode()
}
