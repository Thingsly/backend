package initialize

import (
	"fmt"
	"log"
	"os"
	"time"

	global "github.com/Thingsly/backend/pkg/global"
	utils "github.com/Thingsly/backend/pkg/utils"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type DbConfig struct {
	Host          string
	Port          int
	DbName        string
	Username      string
	Password      string
	TimeZone      string
	LogLevel      int
	SlowThreshold int
	IdleConns     int
	OpenConns     int
}

func PgInit() (*gorm.DB, error) {

	config, err := LoadDbConfig()
	if err != nil {
		logrus.Errorf("Failed to load database configuration: %v", err)
		return nil, err
	}

	db, err := PgConnect(config)
	if err != nil {
		logrus.Errorf("Failed to connect to database: %v", err)
		return nil, err
	}
	global.DB = db

	CasbinInit()

	err = CheckVersion(db)
	if err != nil {
		fmt.Println(err)
	}

	return db, nil
}

func LoadDbConfig() (*DbConfig, error) {
	config := &DbConfig{
		Host:          viper.GetString("db.psql.host"),
		Port:          viper.GetInt("db.psql.port"),
		DbName:        viper.GetString("db.psql.dbname"),
		Username:      viper.GetString("db.psql.username"),
		Password:      viper.GetString("db.psql.password"),
		TimeZone:      viper.GetString("db.psql.time_zone"),
		LogLevel:      viper.GetInt("db.psql.log_level"),
		SlowThreshold: viper.GetInt("db.psql.slow_threshold"),
		IdleConns:     viper.GetInt("db.psql.idle_conns"),
		OpenConns:     viper.GetInt("db.psql.open_conns"),
	}

	if config.Host == "" {
		config.Host = "localhost"
	}
	if config.Port == 0 {
		config.Port = 5432
	}
	if config.TimeZone == "" {
		config.TimeZone = "Asia/Ho_Chi_Minh"
	}
	if config.LogLevel == 0 {
		config.LogLevel = 1
	}
	if config.SlowThreshold == 0 {
		config.SlowThreshold = 200
	}
	if config.IdleConns == 0 {
		config.IdleConns = 10
	}
	if config.OpenConns == 0 {
		config.OpenConns = 50
	}

	if config.DbName == "" || config.Username == "" || config.Password == "" {
		return nil, fmt.Errorf("database configuration is incomplete")
	}

	return config, nil
}

// type Writer struct{}

// func (w Writer) Printf(format string, args ...interface{}) {
// 	log.Println(args...)
// }

func PgConnect(config *DbConfig) (*gorm.DB, error) {
	dataSource := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable TimeZone=%s",
		config.Host, config.Port, config.DbName, config.Username, config.Password, config.TimeZone)

	newLogger := logger.New(
		//Writer{},
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Duration(config.SlowThreshold) * time.Millisecond,
			LogLevel:                  logger.LogLevel(config.LogLevel),
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		})

	var err error
	db, err := gorm.Open(postgres.Open(dataSource), &gorm.Config{
		Logger:                 newLogger,
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: false, // use singular table name, table for `User` would be `user` with this option enabled
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get raw database connection: %v", err)
	}

	sqlDB.SetMaxIdleConns(config.IdleConns)
	sqlDB.SetMaxOpenConns(config.OpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Database connection established successfully...")

	return db, nil
}

/*
Note: Do not have the sys_version table in the SQL
1. Check if the version table exists: Check the database version. If the sys_version table does not exist, create the sys_version table and insert version number 0, version 0.0.0.
2. Program version is lower than data version: Prompt for an upgrade.
3. Data version is lower than program version: Execute the SQL files to update the version.
*/
// Check version in the 'sys_version' table's 'version' field
func CheckVersion(db *gorm.DB) error {
	version := global.VERSION
	versionNumber := global.VERSION_NUMBER // Current program version number
	var dataVersionNumber int              // Database version number

	// Check if the sys_version table exists
	var exists bool
	result := db.Raw("SELECT EXISTS(SELECT 1 FROM information_schema.tables WHERE table_schema='public' AND table_name='sys_version')").Scan(&exists)
	if result.Error != nil {
		return result.Error
	}
	// Start transaction
	logrus.Info("----", exists)
	if !exists { // If the sys_version table does not exist, create the sys_version table
		logrus.Info("Creating sys_version table")
		dataVersionNumber = 0
		t := db.Exec("CREATE TABLE sys_version (version_number INT NOT NULL DEFAULT 0, version varchar(255) NOT NULL, PRIMARY KEY (version_number))")
		if t.Error != nil {
			return t.Error
		}
	}
	tx := db.Begin()
	// Query version number
	result = db.Table("sys_version").Select("version_number").Scan(&dataVersionNumber)
	if result.Error != nil {
		return result.Error
	}
	// If version number is empty, insert the version number
	if dataVersionNumber == 0 {
		t := tx.Exec("INSERT INTO sys_version (version_number, version) VALUES (?, ?)", 0, "0.0.0")
		if t.Error != nil {
			// Rollback
			tx.Rollback()
			return t.Error
		}
	}
	if dataVersionNumber > global.VERSION_NUMBER {
		// Rollback
		tx.Rollback()
		return fmt.Errorf("Current data version is higher than the program version. Please upgrade the program.")
	} else if dataVersionNumber < global.VERSION_NUMBER {
		log.Println("Data version:", dataVersionNumber)
		log.Println("Program version:", global.VERSION_NUMBER)
		log.Println("Starting upgrade...")
		// SQL file names are in the format: version_number.sql, execute SQL files for versions greater than the current data version and less than or equal to the program version
		for i := dataVersionNumber + 1; i <= global.VERSION_NUMBER; i++ {
			fileName := fmt.Sprintf("sql/%d.sql", i)
			// Check if the file exists
			if !utils.FileExist(fileName) {
				// Rollback
				tx.Rollback()
				return fmt.Errorf("SQL file does not exist. Manual upgrade may be required: %s", fileName)
			}
			log.Println("Executing SQL file:", fileName)
			// Read SQL script file
			sqlFile, err := os.ReadFile(fileName)
			if err != nil {
				panic(err)
			}
			fmt.Println("Executing SQL script...")
			// Execute SQL script
			t := tx.Exec(string(sqlFile))
			if t.Error != nil {
				// Rollback
				tx.Rollback()
				return t.Error
			}
		}
		// Update version number
		t := tx.Exec("UPDATE sys_version SET version_number = ?, version = ?", versionNumber, version)
		if t.Error != nil {
			// Rollback
			tx.Rollback()
			return t.Error
		}
		log.Println("Upgrade successful")
	}
	return tx.Commit().Error
}

func ExecuteSQLFile(db *gorm.DB, fileName string) error {
	sqlFile, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}
	t := db.Exec(string(sqlFile))
	if t.Error != nil {
		return t.Error
	}

	return nil
}
