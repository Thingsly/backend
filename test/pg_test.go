package test

import (
	"os"
	"testing"
	"time"

	"github.com/HustIoTPlatform/backend/initialize"

	"github.com/HustIoTPlatform/backend/internal/query"

	"github.com/HustIoTPlatform/backend/internal/model"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

var adr = func(s string) *string { return &s }
var config *initialize.DbConfig
var db *gorm.DB

func TestDatebase(t *testing.T) {
	testConnect(t)
	testDDLInit(t)
	testNotificationGroup(t)
}

func testConnect(t *testing.T) {
	require := require.New(t)
	if os.Getenv("run_env") == "git-actions" {
		initialize.ViperInit("../configs/conf-push-test.yml")
	} else if os.Getenv("run_env") == "localdev" {
		initialize.ViperInit("../configs/conf-localdev.yml")
	} else {
		t.Log("Unknown environment")
		return
	}
	var err error
	config, err = initialize.LoadDbConfig()
	require.Nil(err)
	db, err = initialize.PgConnect(config)
	require.Nil(err)
}

func testDDLInit(t *testing.T) {
	require := require.New(t)

	res := db.Exec("DROP SCHEMA public CASCADE;CREATE SCHEMA public;")
	require.Nil(res.Error)

	db, err := initialize.PgConnect(config)
	require.Nil(err)

	// ts := db.Exec("CREATE TABLE sys_version (version_number INT NOT NULL DEFAULT 0, version varchar(255) NOT NULL, PRIMARY KEY (version_number))")
	// err = ts.Error
	// require.Nilf(err,"CREATE TABLE sys_version error %v",err)

	err = initialize.ExecuteSQLFile(db, "../sql/1.sql")
	require.Nilf(err, "Failed to execute DDL: %v", err)

	require.Nilf(err, "Failed to commit DDL: %v", err)
	t.Log("Database initialization succeeded")

}

func testNotificationGroup(t *testing.T) {
	require := require.New(t)
	require.NotNil(db, "Database connection failed")
	query.SetDefault(db)

	notificationGroup := model.NotificationGroup{
		Name:               "test",
		NotificationType:   "MEMBER",
		Status:             "ON",
		NotificationConfig: adr("{}"),
		Description:        adr("test"),
		TenantID:           "123456",
		Remark:             adr("test"),
		CreatedAt:          time.Now().UTC(),
		UpdatedAt:          time.Now().UTC(),
	}
	err := query.NotificationGroup.Create(&notificationGroup)
	require.Nil(err, "Failed to create notificationGroup")
	db.Commit()
}
